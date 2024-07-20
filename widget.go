package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func spawn_gui() {
	a := app.New()
	w := a.NewWindow("keylogger v1")

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("open",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	keystrokesLabel := widget.NewLabel("keystrokes: 0")
	leftClicksLabel := widget.NewLabel("left mouse clicks: 0")
	rightClicksLabel := widget.NewLabel("right mouse clicks: 0")
	middleClicksLabel := widget.NewLabel("middle mouse clicks: 0")
	mouseTravelledLabel := widget.NewLabel("total mouse travelled: 0.0")
	scrollWheelsLabel := widget.NewLabel("scroll wheels: 0")

	updateStats := func(stats BufferStats) {
		keystrokesLabel.SetText(fmt.Sprintf("keystrokes: %d", stats.Keystrokes))
		leftClicksLabel.SetText(fmt.Sprintf("left mouse clicks: %d", stats.LeftClicks))
		rightClicksLabel.SetText(fmt.Sprintf("right mouse clicks: %d", stats.RightClicks))
		middleClicksLabel.SetText(fmt.Sprintf("middle mouse clicks: %d", stats.MiddleClicks))
		mouseTravelledLabel.SetText(fmt.Sprintf("total mouse travelled: %.3f meters", stats.MouseTravel))
		scrollWheelsLabel.SetText(fmt.Sprintf("scroll wheels: %d", stats.ScrollWheels))
	}

	updateStats(bufferStats)

	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()

		for range ticker.C {
			updateStats(bufferStats)
		}
	}()

	// go func() {
	// 	ticker := time.NewTicker(time.Second * 30)
	// 	defer ticker.Stop()

	// 	for range ticker.C {
	// 		bufferStats = load_daily_stats()
	// 	}
	// }()

	content := container.NewVBox(
		keystrokesLabel,
		leftClicksLabel,
		rightClicksLabel,
		middleClicksLabel,
		mouseTravelledLabel,
		scrollWheelsLabel,
	)

	w.SetContent(content)
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	w.ShowAndRun()
}
