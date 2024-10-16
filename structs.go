package main

type StatusResponse struct {
	Status string `xml:"status"`
}

type SaltResponse struct {
	Salt string `xml:"salt"`
}

type LoginResponse struct {
	Token string `xml:"token"`
}

type FilesResponse struct {
	Folders []struct {
		Ident string `xml:"ident"`
		Name  string `xml:"name"`
		Path  string `xml:"path"`
	} `xml:"folder"`
	Files []struct {
		Ident string `xml:"ident"`
		Name  string `xml:"name"`
		Size  string `xml:"size"`
	} `xml:"file"`
}

type FileLinkResponse struct {
	Link string `xml:"link"`
}

type UrlResponse struct {
	Url string `xml:"url"`
}

type UploadResponse struct {
	Ident string `json:"ident"`
}

type UserResponse struct {
	Id       string `xml:"id"`
	Username string `xml:"username"`
	Email    string `xml:"email"`
	Bytes    string `xml:"bytes"`
	Space    string `xml:"private_space"`
	VipUntil string `xml:"vip_until"`
}
