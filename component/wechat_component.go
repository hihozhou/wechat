package component

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/hihozhou/wechat/http"
	"github.com/hihozhou/wechat/util"
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
	debugMode      bool          //是否开启debug模式
	//componentVerifyTicket string        //内存缓存的component_verify_ticket
	//componentAccessToken  string		//内存缓存的component_access_token
}

//是指是否是debug模式
func (wc *WechatComponent) DebugMode(enable bool) *WechatComponent {
	wc.debugMode = enable
	return wc
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

	requestUrl := fmt.Sprintf(apiComponentTokenUrl, accessToken)
	data := &ComponentAccessTokenData{}

	err = wc.attempt(requestUrl, postData, data)
	if err != nil {
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
// author hihozhou
// 微信文档：https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/pre_auth_code.html
func (wc *WechatComponent) GetComponentPreAuthCode() (data *PreAuthCodeData, err error) {

	accessToken, err := wc.GetAccessToken()
	if err != nil {
		return nil, err
	}

	postData := struct {
		ComponentAppId string `json:"component_appid"`
	}{
		ComponentAppId: wc.AppId,
	}
	requestUrl := fmt.Sprintf(apiCreatePreAuthCodeUrl, accessToken)
	data = &PreAuthCodeData{}

	err = wc.attempt(requestUrl, postData, data)
	if err != nil {
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

//================================================使用授权码获取授权信息=============================================================

// 使用授权码获取授权信息
type AuthorizationInfo struct {
	AuthorizerAppid        string `xml:"authorizer_appid" json:"authorizer_appid"`               //授权方 appid
	AuthorizerAccessToken  string `xml:"authorizer_access_token" json:"authorizer_access_token"` //接口调用令牌（在授权的公众号/小程序具备 API 权限时，才有此返回值）
	ExpiresIn              int64  `xml:"expires_in" json:"expires_in"`                           //authorizer_access_token 的有效期（在授权的公众号/小程序具备API权限时，才有此返回值），单位：秒
	AuthorizerRefreshToken string `xml:"expires_in" json:"authorizer_refresh_token"`             //刷新令牌（在授权的公众号具备API权限时，才有此返回值），
	// 刷新令牌主要用于第三方平台获取和刷新已授权用户的 authorizer_access_token。一旦丢失，只能让用户重新授权，才能再次拿到新的刷新令牌。用户重新授权后，之前的刷新令牌会失效
}

type funcInfo struct {
}

// 使用授权码获取授权信息
func (wc *WechatComponent) GetComponentApiQueryAuth(authorizationCode string) (authorizationInfo *AuthorizationInfo, err error) {

	accessToken, err := wc.GetAccessToken()
	if err != nil {
		return nil, err
	}

	postData := &struct {
		ComponentAppId    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		ComponentAppId:    wc.AppId,
		AuthorizationCode: authorizationCode,
	}
	requestUrl := fmt.Sprintf(apiQueryAuthUrl, accessToken)
	authorizationInfo = &AuthorizationInfo{}
	err = wc.attempt(requestUrl, postData, authorizationInfo)
	if err != nil {
		return nil, err
	}
	return authorizationInfo, nil

}

//请求方法
func (wc *WechatComponent) attempt(requestUrl string, params interface{}, result interface{}) error {
	if wc.debugMode {
		fmt.Println("请求微信接口，url : " + requestUrl)
		fmt.Println("参数 : " + util.GetObjFormatStr(params))
	}

	resultByte, err := http.Post(requestUrl, params, "application/json")
	if err != nil {
		if wc.debugMode {
			fmt.Println("请求接口失败，接口请求错误，err : " + err.Error())
		}
		return err
	}
	apiErr := &ApiError{}
	err = json.Unmarshal(resultByte, apiErr)
	if err != nil {
		if wc.debugMode {
			fmt.Println("解析数据错误1，接口请求返回错误，err : " + err.Error())
		}
		return err
	}
	if apiErr.isError() {
		if wc.debugMode {
			fmt.Println("请求接口失败，接口请求返回错误，err : " + apiErr.Error())
		}
		return errors.New(apiErr.Error())
	}
	err = json.Unmarshal(resultByte, result)
	if (err != nil) {
		if wc.debugMode {
			fmt.Println("解析数据错误2，err : " + err.Error())
		}
		return err
	}
	if wc.debugMode {
		fmt.Println("接口调用成功，返回结果 : " + util.GetObjFormatStr(result))
	}
	return nil
}
