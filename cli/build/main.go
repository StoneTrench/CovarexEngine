package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func Echo(f string, args ...any) {
	fmt.Printf("[Build] %s\n", fmt.Sprintf(f, args...))
}

type Config struct {
	DistDir     string
	OutputDir   string
	ProjectName string
	TargetOS    string
	TargetARCH  string
	Release     bool
}

func main() {
	c := parseFlags()

	Echo("Collecting info")

	platform, err := resolveEngineMeta()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	tag, commit := gitDescribe()
	built := time.Now().UTC().Format("2006-01-02 15:04:05")

	if err := os.MkdirAll(c.OutputDir, 0o755); err != nil {
		os.Exit(1)
	}

	output := "main"
	switch c.TargetOS {
	case "windows":
		output += ".exe"
	case "js", "wasip1":
		output += ".wasm"
	}
	output = filepath.Join(c.OutputDir, output)

	env := append(os.Environ(),
		"GOOS="+c.TargetOS,
		"GOARCH="+c.TargetARCH,
	)
	if c.TargetARCH == "wasm" {
		env = append(env, "CGO_ENABLED=0")
	}

	Echo("Version: v%s %s (%s) %s/%s", tag, commit, built, c.TargetOS, c.TargetARCH)

	pairs := [][2]string{
		{"BUILD_NAME", c.ProjectName},
		{"BUILD_TAG", tag},
		{"BUILD_BUILT", built},
		{"BUILD_COMMIT", commit},
		{"BUILD_RELEASE", strconv.FormatBool(c.Release)},
		{"BUILD_TARGET_OS", c.TargetOS},
		{"BUILD_TARGET_ARCH", c.TargetARCH},
	}
	var b strings.Builder
	for _, p := range pairs {
		fmt.Fprintf(&b, " -X %q", platform+"."+p[0]+"="+p[1])
	}
	ldflags := strings.TrimSpace(b.String())

	Echo("Building: %s", output)

	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", ldflags, "-o", output, ".")
	cmd.Env = env
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}

	// 3. Package
	if !c.Release {
		return
	}
	Echo("Compressing release")
	archive := filepath.Join(c.DistDir, fmt.Sprintf("%s_%s-%s-%s.tar.gz",
		strings.ReplaceAll(c.ProjectName, " ", "_"), tag, c.TargetOS, c.TargetARCH))
	if err := tarball(c.OutputDir, archive); err != nil {
		fmt.Fprintln(os.Stderr, "package failed:", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	var (
		distributable = flag.String("dist", "_dist", "directory to put the released project")
		output        = flag.String("out", "_out", "directory to put the built project")
		projectname   = flag.String("name", "Covarex Project", "application name")
		targetos      = flag.String("os", runtime.GOOS, "windows|linux|darwin|js|wasip1")
		targetarch    = flag.String("arch", runtime.GOARCH, "amd64|wasm|*")
		release       = flag.Bool("release", false, "if omitted build is in debug mode")
	)
	flag.Parse()

	switch *targetos {
	case "windows", "linux", "darwin", "js", "wasip1":
	default:
		fmt.Fprintf(os.Stderr, "error: unsupported os %q\n", *targetos)
		os.Exit(2)
	}

	c := &Config{
		DistDir:     *distributable,
		OutputDir:   *output,
		ProjectName: *projectname,
		TargetOS:    *targetos,
		TargetARCH:  *targetarch,
		Release:     *release,
	}

	return c
}

// resolveEngineMeta asks the Go toolchain what the engine metadata module is.
func resolveEngineMeta() (string, error) {
	out, err := exec.Command("go", "list", "-f", "{{.ImportPath}}", "github.com/StoneTrench/CovarexEngine/engine/platform").Output()
	if err != nil {
		return "", fmt.Errorf("resolve engine platform: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
func gitDescribe() (tag, commit string) {
	if _, err := exec.Command("git", "rev-parse", "--verify", "HEAD").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return "0.0.0-untracked", "no-commit"
	}
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		tag = "0.0.0"
	} else {
		tag = strings.TrimSpace(string(out))
	}
	out, _ = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	commit = strings.TrimSpace(string(out))
	return tag, commit
}

func tarball(dir, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	gz := gzip.NewWriter(f)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		hdr, _ := tar.FileInfoHeader(info, "")
		hdr.Name = rel
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tw.Write(data)
		return err
	})
}
