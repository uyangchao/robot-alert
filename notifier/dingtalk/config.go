package dingtalk

type DingTalkNotification struct {
	MessageType string                        `json:"msgtype"`
	Markdown    *DingTalkNotificationMarkdown `json:"markdown,omitempty"`
	At          *DingTalkNotificationAt       `json:"at,omitempty"`
}

type DingTalkNotificationText struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DingTalkNotificationMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkNotificationAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	AtUserIds []string `json:"atUserIds,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}
