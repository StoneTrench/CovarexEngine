package engine

import (
	"fmt"

	"github.com/StoneTrench/CovarexEngine/engine/core/log"
	"github.com/StoneTrench/CovarexEngine/engine/platform"
)

// Init initializes the platform layer (loads and parses metadata 
// that was embeded during compilation), and initializes the logger.
func Init(virtFilesys *platform.VFS, logDir string) error {
	var err error = nil

	err = platform.InitMeta()
	if err != nil {
		return fmt.Errorf("failed to initialize metadata: %w", err)
	}

	err = log.Init(virtFilesys, logDir, false)
	if err != nil {
		return fmt.Errorf("failed to initialize logging: %w", err)
	}

	return nil
}
