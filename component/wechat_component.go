package component

import (
	"encoding/json"
	"errors"
	"github.com/hihozhou/wechat/http"
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

// 获取component_access_token返回数据
type ComponentAccessTokenData struct {
	ComponentAccessToken string `xml:"component_access_token" json:"component_access_token"`
	ExpiresIn            string `xml:"expires_in" json:"expires_in"`
}

// 获取令牌component_access_token
// author:hihozhou
// param componentVerifyTicket
// 文档https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_access_token.html
func (wc *WechatComponent) GetAccessToken(componentVerifyTicket string) (data *ComponentAccessTokenData, err error) {
	postData := struct {
		Component_appid         string `json:"component_appid"`
		Component_appsecret     string `json:"component_appsecret"`
		Component_verify_ticket string `json:"component_verify_ticket"`
	}{
		Component_appid:         wc.AppId,
		Component_appsecret:     wc.AppSecret,
		Component_verify_ticket: componentVerifyTicket,
	}
	result := http.Post(apiComponentTokenUrl, postData, "application/json")
	apiErr := &ApiError{}
	err = json.Unmarshal([]byte(result), apiErr)
	if err != nil {
		return nil, err
	}
	if apiErr.isError() {
		return nil, errors.New(apiErr.Error())
	}
	err = json.Unmarshal([]byte(result), data)
	if (err != nil) {
		return nil, err
	}
	return data, nil

}
