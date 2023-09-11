package crypto

import (
	"fmt"
	"github.com/yunduansing/gtools/gen"
	"testing"
)

func TestEncryptPwdByPbkdf2(t *testing.T) {
	//p,_:=EncryptPwdByPbkdf2("123456")
	ok := CheckPwdByPbkdf2("123456", "TfQrWD8xfc9fGHdob+c7WuJrZrsxO3ZBefFP/h/TMD8XxITY9qeW9nViQZO1yxjrOfHXG2DuOkd4")
	fmt.Println(ok)
}

func TestGenHmacSha256(t *testing.T) {
	k1 := GenHmacSha256(gen.UUID(), gen.UUID())
	k2 := GenHmacSha256(gen.UUID(), gen.UUID())
	k3 := GenHmacSha256(gen.UUID(), gen.UUID())
	fmt.Println(k1, k2, k3)
}
