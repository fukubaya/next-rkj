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

func TestGenerateImage(t *testing.T) {
	for i, imgInfo := range imageList {
		t.Logf("%+v", imgInfo)
		out := generateTodayImage(imgInfo, "あと 10 日")
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
		songMain()
	}
}
