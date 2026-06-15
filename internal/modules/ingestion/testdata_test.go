package ingestion

import (
	"bytes"
	"encoding/binary"
)

const (
	testTitle       = "Test Song Title"
	testArtist      = "Test Artist"
	testAlbum       = "Test Album"
	testAlbumArtist = "Test Album Artist"
	testGenre       = "Rock"
	testComposer    = "Test Composer"
	testComment     = "Test Comment"
	testLyrics      = "Test lyrics line 1\nTest lyrics line 2"
	testYear        = 2024
	testTrackNum    = 3
	testTrackTotal  = 12
	testDiscNum     = 1
	testDiscTotal   = 2
)

func buildTestMP3(hasTags bool, hasCover bool) []byte {
	buf := &bytes.Buffer{}

	if hasTags {
		id3Size := writeID3v24Tags(buf, hasCover)
		padding := make([]byte, 2048-id3Size)
		buf.Write(padding)
	}

	writeMPEGFrame(buf, 100)

	return buf.Bytes()
}

func buildNoTagMP3() []byte {
	buf := &bytes.Buffer{}

	for i := 0; i < 10; i++ {
		writeMPEGFrame(buf, 100)
	}

	return buf.Bytes()
}

func writeID3v24Tags(buf *bytes.Buffer, hasCover bool) int {
	frameData := &bytes.Buffer{}

	writeTextFrame(frameData, "TIT2", testTitle)
	writeTextFrame(frameData, "TPE1", testArtist)
	writeTextFrame(frameData, "TALB", testAlbum)
	writeTextFrame(frameData, "TPE2", testAlbumArtist)
	writeTextFrame(frameData, "TCON", testGenre)
	writeTextFrame(frameData, "TCOM", testComposer)
	writeCOMMFrame(frameData, testComment)
	writeUSLTFrame(frameData, testLyrics)
	writeTextFrame(frameData, "TDRC", formatInt(testYear))
	writeTrackDiscFrame(frameData, testTrackNum, testTrackTotal, testDiscNum, testDiscTotal)

	if hasCover {
		writeAPICFrame(frameData)
	}

	frameBytes := frameData.Bytes()

	frameSize := len(frameBytes)
	id3Size := 10 + frameSize

	id3Header := make([]byte, 10)
	copy(id3Header, []byte("ID3"))
	id3Header[3] = 4
	id3Header[4] = 0
	id3Header[6] = byte((frameSize >> 21) & 0x7F)
	id3Header[7] = byte((frameSize >> 14) & 0x7F)
	id3Header[8] = byte((frameSize >> 7) & 0x7F)
	id3Header[9] = byte(frameSize & 0x7F)

	buf.Write(id3Header)
	buf.Write(frameBytes)

	return id3Size
}

func writeTextFrame(buf *bytes.Buffer, frameID, value string) {
	encoded := append([]byte{3}, []byte(value)...)
	writeFrameHeader(buf, frameID, len(encoded))
	buf.Write(encoded)
}

func writeNumFrame(buf *bytes.Buffer, frameID string, value int) {
	s := &bytes.Buffer{}
	s.WriteByte(3)
	s.WriteString(formatInt(value))
	encoded := s.Bytes()
	writeFrameHeader(buf, frameID, len(encoded))
	buf.Write(encoded)
}

func writeTrackDiscFrame(buf *bytes.Buffer, trackNum, trackTotal, discNum, discTotal int) {
	writeTextFrame(buf, "TRCK", formatInt(trackNum)+"/"+formatInt(trackTotal))
	writeTextFrame(buf, "TPOS", formatInt(discNum)+"/"+formatInt(discTotal))
}

func formatInt(v int) string {
	if v == 0 {
		return "0"
	}
	digits := []byte{}
	for v > 0 {
		digits = append([]byte{byte('0' + v%10)}, digits...)
		v /= 10
	}
	return string(digits)
}

func writeCOMMFrame(buf *bytes.Buffer, comment string) {
	commBuf := &bytes.Buffer{}
	commBuf.WriteByte(3)
	commBuf.Write([]byte("eng"))
	commBuf.WriteByte(0)
	commBuf.Write([]byte(comment))

	commBytes := commBuf.Bytes()
	writeFrameHeader(buf, "COMM", len(commBytes))
	buf.Write(commBytes)
}

func writeUSLTFrame(buf *bytes.Buffer, lyrics string) {
	usltBuf := &bytes.Buffer{}
	usltBuf.WriteByte(3)
	usltBuf.Write([]byte("eng"))
	usltBuf.WriteByte(0)
	usltBuf.Write([]byte(lyrics))

	usltBytes := usltBuf.Bytes()
	writeFrameHeader(buf, "USLT", len(usltBytes))
	buf.Write(usltBytes)
}

func writeAPICFrame(buf *bytes.Buffer) {
	coverData := buildCoverArt()

	apicBuf := &bytes.Buffer{}
	apicBuf.WriteByte(3)
	apicBuf.Write([]byte("image/jpeg\x00"))
	apicBuf.WriteByte(3)
	apicBuf.Write([]byte("\x00"))
	apicBuf.Write(coverData)

	apicBytes := apicBuf.Bytes()
	writeFrameHeader(buf, "APIC", len(apicBytes))
	buf.Write(apicBytes)
}

func buildCoverArt() []byte {
	img := &bytes.Buffer{}
	img.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})
	img.Write([]byte("JFIF\x00"))
	for i := 0; i < 100; i++ {
		img.WriteByte(byte(i))
	}
	img.Write([]byte{0xFF, 0xD9})
	return img.Bytes()
}

func writeFrameHeader(buf *bytes.Buffer, frameID string, size int) {
	buf.Write([]byte(frameID))
	hdr := make([]byte, 4)
	binary.BigEndian.PutUint32(hdr, uint32(size))
	buf.Write(hdr)
	buf.Write([]byte{0, 0})
}

func writeMPEGFrame(buf *bytes.Buffer, frameCount uint32) {
	header := make([]byte, 4)
	header[0] = 0xFF
	header[1] = 0xFB
	header[2] = 0x90
	header[3] = 0x00
	buf.Write(header)

	sideInfo := make([]byte, 32)
	buf.Write(sideInfo)

	xingStart := 4 + 32
	xingOffset := buf.Len() - xingStart + 4

	xingBuf := &bytes.Buffer{}
	xingBuf.Write([]byte("Xing"))
	flags := uint32(0x0007)
	binary.BigEndian.PutUint32(xingBuf.Bytes()[4:8], flags)

	_ = xingOffset

	xingData := make([]byte, 12)
	copy(xingData[0:4], "Xing")
	binary.BigEndian.PutUint32(xingData[4:8], 0x0007)
	binary.BigEndian.PutUint32(xingData[8:12], frameCount)
	buf.Write(xingData)

	dummyData := make([]byte, 100)
	buf.Write(dummyData)
}

func buildMP3WithTagsAndCover() []byte {
	return buildTestMP3(true, true)
}

func buildMP3WithTagsNoCover() []byte {
	return buildTestMP3(true, false)
}
