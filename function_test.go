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
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 当日0時の1秒前
	now := time.Date(2019, 11, 14, 23, 59, 59, 0, jst)
	days := daysUntil(now, getTargetDate())
	if days != 1 {
		t.Errorf("days=%d", days)
	}
}

func TestDays2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 当日0時の1n秒後
	now := time.Date(2019, 11, 15, 0, 0, 0, 1, jst)
	days := daysUntil(now, getTargetDate())
	if days != 0 {
		t.Errorf("days=%d", days)
	}
}

func TestHours(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の24時間30分前+1n秒後
	now := time.Date(2019, 11, 14, 18, 30, 0, 1, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の23時間30分前
	now := time.Date(2019, 11, 14, 19, 30, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の1時間前
	now := time.Date(2019, 11, 15, 18, 0, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 1 {
		t.Errorf("hours=%d", hours)
	}
}

func TestNearTargetDateTime1(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の59分前
	now := time.Date(2019, 11, 15, 18, 1, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の1分前+1n秒後
	now := time.Date(2019, 11, 15, 18, 59, 0, 1, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の4分59秒後
	now := time.Date(2019, 11, 15, 19, 4, 59, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime4(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の5分後
	now := time.Date(2019, 11, 15, 19, 5, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestCountdownText(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の100時間31分
	now := time.Date(2019, 11, 11, 14, 29, 0, 0, jst)
	text := countdownText(now)
	if text != "あと 4 日" {
		t.Errorf("text=%s", text)
	}
}

func TestCountdownText2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 予定時刻の100時間前
	now := time.Date(2019, 11, 11, 15, 0, 0, 0, jst)
	text := countdownText(now)
	if text != "あと 100 時間" {
		t.Errorf("text=%s", text)
	}
}

func TestGenerateImage(t *testing.T) {
	t.Logf("%+v", lastImage)
	out := generateTodayImage(lastImage, "まもなく\nマウントレーニアホール渋谷の\nステージ!!")
	f, err := os.Create("last.png")
	if err != nil {
		t.Errorf("failed to save file")
	}
	png.Encode(f, out)

	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)
		out := generateTodayImage(imgInfo, "マウントレーニアホール渋谷の\nイベントまで\nあと 18 日")
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
