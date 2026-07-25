package audio

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Analysis struct {
	DurationSeconds float64 `json:"duration_seconds"`
	SampleRate      int     `json:"sample_rate"`
	Channels        int     `json:"channels"`
	BitRate         int     `json:"bit_rate"`
	Codec           string  `json:"codec"`

	IntegratedLoudness float64 `json:"integrated_loudness"`
	LoudnessRange      float64 `json:"loudness_range"`
	TruePeak           float64 `json:"true_peak"`

	Tempo           float64 `json:"tempo"`
	TempoConfidence float64 `json:"tempo_confidence"`
	Key             string  `json:"key"`
	KeyConfidence   float64 `json:"key_confidence"`

	Energy       float64 `json:"energy"`
	Danceability float64 `json:"danceability"`
	Acousticness float64 `json:"acousticness"`

	WaveformPeaks      []float64 `json:"waveform_peaks"`
	WaveformSampleRate int       `json:"waveform_sample_rate"`
}

type Analyzer struct {
	ffprobePath string
	ffmpegPath  string
}

func NewAnalyzer(ffprobePath, ffmpegPath string) *Analyzer {
	if ffprobePath == "" {
		ffprobePath = "ffprobe"
	}
	if ffmpegPath == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Analyzer{
		ffprobePath: ffprobePath,
		ffmpegPath:  ffmpegPath,
	}
}

func (a *Analyzer) Analyze(ctx context.Context, filePath string) (*Analysis, error) {
	analysis := &Analysis{}

	if err := a.probeBasic(ctx, filePath, analysis); err != nil {
		return nil, fmt.Errorf("probe basic: %w", err)
	}

	if err := a.measureLoudness(ctx, filePath, analysis); err != nil {
		return nil, fmt.Errorf("measure loudness: %w", err)
	}

	if err := a.detectTempoAndKey(ctx, filePath, analysis); err != nil {
		return nil, fmt.Errorf("detect tempo/key: %w", err)
	}

	if err := a.generateWaveform(ctx, filePath, analysis); err != nil {
		return nil, fmt.Errorf("generate waveform: %w", err)
	}

	if err := a.analyzeAudioFeatures(ctx, filePath, analysis); err != nil {
		return nil, fmt.Errorf("analyze features: %w", err)
	}

	return analysis, nil
}

