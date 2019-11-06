package component

type ComponentVerifyTicketData struct {
	AppId                 string `xml:"AppId"					json:"AppId"`                // 第三方平台 appid
	CreateTime            int64  `xml:"CreateTime"       		json:"CreateTime"`           // 时间戳，单位：s
	InfoType              string `xml:"InfoType"       			json:"InfoType"`           // 固定为："component_verify_ticket"
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"   json:"ComponentVerifyTicket"` // Ticket 内容
}

func GetAccessToken() {
	//http client
	//http 请求
	//返回结果
	//储存token
	//appId := "wxdc4e3ad21308ea2a"
	//appSecret := "3eb6c01e7a49fc64a79c12057420e6c1"

}
