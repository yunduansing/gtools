package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/youmark/pkcs8"
	"github.com/yunduansing/gtools/utils"
	"os"
	"reflect"
	"sort"
)

func getPrivateKeyFromRaw(raw string) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode(utils.StringToByte(raw))
	priKey, err := x509.ParsePKCS1PrivateKey(b.Bytes)
	return priKey, err
}

func SignRsaSHA256(privateKey string, data interface{}) (string, error) {
	priKey, err := getPrivateKeyFromRaw(privateKey)
	if err != nil {
		return "", err
	}
	signData, err := getSortedSignData(data)
	if err != nil {
		return "", err
	}
	//hash := sha256.New()
	//hash.Write(gen.StringToByte(signData))
	h := sha256.Sum256(utils.StringToByte(signData))
	sign, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, h[:])
	return hex.EncodeToString(sign), err
}

func VerifySignRsaSha256(key, sign string, data interface{}) (bool, error) {
	priKey, err := getPrivateKeyFromRaw(key)
	if err != nil {
		return false, err
	}
	signData, err := getSortedSignData(data)
	if err != nil {
		return false, err
	}
	//hash := sha256.New()
	//hash.Write(gen.StringToByte(signData))
	//signDataHashed := hash.Sum(nil)
	signature, _ := hex.DecodeString(sign)
	h := sha256.Sum256(utils.StringToByte(signData))
	err = rsa.VerifyPKCS1v15(&priKey.PublicKey, crypto.SHA256, h[:], signature)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getSortedSignData(data interface{}) (string, error) {
	var res string
	t := reflect.TypeOf(data)
	vals := reflect.ValueOf(data)
	switch t.Kind() {
	case reflect.Struct:
		var keys []string

		var m = make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			keys = append(keys, t.Field(i).Name)
			m[t.Field(i).Name] = vals.Field(i).Interface()
		}
		sort.Strings(keys)
		for _, field := range keys {
			res += getFieldValueString(m[field])
		}
	case reflect.Map:
		var keys []string
		var m = data.(map[string]interface{})
		for k, _ := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			res += fmt.Sprint(m[k])
		}
	default:
		return "", errors.New("不支持的数据类型")
	}
	return res, nil

}

func getFieldValueString(i interface{}) string {
	t := reflect.ValueOf(i)
	switch t.Kind() {
	case reflect.Int, reflect.Uint, reflect.Uint64, reflect.Uint16, reflect.Uint8, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return fmt.Sprint(i)
	case reflect.String:
		return i.(string)
	case reflect.Map, reflect.Array, reflect.Struct, reflect.Slice:
		return utils.ToJsonString(i)
	}
	return ""
}

func EncryptRsaPrivateKey() {
	// Generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Define the encryption password
	password := "your-password"

	// Convert the private key to PKCS8
	pkcs8Key, err := pkcs8.ConvertPrivateKeyToPKCS8(privateKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Encrypt the PKCS8 private key
	encryptedKey, err := pkcs8.EncryptPKCS8PrivateKey(pkcs8Key, password, pkcs8.PBKDF2Count, x509.PEMCipherAES256)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write the encrypted private key to a file
	err = os.WriteFile("encrypted_private_key.pem", encryptedKey, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
