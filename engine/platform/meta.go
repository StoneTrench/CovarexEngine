package platform

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
)

type ENGINE_OS uint8

const (
	OS_WINDOWS ENGINE_OS = iota
	OS_LINUX
	OS_DARWIN
	OS_JS
	OS_WASIP1
)

type ENGINE_ARCH uint8

const (
	ARCH_AMD64 ENGINE_ARCH = iota
	ARCH_WASM
)

var (
	BUILD_NAME        string = "" // Do not read from these!
	BUILD_TAG         string = "" // Do not read from these!
	BUILD_BUILT       string = "" // Do not read from these!
	BUILD_COMMIT      string = "" // Do not read from these!
	BUILD_RELEASE     string = "" // Do not read from these!
	BUILD_TARGET_OS   string = "" // Do not read from these!
	BUILD_TARGET_ARCH string = "" // Do not read from these!
)

var (
	internal_is_initialized bool = false
	internal_name           string
	internal_tag            *semver.Version
	internal_built          time.Time
	internal_commit         string
	internal_release        bool
	internal_target         string
	internal_target_os      ENGINE_OS
	internal_target_arch    ENGINE_ARCH
)

func InitMeta() error {
	if internal_is_initialized {
		return fmt.Errorf("metadata already initialized")
	}
	if len(BUILD_NAME) == 0 {
		return fmt.Errorf("name was nil")
	}
	if len(BUILD_TAG) == 0 {
		return fmt.Errorf("git tag was nil")
	}
	if len(BUILD_COMMIT) == 0 {
		return fmt.Errorf("git commit was nil")
	}
	if len(BUILD_BUILT) == 0 {
		return fmt.Errorf("built timestamp was nil")
	}
	if len(BUILD_RELEASE) == 0 {
		return fmt.Errorf("engine release was nil")
	}
	if len(BUILD_TARGET_OS) == 0 {
		return fmt.Errorf("target os was nil")
	}
	if len(BUILD_TARGET_ARCH) == 0 {
		return fmt.Errorf("target arch was nil")
	}

	var err error
	internal_name = BUILD_NAME
	internal_tag, err = semver.NewVersion(BUILD_TAG)
	if err != nil {
		return fmt.Errorf("failed to parse version (build tag): %w", err)
	}
	internal_commit = BUILD_COMMIT
	internal_built, err = time.Parse(time.DateTime, BUILD_BUILT)
	if err != nil {
		return fmt.Errorf("failed to parse build time: %w", err)
	}
	switch BUILD_RELEASE {
	case "true":
		internal_release = true
	case "false":
		internal_release = false
	default:
		return fmt.Errorf("invalid engine mode: %s", BUILD_RELEASE)
	}
	internal_target = fmt.Sprintf("%s/%s", BUILD_TARGET_OS, BUILD_TARGET_ARCH)
	switch BUILD_TARGET_OS {
	case "windows":
		internal_target_os = OS_WINDOWS
	case "linux":
		internal_target_os = OS_LINUX
	case "darwin":
		internal_target_os = OS_DARWIN
	case "js":
		internal_target_os = OS_JS
	case "wasip1":
		internal_target_os = OS_WASIP1
	default:
		return fmt.Errorf("invalid os build target: %s", BUILD_TARGET_OS)
	}
	switch BUILD_TARGET_ARCH {
	case "amd64":
		internal_target_arch = ARCH_AMD64
	case "wasm":
		internal_target_arch = ARCH_WASM
	default:
		return fmt.Errorf("invalid arch build target: %s", BUILD_TARGET_ARCH)
	}

	BUILD_NAME = ""
	BUILD_TAG = ""
	BUILD_COMMIT = ""
	BUILD_BUILT = ""
	BUILD_RELEASE = ""
	BUILD_TARGET_OS = ""
	BUILD_TARGET_ARCH = ""

	internal_is_initialized = true
	return nil
}

func verifyIsInit() {
	if !internal_is_initialized {
		panic("metadata not initialized")
	}
}

func GetFormattedApplicationLabel() string {
	verifyIsInit()

	var suffix = "Release"
	if IsInDebugMode() {
		suffix = "Debug"
	}
	return fmt.Sprintf("%s v%s %s (%s) %s %s",
		internal_name,
		internal_tag,
		internal_commit,
		internal_built.Format(time.DateTime),
		internal_target,
		suffix,
	)
}

func GetTarget() string {
	verifyIsInit()
	return internal_target
}

func IsInDebugMode() bool {
	verifyIsInit()
	return !internal_release
}

func HasFileSystem() bool {
	return internal_target_os != OS_JS && internal_target_os != OS_WASIP1
}
