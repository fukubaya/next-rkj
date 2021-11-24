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
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"

	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	location                   = "Asia/Tokyo"
	fontFilePath               = "font/mplus-1p-bold.ttf"
	youtubeChannelId           = "UCNsGYZjlivJYZdbxqrR-26g"
	uploadImageRetryCount  int = 3
	uploadImageRetrySecond int = 30
)

var (
	imageList  []ImageInfo
	songsList  []SongInfo
	eventsList []EventInfo
	tour       TourInfo
	colorList  = [...]color.RGBA{{215, 46, 42, 255}, {151, 95, 162, 255}, {254, 246, 155, 255}, {11, 83, 148, 255}}
	fontData   *truetype.Font
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

// EventInfo struct
type EventInfo struct {
	Title string    `json:"title"`
	Time  time.Time `json:"time"`
}

// TourInfo struct
type TourInfo struct {
	Title  string      `json:"title"`
	Events []EventInfo `json:"events"`
}

// YouTubeInfo struct
type YouTubeInfo struct {
	Title           string
	SubscriberCount int
	VideoCount      int
	ViewCount       int
}

func (e EventInfo) IsZero() bool {
	return e.Title == "" && e.Time.IsZero()
}

func (e EventInfo) DaysUntil(now time.Time) int {
	// 時間より小さい単位を落とす
	d := e.Time.Truncate(time.Hour).Add(-time.Duration(e.Time.Hour()) * time.Hour)

	if now.After(d) {
		return 0
	}
	h := int(d.Sub(now).Hours())
	return (h / 24) + 1
}

func (e EventInfo) HoursUntil(now time.Time) int {
	if now.After(e.Time) {
		return 0
	}
	m := int(e.Time.Sub(now).Minutes())
	// n時間30分前〜n-1時間30分前はn
	return (m + 30) / 60
}

func (e EventInfo) NearTargetDateTime(now time.Time) bool {
	s := e.Time.Sub(now).Seconds()
	// 1分前から5分後まで
	return s < 60 && s > -300
}

func (e EventInfo) GetCountdownText(now time.Time) (string, string) {
	if e.NearTargetDateTime(now) {
		return fmt.Sprintf("%s！", e.Title),
			fmt.Sprintf("%s！", strings.Replace(strings.Replace(e.Title, "\n", " ", -1), "@", "@ ", -1))
	}

	var countdown string

	hours := e.HoursUntil(now)
	if hours <= 100 {
		countdown = fmt.Sprintf("あと %d 時間", hours)
	} else {
		days := e.DaysUntil(now)
		countdown = fmt.Sprintf("あと %d 日", days)
	}

	text := fmt.Sprintf("%s\n%sまで\n%s", e.Time.Format("2006/01/02"), e.Title, countdown)
	return text, strings.Replace(strings.Replace(text, "\n", " ", -1), "@", "@ ", -1)
}

func (t TourInfo) Finished(now time.Time) bool {
	return len(t.Remained(now)) == 0
}

func (t TourInfo) GetCountdownText(now time.Time) (string, string) {
	texts := make([]string, 0, len(t.Events)+1)
	texts = append(texts, t.Title)

	for _, e := range t.Remained(now) {
		var countdown string

		hours := e.HoursUntil(now)
		if hours <= 100 {
			countdown = fmt.Sprintf("%d時間", hours)
		} else {
			days := e.DaysUntil(now)
			countdown = fmt.Sprintf("%d日", days)
		}
		texts = append(texts, fmt.Sprintf("%s %sまで%s", e.Time.Format("2006/01/02"), e.Title, countdown))
	}
	text := strings.Join(texts, "\n")
	return text, strings.Replace(text, "@", "@ ", -1)
}

func (t TourInfo) Remained(now time.Time) []EventInfo {
	filtered := make([]EventInfo, 0, len(t.Events))
	for _, e := range t.Events {
		if e.Time.After(now) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// PubSubMessage struct
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func init() {
	loadImageList()
	loadSongsList()
	loadEventsList()
	loadTourEventsList()
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

func loadEventsList() {
	f, _ := os.Open("events.json")
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&eventsList)
	sort.Slice(eventsList, func(i, j int) bool {
		return eventsList[i].Time.Before(eventsList[j].Time)
	})
}

func loadTourEventsList() {
	f, _ := os.Open("tour_events.json")
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&tour)
	sort.Slice(tour.Events, func(i, j int) bool {
		return tour.Events[i].Time.Before(tour.Events[j].Time)
	})
}

func selectRandomImage() ImageInfo {
	return imageList[rand.Intn(len(imageList))]
}

func selectRandomColor() color.RGBA {
	return colorList[rand.Intn(len(colorList))]
}

func selectRandomSong() SongInfo {
	return songsList[rand.Intn(len(songsList))]
}

func getTargetEvent(now time.Time) EventInfo {
	for _, e := range eventsList {
		if e.Time.After(now) {
			return e
		}
	}

	return EventInfo{}
}

func getNow() time.Time {
	jst, _ := time.LoadLocation(location)
	return time.Now().In(jst)
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

func encodeJpg(img image.Image) string {
	var buff bytes.Buffer
	jpeg.Encode(&buff, img, &jpeg.Options{Quality: 80})
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
	color := selectRandomColor()
	for i, l := range lines {
		face := truetype.NewFace(fontData, &opt)
		dr := &font.Drawer{
			Dst:  out,
			Src:  image.NewUniform(color),
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

func uploadImage(api *anaconda.TwitterApi, imgString string) (anaconda.Media, error) {
	var media anaconda.Media
	var err error

	for i := 0; i < uploadImageRetryCount; i++ {
		media, err = api.UploadMedia(imgString)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * time.Duration(uploadImageRetrySecond))
		} else {
			return media, err
		}
	}
	return media, err
}

// Tweet daily comment
func Tweet(ctx context.Context, m PubSubMessage) error {
	loadImageList()
	initRand()
	fontData = loadFont(fontFilePath)
	tweetMain()
	return nil

}

func tweetMain() {
	now := getNow()

	// 現在から5分前まで
	event := getTargetEvent(now.Add(-time.Duration(5) * time.Minute))

	// 期限後は実行しない
	if event.IsZero() || event.HoursUntil(now) <= 0 && !event.NearTargetDateTime(now) {
		return
	}

	// text, image
	text, textTw := event.GetCountdownText(now)
	out := generateTodayImage(selectRandomImage(), text)

	// encode image to base64
	encodeString := encodeJpg(out)

	// upload media
	api := getTwitterAPI()
	media, err := uploadImage(api, encodeString)

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
	song := selectRandomSong()

	// tweet
	tweetText := fmt.Sprintf(
		"%d時の1曲: %s\n%s\n%s\n#内藤るな #白浜あや #高井千帆 #青山菜花\n#BOLT #ボルト",
		now.Hour(), song.Title, song.Link.Apple, song.Link.Spotify)

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

// YouTube channel
func TweetYouTubeChannel(ctx context.Context, m PubSubMessage) error {
	youtubeMain()
	return nil
}

func youtubeMain() {
	now := getNow()

	// tweet
	youtubeInfo := getYouTubeInfo(youtubeChannelId)
	tweetText := fmt.Sprintf(
		"%s\n(%s)\nチャンネル登録者数:%d\n総視聴回数:%d\n公開動画数:%d\nhttps://www.youtube.com/channel/%s\n#BOLT #ボルト",
		youtubeInfo.Title,
		now.Format("2006年01月02日"),
		youtubeInfo.SubscriberCount,
		youtubeInfo.ViewCount,
		youtubeInfo.VideoCount,
		youtubeChannelId,
	)

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

func getYouTubeInfo(channelId string) YouTubeInfo {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	call := service.Channels.List([]string{"id", "snippet", "statistics"}).Id(channelId).MaxResults(1)
	response, err := call.Do()

	if err != nil {
		log.Fatalf("Error in retreiving from YouTube API: %#v", err)
	}

	ch := response.Items[0]

	return YouTubeInfo{
		Title:           ch.Snippet.Title,
		SubscriberCount: int(ch.Statistics.SubscriberCount),
		VideoCount:      int(ch.Statistics.VideoCount),
		ViewCount:       int(ch.Statistics.ViewCount),
	}
}

func TweetFirstTour(ctx context.Context, m PubSubMessage) error {
	loadImageList()
	initRand()
	fontData = loadFont(fontFilePath)
	firstTourMain()
	return nil
}

func firstTourMain() {
	now := getNow()

	// 期限後は実行しない
	if tour.Finished(now) {
		return
	}

	// text, image
	text, textTw := tour.GetCountdownText(now)
	out := generateTodayImage(selectRandomImage(), text)

	// encode image to base64
	encodeString := encodeJpg(out)

	// upload media
	api := getTwitterAPI()
	media, err := uploadImage(api, encodeString)

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
