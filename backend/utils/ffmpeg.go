package utils

import (
	"fmt"
	"os/exec"
)

// IsFFmpegAvailable checks if ffmpeg is installed and available in the system PATH.
func IsFFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

type AudioMetadata struct {
	Title      string
	Artist     string
	Album      string
	Year       string
	Genre      string
	ArtworkURL string // Used if coverPath is not provided
}

// ConvertToMp3 converts an audio file to MP3 with ID3 tags and optional cover art.
// It uses -q:a 0 for best variable bitrate quality (approx 220-260kbps).
func ConvertToMp3(inputPath, outputPath, coverPath string, meta AudioMetadata) error {
	// ffmpeg -i input.m4a -i cover.jpg -map 0:a -map 1:0 -c:a libmp3lame -q:a 0 -id3v2_version 3
	// -metadata title="..." -metadata artist="..." -metadata album="..." output.mp3

	args := []string{
		"-y", // Overwrite output
		"-i", inputPath,
	}

	// Add cover art if provided
	if coverPath != "" {
		args = append(args, "-i", coverPath)
		args = append(args, "-map", "0:a")
		args = append(args, "-map", "1:0")
		args = append(args, "-c:v", "copy") // Copy the image data
		// Ensure it's attached as front cover
		args = append(args, "-disposition:v:0", "attached_pic")
	} else {
		// Just map audio
		args = append(args, "-map", "0:a")
	}

	args = append(args,
		"-c:a", "libmp3lame",
		"-q:a", "0", // Best quality VBR
		"-id3v2_version", "3", // Widely compatible ID3v2.3
		"-metadata", fmt.Sprintf("title=%s", meta.Title),
		"-metadata", fmt.Sprintf("artist=%s", meta.Artist),
		"-metadata", fmt.Sprintf("album=%s", meta.Album),
	)

	if meta.Year != "" {
		args = append(args, "-metadata", fmt.Sprintf("date=%s", meta.Year))
	}
	if meta.Genre != "" {
		args = append(args, "-metadata", fmt.Sprintf("genre=%s", meta.Genre))
	}

	args = append(args, outputPath)

	cmd := exec.Command("ffmpeg", args...)
	// Capture output for debugging if needed, but for now just run it
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}

	return nil
}
