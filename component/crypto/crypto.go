package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

// MessageCrypter 封装了生成签名和消息加解密的方法
type Decryptor struct {
	token string
	appId string
	key   []byte
	iv    []byte
}

// NewMessageCrypter 方法用于创建 MessageCrypter 实例
//
// token 为开发者在微信开放平台上设置的 Token，
// encodingAESKey 为开发者在微信开放平台上设置的 EncodingAESKey，
// AppId
func NewDecryptor(appId, token, encodingAESKey string) (decryptor *Decryptor, err error) {
	var key []byte

	if key, err = base64.StdEncoding.DecodeString(encodingAESKey + "="); err != nil {
		return nil, err
	}

	if len(key) != 32 {
		return nil, ENCODING_AES_KEY_INVALID
	}

	iv := key[:16]
	decryptor = &Decryptor{token, appId, key, iv,}
	return decryptor, nil
}

// Decrypt 方法用于对密文进行解密
//
// 返回解密后的消息，AppId, 或者错误信息
func (decryptor Decryptor) Decrypt(text string) (decryptData []byte, appId string, err error) {

	deciphered, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, "", err
	}

	c, err := aes.NewCipher(decryptor.key)
	if err != nil {
		return nil, "", err
	}

	cbc := cipher.NewCBCDecrypter(c, decryptor.iv)
	cbc.CryptBlocks(deciphered, deciphered)

	decoded := PKCS7Decode(deciphered)

	buf := bytes.NewBuffer(decoded[16:20])

	var msgLen int32
	binary.Read(buf, binary.BigEndian, &msgLen)

	decryptData = decoded[20 : 20+msgLen]
	appId = string(decoded[20+msgLen:])

	return decryptData, appId, nil
}

//Signature sha1签名
func Signature(params ...string) string {

	sort.Strings(params)

	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 请求验证
func RequestSignature(token, timestamp, nonce, signature string) bool {
	return Signature(token, timestamp, nonce) == signature
}
