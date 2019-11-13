package wechat

import (
	"github.com/go-redis/redis/v7"
	"github.com/hihozhou/wechat/component"
)

type WechatFactory struct {
	redisClient *redis.Client
}

func New(options *redis.Options) *WechatFactory {
	factory := &WechatFactory{
		redisClient: redis.NewClient(options),
	}
	return factory
}

func NewByRedis(client *redis.Client) *WechatFactory {
	factory := &WechatFactory{
		redisClient: client,
	}
	return factory
}

func (factory *WechatFactory) NewWechatComponent(appId, appSecret, token, encodingAESKey string) *component.WechatComponent {
	wechatComponent := &component.WechatComponent{
		AppId:          appId,
		AppSecret:      appSecret,
		Token:          token,
		EncodingAESKey: encodingAESKey,
		RedisClient:    factory.redisClient,
	}
	return wechatComponent
}
