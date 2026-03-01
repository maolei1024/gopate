// Package i18n provides internationalization support for gopate CLI.
// It auto-detects the terminal language from environment variables
// and defaults to English.
package i18n

import (
	"fmt"
	"os"
	"strings"
)

var currentLang = "en"

func init() {
	currentLang = DetectLanguage()
}

// DetectLanguage detects the user's language from environment variables.
// Checks LANGUAGE, LC_ALL, LC_MESSAGES, LANG in order.
// Returns "zh" for Chinese, "en" for everything else.
func DetectLanguage() string {
	for _, env := range []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		val := os.Getenv(env)
		if val == "" || val == "C" || val == "POSIX" {
			continue
		}
		lang := strings.ToLower(val)
		if strings.HasPrefix(lang, "zh") {
			return "zh"
		}
		// Any non-empty, non-Chinese value means non-Chinese
		return "en"
	}
	return "en"
}

// SetLanguage manually sets the current language.
// Supported values: "en", "zh".
func SetLanguage(lang string) {
	lang = strings.ToLower(lang)
	if lang == "zh" || strings.HasPrefix(lang, "zh") {
		currentLang = "zh"
	} else {
		currentLang = "en"
	}
}

// GetLanguage returns the current language code.
func GetLanguage() string {
	return currentLang
}

// T returns the translated string for the given key.
// Falls back to English if the key is not found in the current language.
// Falls back to the key itself if not found in any language.
func T(key string) string {
	var msgs map[string]string
	if currentLang == "zh" {
		msgs = messagesZh
	} else {
		msgs = messagesEn
	}

	if val, ok := msgs[key]; ok {
		return val
	}
	// Fallback to English
	if val, ok := messagesEn[key]; ok {
		return val
	}
	return key
}

// Tf returns the translated and formatted string for the given key.
func Tf(key string, args ...any) string {
	return fmt.Sprintf(T(key), args...)
}
