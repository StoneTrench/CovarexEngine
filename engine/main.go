package engine

import (
	"fmt"

	"github.com/StoneTrench/CovarexEngine/engine/core/log"
	"github.com/StoneTrench/CovarexEngine/engine/platform"
)

func Init(logDir string) error {
	var err error = nil

	err = platform.InitMeta()
	if err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	err = log.Init(logDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}

	return nil
}
