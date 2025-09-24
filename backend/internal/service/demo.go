package service

import (
	"fmt"
	log "insight/internal/pkg/logger"

	"go.uber.org/zap"
)

type DemoService interface {
	Demo(path string) (string, error)
}

type demoServiceImpl struct{}

func NewDemoService() DemoService {
	return &demoServiceImpl{}
}

func (s *demoServiceImpl) Demo(path string) (string, error) {
	log.Logger.Info("DemoService processed path",
		zap.String("path", path),
	)
	return fmt.Sprintf("path is %s", path), nil
}
