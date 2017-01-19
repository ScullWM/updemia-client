package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/howeyc/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const defaultDestinationPath string = "/tmp/screendirect/"
const apiPostEndpoint string = "http://www.updemia.com/api/v1/post"

type UpdemiaResponse struct {
	Key string
	Url string
}

func main() {
	saveNotificationLogo()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					sendFile(getNewFilePath(ev))
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	dest := getDestinationPath()
	err = watcher.Watch(dest)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}

func getDestinationPath() string {
	destination_path := defaultDestinationPath

	if len(os.Args) > 1 {
		if len(os.Args[1]) > 0 {
			destination_path = os.Args[1]
		}
	}

	os.MkdirAll(destination_path, os.ModePerm)
	log.Printf("Send to destination_path: %+v", destination_path)

	return destination_path
}

func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func sendFile(filename string) {
	request, err := newfileUploadRequest(apiPostEndpoint, "file", filename)
	if err == nil {
		client := &http.Client{}

		resp, err := client.Do(request)
		if err != nil {
			notifyUserFail()
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var response UpdemiaResponse
		json.Unmarshal(body, &response)

		err = clipboard.WriteAll(response.Url)
		notifyUserSuccess(response.Url)

		log.Printf("New upload: %+v", response.Url)
	}
}

func getNewFilePath(ev *fsnotify.FileEvent) string {
	strEvent := fmt.Sprintf("%s", ev)
	endingIndex := strings.Index(strEvent, "\":")
	filename := strEvent[1:endingIndex]
	beginningIndex := strings.LastIndex(filename, "/") + 1

	if string(filename[beginningIndex]) != "." {
		return filename
	}

	return ""
}
