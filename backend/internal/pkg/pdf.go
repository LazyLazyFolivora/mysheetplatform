package pkg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const pdfConvertTimeout = 5 * time.Minute

var pageImagePattern = regexp.MustCompile(`^page-0*(\d+)\.png$`)

// ConvertPDFToImages 用 pdftoppm 把 PDF 逐页转成 1.png..N.png，返回页数
func ConvertPDFToImages(pdfPath, outputDir string) (int, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return 0, fmt.Errorf("创建图片目录失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pdfConvertTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pdftoppm",
		"-png",
		"-r", "120",
		pdfPath,
		filepath.Join(outputDir, "page"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(out)
		if len(msg) > 500 {
			msg = msg[len(msg)-500:]
		}
		return 0, fmt.Errorf("pdftoppm 转换失败: %w, 输出: %s", err, msg)
	}

	// pdftoppm 输出 page-01.png 这类零填充命名，统一重命名为 1.png..N.png
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return 0, fmt.Errorf("读取图片目录失败: %w", err)
	}

	count := 0
	for _, entry := range entries {
		match := pageImagePattern.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}
		pageNum, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		oldPath := filepath.Join(outputDir, entry.Name())
		newPath := filepath.Join(outputDir, fmt.Sprintf("%d.png", pageNum))
		if err := os.Rename(oldPath, newPath); err != nil {
			return 0, fmt.Errorf("重命名图片失败: %w", err)
		}
		count++
	}

	if count == 0 {
		return 0, fmt.Errorf("pdftoppm 未生成任何页面图片")
	}
	return count, nil
}
