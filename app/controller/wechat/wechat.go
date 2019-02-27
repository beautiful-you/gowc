package wechat

import (
	"errors"
	"fmt"

	"github.com/beautiful-you/gowc/config"

	"github.com/beautiful-you/wechat/cache"

	"github.com/beautiful-you/wechat"
	"github.com/beautiful-you/wechat/context"
	"github.com/beautiful-you/wechat/message"
	"github.com/gin-gonic/gin"
)

// WeChat 控制器
type WeChat struct {
}

// 接口信息
const (
	RedirectURL           = "http://am.jyacad.cc/wechat/public/account/auth_call"
	ComponentloginPageURL = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=3"
)

var memCache = cache.NewMemcache("127.0.0.1:11211")

//配置微信参数
var cfg = &wechat.Config{
	AppID:          "xxxx",
	AppSecret:      "xxxx",
	Token:          "xxxx",
	EncodingAESKey: "xxxx",
	Cache:          memCache,
}

// AuthCall ... 授权后回调地址
func (w *WeChat) AuthCall(c *gin.Context) {

}

// AuthURL ... 授权地址
func (w *WeChat) AuthURL(c *gin.Context) {
	ctx := new(context.Context)
	// 获取预授权码
	PreCode, err := ctx.GetPreCode()
	if err != nil {
		fmt.Println(err)
		c.Writer.WriteString("GetPreCode() error")
		return
	}
	URL := fmt.Sprintf(ComponentloginPageURL, cfg.AppID, PreCode, RedirectURL)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, "<a href='"+URL+"'>点击授权</a>")

}

// MessageWithEvent ... 消息与事件接收url
func (w *WeChat) MessageWithEvent(c *gin.Context) {
	wc := wechat.NewWechat(cfg)
	server := wc.GetServer(c.Request, c.Writer)

	//设置接收消息的处理方法
	server.SetMessageHandler(messageWithEventHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(errors.New("server.Serve() error: "))
		fmt.Println(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		fmt.Println(errors.New("server.Send() error: "))
		fmt.Println(err)
		return
	}
}

// messageWithEventHandler  消息与事件处理
func messageWithEventHandler(msg message.MixMessage) *message.Reply {
	if len(msg.ComponentVerifyTicket) > 0 {
		// 缓存这个参数
		sc := new(config.Cache)
		err := sc.Set("ComponentVerifyTicket", msg.ComponentVerifyTicket)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		return nil
	}
	return nil
}

// Test 执行测试任务
func (w *WeChat) Test(c *gin.Context) {

}

// AuthEvent 授权事件接收URL
func (w *WeChat) AuthEvent(c *gin.Context) {

	wc := wechat.NewWechat(cfg)
	server := wc.GetServer(c.Request, c.Writer)
	//设置接收消息的处理方法
	server.SetMessageHandler(authEventHandler)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(errors.New("server.Serve() error: "))
		fmt.Println(err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		fmt.Println(errors.New("server.Send() error: "))
		fmt.Println(err)
		return
	}
}

// authEventHandler 授权事件处理
func authEventHandler(msg message.MixMessage) *message.Reply {
	return nil
}

// VerifyFile 微信文件效验
func (w *WeChat) VerifyFile(c *gin.Context) {
	c.Writer.WriteString("65e10e6cda0f37d81cdffaf5a3441979")
}
