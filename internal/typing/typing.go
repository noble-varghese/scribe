// internal/typing/typing.go
package typing

import (
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
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
	for _, char := range text {
		robotgo.TypeStr(string(char))
		time.Sleep(time.Duration(config.CharacterDelay) * time.Millisecond)
	}
}

// SafeNewlineWithConfig inserts a newline using configuration
func SafeNewlineWithConfig(config TypingConfig) {

	robotgo.KeyTap("enter", "shift")
	time.Sleep(time.Duration(config.ModifierDelay) * time.Millisecond)

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
