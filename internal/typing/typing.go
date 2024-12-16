// internal/typing/typing.go
package typing

import (
	"github.com/go-vgo/robotgo"
	"strings"
	"time"
)

type TypingConfig struct {
	CharacterDelay int
	ModifierDelay  int
	OperationDelay int
}

var FastConfig = TypingConfig{
	CharacterDelay: 0,
	ModifierDelay:  5,
	OperationDelay: 10,
}

var DefaultConfig = TypingConfig{
	CharacterDelay: 0,
	ModifierDelay:  1,
	OperationDelay: 1,
}

// TypeWithConfig types text using specified timing configuration
func TypeWithConfig(text string, config TypingConfig) {
	robotgo.TypeStr(text)
}

// SafeNewlineWithConfig inserts a newline using configuration
func SafeNewlineWithConfig(config TypingConfig) {
	robotgo.KeyToggle("shift", "down")
	time.Sleep(time.Duration(config.ModifierDelay) * time.Millisecond)

	robotgo.KeyTap("enter")
	time.Sleep(time.Duration(config.ModifierDelay) * time.Millisecond)

	robotgo.KeyToggle("shift", "up")
	time.Sleep(time.Duration(config.OperationDelay) * time.Millisecond)
}

// TypeMultilineTextWithConfig types multiline text with specified config
func TypeMultilineTextWithConfig(text string, config TypingConfig) {
	lines := strings.Split(strings.TrimSpace(text), "\n")

	for i, line := range lines {
		TypeWithConfig(line, config)
		if i < len(lines)-1 {
			SafeNewlineWithConfig(config)
		}
	}
}
