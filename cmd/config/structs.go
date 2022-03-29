package config

import (
	"net"
)

type TcpBeacon struct {
	Id int
	Conn net.Conn
	Ppid int
	EncryptedMetaInfo []byte
}

// type HostedFile struct {
// 	Id int
// 	BaseEncFileContent string
// }