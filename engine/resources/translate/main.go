package translate

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type LanguageSeverity uint8
const (
	LanguageSeverity_Warning LanguageSeverity = iota
	LanguageSeverity_Error
)

func NewLanguageIssue(op string, msg string, severity LanguageSeverity) *LanguageIssue {
	return &LanguageIssue{
		Op:       op,
		Message:  msg,
		Severity: severity,
	}
}

type LanguageIssue struct {
	Op       string
	Message  string
	Severity LanguageSeverity
}

func (t *LanguageIssue) Error() string {
	return fmt.Sprintf("%s: %s", t.Op, t.Message)
}

var language_fallback_map map[string]string = make(map[string]string)
var language_map map[string]string = make(map[string]string)

// Translate tries to translate using the main map, and uses the fallback if it is missing, otherwise returns an error issue.
func Translate(key string, args ...any) (text string, issue *LanguageIssue) {
	if format, exists := language_map[key]; exists {
		return fmt.Sprintf(format, args...), nil
	}
	if format, exists := language_fallback_map[key]; exists {
		return fmt.Sprintf(format, args...), NewLanguageIssue("translate", fmt.Sprintf("using fallback translation for key %s", key), LanguageSeverity_Warning)
	}
	return key, NewLanguageIssue("translate", fmt.Sprintf("missing translation for key %s", key), LanguageSeverity_Error)
}

// TranslateLog translates the key, and instead of returning any issues, just directly emits them using slog.
func TranslateLog(key string, args ...any) string {
	value, err := Translate(key, args...)
	switch err.Severity {
	case LanguageSeverity_Warning:
		slog.Warn(err.Error())
	case LanguageSeverity_Error:
		slog.Error(err.Error())
	}
	return value
}

func append_lang(json_data []byte, m map[string]string) error {
	lang := make(map[string]string)
	err := json.Unmarshal(json_data, &lang)
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
func Append(json_data []byte) error {
	return append_lang(json_data, language_map)
}
func AppendFallback(json_data []byte) error {
	return append_lang(json_data, language_fallback_map)
}

// ClearAll clears the fallback map and the main map, after which everything has to be re-appended.
func ClearAll() {
	language_map = make(map[string]string)
	language_fallback_map = make(map[string]string)
}
