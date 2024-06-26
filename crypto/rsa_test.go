package crypto

import (
	"fmt"
	"testing"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAnEMkukfTyJRz57abcMkUCr7shzltyqfh/Qpu9K8h6LfFK+Xj
MS5tWPsHCTwo2sBdrGJaBBSsdtiIkIuZ/lZc+zqZzEb0iA32ydMxq2JDsqvpjqTG
CV5i/uAAw6enVMVTq74FjWZD1EA5wvj+cuvI3fb8RFQtbVMNX+WHLDK0tdV9oS3y
auY3OwA7lE4hMbNfBEOCU6flF3g4HHE4Atw5llJ81pR7C5q38v5BiHpaO2wpgC3Y
dPV+fT9XNBsAk4813H0junx+xkh2At/7bfRFgxWoJffTNxREjiYCNj5mbU1TWs5w
kZMMM8aXrnK1MW6ieT1aTqCwnuv9y4lh8BHRmQIDAQABAoIBAB8ZeN+PKI4Ax7Gb
5QNVLaU22uXNzvVYpNBC6ZLOoTnLC1Whhu40noQpT7ztiXcB9pK2r1IlGC7+CO31
rsQRBaflgZMHoybJ2q5P11CT+cO/VbgzQIvASoUN6XEHNrUXIcAqTToXWpQNZUGR
4zhVh3BvtcTWmQTrVNvbV2P7Qn52qS4ltOKDbKgmqR2rPPcJlv6PBNGC8eScnVWb
zgBn7OBDmJQouRkA4qm+JQvMXgWl4ZP4ILJ3QL2D+y4Lp+vxxvyNoagpTJep9Kqu
rkFXM4zd25McD/FCQ/5GwRlMTrWQlbx0RuVNoqJOSaA76wIt6mKFjaSJ1L8HKYit
emt2nyUCgYEAnKaPkEA6RjUTJu5j/N2k6cub6IchRGDSxbIh+zwVnWB0AQ/8P/23
F7JdI/Y6d/VPQ5CVc5sqkWq7gFTreHlLYW8zsrasH0oIqq4HjJioIGmf6EFkG+gS
XRxoLZd5a9hJwl9GidPEkRDBCo1UNGvOGLTfmT/Dy1zQBSR41MAajl0CgYEA/12H
/kpS7f6D/My62lgFImq81vD+G+0CDtuYCWdOe4ITjtpltKnmUklbemX6RTkgaXHZ
zrjPevrrQpHRuFgEdyDaIe/UYbk6uBbfCeVeLy0036epcnDibV5OUTYf/nnED+je
x0DTeO9Q3IjwtVxTJIb6o7B44WoDRP12Z2qUxG0CgYAj6K56I674DwyP3Q/AlFJM
Yacgm3FNymAmf1n4zGIsDEMrRy3MwW0qGnHtX3ExtGLyGv+XW5dbdDr9BPphSHS2
JhaUvfkgxwjLj08dtJS9sVi9ldfL6dvjyjI7WUC1LwRhcROrXobttgh6UcFeruO0
mCm1aH+Ka4En3J9yLYMo5QKBgQD6wnq/b3ad8/bh8Jxvuk1YbZ4AJ2gwPf1uX56d
ZU+BIstX1QhYbxoXpadpuugWv7EfOkhaocpEBk+s8AoTzoKcBKuO9bYFPGH30aW1
TwG8q9Sm5zoSgd0m7ptTTiX+us2XneHdob/HmqE53lMYlI+kZ9OJl8or9/E79vUl
OljZxQKBgC5N+E25B2KEZJowTRnTSEI/Ha3UwpU6JLvGOkzFOh01uzpsSFPHbEoz
cPHMR9emOnvinpPazV2VxqlTqq2A9o5ipCQsbmD2DBHLfIuT6gqPRBPlkDXAMjQB
g32lXSkq7S8XxSin7fL9gQQupSEeA62sxDWWM941j+PFGyjds4BO
-----END RSA PRIVATE KEY-----`

type G struct {
	Codes []string `json:"codes"`
	Token string   `json:"token"`
}

func TestSignRsaSHA256(t *testing.T) {
	var m = map[string]interface{}{
		"id": 1,
		"class": map[string]interface{}{
			"id":   1,
			"name": "二班",
		},
	}
	sign, err := SignRsaSHA256(privateKey, m)
	fmt.Println(sign, err)

	fmt.Println(VerifySignRsaSha256(privateKey, sign+"2", m))

	c := G{
		Codes: []string{"A", "B"},
		Token: "eee",
	}

	sign, err = SignRsaSHA256(privateKey, c)
	fmt.Println(sign, err)

	fmt.Println(VerifySignRsaSha256(privateKey, sign+"2", c))
}
