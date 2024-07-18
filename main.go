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
	saveInterval   = time.Second * 10
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
	fmt.Println("--- press 'ctrl + space + b' to stop tracking ---")

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
			} else if ev.Kind == hook.MouseDown {
				switch ev.Button {
				case 1:
					bucket.LeftClicks++
				case 2:
					bucket.MiddleClicks++
				case 3:
					bucket.RightClicks++
				}
			} else if ev.Kind == hook.MouseWheel {
				bucket.ScrollWheels++
			} else if ev.Kind == hook.MouseMove {
				if firstMove {
					distance := calculate_distance(lastX, lastY, ev.X, ev.Y)
					bucket.MouseTravel += distance
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