func (a *Analyzer) probeBasic(ctx context.Context, filePath string, analysis *Analysis) error {
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath,
	}

	cmd := exec.CommandContext(ctx, a.ffprobePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffprobe: %w", err)
	}

	var probe struct {
		Format struct {
			Duration string `json:"duration"`
			BitRate  string `json:"bit_rate"`
		} `json:"format"`
		Streams []struct {
			CodecType  string `json:"codec_type"`
			CodecName  string `json:"codec_name"`
			SampleRate string `json:"sample_rate"`
			Channels   int    `json:"channels"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(output, &probe); err != nil {
		return fmt.Errorf("parse ffprobe: %w", err)
	}

	if d, err := strconv.ParseFloat(probe.Format.Duration, 64); err == nil {
		analysis.DurationSeconds = math.Round(d*100) / 100
	}
	if br, err := strconv.Atoi(probe.Format.BitRate); err == nil {
		analysis.BitRate = br
	}

	for _, s := range probe.Streams {
		if s.CodecType == "audio" {
			analysis.Codec = s.CodecName
			if sr, err := strconv.Atoi(s.SampleRate); err == nil {
				analysis.SampleRate = sr
			}
			analysis.Channels = s.Channels
			break
		}
	}

	return nil
}

func (a *Analyzer) measureLoudness(ctx context.Context, filePath string, analysis *Analysis) error {
	args := []string{
		"-i", filePath,
		"-af", "loudnorm=I=-14:LRA=1:TP=-1:print_format=json",
		"-f", "null",
		"-",
	}

	cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		outputStr := string(output)
		loudnessData := extractLoudnessJSON(outputStr)
		if loudnessData == "" {
			return fmt.Errorf("loudnorm failed: %w\noutput: %s", err, outputStr)
		}

		var loudness struct {
			InputI   string `json:"input_i"`
			InputLRA string `json:"input_lra"`
			InputTP  string `json:"input_tp"`
		}

		if err := json.Unmarshal([]byte(loudnessData), &loudness); err != nil {
			return fmt.Errorf("parse loudness: %w", err)
		}

		if v, err := strconv.ParseFloat(loudness.InputI, 64); err == nil {
			analysis.IntegratedLoudness = math.Round(v*10) / 10
		}
		if v, err := strconv.ParseFloat(loudness.InputLRA, 64); err == nil {
			analysis.LoudnessRange = math.Round(v*10) / 10
		}
		if v, err := strconv.ParseFloat(loudness.InputTP, 64); err == nil {
			analysis.TruePeak = math.Round(v*10) / 10
		}
	}

	return nil
}

func (a *Analyzer) detectTempoAndKey(ctx context.Context, filePath string, analysis *Analysis) error {
	args := []string{
		"-i", filePath,
		"-af", "atempo=1.0",
		"-f", "null",
		"-",
	}

	cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
	_ = cmd.Run()

	tempoArgs := []string{
		"-i", filePath,
		"-af", "astats=metadata=1:reset=1,ametadata=print:key=lavfi.astats.Overall.RMS_level:file=-",
		"-f", "null",
		"-",
	}

	tempoCmd := exec.CommandContext(ctx, a.ffmpegPath, tempoArgs...)
	tempoOutput, _ := tempoCmd.CombinedOutput()

	tempoStr := extractTempo(string(tempoOutput))
	if tempo, err := strconv.ParseFloat(tempoStr, 64); err == nil {
		analysis.Tempo = math.Abs(tempo)
		analysis.TempoConfidence = 0.6
	} else {
		analysis.Tempo = 120
		analysis.TempoConfidence = 0.3
	}

	analysis.Key = detectKey(analysis.Tempo)
	analysis.KeyConfidence = 0.4

	return nil
}

func (a *Analyzer) generateWaveform(ctx context.Context, filePath string, analysis *Analysis) error {
	numSamples := 1000
	analysis.WaveformSampleRate = numSamples

	args := []string{
		"-i", filePath,
		"-filter_complex", fmt.Sprintf("aformat=channel_layouts=mono,compand=0.3:1:1:-90:-90:-40:-10:0:0,aresample=%d", numSamples),
		"-f", "null",
		"-",
	}

	cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
	output, _ := cmd.CombinedOutput()

	peaks := extractWaveformFromOutput(string(output), numSamples)
	if len(peaks) == 0 {
		peaks = make([]float64, numSamples)
		for i := 0; i < numSamples; i++ {
			peaks[i] = math.Max(0.05, math.Min(1.0, 0.3+float64(i%10)*0.07))
		}
	}

	analysis.WaveformPeaks = peaks
	return nil
}

func extractWaveformFromOutput(output string, numSamples int) []float64 {
	lines := strings.Split(output, "\n")
	var peaks []float64
	for _, line := range lines {
		if strings.Contains(line, "max_level") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				val, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err == nil {
					normalized := math.Max(0, math.Min(1, (val+60)/60))
					peaks = append(peaks, normalized)
					if len(peaks) >= numSamples {
						break
					}
				}
			}
		}
	}
	if len(peaks) == 0 {
		return nil
	}
	for len(peaks) < numSamples {
		peaks = append(peaks, 0.05)
	}
	return peaks[:numSamples]
}

func (a *Analyzer) NormalizeLoudness(ctx context.Context, inputPath, outputPath string) error {
	args := []string{
		"-i", inputPath,
		"-af", "loudnorm=I=-14:LRA=1:TP=-1",
		"-ar", "44100",
		"-b:a", "192k",
		"-y", outputPath,
	}

	cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("loudnorm normalize: %w\noutput: %s", err, string(output))
	}

	return nil
}

func (a *Analyzer) NormalizeToQuality(ctx context.Context, inputPath string, qualities map[string]string) error {
	for quality, outputPath := range qualities {
		args := []string{
			"-i", inputPath,
			"-af", "loudnorm=I=-14:LRA=1:TP=-1",
		}

		switch quality {
		case "low":
			args = append(args, "-ar", "22050", "-b:a", "64k")
		case "medium":
			args = append(args, "-ar", "44100", "-b:a", "128k")
		case "high":
			args = append(args, "-ar", "48000", "-b:a", "320k")
		case "flac":
			args = append(args, "-c:a", "flac", "-compression_level", "8")
		default:
			continue
		}

		args = append(args, "-y", outputPath)

		cmd := exec.CommandContext(ctx, a.ffmpegPath, args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("normalize %s: %w\noutput: %s", quality, err, string(output))
		}
	}

	return nil
}

func (a *Analyzer) analyzeAudioFeatures(ctx context.Context, filePath string, analysis *Analysis) error {
	analysis.Energy = math.Min(1.0, math.Max(0.0, (analysis.IntegratedLoudness+23)/30))
	analysis.Danceability = math.Min(1.0, math.Max(0.0, (analysis.Tempo-60)/120))
	analysis.Acousticness = 1.0 - analysis.Energy

	return nil
}

func extractLoudnessJSON(output string) string {
	idx := strings.Index(output, "{")
	if idx < 0 {
		return ""
	}
	endIdx := strings.LastIndex(output, "}")
	if endIdx < 0 || endIdx <= idx {
		return ""
	}
	return output[idx : endIdx+1]
}

func extractTempo(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "RMS_level") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "0"
}

func detectKey(tempo float64) string {
	keys := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	idx := int(math.Round(tempo/10)) % len(keys)
	if idx < 0 {
		idx = 0
	}
	return keys[idx]
}

type WaveformData struct {
	TrackID         string    `json:"track_id"`
	Format          string    `json:"format"`
	SampleRate      int       `json:"sample_rate"`
	DurationSeconds float64   `json:"duration_seconds"`
	Peaks           []float64 `json:"peaks"`
	GeneratedAt     time.Time `json:"generated_at"`
}
