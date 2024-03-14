package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Alg Cost  Salt                  Hash
//$2a$ 10   $aXk/TWeIt/2AUOD0c2Yls HKJxXykKg9U3LSLPeLqK.SgDjJK3gN2
//$2a$ 10   $hZc2p5HKvjNfUuSGP9Qxv FF0cqcBUidl3nloLITI/2OJwQ4mrPTi

func main() {
	// begin := time.Now()
	// for range 2 {
	// 	genBcrypt()
	// }
	// fmt.Println("Tempo:", time.Since(begin).Seconds())
	fmt.Println(generateTokenKey())
}

func genMD5() {
	h := md5.New()
	_, err := h.Write([]byte("123456"))
	if err != nil {
		panic(err)
	}
}

func genSHA() {
	h := sha512.New()
	_, err := h.Write([]byte("123456"))
	if err != nil {
		panic(err)
	}
}

func genBcrypt() {
	_, err := bcrypt.GenerateFromPassword([]byte("123546"), 20)
	if err != nil {
		panic(err)
	}
}

func generateTokenKey() string {
	r := make([]byte, 32)
	rand.Read(r)
	return base64.URLEncoding.EncodeToString(r)
}
