package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

const transcodeTimeout = 5 * time.Minute

// TranscodeToHLS 用 ffmpeg 把音频转成 HLS 分片（playlist.m3u8 + ts），参数与原 Java 版一致
func TranscodeToHLS(inputPath, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建 HLS 目录失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), transcodeTimeout)
	defer cancel()

	playlist := filepath.Join(outputDir, "playlist.m3u8")
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-i", inputPath,
		"-vn",
		"-c:a", "aac",
		"-b:a", "128k",
		"-hls_time", "10",
		"-hls_list_size", "0",
		playlist,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(out)
		if len(msg) > 500 {
			msg = msg[len(msg)-500:]
		}
		return fmt.Errorf("ffmpeg 转码失败: %w, 输出: %s", err, msg)
	}
	return nil
}

// ProbeDuration 用 ffprobe 读取音频时长（秒），失败返回 0
func ProbeDuration(inputPath string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		inputPath,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	var probe struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if json.Unmarshal(out, &probe) != nil {
		return 0
	}
	seconds, _ := strconv.ParseFloat(probe.Format.Duration, 64)
	return int(seconds)
}
