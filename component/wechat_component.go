package component

import (
	"fmt"
	"github.com/hihozhou/wechat/cache"
	"github.com/hihozhou/wechat/http"
	"time"
)

const VerifyTicketCacheKeyPrefix = "wechat:component:verify_ticket:"

var (
	weixinComponentHost       = "https://api.weixin.qq.com/cgi-bin/component"
	apiComponentTokenUrl      = weixinComponentHost + "/api_component_token"
	apiCreatePreAuthCodeUrl   = weixinComponentHost + "/api_create_preauthcode?component_access_token=%s"
	apiQueryAuthUrl           = weixinComponentHost + "/api_query_auth?component_access_token=%s"
	apiAuthorizerTokenUrl     = weixinComponentHost + "/api_authorizer_token?component_access_token=%s"
	apiGetAuthorizerInfoUrl   = weixinComponentHost + "/api_get_authorizer_info?component_access_token=%s"
	apiGetAuthorizerOptionUrl = weixinComponentHost + "/api_get_authorizer_option?component_access_token=%s"
	apiSetAuthorizerOptionUrl = weixinComponentHost + "/api_set_authorizer_option?component_access_token=%s"

	oauthUrl = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
)

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

// 获取令牌component_access_token
// author:hihozhou
// 文档https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_access_token.html
func (wc *WechatComponent) GetAccessToken() {
	//http client
	//http 请求
	ticket, exist := wc.GetVerifyTicket()
	if !exist {
		//返回错误
	}
	postData := struct {
		Component_appid         string `json:"component_appid"`
		Component_appsecret     string `json:"component_appsecret"`
		Component_verify_ticket string `json:"component_verify_ticket"`
	}{
		Component_appid:         wc.AppId,
		Component_appsecret:     wc.AppSecret,
		Component_verify_ticket: ticket,
	}
	result := http.Post(apiComponentTokenUrl, postData, "application/json")
	fmt.Println(result)

	//返回结果
	//储存token
	//appId := "wxdc4e3ad21308ea2a"
	//appSecret := "3eb6c01e7a49fc64a79c12057420e6c1"

}
