package main

import (
	"crypto/md5"
	crypt_rand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func MAChash(TestString string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(TestString))
	Result := Md5Inst.Sum([]byte(""))
	md5Str := string(Result[:])
	//fmt.Printf("%x\n\n", Result)
	rand.Seed(time.Now().Unix())
	x := rand.Int63()
	y := rand.Int63()
	S, _ := json.Marshal(deviceInfoStr)

	timeStr := time.Now().Format("2006-01-02_150405")
	Result1 := md5Str + string(S) + strconv.FormatInt(x, 10) + strconv.FormatInt(y, 10) + timeStr

	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(Result1))
	Result = Sha1Inst.Sum([]byte(""))
	hashStr := fmt.Sprintf("%x", Result) //将[]byte转成16进制
	//fmt.Printf("%x\n\n", Result)
	return hashStr
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(crypt_rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	timeStr := time.Now().Format("2006-01-02_15_04_05")
	file, err := os.Create("private" + timeStr + ".pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	pubkey = fmt.Sprintf("%x", derPkix)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public" + timeStr + ".pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
