package component

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hihozhou/wechat/http"
	"net/url"
	"time"
)

const VerifyTicketCacheKeyPrefix = "wechat:component:verify_ticket:"
const AccessTokenCacheKeyPrefix = "wechat:component:access_token:"

var (
	weixinComponentHost       = "https://api.weixin.qq.com/cgi-bin/component"
	apiComponentTokenUrl      = weixinComponentHost + "/api_component_token"
	apiCreatePreAuthCodeUrl   = weixinComponentHost + "/api_create_preauthcode?component_access_token=%s"
	apiQueryAuthUrl           = weixinComponentHost + "/api_query_auth?component_access_token=%s"
	apiAuthorizerTokenUrl     = weixinComponentHost + "/api_authorizer_token?component_access_token=%s"
	apiGetAuthorizerInfoUrl   = weixinComponentHost + "/api_get_authorizer_info?component_access_token=%s"
	apiGetAuthorizerOptionUrl = weixinComponentHost + "/api_get_authorizer_option?component_access_token=%s"
	apiSetAuthorizerOptionUrl = weixinComponentHost + "/api_set_authorizer_option?component_access_token=%s"

	oauthUrl = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%d"
)

type WechatComponent struct {
	AppId          string        //appId
	AppSecret      string        //AppSecret
	Token          string        //token
	EncodingAESKey string        //aesKey
	RedisClient    *redis.Client //redis客户端，用于缓存component_verify_ticket和component_access_token
	//componentVerifyTicket string        //内存缓存的component_verify_ticket
	//componentAccessToken  string		//内存缓存的component_access_token
}

// 获取微信开放平台账号对应的component_verify_ticket的缓存key
// author : hihozhou
func (wc *WechatComponent) GetVerifyTicketCacheKey() (string) {
	return VerifyTicketCacheKeyPrefix + wc.AppId
}

// 设置component_verify_ticket
func (wc *WechatComponent) SetVerifyTicketCache(ticket string) {
	h, _ := time.ParseDuration("2h")
	wc.RedisClient.Set(wc.GetVerifyTicketCacheKey(), ticket, h)
}

// 获取component_verify_ticket
// todo 获取失败或不存在
func (wc *WechatComponent) GetVerifyTicket() (string, error) {
	return wc.RedisClient.Get(wc.GetVerifyTicketCacheKey()).Result()
}

//================================================获取access_token=============================================================

func (wc *WechatComponent) GetAccessTokenCacheKey() (string) {
	return AccessTokenCacheKeyPrefix + wc.AppId
}

func (wc *WechatComponent) SetAccessTokenCache(accessToken string) error {
	h, _ := time.ParseDuration("1h50m")
	statusCmd := wc.RedisClient.Set(wc.GetAccessTokenCacheKey(), accessToken, h)
	return statusCmd.Err()
}

// 获取component_verify_ticket
func (wc *WechatComponent) GetAccessTokenOnCache() (string, error) {
	return wc.RedisClient.Get(wc.GetAccessTokenCacheKey()).Result()
}

// 获取component_access_token返回数据
type ComponentAccessTokenData struct {
	ComponentAccessToken string `xml:"component_access_token" json:"component_access_token"`
	ExpiresIn            int64  `xml:"expires_in" json:"expires_in"`
}

// 获取令牌component_access_token
// author:hihozhou
// param componentVerifyTicket
// 文档https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_access_token.html
func (wc *WechatComponent) GetAccessToken() (accessToken string, err error) {

	accessToken, err = wc.GetAccessTokenOnCache()
	if err == nil {
		return accessToken, err
	}

	componentVerifyTicket, err := wc.GetVerifyTicket()
	if err != nil {
		return "", err
	}

	//通过缓存获取ticket
	postData := struct {
		Component_appid         string `json:"component_appid"`
		Component_appsecret     string `json:"component_appsecret"`
		Component_verify_ticket string `json:"component_verify_ticket"`
	}{
		Component_appid:         wc.AppId,
		Component_appsecret:     wc.AppSecret,
		Component_verify_ticket: componentVerifyTicket,
	}
	//请求接口
	resultByte, err := http.Post(apiComponentTokenUrl, postData, "application/json")
	if err != nil {
		return "", err
	}
	//判断请求接口返回数据是否错误
	apiErr := &ApiError{}
	err = json.Unmarshal(resultByte, apiErr)
	if err != nil {
		return "", err
	}
	if apiErr.isError() {
		return "", errors.New(apiErr.Error())
	}
	//解析正常数据
	data := &ComponentAccessTokenData{}
	err = json.Unmarshal(resultByte, data)
	if (err != nil) {
		return "", err
	}
	wc.SetAccessTokenCache(data.ComponentAccessToken)
	//缓存token
	return data.ComponentAccessToken, nil

}

//================================================获取授权码pre_auth_code=============================================================

// 获取预授权码返回结果数据
type PreAuthCodeData struct {
	PreAuthCode string `xml:"pre_auth_code" json:"pre_auth_code"`
	ExpiresIn   int64  `xml:"expires_in" json:"expires_in"`
}

// 获取预授权码pre_auth_code
func (wc *WechatComponent) GetComponentPreAuthCode() (data *PreAuthCodeData, err error) {

	accessToken, err := wc.GetAccessToken()
	if err != nil {
		return nil, err
	}

	postData := struct {
		Component_appid string `json:"component_appid"`
	}{
		Component_appid: wc.AppId,
	}
	requestUrl := fmt.Sprintf(apiCreatePreAuthCodeUrl, accessToken)

	resultByte, err := http.Post(requestUrl, postData, "application/json")
	if err != nil {
		return nil, err
	}
	apiErr := &ApiError{}
	err = json.Unmarshal(resultByte, apiErr)
	if err != nil {
		return nil, err
	}
	if apiErr.isError() {
		return nil, errors.New(apiErr.Error())
	}
	data = &PreAuthCodeData{}
	err = json.Unmarshal(resultByte, data)
	if (err != nil) {
		return nil, err
	}
	return data, nil
}

//================================================获取授权url=============================================================

// 获取授权url
// param redirectUri 回调url
// param auth_type 要授权的帐号类型， 1则商户扫码后，手机端仅展示公众号、2表示仅展示小程序，3表示公众号和小程序都展示。如果为未制定，则默认小程序和公众号都展示。第三方平台开发者可以使用本字段来控制授权的帐号类型。
func (wc *WechatComponent) GetComponentOauthUrl(redirectUri string, authType int) (string, error) {
	code, err := wc.GetComponentPreAuthCode()
	if (err != nil) {
		return "", err
	}
	//拼接url
	result := fmt.Sprintf(oauthUrl, wc.AppId, code.PreAuthCode, url.QueryEscape(redirectUri), authType)
	return result, nil
}
