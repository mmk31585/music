package audioinfo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Info struct {
	Format   string
	Duration float64
	Bitrate  int
}

var mp3BitrateTable = [5][4][16]int{
	{},
	{
		{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0},
		{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0},
		{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0},
	},
	{
		{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0},
		{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},
	},
	{
		{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0},
		{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},
	},
}

var mp3SampleRateTable = [4][4]int{
	{44100, 48000, 32000, 0},
	{22050, 24000, 16000, 0},
	{11025, 12000, 8000, 0},
}

func Extract(r io.ReadSeeker) (*Info, error) {
	header := make([]byte, 10)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, fmt.Errorf("cannot read header: %w", err)
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek error: %w", err)
	}

	if bytes.HasPrefix(header, []byte("ID3")) {
		return extractMP3WithID3(r, header)
	}

	if bytes.HasPrefix(header, []byte("fLaC")) {
		return extractFLAC(r)
	}

	if bytes.HasPrefix(header, []byte("OggS")) {
		return extractOGG(r)
	}

	if len(header) >= 8 && (bytes.Equal(header[4:8], []byte("ftyp")) || bytes.Equal(header[4:8], []byte("moov"))) {
		return extractMP4(r)
	}

	if header[0] == 0xFF && (header[1]&0xE0) == 0xE0 {
		return extractMP3Raw(r)
	}

	return nil, fmt.Errorf("unsupported or unknown audio format")
}

func extractMP3WithID3(r io.ReadSeeker, header []byte) (*Info, error) {
	tagSize := int(header[6])<<21 | int(header[7])<<14 | int(header[8])<<7 | int(header[9])
	tagSize += 10

	if _, err := r.Seek(int64(tagSize), io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek past ID3 tag: %w", err)
	}

	return parseMPEG(r, tagSize)
}

func extractMP3Raw(r io.ReadSeeker) (*Info, error) {
	return parseMPEG(r, 0)
}

func parseMPEG(r io.ReadSeeker, offset int) (*Info, error) {
	frameHeader := make([]byte, 4)
	if _, err := io.ReadFull(r, frameHeader); err != nil {
		return nil, fmt.Errorf("cannot read MPEG frame header: %w", err)
	}

	if frameHeader[0] != 0xFF || (frameHeader[1]&0xE0) != 0xE0 {
		return nil, fmt.Errorf("invalid MPEG sync word")
	}

	h := uint32(frameHeader[0])<<24 | uint32(frameHeader[1])<<16 | uint32(frameHeader[2])<<8 | uint32(frameHeader[3])

	version := int((h >> 19) & 0x3)
	layer := int((h >> 17) & 0x3)
	bitrateIdx := int((h >> 12) & 0xF)
	sampleRateIdx := int((h >> 10) & 0x3)
	padding := int((h >> 9) & 0x1)

	if version == 0 || version == 1 || layer == 0 {
		return nil, fmt.Errorf("unsupported MPEG version/layer")
	}

	verIdx := mpegVersionIndex(version)
	layerIdx := mpegLayerIndex(layer)

	bitrate := mp3BitrateTable[verIdx][layerIdx-1][bitrateIdx]
	sampleRate := mp3SampleRateTable[verIdx][sampleRateIdx]

	if bitrate == 0 || sampleRate == 0 {
		return nil, fmt.Errorf("invalid bitrate or sample rate in MPEG header")
	}

	frameLen := mpegFrameLength(layer, bitrate, sampleRate, padding)
	if frameLen == 0 {
		return nil, fmt.Errorf("could not compute MPEG frame length")
	}

	xingBuf := make([]byte, 4)
	if _, err := r.Read(xingBuf); err != nil {
		return nil, fmt.Errorf("cannot read Xing header: %w", err)
	}

	if string(xingBuf) == "Xing" || string(xingBuf) == "Info" {
		return parseXingHeader(r, sampleRate, bitrate, frameLen)
	}

	totalSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seek end: %w", err)
	}

	audioSize := totalSize - int64(offset)
	if audioSize <= 0 {
		return nil, fmt.Errorf("no audio data found")
	}

	duration := float64(audioSize) * 8 / float64(bitrate*1000)

	return &Info{
		Format:   "mp3",
		Duration: duration,
		Bitrate:  bitrate,
	}, nil
}

func mpegVersionIndex(version int) int {
	switch version {
	case 3:
		return 1
	case 2:
		return 2
	default:
		return 3
	}
}

func mpegLayerIndex(layer int) int {
	switch layer {
	case 3:
		return 1
	case 2:
		return 2
	case 1:
		return 3
	default:
		return 0
	}
}

func mpegFrameLength(layer, bitrate, sampleRate, padding int) int {
	switch layer {
	case 1:
		return (12*bitrate*1000/sampleRate + padding) * 4
	case 2, 3:
		return 144*bitrate*1000/sampleRate + padding
	default:
		return 0
	}
}

