package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

const (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	contentType = "application/json"
)

// Text 文本消息
type Text struct {
	ToUser string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentID int `json:"agentid"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

// TextCard 文本卡片消息
type TextCard struct {
	ToUser string `json:"touser"`
	MsgType string `json:"msgtype"`
	AgentID int `json:"agentid"`
	TextCard struct {
		Title string `json:"title"`
		Description string `json:"description"`
	} `json:"textcard"`
	Safe int `json:"safe"`
}

// Send 发送消息
func (t *Text) Send(accessToken string) error {
	err := PostAPI(accessToken, t)
	if err != nil {
		return err
	}
	return nil
}

// Send 文本卡片消息发送
func (t *TextCard) Send(accessToken string) error {
	err := PostAPI(accessToken, t)
	if err != nil {
		return err
	}
	return nil
}

// PostAPI 请求企业微信发送消息接口
func PostAPI(accessToken string, d interface{}) error {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", accessToken)
	jsonBytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	client := &http.Client{}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	// req.Header.Set("Content-Type", contentType)
	// req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	errCode := gjson.GetBytes(body, "errcode")
	if errCode.Int() != 0 {
		errMsg := gjson.GetBytes(body, "errmsg").String()
		err := fmt.Errorf("%v", errMsg)
		return err
	}

	return nil
}