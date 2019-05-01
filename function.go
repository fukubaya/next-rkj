package function

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	location = "Asia/Tokyo"
	fontsize = 75
)

// PubSubMessage struct
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func getTargetDate() time.Time {
	jst, _ := time.LoadLocation(location)
	return time.Date(2019, 7, 15, 0, 0, 0, 0, jst)
}

func getNow() time.Time {
	jst, _ := time.LoadLocation(location)
	return time.Now().In(jst)
}

func daysUntil(from time.Time, to time.Time) int {
	if from.After(to) {
		return 0
	}
	h := int(to.Sub(from).Hours())
	return (h / 24) + 1
}

func loadImg(imgPath string) image.Image {
	f, _ := os.Open(imgPath)
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}
	return img
}

func loadFont(ttfPath string) *truetype.Font {
	ttf, err := ioutil.ReadFile(ttfPath)
	if err != nil {
		log.Fatalln(err)
	}

	ft, err := truetype.Parse(ttf)
	if err != nil {
		log.Fatalln(err)
	}
	return ft
}

func encodePng(img image.Image) string {
	var buff bytes.Buffer
	png.Encode(&buff, img)
	return base64.StdEncoding.EncodeToString(buff.Bytes())
}

// getTwitterApi creates api client
func getTwitterAPI() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"))
}

func generateTodayImage(baseImgPath string, text string) image.Image {
	// load image
	img := loadImg(baseImgPath)
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{0, 0}, draw.Over)

	// load font
	ft := loadFont("font/mplus-1p-bold.ttf")
	opt := truetype.Options{
		Size: fontsize,
	}
	face := truetype.NewFace(ft, &opt)
	dr := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.RGBA{215, 46, 42, 255}),
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	dr.Dot.X = (fixed.I(out.Bounds().Dx()) - dr.MeasureString(text)) / 2
	dr.Dot.Y = fixed.I(out.Bounds().Dy() - int(fontsize/2))
	dr.DrawString(text)
	return out
}

// Tweet daily comment
func Tweet(ctx context.Context, m PubSubMessage) error {
	main()
	return nil
}

func main() {

	days := daysUntil(getNow(), getTargetDate())
	text := fmt.Sprintf("あと %d 日", days)

	// create image
	out := generateTodayImage("img/next-rkj-16-9.jpg", text)

	// encode image to base64
	encodeString := encodePng(out)

	// upload media
	api := getTwitterAPI()
	media, _ := api.UploadMedia(encodeString)

	// tweet
	v := url.Values{}
	v.Add("media_ids", media.MediaIDString)
	tweetText := fmt.Sprintf("%s\n#内藤るな #高井千帆 #平瀬美里\n#ELRFES", text)

	tweet, err := api.PostTweet(tweetText, v)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(tweet.Text)

}