func parseXingHeader(r io.ReadSeeker, sampleRate, bitrate, frameLen int) (*Info, error) {
	hdrBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, hdrBuf); err != nil {
		return nil, fmt.Errorf("reading Xing flags: %w", err)
	}

	flags := binary.BigEndian.Uint32(hdrBuf)
	frames := 0
	if flags&1 != 0 {
		frameBuf := make([]byte, 4)
		if _, err := io.ReadFull(r, frameBuf); err != nil {
			return nil, fmt.Errorf("reading Xing frame count: %w", err)
		}
		frames = int(binary.BigEndian.Uint32(frameBuf))
	}

	if frames > 0 && sampleRate > 0 {
		duration := float64(frames) * float64(frameLen) * 8 / float64(bitrate*1000)
		return &Info{
			Format:   "mp3",
			Duration: duration,
			Bitrate:  bitrate,
		}, nil
	}

	totalSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seek end: %w", err)
	}

	duration := float64(totalSize) * 8 / float64(bitrate*1000)

	return &Info{
		Format:   "mp3",
		Duration: duration,
		Bitrate:  bitrate,
	}, nil
}

func extractFLAC(r io.ReadSeeker) (*Info, error) {
	streamInfo := make([]byte, 42)
	if _, err := io.ReadFull(r, streamInfo); err != nil {
		return nil, fmt.Errorf("cannot read FLAC STREAMINFO: %w", err)
	}

	if string(streamInfo[:4]) != "fLaC" {
		return nil, fmt.Errorf("invalid FLAC header")
	}

	blockHeader := streamInfo[4:8]
	blockType := blockHeader[0] & 0x7F

	if blockType != 0 {
		return nil, fmt.Errorf("first block is not STREAMINFO")
	}

	info := streamInfo[8:42]
	sampleRate := int(binary.BigEndian.Uint32([]byte{info[10] & 0x0F, info[11], info[12], info[13]}) >> 12)
	totalSamples := int64(binary.BigEndian.Uint64(append([]byte{0, 0}, info[13:19]...)) & 0x0FFFFFFFFFFFFF)

	if sampleRate == 0 {
		return nil, fmt.Errorf("invalid FLAC sample rate")
	}

	var duration float64
	if totalSamples > 0 {
		duration = float64(totalSamples) / float64(sampleRate)
	}

	return &Info{
		Format:   "flac",
		Duration: duration,
	}, nil
}

func extractOGG(r io.ReadSeeker) (*Info, error) {
	ident := make([]byte, 58)
	if _, err := io.ReadFull(r, ident); err != nil {
		return nil, fmt.Errorf("cannot read OGG identification header: %w", err)
	}

	if string(ident[:4]) != "OggS" {
		return nil, fmt.Errorf("invalid OGG header")
	}

	if len(ident) < 58 {
		return nil, fmt.Errorf("OGG header too short")
	}

	sampleRate := int(binary.LittleEndian.Uint32(ident[50:54]))
	if sampleRate == 0 {
		return nil, fmt.Errorf("invalid OGG sample rate")
	}

	totalSize, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seek end: %w", err)
	}

	estDuration := float64(totalSize) / float64(sampleRate) * 2

	return &Info{
		Format:   "ogg",
		Duration: estDuration,
	}, nil
}

func extractMP4(r io.ReadSeeker) (*Info, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek start: %w", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading MP4: %w", err)
	}

	fileSize := len(data)

	sampleRate, duration := findMoovDuration(data)

	if sampleRate == 0 {
		totalSize, err := r.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, fmt.Errorf("seek end: %w", err)
		}
		estDuration := float64(totalSize) / 44100 * 2
		return &Info{
			Format:   "m4a",
			Duration: estDuration,
			Bitrate:  int(float64(fileSize) * 8 / estDuration / 1000),
		}, nil
	}

	durSec := float64(duration) / float64(sampleRate)
	bitrate := 0
	if durSec > 0 {
		bitrate = int(float64(fileSize) * 8 / durSec / 1000)
	}

	return &Info{
		Format:   "m4a",
		Duration: durSec,
		Bitrate:  bitrate,
	}, nil
}

func findMoovDuration(data []byte) (sampleRate int, duration int) {
	i := 0
	for i < len(data)-8 {
		if i+8 > len(data) {
			break
		}
		boxSize := int(binary.BigEndian.Uint32(data[i : i+4]))
		if boxSize < 8 {
			break
		}
		boxType := string(data[i+4 : i+8])

		switch boxType {
		case "moov":
			// Recurse into the moov box body (skip 8-byte header) to find mvhd
			if i+8 >= len(data) {
				return 0, 0
			}
			return findMoovDuration(data[i+8:])
		case "mvhd":
			if i+20 <= len(data) {
				version := data[i+8]
				if version == 0 {
					if i+24 <= len(data) {
						duration = int(binary.BigEndian.Uint32(data[i+20 : i+24]))
					}
					sampleRate = int(binary.BigEndian.Uint32(data[i+16 : i+20]))
				} else {
					if i+36 <= len(data) {
						duration = int(binary.BigEndian.Uint32(data[i+28 : i+32]))
					}
					sampleRate = int(binary.BigEndian.Uint32(data[i+24 : i+28]))
				}
			}
			return sampleRate, duration
		}

		i += boxSize
	}
	return 0, 0
}
