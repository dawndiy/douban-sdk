package doubanfm

import ()

// User is Douban account
type User struct {
	UID        string     `json:"uid"`
	Name       string     `json:"name"`
	ID         string     `json:"id"`
	IsPro      bool       `json:"is_pro"`
	PlayRecord PlayRecord `json:"play_record"`
	DBCL2      string
	Expires    int
}

// PlayRecord of Douban account
type PlayRecord struct {
	Liked  int `json:"liked"`
	Played int `json:"played"`
	Banned int `json:"banned"`
}

// Channel for Douban FM
type Channel struct {
	Name   string      `json:"name"`
	SeqID  int         `json:"seq_id"`
	AbbrEN string      `json:"AbbrEN"`
	ID     interface{} `json:"channel_id"` // channel_id 可能是字符或数字 :(
	NameEN string      `json:"name_en"`
}

// Song from Douban FM
type Song struct {
	Picture    string `json:"picture"`
	AlbumTitle string `json:"albumtitle"`
	Like       int    `json:"like"`
	Album      string `json:"album"`
	SSID       string `json:"ssid"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Artist     string `json:"artist"`
	SubType    string `json:"subtype"`
	Length     int    `json:"length"`
	SID        string `json:"sid"`
	AID        string `json:"aid"`
	Company    string `json:"company"`
	PublicTime string `json:"public_time"`
	SHA256     string `json:"sha256"`
	Kbps       string `json:"kbps"`
}
