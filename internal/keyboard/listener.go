package keyboard

import (
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/noble-varghese/scribe/internal/logger"
	"github.com/noble-varghese/scribe/internal/types"
	hook "github.com/robotn/gohook"
)

type Listener struct {
	expander    types.Expander // Using the interface instead of concrete type
	buffer      []string
	bufferMutex sync.Mutex
	running     bool
	stopChan    chan struct{}
}

func NewListener(exp types.Expander) *Listener {
	return &Listener{
		expander: exp,
		buffer:   make([]string, 0, 100),
		stopChan: make(chan struct{}),
	}
}

func (l *Listener) Start() error {
	l.running = true

	// Start listening for events
	evChan := hook.Start()

	go func() {
		for l.running {
			select {
			case <-l.stopChan:
				return
			case ev := <-evChan:
				// We only care about key press events
				if ev.Kind != hook.KeyDown {
					continue
				}

				// Get the character representation of the key
				key := hook.RawcodetoKeychar(ev.Rawcode)
				// logger.Infof("Key event - Keycode: %v, Rawcode: %v, Char: %s", ev.Keycode, ev.Rawcode, key)

				// Handle special keys based on their keycodes
				// logger.Info(hook.RawcodetoKeychar(ev.Rawcode))
				switch ev.Rawcode {
				case 51: // Common backspace keycode
					l.removeLastChar()
				case 36: // Common enter keycode
					l.processKeypress("\n")
				case 32: // Common space keycode
					l.processKeypress(" ")
				case 48: // Common tab keycode
					l.processKeypress("\t")
				default:
					// For regular characters, only process if we have a valid character
					if len(key) == 1 {
						// Convert to lowercase for consistent matching
						l.processKeypress(strings.ToLower(key))
					}
				}
			}
		}
	}()

	return nil
}

func (l *Listener) Stop() error {
	l.running = false
	close(l.stopChan)
	return nil
}

func (l *Listener) processKeypress(key string) {
	l.bufferMutex.Lock()
	defer l.bufferMutex.Unlock()
	l.buffer = append(l.buffer, key)
	currentText := strings.Join(l.buffer, "")

	if expanded, trigger, ok := l.expander.Expand(currentText); ok {
		backspaceLen := len(trigger)
		for i := 0; i < backspaceLen; i++ {
			robotgo.KeyTap("backspace")
		}
		l.expander.TypeExpansion(expanded)
		l.buffer = l.buffer[:len(l.buffer)-backspaceLen]
	}
}

func (l *Listener) removeLastChar() {
	l.bufferMutex.Lock()
	defer l.bufferMutex.Unlock()

	if len(l.buffer) > 0 {
		l.buffer = l.buffer[:len(l.buffer)-1]
	}
}

func (l *Listener) cleanBuffer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.bufferMutex.Lock()
			defer l.bufferMutex.Unlock()

			if len(l.buffer) > 50 { // If buffer gets too long, clear it
				l.buffer = l.buffer[:0]
				logger.Info("Cleared keyboard buffer due to length")
			}
			return
		case <-l.stopChan:
			return
		}
	}
}
