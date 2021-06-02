package function

import (
	"fmt"
	"image/png"
	_ "image/png"
	"os"
	"testing"
	"time"
)

func getTestTargetEvent() EventInfo {
	jst, _ := time.LoadLocation(location)
	return EventInfo{
		Title: "1行目\n2行目\n@3行目",
		Time:  time.Date(2021, 1, 17, 17, 10, 0, 0, jst)}
}

func TestEventInfoZero(t *testing.T) {
	e := EventInfo{}
	if !e.IsZero() {
		t.Errorf("event=%+v", e)
	}
}

func TestDays(t *testing.T) {
	// 当日0時の1秒前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 0, 0, 0, 0, jst).Add(-1 * time.Second)

	e := getTestTargetEvent()
	days := e.DaysUntil(now)
	if days != 1 {
		t.Errorf("days=%d", days)
	}
}

func TestDays2(t *testing.T) {
	// 当日0時の1n秒後
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 0, 0, 0, 0, jst).Add(1 * time.Second)

	e := getTestTargetEvent()
	days := e.DaysUntil(now)
	if days != 0 {
		t.Errorf("days=%d", days)
	}
}

func TestHours(t *testing.T) {
	// 予定時刻の24時間30分前+1n秒後
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-1470*time.Minute + 1*time.Nanosecond)

	e := getTestTargetEvent()
	hours := e.HoursUntil(now)
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours2(t *testing.T) {
	// 予定時刻の23時間30分前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-1410 * time.Minute)

	e := getTestTargetEvent()
	hours := e.HoursUntil(now)
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours3(t *testing.T) {
	// 予定時刻の1時間前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-1 * time.Hour)

	e := getTestTargetEvent()
	hours := e.HoursUntil(now)
	if hours != 1 {
		t.Errorf("hours=%d", hours)
	}
}

func TestNearTargetDateTime1(t *testing.T) {
	// 予定時刻の59分前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-59 * time.Minute)

	e := getTestTargetEvent()
	near := e.NearTargetDateTime(now)
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime2(t *testing.T) {
	// 予定時刻の1分前+1n秒後
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-1*time.Minute + 1*time.Nanosecond)

	e := getTestTargetEvent()
	near := e.NearTargetDateTime(now)
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime3(t *testing.T) {
	// 予定時刻の4分59秒後
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(299 * time.Second)

	e := getTestTargetEvent()
	near := e.NearTargetDateTime(now)
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime4(t *testing.T) {
	// 予定時刻の5分後
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(5 * time.Minute)

	e := getTestTargetEvent()
	near := e.NearTargetDateTime(now)
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestCountdownText(t *testing.T) {
	// 予定時刻の100時間31分前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-100*time.Hour - 31*time.Minute)

	e := getTestTargetEvent()
	text, textTw := e.GetCountdownText(now)
	if text != "2021/01/17\n1行目\n2行目\n@3行目まで\nあと 4 日" {
		t.Errorf("text=%s", text)
	}
	if textTw != "2021/01/17 1行目 2行目 @ 3行目まで あと 4 日" {
		t.Errorf("textTw=%s", textTw)
	}
}

func TestCountdownText2(t *testing.T) {
	// 予定時刻の100時間前
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst).Add(-100 * time.Hour)

	e := getTestTargetEvent()
	text, textTw := e.GetCountdownText(now)
	if text != "2021/01/17\n1行目\n2行目\n@3行目まで\nあと 100 時間" {
		t.Errorf("text=%s", text)
	}
	if textTw != "2021/01/17 1行目 2行目 @ 3行目まで あと 100 時間" {
		t.Errorf("textTw=%s", text)
	}
}

func TestCountdownText3(t *testing.T) {
	// 予定時刻ちょうど
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 17, 17, 10, 0, 0, jst)

	e := getTestTargetEvent()
	text, textTw := e.GetCountdownText(now)
	if text != "1行目\n2行目\n@3行目！" {
		t.Errorf("text=%s", text)
	}
	if textTw != "1行目 2行目 @ 3行目！" {
		t.Errorf("textTw=%s", text)
	}
}

func TestGetTargetEvent(t *testing.T) {
	jst, _ := time.LoadLocation(location)
	now := time.Date(2021, 1, 15, 1, 2, 3, 0, jst)

	e := getTargetEvent(now)
	if e.Title != "「Don’t Blink」発売記念\nインターネットサイン会" {
		t.Errorf("event(%+v) is not expected", e)
	}
}

func TestGenerateImage(t *testing.T) {
	now := getNow()
	event := getTargetEvent(now)
	t.Logf("e=%+v", event)
	text, _ := event.GetCountdownText(now)
	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)

		out := generateTodayImage(imgInfo, text)
		f, err := os.Create(fmt.Sprintf("./output%02d.png", i))
		if err != nil {
			t.Errorf("failed to save file")
		}
		png.Encode(f, out)
	}
}

func TestLoadEvents(t *testing.T) {
	if len(eventsList) == 0 {
		t.Errorf("len(eventsList) == 0")
	}
	for _, e := range eventsList {
		if e.Time.IsZero() {
			t.Errorf("time is invalid: event=%+v", e)
		}
	}
}

func TestSelectRandomImage(t *testing.T) {
	for i := 0; i < 10; i++ {
		img := selectRandomImage()
		t.Logf("%+v", img)
	}
}

func TestSelectRandomSong(t *testing.T) {
	for i := 0; i < 10; i++ {
		song := selectRandomSong()
		t.Logf("%+v", song)
	}
}
