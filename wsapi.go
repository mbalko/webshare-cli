package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/GehirnInc/crypt/md5_crypt"
)

func find_ident(token string, filename string) string {
	dir := path.Dir(filename)

	file_list := files(token, dir, true)

	for i := 0; i < len(file_list.Files); i++ {
		if path.Base(filename) == file_list.Files[i].Name {
			return file_list.Files[i].Ident
		}
	}
	return ""
}

func login(username string, password string) string {
	var salt_url string = WEBSHARE + "/salt/"
	var login_url string = WEBSHARE + "/login/"
	var data string = "username_or_email=" + username
	var salt_response SaltResponse
	var login_response LoginResponse
	var hash_password string

	salt_raw := post_request(salt_url, data)
	if !verify_request(salt_raw) {
		return "FATAL - Salt"
	}
	xml.Unmarshal(salt_raw, &salt_response)
	hash_password = get_passwd_hash(salt_response.Salt, password)

	data += "&password=" + hash_password + "&keep_logged_in=0"
	login_raw := post_request(login_url, data)
	if !verify_request(login_raw) {
		return "FATAL - Login"
	}
	xml.Unmarshal(login_raw, &login_response)
	return login_response.Token
}

func user_data(token string) UserResponse {
	var url string = WEBSHARE + "/user_data/"
	var data string = "wst=" + token
	var user_response UserResponse
	req := post_request(url, data)
	if !verify_request(req) {
		fmt.Println("FATAL - user_data")
	}
	xml.Unmarshal(req, &user_response)
	return user_response
}

func files(token string, path string, private bool) FilesResponse {
	var url string = WEBSHARE + "/files/"
	var files_response FilesResponse
	data := "wst=" + token + "&path=" + path + "&private="
	if private {
		data += "1"
	} else {
		data += "0"
	}

	req := post_request(url, data)
	if !verify_request(req) {
		fmt.Println("FATAL - files")
	}
	xml.Unmarshal(req, &files_response)
	return files_response
}

func file_link(token string, ident string, download bool, ident_type int8) string {
	var url string = WEBSHARE + "/file_link/"
	var file_link_response FileLinkResponse

	switch ident_type {
	case IDENT_TYPE_URL:
		ident = strings.Split(strings.Split(ident, "file/")[1], "/")[0]
	case IDENT_TYPE_FILENAME:
		ident = find_ident(token, ident)
		if ident == "" {
			return "FATAL - File not found"
		}
	}

	data := "wst=" + token + "&ident=" + ident
	req := post_request(url, data)
	if !verify_request(req) {
		return "FATAL - file_link"
	}
	xml.Unmarshal(req, &file_link_response)

	if download {
		f := path.Base(file_link_response.Link)
		fmt.Println("Downloading " + f)
		out, err := os.Create(f)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		resp, err := http.Get(file_link_response.Link)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != http.StatusOK {
			panic(resp.StatusCode)
		}
		defer resp.Body.Close()
		io.Copy(out, resp.Body)
	}
	return file_link_response.Link
}

func normal_link(token string, ident string, ident_type int8) string {
	if ident_type == IDENT_TYPE_FILENAME {
		ident = find_ident(token, ident)
	}
	return "https://webshare.cz/#/file/" + ident
}

func remove_file(token string, ident string, ident_type int8) bool {
	var url string = WEBSHARE + "/remove_file/"
	if ident_type == IDENT_TYPE_FILENAME {
		ident = find_ident(token, ident)
	}
	data := "wst=" + token + "&ident=" + ident
	req := post_request(url, data)
	return verify_request(req)
}

func upload(token string, file_path string, folder string) string {
	var url string = WEBSHARE + "/upload_url/"
	var url_response UrlResponse
	req := post_request(url, "")
	if !verify_request(req) {
		return "FATAL - upload_url"
	}
	xml.Unmarshal(req, &url_response)
	url = url_response.Url

	file, _ := os.Open(file_path)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("wst", token)
	writer.WriteField("folder", folder)
	writer.WriteField("private", "1")
	writer.WriteField("adult", "0")
	file_stat, _ := file.Stat()
	writer.WriteField("total", strconv.FormatInt(file_stat.Size(), 10))
	writer.WriteField("offset", "0")
	part, err := writer.CreateFormFile("name", filepath.Base(file_path))
	if err != nil {
		panic(err)
	}
	io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	response_body, _ := io.ReadAll(response.Body)

	var upload_response UploadResponse
	json.Unmarshal(response_body, &upload_response)
	return upload_response.Ident
}
