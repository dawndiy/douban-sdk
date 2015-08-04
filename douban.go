package doubanfm

type DoubanFM struct {
	host    string
	appName string
	version int
}

// 用户
type User struct {
	ID     string `json:"user_id"`
	ERR    string `json:"err"`
	Token  string `json:"token"`
	Expire string `json:"expire"`
	Result int    `json:"r"`
	Name   string `json:"user_name"`
	Email  string `json:"email"`
}

// 频道
type Channel struct {
	Name   string      `json:"name"`
	SeqID  int         `json:"seq_id"`
	AbbrEN string      `json:"AbbrEN"`
	ID     interface{} `json:"channel_id"` // channel_id 可能是字符或数字 :(
	NameEN string      `json:"name_en"`
}

// 歌曲
type Song struct {
	Album      string  `json:"album"`
	Picture    string  `json:"picture"`
	SSID       string  `json:"ssid"`
	Artist     string  `json:"artist"`
	URL        string  `json:"url"`
	Company    string  `json:"company"`
	Title      string  `json:"title"`
	RatingAvg  float64 `json:"rating_avg"`
	Length     int     `json:"length"`
	SubType    string  `json:"subtype"`
	PublicTime string  `json:"public_time"`
	SID        string  `json:"sid"`
	AID        string  `json:"aid"`
	SHA256     string  `json:"sha256"`
	Kbps       string  `json:"kbps"`
	AlbumTitle string  `json:"albumtitle"`
	Like       int     `json:"like"`
}

// 默认参数
var DefaultOptions = map[string]string{
	"app_name": "radio_desktop_win",
	"version":  "100",
	"type":     "n",
	"channel":  "0",
}

func NewDoubanFM() *DoubanFM {
	return &DoubanFM{
		host:    "http://www.douban.com/",
		appName: "radio_desktop_win",
		version: 100,
	}
}
