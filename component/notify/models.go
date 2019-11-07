package notify

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
