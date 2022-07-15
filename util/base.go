package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/big"
)

func RsaEncode(scr, pk []byte) []byte {
	//pem解码
	block, _ := pem.Decode(pk)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, scr)
	if err != nil {
		panic(err)
	}
	//返回密文
	return cipherText

}

func Str2byte(scr string) ([]byte, error) {
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(scr)), simplifiedchinese.GBK.NewEncoder()))
}

func RandMAC() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	// Set the local bit
	buf[0] |= 2
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5]), err
}

func IntToBytes(n int) []byte {
	data := int64(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, data)
	return bytebuf.Bytes()
}

func BytesToInt32(bys []byte) int32 {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

func BytesToInt16(bys []byte) int16 {
	bytebuff := bytes.NewBuffer(bys)
	var data int16
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

func BytesToInt64(bys []byte) int64 {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.LittleEndian, &data)
	return data
}

func RandDigitString(n int) string {
	var digits = "0123456789"
	var d2 = "123456789"
	if n <= 1 {
		idx, _ := rand.Int(rand.Reader, big.NewInt(10))
		return idx.String()
	}

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(10))
		b[i] = digits[int(idx.Int64())]
	}

	idx, _ := rand.Int(rand.Reader, big.NewInt(10))
	b[0] = d2[idx.Int64()%int64(len(d2))]

	return string(b)
}
