package component

import (
	"encoding/xml"
	"errors"
	"github.com/hihozhou/wechat/component/crypto"
)

// 微信回调数据
type NotifyData struct {
	AppId   string `xml:"AppId"`   //开发平台第三方appid
	Encrypt string `xml:"Encrypt"` //加密的数据
}

// 微信回调加密解密数据
type NotifyInfo struct {
	//基本数据
	AppId      string `xml:"AppId"`      //开发平台第三方appid
	CreateTime int64  `xml:"CreateTime"` //时间戳，单位：s
	InfoType   string `xml:"InfoType"`   //固定为："component_verify_ticket"

	//component_verify_ticket 回到数据
	*ComponentVerifyTicketInfo
	// 授权变更数据
	*UpdateAuthorizedInfo
}

// 验证票据数据，推送内容解密后数据
// author:hihozhou
// 微信文档： https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_verify_ticket.html
type ComponentVerifyTicketInfo struct {
	//`xml:"ComponentVerifyTicketInfo,omitempty" json:"ComponentVerifyTicketInfo,omitempty"`
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket,omitempty"` //component_verify_ticket 内容
}

// 公众号或小程序授权后数据
// author :hihozhou
// 微信文档：https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/authorize_event.html#%E5%AD%97%E6%AE%B5%E8%AF%B4%E6%98%8E
// InfoType等于unauthorized，updateauthorized，authorized
type UpdateAuthorizedInfo struct {
	AuthorizerAppid              string `xml:"AuthorizerAppid,omitempty"`              //公众号或小程序的 appid
	AuthorizationCode            string `xml:"AuthorizationCode,omitempty"`            //授权码，可用于获取授权信息
	AuthorizationCodeExpiredTime int64  `xml:"AuthorizationCodeExpiredTime,omitempty"` //授权码过期时间 单位秒
	PreAuthCode                  string `xml:"PreAuthCode,omitempty"`                  //预授权码
}

//回调验证
//todo 直接传入request
//todo error整理
func (wc *WechatComponent) NotifyValid(notifyData *NotifyData, timestamp, nonce, signature string) (*NotifyInfo, error) {
	//验证签名
	if crypto.Signature(wc.Token, timestamp, nonce) != signature {
		return nil, errors.New("请求签名signature验证错误")
	}
	//获取
	//创建解密struct
	decryptor, err := crypto.NewDecryptor(wc.AppId, wc.Token, wc.EncodingAESKey)
	if err != nil {
		return nil, errors.New("解密失败，err:" + err.Error())
	}

	//数据解密
	msgDecrypt, decryptAppId, err := decryptor.Decrypt(notifyData.Encrypt)
	if err != nil {
		return nil, errors.New("decrypt err:" + err.Error())
	}

	//判断解密的appid和调用appid
	if (decryptAppId != wc.AppId) {
		return nil, errors.New("解密appid不同")
	}

	//解密数据绑定到struct
	var notifyInfo NotifyInfo
	err = xml.Unmarshal(msgDecrypt, &notifyInfo)
	if (err != nil) {
		return nil, errors.New("绑定解密后数据失败，err：" + err.Error())
	}
	//缓存ticket

	return &notifyInfo, nil

}
