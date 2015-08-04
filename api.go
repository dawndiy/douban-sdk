package doubanfm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// 用户登录
func (fm *DoubanFM) Login(email, password string) (*User, error) {
	url := fm.host + "j/app/login"
	params := map[string][]string{
		"email":    {email},
		"password": {password},
		"app_name": {fm.appName},
		"version":  {strconv.Itoa(fm.version)},
	}

	res, err := http.PostForm(url, params)
	if err != nil {
		return &User{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &User{}, err
	}

	user := &User{}
	json.Unmarshal(body, &user)
	if user.Result != 0 {
		return &User{}, errors.New(user.ERR)
	}

	return user, nil
}

// 频道
func (fm *DoubanFM) Channels() ([]Channel, error) {

	url := fm.host + "j/app/radio/channels"

	res, err := http.Get(url)
	if err != nil {
		return []Channel{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Channel{}, err
	}

	data := map[string][]Channel{}
	json.Unmarshal(body, &data)
	channels := data["channels"]

	return channels, nil
}

// 歌曲
func (fm *DoubanFM) Songs(opts ...map[string]string) ([]Song, error) {

	url := fm.host + "j/app/radio/people"

	var opt map[string]string

	if opts == nil {
		opt = DefaultOptions
	} else {
		opt = opts[0]
		opt["app_name"] = fm.appName
		opt["version"] = strconv.Itoa(fm.version)
	}

	if len(opt) > 0 {
		url += "?"
	}
	for k, v := range opt {
		url += k + "=" + v + "&"
	}

	res, err := http.Get(url)
	if err != nil {
		return []Song{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Song{}, err
	}

	data := struct {
		Result     int    `json:"r"`
		VersionMax int    `json:"version_max"`
		Songs      []Song `json:"song"`
		ERR        string `json:"err"`
	}{}
	json.Unmarshal(body, &data)

	if data.Result != 0 {
		return []Song{}, errors.New(data.ERR)
	}

	songs := data.Songs
	return songs, nil

}
