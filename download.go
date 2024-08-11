package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// func followRedirects(url string) (string, error) {
// 	var resp http
// }

func LibGenDownload(md5_hash_list []string) error {
	for _, hash := range md5_hash_list {

		get_request_url := "https://libgen.li/get.php?" + hash

		resp, err := http.Get(get_request_url)

		if err != nil {
			log.Fatal("Failed to download book", err)
		}
		fmt.Println(resp.Header.Get("Location"))
		filename := "test.pdf"
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil

	// https: //libgen.li/get.php?md5=311a42b9ba4d0e77d6bca615b1b333b1&key=BQ54PZGW0VWLOL02
}
