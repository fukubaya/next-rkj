package function

import (
	"fmt"
	"image/png"
	_ "image/png"
	"os"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	days := daysUntil(getNow(), getTargetDate())
	t.Logf("days=%d", days)
}

func TestDays(t *testing.T) {
	// 当日0時の1秒前
	now := getTargetDate().Add(-1 * time.Second)
	days := daysUntil(now, getTargetDate())
	if days != 1 {
		t.Errorf("days=%d", days)
	}
}

func TestDays2(t *testing.T) {
	// 当日0時の1n秒後
	now := getTargetDate().Add(1 * time.Second)
	days := daysUntil(now, getTargetDate())
	if days != 0 {
		t.Errorf("days=%d", days)
	}
}

func TestHours(t *testing.T) {
	// 予定時刻の24時間30分前+1n秒後
	now := getTargetDateTime().Add(-1470*time.Minute + 1*time.Nanosecond)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours2(t *testing.T) {
	// 予定時刻の23時間30分前
	now := getTargetDateTime().Add(-1410 * time.Minute)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours3(t *testing.T) {
	// 予定時刻の1時間前
	now := getTargetDateTime().Add(-1 * time.Hour)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 1 {
		t.Errorf("hours=%d", hours)
	}
}

func TestNearTargetDateTime1(t *testing.T) {
	// 予定時刻の59分前
	now := getTargetDateTime().Add(-59 * time.Minute)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime2(t *testing.T) {
	// 予定時刻の1分前+1n秒後
	now := getTargetDateTime().Add(-1*time.Minute + 1*time.Nanosecond)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime3(t *testing.T) {
	// 予定時刻の4分59秒後
	now := getTargetDateTime().Add(299 * time.Second)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime4(t *testing.T) {
	// 予定時刻の5分後
	now := getTargetDateTime().Add(5 * time.Minute)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestCountdownText(t *testing.T) {
	// 予定時刻の100時間31分前
	now := getTargetDateTime().Add(-100*time.Hour - 31*time.Minute)
	text := countdownText(now)
	if text != fmt.Sprintf("あと %d 日", daysUntil(now, getTargetDate())) {
		t.Errorf("text=%s", text)
	}
}

func TestCountdownText2(t *testing.T) {
	// 予定時刻の100時間前
	now := getTargetDateTime().Add(-100 * time.Hour)
	text := countdownText(now)
	if text != "あと 100 時間" {
		t.Errorf("text=%s", text)
	}
}

func TestGenerateImage(t *testing.T) {
	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)
		out := generateTodayImage(imgInfo, "2020/10/24\n1st single 「Don’t Blink」 \n発売決定記念インターネットサイン会まで\nあと 18 日")
		f, err := os.Create(fmt.Sprintf("./output%02d.png", i))
		if err != nil {
			t.Errorf("failed to save file")
		}
		png.Encode(f, out)
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

func TestGetSong(t *testing.T) {
	if s := getSong("竜人くんが大好きです♡"); s.Title != "竜人くんが大好きです♡" {
		t.Errorf("invalid")
	}
}

func TestSelectPOPSong(t *testing.T) {
	jst, _ := time.LoadLocation(location)

	if s, l := selectPOPSong(time.Date(2020, 7, 15, 3, 1, 2, 3, jst)); s.Title != "星が降る街 (ALBUM ver.)" || l != "夜明け" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 7, 1, 2, 3, jst)); s.Title != "足音" || l != "朝" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 10, 1, 2, 3, jst)); s.Title != "BON-NO BORN" || l != "午前中" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 12, 1, 2, 3, jst)); s.Title != "宙に浮くぐらい" || l != "正午" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 15, 1, 2, 3, jst)); s.Title != "SLEEPY BUSTERS" || l != "昼" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 18, 1, 2, 3, jst)); s.Title != "わたし色のトビラ" || l != "夕方" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 20, 1, 2, 3, jst)); s.Title != "axis" || l != "夜" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 23, 1, 2, 3, jst)); s.Title != "スーパースター" || l != "深夜" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 1, 1, 2, 3, jst)); s.Title != "寝具でSING A SONG" || l != "就寝" {
		t.Errorf("invalid")
	}
	if s, l := selectPOPSong(time.Date(2020, 7, 15, 5, 1, 2, 3, jst)); s.Title != "夜更けのプロローグ (ALBUM ver.)" || l != "夜明け" {
		t.Errorf("invalid")
	}
}
