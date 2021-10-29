package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// func SendAlert(c *gin.Context) {
// 	var Data template.Data
// 	err := c.BindJSON(&Data)
// 	if err != nil {

// 	}

// }

type NotificationResponse struct {
	ErrorMessage string `json:"errmsg"`
	ErrorCode    int    `json:"errcode"`
}

type Notifier struct {
	Url string `json:"url"`
}

func (a *Notifier) Send(context []byte) (*NotificationResponse, error) {
	httpReq, err := http.NewRequest("POST", a.Url, bytes.NewReader(context))
	if err != nil {
		return nil, errors.Wrap(err, "error building DingTalk request")
	}
	// 设置请求头{"Content-Type":"application/json"}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	var httpClient = &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, errors.Wrap(err, "error sending notification")
	}
	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("unacceptable response code %d", resp.StatusCode)
	}

	var robotResp NotificationResponse
	enc := json.NewDecoder(resp.Body)
	if err := enc.Decode(&robotResp); err != nil {
		return nil, errors.Wrap(err, "error decoding response")
	}

	return &robotResp, nil
}

func NewNotifier(url string) *Notifier {
	return &Notifier{
		Url: url,
	}
}
