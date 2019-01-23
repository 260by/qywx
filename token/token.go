package token

import (
	"io/ioutil"
	"fmt"
	"time"
	"net/http"
	"github.com/tidwall/gjson"
	"github.com/sirupsen/logrus"
)

const (
	tokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	contentType = "application/json"
)

// AccessToken 凭证
type AccessToken struct {
	Ticket string	// 凭证
	ExpiresIn int64	// 凭证有效时间,单位: 秒
	NextGet int64	// 下次获取凭证时间
	CreateAt int64	// 取得凭证时间
}

// HTTPClient 客户端
type HTTPClient struct {
	UserAgent string
	ContentType string
}

// Get 获取token
func Get(corpID, appSecret string) (*AccessToken, error) {
	url := fmt.Sprintf(tokenURL, corpID, appSecret)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token AccessToken
	errCode := gjson.GetBytes(body, "errcode")
	if errCode.Int() != 0 {
		errMsg := gjson.GetBytes(body, "errmsg").String()
		err := fmt.Errorf("%v", errMsg)
		return nil, err
	}
	
	token.Ticket = gjson.GetBytes(body, "access_token").String()
	token.ExpiresIn = gjson.GetBytes(body, "expires_in").Int()
	token.CreateAt = time.Now().Unix()
	token.NextGet = token.CreateAt + token.ExpiresIn

	return &token, nil
}

// Loop 定时器
func Loop(corpID, appSecret string, token *AccessToken) {
	var refreshInterval time.Duration

	t, err := Get(corpID, appSecret)
	if err != nil {
		logrus.Errorln(err)
	}
	*token = *t
	refreshInterval = time.Duration(t.ExpiresIn) * time.Second
	logrus.Infof("Next get access-token time %s", time.Unix(t.NextGet, 0).Format("2006-01-02 15:04:05"))
	timeTicker := time.NewTicker(refreshInterval)

	for {
		select {
		case <- timeTicker.C:
			
			t, err := Get(corpID, appSecret)
			if err != nil {
				logrus.Errorln(err)
			}
			*token = *t
			refreshInterval = time.Duration(t.ExpiresIn) * time.Second
			logrus.Infof("Next get access-token time %s", time.Unix(t.NextGet, 0).Format("2006-01-02 15:04:05"))
			// timeTicker.Stop()
			timeTicker = time.NewTicker(refreshInterval)
		}
	}
}