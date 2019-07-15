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
	now := time.Date(2019, 8, 2, 23, 59, 59, 0, jst)
	days := daysUntil(now, getTargetDate())
	if days != 1 {
		t.Errorf("days=%d", days)
	}
}

func TestDays2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 0, 0, 0, 1, jst)
	days := daysUntil(now, getTargetDate())
	if days != 0 {
		t.Errorf("days=%d", days)
	}
}

func TestHours(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 2, 10, 40, 0, 1, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 2, 11, 40, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 24 {
		t.Errorf("hours=%d", hours)
	}
}

func TestHours3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 10, 10, 0, 0, jst)
	hours := hoursUntil(now, getTargetDateTime())
	if hours != 1 {
		t.Errorf("hours=%d", hours)
	}
}

func TestNearTargetDateTime1(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 11, 9, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 11, 9, 0, 1, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime3(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 11, 14, 59, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if !near {
		t.Errorf("near=%v", near)
	}
}

func TestNearTargetDateTime4(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Date(2019, 8, 3, 11, 15, 0, 0, jst)
	near := nearTargetDateTime(now, getTargetDateTime())
	if near {
		t.Errorf("near=%v", near)
	}
}

func TestCountdownText(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 100.4999h
	now := time.Date(2019, 7, 30, 6, 40, 0, 0, jst)
	text := countdownText(now)
	if text != "あと 4 日" {
		t.Errorf("text=%s", text)
	}
}

func TestCountdownText2(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	// 100.4999h
	now := time.Date(2019, 7, 30, 6, 40, 0, 1, jst)
	text := countdownText(now)
	if text != "あと 100 時間" {
		t.Errorf("text=%s", text)
	}
}

func TestGenerateImage(t *testing.T) {
	t.Logf("%+v", lastImage)
	out := generateTodayImage(lastImage, "まもなく\nTIFに登場!!")
	f, err := os.Create("last.png")
	if err != nil {
		t.Errorf("failed to save file")
	}
	png.Encode(f, out)

	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)
		out := generateTodayImage(imgInfo, "TIFのステージまで\nあと 18 日")
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
