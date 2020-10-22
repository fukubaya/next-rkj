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
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	location     = "Asia/Tokyo"
	fontFilePath = "font/mplus-1p-bold.ttf"
)

var (
	imageList []ImageInfo
	songsList []SongInfo
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

// SongInfo struct
type SongInfo struct {
	Title string `json:"title"`
	Link  struct {
		Spotify string `json:"spotify"`
		Apple   string `json:"apple"`
	} `json:"link"`
}

// PubSubMessage struct
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func init() {
	loadImageList()
	loadSongsList()
	fontData = loadFont(fontFilePath)
	initRand()
}

func initRand() {
	// random seed
	t := time.Now().UnixNano() % 1000
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64-1000))
	rand.Seed(seed.Int64() + t)
}

func loadImageList() {
	f, _ := os.Open("image.json")
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&imageList)
}

func loadSongsList() {
	f, _ := os.Open("songs.json")
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&songsList)
}

func selectRandomImage() ImageInfo {
	return imageList[rand.Intn(len(imageList))]
}

func selectRandomSong() SongInfo {
	return songsList[rand.Intn(len(songsList))]
}

func getSong(name string) SongInfo {
	for _, s := range songsList {
		if s.Title == name {
			return s
		}
	}
	return selectRandomSong()
}

func selectPOPSong(t time.Time) (SongInfo, string) {
	h := t.Hour()

	switch h {
	case 3, 4:
		return getSong("星が降る街 (ALBUM ver.)"), "夜明け"
	case 7, 8, 9:
		return getSong("足音"), "朝"
	case 10, 11:
		return getSong("BON-NO BORN"), "午前中"
	case 12:
		return getSong("宙に浮くぐらい"), "正午"
	case 13, 14, 15:
		return getSong("SLEEPY BUSTERS"), "昼"
	case 16, 17, 18:
		return getSong("わたし色のトビラ"), "夕方"
	case 19, 20, 21, 22:
		return getSong("axis"), "夜"
	case 23, 0:
		return getSong("スーパースター"), "深夜"
	case 1, 2:
		return getSong("寝具でSING A SONG"), "就寝"
	case 5, 6:
		return getSong("夜更けのプロローグ (ALBUM ver.)"), "夜明け"
	default:
		return getSong("axis"), "リード曲"
	}
}

func getTargetDate() time.Time {
	jst, _ := time.LoadLocation(location)
	return time.Date(2020, 10, 24, 0, 0, 0, 0, jst)
}

func getTargetDateTime() time.Time {
	jst, _ := time.LoadLocation(location)
	return time.Date(2020, 10, 24, 12, 0, 0, 0, jst)
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

func hoursUntil(from time.Time, to time.Time) int {
	if from.After(to) {
		return 0
	}
	m := int(to.Sub(from).Minutes())
	// n時間30分前〜n-1時間30分前はn
	return (m + 30) / 60
}

func nearTargetDateTime(from time.Time, to time.Time) bool {
	s := to.Sub(from).Seconds()
	// 1分前から5分後まで
	return s < 60 && s > -300
}

func countdownText(from time.Time) string {
	hours := hoursUntil(from, getTargetDateTime())
	if hours <= 100 {
		return fmt.Sprintf("あと %d 時間", hours)
	}
	days := daysUntil(from, getTargetDate())
	return fmt.Sprintf("あと %d 日", days)
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

	// split by newline
	lines := strings.Split(text, "\n")

	// draw
	maxIndex := -1
	maxLen := 0
	for i, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
			maxIndex = i
		}
	}
	fontsize := float64(calcFontSize(imgInfo, lines[maxIndex], len(lines)))
	opt := truetype.Options{
		Size: fontsize,
	}
	lineHeight := (imgInfo.BottomLeft.Y - imgInfo.TopLeft.Y) / len(lines)
	for i, l := range lines {
		face := truetype.NewFace(fontData, &opt)
		dr := &font.Drawer{
			Dst:  out,
			Src:  image.NewUniform(color.RGBA{0, 0, 128, 255}),
			Face: face,
			Dot:  fixed.Point26_6{},
		}
		dr.Dot.X = fixed.I((imgInfo.BottomRight.X+imgInfo.BottomLeft.X)/2) - dr.MeasureString(l)/2
		dr.Dot.Y = fixed.I(imgInfo.TopLeft.Y + i*lineHeight + int(lineHeight/2) + int(fontsize/2))
		dr.DrawString(l)
	}
	return out
}

func calcFontSize(imgInfo ImageInfo, text string, n int) int {
	var width = imgInfo.TopRight.X - imgInfo.TopLeft.X
	var height = (imgInfo.BottomLeft.Y - imgInfo.TopLeft.Y) / n

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
	loadImageList()
	initRand()
	fontData = loadFont(fontFilePath)
	main()
	return nil

}

func main() {
	now := getNow()

	hours := hoursUntil(now, getTargetDateTime())
	near := nearTargetDateTime(now, getTargetDateTime())

	// 期限後は実行しない
	if hours <= 0 && !near {
		return
	}

	var out image.Image
	var text string
	var textTw string
	if near {
		text = "まもなく\n1st single 「Don’t Blink」 \n発売決定記念インターネットサイン会！"
		textTw = "まもなく1st single 「Don’t Blink」 発売決定記念インターネットサイン会！"
		out = generateTodayImage(selectRandomImage(), text)
	} else {
		text = fmt.Sprintf("2020/10/24\n1st single 「Don’t Blink」\n発売決定記念インターネットサイン会まで\n%s!!", countdownText(now))
		textTw = fmt.Sprintf("2020/10/24 1st single 「Don’t Blink」 発売決定記念インターネットサイン会まで%s!!", countdownText(now))
		out = generateTodayImage(selectRandomImage(), text)
	}
	// encode image to base64
	encodeString := encodePng(out)

	// upload media
	api := getTwitterAPI()
	media, _ := api.UploadMedia(encodeString)

	// tweet
	v := url.Values{}
	v.Add("media_ids", media.MediaIDString)
	tweetText := fmt.Sprintf("%s\n#内藤るな #白浜あや #高井千帆 #青山菜花\n#BOLT #ボルト #BOLTデマス", textTw)

	tweet, err := api.PostTweet(tweetText, v)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(tweet.Text)

}

// Tweet daily song
func TweetSong(ctx context.Context, m PubSubMessage) error {
	loadSongsList()
	initRand()
	songMain()
	return nil
}

func songMain() {
	// select random song
	now := getNow()
	song, label := selectPOPSong(now)

	// tweet
	tweetText := fmt.Sprintf(
		"%d時『%s』: %s\n\n%s\n%s\n\n#BOLT #ボルト",
		now.Hour(), label, song.Title, song.Link.Apple, song.Link.Spotify)

	// api
	api := getTwitterAPI()

	// tweet
	v := url.Values{}
	tweet, err := api.PostTweet(tweetText, v)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(tweet.Text)
}
