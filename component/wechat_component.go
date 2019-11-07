package component

import (
	"github.com/hihozhou/wechat/cache"
	"time"
)

const VerifyTicketCacheKeyPrefix = "wechat:component:verify_ticket:"

type WechatComponent struct {
	AppId          string //appId
	AppSecret      string //AppSecret
	Token          string //token
	EncodingAESKey string //aesKey
}

// 创建开放平台操作对象
// author : hihozhou
func New(appId, appSecret, token, encodingAESKey string) (*WechatComponent) {
	return &WechatComponent{
		AppId:          appId,
		AppSecret:      appSecret,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
}

// 获取微信开放平台账号对应的component_verify_ticket的缓存key
// author : hihozhou
func (wc *WechatComponent) GetVerifyTicketCacheKey() (string) {
	return VerifyTicketCacheKeyPrefix + wc.AppId
}

// 设置component_verify_ticket
func (wc *WechatComponent) SetVerifyTicketCache(ticket string) {
	h, _ := time.ParseDuration("1h")
	cache.Set(wc.GetVerifyTicketCacheKey(), ticket, h)
}

// 获取component_verify_ticket
func (wc *WechatComponent) GetVerifyTicket() (string, bool) {
	ticket, exist := cache.Get(wc.GetVerifyTicketCacheKey())
	if (exist) {
		return ticket.(string), exist
	}
	return "", exist
}

// 请求验证

func GetAccessToken() {
	//http client
	//http 请求
	//返回结果
	//储存token
	//appId := "wxdc4e3ad21308ea2a"
	//appSecret := "3eb6c01e7a49fc64a79c12057420e6c1"

}
