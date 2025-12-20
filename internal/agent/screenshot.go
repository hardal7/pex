package agent

import (
	"fmt"
	"image/png"
	"log/slog"
	"os"

	"github.com/kbinani/screenshot"
)

func CaptureScreen() []string {
	var screenshots []string

	n := screenshot.NumActiveDisplays()

	for i := range n {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			slog.Info("Failed capturing screenshot")
			return []string{}
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		file, _ := os.Create(fileName)
		defer file.Close()
		png.Encode(file, img)

		screenshots[i] = fileName
		slog.Info("Captured screenshot: #%d : %v \"%s\"\n", i, bounds, fileName)
	}

	return screenshots
}
