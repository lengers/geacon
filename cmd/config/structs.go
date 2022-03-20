package config

import (
	"net"
	"io"
)

type TcpBeacon struct {
    Id int
    Conn net.Conn
	Reader io.Reader
	EncryptedMetaInfo []byte
}