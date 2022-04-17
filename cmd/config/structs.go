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

type SocksSession struct {
	Id int
	Address string
	Port int
	Conn net.Conn
	OutBuf []byte
}