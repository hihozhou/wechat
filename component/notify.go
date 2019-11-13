package component

import (
	"encoding/xml"
	"errors"
	"github.com/hihozhou/wechat/component/crypto"
)

// 微信回调数据
type ComponentNotifyData struct {
	AppId   string `xml:"AppId"`   //开发平台第三方appid
	Encrypt string `xml:"Encrypt"` //加密的数据
}

// 验证票据数据，推送内容解密后数据
// author:hihozhou
// 文档： https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/api/component_verify_ticket.html
type ComponentVerifyTicketData struct {
	AppId                 string `xml:"AppId"`                 //开发平台第三方appid
	CreateTime            int64  `xml:"CreateTime"`            //时间戳，单位：s
	InfoType              string `xml:"InfoType"`              //固定为："component_verify_ticket"
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"` //Ticket 内容
}

//回调验证
//todo 直接传入request
//todo error整理
func (wc *WechatComponent) NotifyValid(componentNotifyData *ComponentNotifyData, timestamp, nonce, signature string) (data *ComponentVerifyTicketData, err error) {
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
	msgDecrypt, decryptAppId, err := decryptor.Decrypt(componentNotifyData.Encrypt)
	if err != nil {
		return nil, errors.New("decrypt err:" + err.Error())
	}

	//判断解密的appid和调用appid
	if (decryptAppId != wc.AppId) {
		return nil, errors.New("解密appid不同")
	}

	//解密数据绑定到struct
	var componentVerifyTicketData ComponentVerifyTicketData
	err = xml.Unmarshal(msgDecrypt, &componentVerifyTicketData)
	if (err != nil) {
		return nil, errors.New("绑定解密后数据失败，err：" + err.Error())
	}
	//缓存ticket

	return &componentVerifyTicketData, nil

}
