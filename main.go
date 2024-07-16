package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	_ "github.com/mattn/go-sqlite3"
	hook "github.com/robotn/gohook"
)

var (
	// keyPressCount    int
	// mouseClickCount  int
	// totalMouseTravel float64
	lastX, lastY     int16
	firstMove        bool
    // this is my dpi => check this site to suit it https://www.mouse-sensitivity.com/dpianalyzer/
	dpi              = 475.47
	inchesToMeters   = 0.0254
	db               *sql.DB
	saveInterval     = time.Second * 10
)

type EventBucket struct {
	LeftClicks   int
	RightClicks  int
	MiddleClicks int
	Keystrokes   int
	MouseTravel   float64
	ScrollWheels  int
}

var bucket EventBucket

func initialize_database() {
	var err error
	db, err = sql.Open("sqlite3", "events.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS event_counts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		left_clicks INTEGER DEFAULT 0,
		right_clicks INTEGER DEFAULT 0,
		middle_clicks INTEGER DEFAULT 0,
		keystrokes INTEGER DEFAULT 0,
		mouse_travel_distance REAL DEFAULT 0,
		scroll_wheel_movements INTEGER DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO event_counts (left_clicks, right_clicks, middle_clicks, keystrokes, mouse_travel_distance, scroll_wheel_movements) VALUES (0, 0, 0, 0, 0, 0)`)
	if err != nil {
		log.Println("error initializing counts:", err)
	}
}

func save_bucket() {
	_, err := db.Exec(`INSERT INTO event_counts (left_clicks, right_clicks, middle_clicks, keystrokes, mouse_travel_distance, scroll_wheel_movements) VALUES (?, ?, ?, ?, ?, ?)`,
		bucket.LeftClicks,
		bucket.RightClicks,
		bucket.MiddleClicks,
		bucket.Keystrokes,
		bucket.MouseTravel,
		bucket.ScrollWheels,
	)
	if err != nil {
		log.Println("error saving bucket:", err)
	}
	bucket = EventBucket{} // reset the bucket
}

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
				keyPressCount++
				bucket.Keystrokes++
				// fmt.Printf("key pressed: %v (Total: %d)\n", ev.Rawcode, keyPressCount)
			} else if ev.Kind == hook.MouseDown {
				mouseClickCount++
				switch ev.Button {
				case 1:
					bucket.LeftClicks++
					// fmt.Printf("left mouse button pressed: %v (Total: %d)\n", ev.Button, mouseClickCount)
				case 2:
					bucket.MiddleClicks++
					// fmt.Printf("middle mouse button pressed: %v (Total: %d)\n", ev.Button, mouseClickCount)
				case 3:
					bucket.RightClicks++
					// fmt.Printf("right mouse button pressed: %v (Total: %d)\n", ev.Button, mouseClickCount)
				}
			} else if ev.Kind == hook.MouseWheel {
				bucket.ScrollWheels++
				// fmt.Printf("scroll direction: (Total: %d)\n", ev.Rotation)
			} else if ev.Kind == hook.MouseMove {
				if firstMove {
					distance := calculate_distance(lastX, lastY, ev.X, ev.Y)
					totalMouseTravel += distance
					bucket.MouseTravel += distance
				} else {
					firstMove = true
				}
				lastX, lastY = ev.X, ev.Y
				// fmt.Printf("mouse movement: (X: %v, Y: %d) Total Distance: %f meters\n", ev.X, ev.Y, totalMouseTravel)
			}
		}
	}
}

func calculate_distance(x1, y1, x2, y2 int16) float64 {
	distanceInPixels := math.Sqrt(math.Abs(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))))
	distanceInInches := distanceInPixels / dpi
	distanceInMeters := distanceInInches * inchesToMeters
	return distanceInMeters
}

func main() {
	firstMove = false
	initialize_database()
	defer db.Close()

	track_events()
}
