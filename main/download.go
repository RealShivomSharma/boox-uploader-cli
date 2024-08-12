package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/schollz/progressbar/v3"
)

type title_and_hash struct {
	title string
	hash  string
}

func UploadFromMemory(file_to_download title_and_hash, client http.Client) {
	get_request_url := "https://cdn3.booksdl.org/get.php?" + file_to_download.hash

	post_request_url := getRoutes(getBooxURL()).Upload

	fmt.Println(post_request_url)

	resp, err := client.Get(get_request_url)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	filename := file_to_download.title + ".pdf"

	var buffer bytes.Buffer
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(&buffer, bar), resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, &buffer)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", post_request_url, body)
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Println("testing req", req.Header)

	uploadResp, err := client.Do(req)
	if err != nil {
		log.Fatal("testing request", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		log.Fatal("request status", uploadResp.Status)
	}
	fmt.Println(uploadResp)

}
func LibGenDownload(md5_hash_list []title_and_hash) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	for _, title_and_hash := range md5_hash_list {

		UploadFromMemory(title_and_hash, *client)

		// 	get_request_url := "https://cdn3.booksdl.org/get.php?" + title_and_hash.hash

		// 	resp, err := client.Get(get_request_url)

		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	defer resp.Body.Close()

		// 	filename := title_and_hash.title + ".pdf"

		// 	bar := progressbar.DefaultBytes(
		// 		resp.ContentLength,
		// 		"downloading",
		// 	)

		// 	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		// return nil
	}
	return nil
}
