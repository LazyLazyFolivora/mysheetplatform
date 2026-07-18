package service

import (
	"os/exec"

	"go.uber.org/zap"
)

type PDFService struct {
	logger *zap.Logger
}

func NewPDFService(logger *zap.Logger) *PDFService {
	return &PDFService{logger: logger}
}

func (s *PDFService) GetPageCount(filePath string) (int, error) {
	cmd := exec.Command("pdfcpu", "info", filePath)
	output, err := cmd.Output()
	if err != nil {
		s.logger.Warn("pdfcpu info failed, try fallback", zap.String("path", filePath), zap.Error(err))
		return 0, nil
	}
	_ = output
	return 0, nil
}
