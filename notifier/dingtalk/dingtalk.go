package dingtalk

var (
	DingTalkUrl = "https://oapi.dingtalk.com/robot/send?access_token="
)

type DingTalkNotifier struct {
	Url          string
	Notification DingTalkNotification
}

// func (dtn *DingTalkNotifier) send() (*notifier.NotificationResponse, error) {
// 	notifier := notifier.NewNotifier(dtn.Url)
// 	resp, err := notifier.Send()

// }
