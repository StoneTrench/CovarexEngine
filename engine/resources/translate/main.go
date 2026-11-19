package translate

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/StoneTrench/CovarexEngine/engine/platform"
)

var language_fallback_map map[string]string = make(map[string]string)
var language_map map[string]string = make(map[string]string)

func Translate(key string, args ...any) (text string, err error, warn error) {
	if format, exists := language_map[key]; exists {
		return fmt.Sprintf(format, args...), nil, nil
	}
	if format, exists := language_fallback_map[key]; exists {
		return fmt.Sprintf(format, args...), nil, fmt.Errorf("using fallback for key %s", key)
	}
	return key, fmt.Errorf("no translation found for key %s", key), nil
}
func TryTranslate(key string, args ...any) string {
	value, err, warn := Translate(key, args...)
	if err != nil {
		slog.Error(err.Error())
	}
	if warn != nil {
		slog.Warn(warn.Error())
	}
	return value
}

func append_lang(fs *platform.VFS, path string, m map[string]string) error {
	data, err := fs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open json language file: %w", err)
	}

	lang := make(map[string]string)
	err = json.Unmarshal(data, &lang)
	if err != nil {
		return fmt.Errorf("failed to parse json language file: %w", err)
	}
	for k, v := range lang {
		if other_v, exists := m[k]; exists {
			slog.Warn(fmt.Sprintf("overlapping language keys (%s) \"%s\" and \"%s\"", k, v, other_v))
		}
		m[k] = v
	}

	return nil
}
func Append(fs *platform.VFS, path string) error {
	return append_lang(fs, path, language_map)
}
func AppendFallback(fs *platform.VFS, path string) error {
	return append_lang(fs, path, language_fallback_map)
}
func ClearAll() {
	language_map = make(map[string]string)
	language_fallback_map = make(map[string]string)
}
