package function

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	location     = "Asia/Tokyo"
	fontFilePath = "font/mplus-1p-bold.ttf"
	fontsize     = 75
)

var (
	imageList []ImageInfo
	fontData  *truetype.Font
)

// Point struct
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ImageInfo struct
type ImageInfo struct {
	Path        string `json:"path"`
	TopLeft     Point  `json:"topLeft"`
	TopRight    Point  `json:"topRight"`
	BottomLeft  Point  `json:"bottomLeft"`
	BottomRight Point  `json:"bottomRight"`
}

// PubSubMessage struct
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func init() {
	loadImageList()
	fontData = loadFont(fontFilePath)

	// random seed
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
}

func loadImageList() {
	f, _ := os.Open("image.json")
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&imageList)
	fmt.Printf("%+v\n", imageList)
}

func selectRandomImage() ImageInfo {
	return imageList[rand.Intn(len(imageList))]
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

func generateTodayImage(imgInfo ImageInfo, text string) image.Image {
	// load image
	img := loadImg(imgInfo.Path)
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.Point{0, 0}, draw.Over)

	// draw
	opt := truetype.Options{
		Size: float64(calcFontSize(imgInfo, text)),
	}
	face := truetype.NewFace(fontData, &opt)
	dr := &font.Drawer{
		Dst:  out,
		Src:  image.NewUniform(color.RGBA{215, 46, 42, 255}),
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	dr.Dot.X = fixed.I((imgInfo.BottomRight.X+imgInfo.BottomLeft.X)/2) - dr.MeasureString(text)/2
	dr.Dot.Y = fixed.I((imgInfo.BottomLeft.Y+imgInfo.TopLeft.Y)/2 + int(fontsize/2))
	dr.DrawString(text)
	return out
}

func calcFontSize(imgInfo ImageInfo, text string) int {
	var width = imgInfo.TopRight.X - imgInfo.TopLeft.X
	var height = imgInfo.BottomLeft.Y - imgInfo.TopLeft.Y

	for i := 0; i < 200; i += 5 {
		if i > height {
			return i - 5
		}
		face := truetype.NewFace(fontData, &truetype.Options{Size: float64(i)})
		textWidth := font.MeasureString(face, text)
		if i > height || textWidth > fixed.I(width) {
			return i - 5
		}
	}
	return 200
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
	out := generateTodayImage(selectRandomImage(), text)

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
