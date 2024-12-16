// internal/expander/expander.go
package expander

import (
	"strings"
	"time"

	"github.com/noble-varghese/expandr/internal/config"
	"github.com/noble-varghese/expandr/internal/typing"
)

type Expander struct {
	config       *config.Config
	typingConfig typing.TypingConfig
}

func New(cfg *config.Config) *Expander {
	return &Expander{
		config:       cfg,
		typingConfig: typing.DefaultConfig,
	}
}

func (e *Expander) SetTypingSpeed(config typing.TypingConfig) {
	e.typingConfig = config
}

func (e *Expander) Expand(text string) (string, string, bool) {
	if expansion, ok := e.config.Expansions[text]; ok {
		return e.processTemplate(expansion), text, true
	}
	for key, expansion := range e.config.Expansions {
		if strings.HasSuffix(text, key) {
			// logger.Infof("Found expansion for pattern '%s'", key)
			return e.processTemplate(expansion), key, true
		}
	}
	return "", "", false
}

func (e *Expander) processTemplate(text string) string {
	// Replace common template patterns
	replacements := map[string]string{
		"{{.Now}}":     time.Now().Format(time.RFC3339),
		"{{.Date}}":    time.Now().Format("2006-01-02"),
		"{{.Time}}":    time.Now().Format("15:04:05"),
		"{{.Year}}":    time.Now().Format("2006"),
		"{{.Month}}":   time.Now().Format("January"),
		"{{.Day}}":     time.Now().Format("Monday"),
		"{{.Utc_Now}}": time.Now().UTC().Format(time.RFC3339),
	}

	result := text
	for pattern, replacement := range replacements {
		result = strings.ReplaceAll(result, pattern, replacement)
	}

	return result
}

func (e *Expander) TypeExpansion(text string) {
	typing.TypeMultilineTextWithConfig(text, e.typingConfig)
}
