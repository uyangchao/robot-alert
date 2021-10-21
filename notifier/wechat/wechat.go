package notifier

import "github.com/gin-gonic/gin"

var (
	wechatHook = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
)

type WechatHook struct {
	HookUrl string
}

type MarkdownData struct {
	Msgtype  string            `json:"msgtype"`
	Markdown map[string]string `json:"Markdown"`
}

func WechatAlert(c *gin.Context) {

}

func (w WechatHook) NewMarkdownMsg(context string) *MarkdownData {
	return &MarkdownData{
		Msgtype: "markdown",
		Markdown: map[string]string{
			"content": "context",
		},
	}
}
