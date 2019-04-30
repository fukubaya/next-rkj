package main

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

// getTwitterApi creates api client
func getTwitterAPI() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"))
}

// Tweet daily comment
func Tweet(ctx context.Context, m PubSubMessage) error {
	main()
	return nil
}

func main() {
	api := getTwitterAPI()

	target := getTargetDate()
	now := getNow()
	duration := target.Sub(now)
	days := (int(duration.Hours()) / 24) + 1
	text := fmt.Sprintf("あと %d 日", days)

	// load image
	f, _ := os.Open("img/next-rkj.jpg")
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{0, 0}, draw.Over)

	// load font
	ttf, err := ioutil.ReadFile("font/mplus-1p-bold.ttf")
	if err != nil {
		log.Fatalln(err)
	}

	ft, err := truetype.Parse(ttf)
	if err != nil {
		log.Fatalln(err)
	}
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

	// encode image to base64
	var buff bytes.Buffer
	png.Encode(&buff, out)
	encodeString := base64.StdEncoding.EncodeToString(buff.Bytes())

	// upload media
	media, _ := api.UploadMedia(encodeString)

	// tweet
	if true {
		v := url.Values{}
		v.Add("media_ids", media.MediaIDString)

		tweet, err := api.PostTweet(text, v)

		if err != nil {
			log.Fatalln(err)
		}

		log.Println(tweet.Text)
	}

}
