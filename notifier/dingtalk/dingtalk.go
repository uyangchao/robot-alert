package notifier

var (
	DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingTalkNotification struct {
	MessageType string            `json:"msgtype"`
	Markdown    map[string]string `json:"markdown"`
}

type DingTalkNotificationMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkNotificationAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// func Build(msgType string, ctx string) {
// 	return
// }

// func SendNotification(token string) error {

// }
