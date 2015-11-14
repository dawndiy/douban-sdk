package doubanfm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	httpurl "net/url"
	"strings"
)

// douban url
const DoubanURL = "http://douban.fm"

var client *http.Client

// init http client
func init() {
	jar, _ := cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
	}
}

// SetUser set douban account
func SetUser(user User) {
	dbcl2 := user.DBCL2
	u, _ := httpurl.Parse(DoubanURL)
	cookies := []*http.Cookie{
		&http.Cookie{Name: "dbcl2", Value: dbcl2},
	}
	client.Jar.SetCookies(u, cookies)
}

// PlayListNew get a new play list
func PlayListNew(channel string) ([]Song, error) {
	return PlayList("n", "", "0.0", channel, "")
}

// PlayListPlaying get a play list when playing
func PlayListPlaying(sid, channel string) ([]Song, error) {
	return PlayList("p", sid, "0.0", channel, "128")
}

// PlayListSkip get a play list by skip current playing
func PlayListSkip(sid, pt, channel string) ([]Song, error) {
	return PlayList("s", sid, pt, channel, "128")
}

// PlayListBan stop play current list and get a new list
func PlayListBan(sid, pt, channel string) ([]Song, error) {
	return PlayList("b", sid, pt, channel, "128")
}

// PlayListEnd mark finished playing a song
func PlayListEnd(sid, pt, channel string) ([]Song, error) {
	return PlayList("e", sid, pt, channel, "128")
}

// PlayListRate mark like this song
func PlayListRate(sid, pt, channel string) ([]Song, error) {
	return PlayList("r", sid, pt, channel, "128")
}

// PlayListUnrate mark dislike this song
func PlayListUnrate(sid, pt, channel string) ([]Song, error) {
	return PlayList("u", sid, pt, channel, "128")
}

// PlayList get play list from douban.fm
func PlayList(_type, sid, pt, channel, pb string) ([]Song, error) {
	api := "/j/mine/playlist"
	params := map[string]interface{}{
		"type":    _type,
		"sid":     sid,
		"pt":      pt,
		"channel": channel,
		"pb":      pb,
		"from":    "mainsite",
		"r":       fmt.Sprintf("%x", (int64(rand.Float64()*0xEFFFFFFFFF))+0x1000000000),
	}
	url := makeURL(api, params)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	data := struct {
		Result int    `json:"r"`
		Songs  []Song `json:"song"`
	}{}
	json.Unmarshal(body, &data)

	for _, v := range res.Cookies() {
		// logged out
		if v.Name == "dbcl2" && v.Value == "deleted" {
			return data.Songs, errors.New("logged out")
		}
	}

	return data.Songs, nil
}

// ChannelChange report chanel change
func ChannelChange(fromCID, toCID, area string) error {
	api := "/j/change_channel"
	params := map[string]interface{}{
		"fcid": fromCID,
		"tcid": toCID,
		"area": area,
	}
	url := makeURL(api, params)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	for _, v := range res.Cookies() {
		// logged out
		if v.Name == "dbcl2" && v.Value == "deleted" {
			return errors.New("logged out")
		}
	}
	return nil
}

// NewCaptcha get captcha id
func captchaID() (string, error) {
	api := "/j/new_captcha"
	params := map[string]interface{}{}
	url := makeURL(api, params)
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	str := string(body)
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		str = strings.Replace(str, "\"", "", -1)
	}
	return str, err
}

// Captcha get captcha base64 string and id
func Captcha() (string, string, error) {
	captchaID, err := captchaID()
	if err != nil {
		return "", "", err
	}
	api := "/misc/captcha"
	params := map[string]interface{}{
		"size": "m",
		"id":   captchaID,
	}
	imageURL := makeURL(api, params)
	res, err := client.Get(imageURL)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	str := base64.StdEncoding.EncodeToString(body)

	return str, captchaID, nil
}

// Login login douban account
func Login(name, password, captcha, captchaID string) (User, error) {
	api := "/j/login"
	url := makeURL(api, map[string]interface{}{})

	params := httpurl.Values{}
	params.Set("source", "radio")
	params.Set("alias", name)
	params.Set("form_password", password)
	params.Set("captcha_solution", captcha)
	params.Set("captcha_id", captchaID)

	req, err := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	if err != nil {
		return User{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	data := struct {
		Result int    `json:"r"`
		ErrMsg string `json:"err_msg"`
		User   User   `json:"user_info"`
	}{}
	json.Unmarshal(body, &data)

	if data.Result != 0 {
		return User{}, errors.New(data.ErrMsg)
	}

	user := data.User

	for _, v := range res.Cookies() {
		if v.Name == "dbcl2" {
			user.DBCL2 = v.Value
		}
		if v.Domain == ".douban.fm" && v.RawExpires != "" {
			user.Expires = int(v.Expires.Unix())
		}
	}

	return user, nil
}

// Logout douban account
func Logout() {
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
}

// makeURL make url with query string
func makeURL(api string, kwargs map[string]interface{}) string {
	params := []string{}
	for k, v := range kwargs {
		params = append(params, k+"="+fmt.Sprint(v))
	}

	if len(params) == 0 {
		return DoubanURL + api
	}

	return DoubanURL + api + "?" + strings.Join(params, "&")
}
