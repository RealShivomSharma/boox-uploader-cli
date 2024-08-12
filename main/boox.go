package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Routes struct {
	Library  string
	Images   string
	Recent   string
	Music    string
	Audios   string
	Download string
	Storage  string
	Device   string
	Upload   string
}

type DeviceDetails struct {
	Host         string `json:"host"`
	ID           string `json:"id"`
	MAC          string `json:"mac"`
	Model        string `json:"model"`
	StorageTotal string `json:"storageTotal"`
	StorageUsed  string `json:"storageUsed"`
	DeviceType   string `json:"type"`
}

type LibraryResponse struct {
	BookCount          int           `json:"bookCount"`
	LibraryCount       int           `json:"libraryCount"`
	VisibleBookList    []LibraryBook `json:"visibleBookList"`
	VisibleLibraryList []Library     `json:"visibleLibraryList"`
}

type LibraryBook struct {
	Title string `json:"title"`
}

type Library struct {
	Title string `json:"title"`
}

type LibraryQueryParams struct {
	Limit           int
	Offset          int
	SortBy          string
	Order           string
	LibraryUniqueID string
}

func getBooxURL() string {
	config := create_config()
	return config.boox_url
}

func getRoutes(url string) Routes {

	libraryRoute := url + "/api/library"
	imagesRoute := url + "/#/pc/image"
	recentRoute := url + "/#/pc/recent"
	musicRoute := url + "/#/pc/music"
	audiosRoute := url + "/#/pc/audio"
	downloadRoute := url + "/#/pc/download"
	storageRoute := url + "/api/storage"
	deviceRoute := url + "/api/device"
	uploadRoute := libraryRoute + "/upload"

	return Routes{
		Library:  libraryRoute,
		Images:   imagesRoute,
		Recent:   recentRoute,
		Music:    musicRoute,
		Audios:   audiosRoute,
		Download: downloadRoute,
		Storage:  storageRoute,
		Device:   deviceRoute,
		Upload:   uploadRoute,
	}

}

func getBooxDetails() {
	url := getBooxURL()
	fmt.Println("URL", url)
	deviceRoute := getRoutes(url).Device
	resp, err := http.Get(deviceRoute)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var device DeviceDetails
	if err := json.Unmarshal(body, &device); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println(device)

}

func constructLibraryURL(params LibraryQueryParams) (string, error) {
	library_url := getRoutes(getBooxURL()).Library

	u, err := url.Parse(library_url)

	if err != nil {
		return "", err
	}

	q := u.Query()
	args := make(map[string]interface{})
	args["limit"] = params.Limit
	args["offset"] = params.Limit
	args["sortBy"] = params.SortBy
	args["order"] = params.Order
	if params.LibraryUniqueID != "" {
		args["libraryUniqueId"] = params.LibraryUniqueID
	} else {
		args["libraryUniqueID"] = nil
	}

	argsJSON, err := json.Marshal(args)

	if err != nil {
		return "", err
	}

	q.Set("args", string(argsJSON))
	u.RawQuery = q.Encode()

	return u.String(), nil

}

func getLibraryTitlesWithParams(params LibraryQueryParams) ([]string, error) {
	queryURL, err := constructLibraryURL(params)
	if err != nil {
		return nil, fmt.Errorf("error constructing URL: %v", err)
	}
	resp, err := http.Get(queryURL)

	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var libraryResp LibraryResponse
	err = json.Unmarshal(body, &libraryResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON %v", err)
	}

	var titles []string

	for _, book := range libraryResp.VisibleBookList {
		titles = append(titles, book.Title)
	}

	for _, library := range libraryResp.VisibleLibraryList {
		titles = append(titles, library.Title)
	}

	return titles, nil

}

func UploadFromFile(uploadURL string, filePath string) error {
	filePath = "/Users/shivom/boox-uploader-cli/boox-uploader-cli/main/Testing Lifting Bodies at Edwards.pdf"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileName := filepath.Base(filePath)

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("error creating form file %v", err)
	}

	_, err = io.Copy(part, bytes.NewReader(fileContent))
	if err != nil {
		return fmt.Errorf("failed to copy file content %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer %v", err)
	}

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to submit request %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	fmt.Printf("Upload response: %s\n", string(responseBody))

	return nil

}
