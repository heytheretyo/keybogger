package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	hook "github.com/robotn/gohook"
)

var (
	// keyPressCount    int
	// mouseClickCount  int
	// totalMouseTravel float64
	lastX, lastY int16
	firstMove    bool
	// this is my dpi => check this site to suit it https://www.mouse-sensitivity.com/dpianalyzer/
	dpi            = 475.47
	inchesToMeters = 0.0254
	db             *sql.DB
	saveInterval   = time.Second * 45
	bufferStats    = BufferStats{
		LeftClicks:   0,
		RightClicks:  0,
		MiddleClicks: 0,
		Keystrokes:   0,
		MouseTravel:  0.0,
		ScrollWheels: 0,
	}
)

type EventBucket struct {
	LeftClicks   int
	RightClicks  int
	MiddleClicks int
	Keystrokes   int
	MouseTravel  float64
	ScrollWheels int
}

var bucket EventBucket

func track_events() {
	hook.Register(hook.KeyDown, []string{"ctrl", "space", "b"}, func(e hook.Event) {
		fmt.Println("stopping tracking...")
		hook.End()
	})

	evChan := hook.Start()
	defer hook.End()

	ticker := time.NewTicker(saveInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			save_bucket()
		}
	}()

	for ev := range evChan {
		if ev.Kind == hook.KeyDown || ev.Kind == hook.MouseDown || ev.Kind == hook.MouseWheel || ev.Kind == hook.MouseMove {
			if ev.Kind == hook.KeyDown {
				bucket.Keystrokes++
				bufferStats.Keystrokes++
			} else if ev.Kind == hook.MouseDown {
				switch ev.Button {
				case 1:
					bucket.LeftClicks++
					bufferStats.LeftClicks++
				case 2:
					bucket.MiddleClicks++
					bufferStats.MiddleClicks++
				case 3:
					bucket.RightClicks++
					bufferStats.RightClicks++
				}
			} else if ev.Kind == hook.MouseWheel {
				bucket.ScrollWheels++
				bufferStats.ScrollWheels++
			} else if ev.Kind == hook.MouseMove {
				if firstMove {
					distance := calculate_distance(lastX, lastY, ev.X, ev.Y)
					bucket.MouseTravel += distance
					bufferStats.MouseTravel += distance
				} else {
					firstMove = true
				}
				lastX, lastY = ev.X, ev.Y
			}
		}
	}
}

func main() {
	firstMove = false
	initialize_database()
	defer db.Close()

	go track_events()

	spawn_gui()
}
