package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GehirnInc/crypt"
	"gopkg.in/ini.v1"
)

func get_passwd_hash(salt string, password string) string {
	md5_crypt, _ := crypt.MD5.New().Generate([]byte(password), []byte("$1$"+salt))
	hash := sha1.Sum([]byte(md5_crypt))
	return hex.EncodeToString(hash[:])
}

func post_request(url string, data string) []byte {
	client := &http.Client{}
	resp, err := client.Post(url, "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewReader([]byte(data)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

func verify_request(raw []byte) bool {
	var response StatusResponse
	err := xml.Unmarshal(raw, &response)
	if err != nil {
		panic(err)
	}
	return response.Status == "OK"
}

func load_config() *ini.File {
	user_home, _ := os.UserHomeDir()
	config_path := user_home + "/.wscli"

	cfg, err := ini.Load(config_path)
	if err != nil {
		fmt.Println("Config file (~/.wscli) not found")
		panic(err)
	}
	return cfg
}
