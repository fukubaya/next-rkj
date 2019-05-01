package function

import (
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

func TestGenerateImage(t *testing.T) {
	out := generateTodayImage("img/next-rkj-16-9.jpg", "あと 10 日")
	f, err := os.Create("./output.png")
	if err != nil {
		t.Errorf("failed to save file")
	}
	png.Encode(f, out)
}
