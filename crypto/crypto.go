package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	mathrand "math/rand"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
	iter       = 1000
	keyLen     = 32
)

// EncryptPwdByPbkdf2 EncryptPwd 加密密码
func EncryptPwdByPbkdf2(pwd string) (encrypt string, err error) {
	// 1、生成随机长度的盐值
	salt, err := randSalt()
	if err != nil {
		return
	}

	// 2、生成加密串
	en := encryptPwdWithSalt([]byte(pwd), salt)
	en = append(en, salt...)

	// 3、合并盐值
	encrypt = base64.StdEncoding.EncodeToString(en)

	return
}

func randSalt() ([]byte, error) {
	// 生成8-32之间的随机数字
	salt := make([]byte, mathrand.Intn(saltMaxLen-saltMinLen)+saltMinLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptPwdWithSalt(pwd, salt []byte) (pwdEn []byte) {
	pwd = append(pwd, salt...)
	pwdEn = pbkdf2.Key(pwd, salt, iter, keyLen, sha256.New)
	return
}

// CheckPwdByPbkdf2 CheckEncryptPwdMatch 验证密码是否与加密串匹配
func CheckPwdByPbkdf2(pwd, encrypt string) (ok bool) {
	// 1、参数校验
	if len(encrypt) == 0 {
		return
	}

	enDecode, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return
	}

	// 2、截取加密串 固定长度
	salt := enDecode[keyLen:]

	// 3、比对
	enBase64 := base64.StdEncoding.EncodeToString(enDecode[0:keyLen])
	pwdEnBase64 := base64.StdEncoding.EncodeToString(encryptPwdWithSalt([]byte(pwd), salt))
	ok = enBase64 == pwdEnBase64

	return
}

func GenHmacSha256(key string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
