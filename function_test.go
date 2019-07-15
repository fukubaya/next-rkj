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
	now := time.Date(2019, 7, 14, 23, 59, 59, 0, jst)
	days := daysUntil(now, getTargetDate())
	if days != 1 {
		t.Errorf("days=%d", days)
	}
}

func TestDays2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 0, 0, 0, 1, jst)
	days := daysUntil(now, getTargetDate())
	if days != 0 {
		t.Errorf("days=%d", days)
	}
}

func TestHours(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 14, 17, 30, 0, 1, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 14, 18, 30, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 17, 0, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 1 {
		t.Errorf("hours=%d", hours)
	}
}

func TestNearTargetDateTime1(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 17, 59, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 17, 59, 0, 1, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 18, 4, 59, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime4(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 7, 15, 18, 5, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestCountdownText(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 100.4999h
	now := time.Date(2019, 7, 11, 13, 30, 0, 0, jst)
	text := countdownText(now)
	if text != "あと 4 日" {
		t.Errorf("text=%s", text)
	}
}

func TestCountdownText2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 100.4999h
	now := time.Date(2019, 7, 11, 13, 30, 0, 1, jst)
	text := countdownText(now)
	if text != "あと 100 時間" {
		t.Errorf("text=%s", text)
	}
}

func TestGenerateImage(t *testing.T) {
	t.Logf("%+v", lastImage)
	out := generateTodayImage(lastImage, "まもなく\n再始動!!")
	f, err := os.Create("last.png")
	if err != nil {
		t.Errorf("failed to save file")
	}
	png.Encode(f, out)

	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)
		out := generateTodayImage(imgInfo, "再始動まで\nあと 10 日")
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
