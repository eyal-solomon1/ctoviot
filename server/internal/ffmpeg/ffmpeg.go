package ffmpeg

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ffmpeG "github.com/u2takey/ffmpeg-go"
)

type FFMPEG interface {
	CreateAudioFile(filePath string) (string, error)
	MatchTranscriptionsToVideo(filePath string, subtitlesPath string) (string, error)
	GetFileDuration(filePath string) (float64, error)
}

type FFMPEGService struct {
}

func Initialize() FFMPEG {
	return &FFMPEGService{}
}

// CreateAudioFile recives a filePath to the video, and creates a .mp3 (audio) version of it
// Returning the "end result" audio file path
func (ffmpeg FFMPEGService) CreateAudioFile(filePath string) (string, error) {
	basePath := strings.Split(filePath, ".")
	outputFilePath := basePath[0] + ".mp3"
	err := ffmpeG.Input(filePath).
		Output(outputFilePath, ffmpeG.KwArgs{"b:a": "392K"}).
		OverWriteOutput().Run()

	if err != nil {
		return "", err
	}

	return outputFilePath, nil
}

// MatchTranscriptionsToVideo recives a filePath to the video, and a file path to the subs to match
// It returns the "matched" video file path
func (ffmpeg FFMPEGService) MatchTranscriptionsToVideo(filePath string, subtitlesPath string) (string, error) {
	basePath := strings.Split(filePath, ".")
	outputFilePath := basePath[0] + "-ready" + ".mp4"

	err := ffmpeG.Input(filePath).
		Output(outputFilePath, ffmpeG.KwArgs{"vf": fmt.Sprintf("subtitles=%s:force_style='Alignment=2,MarginV=50", subtitlesPath)}).OverWriteOutput().Run()

	if err != nil {
		return "", err
	}
	return outputFilePath, err
}

// GetFileDuration recives a filePath to the video and returns the it's length in secods
func (ffmpeg FFMPEGService) GetFileDuration(filePath string) (float64, error) {
	result, err := ffmpeG.Probe(filePath)
	if err != nil {
		return 0, err
	}

	var data fileOutput

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(data.Format.Duration, 64)

	if err != nil {
		return 0, err
	}
	return duration, nil
}

type fileOutput struct {
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}
