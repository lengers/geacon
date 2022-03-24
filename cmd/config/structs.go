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