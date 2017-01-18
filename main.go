package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/deckarep/gosx-notifier"
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

type UpdemiaResponse struct {
	Key string
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
	destination_path := "/tmp/screendirect/"

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
	request, err := newfileUploadRequest("http://www.updemia.com/api/v1/post", "file", filename)
	if err == nil {
		client := &http.Client{}

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var response UpdemiaResponse
		json.Unmarshal(body, &response)
		url := "http://www.updemia.com/media/" + response.Key

		err = clipboard.WriteAll(url)
		notifyUserSuccess(url)

		log.Printf("New upload: %+v", url)
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

func notifyUserSuccess(url string) {
	note := gosxnotifier.NewNotification("Upload complete")
	note.AppIcon = "/tmp/m1UIjW1.png"
	note.Title = "On updemia.com"
	// note.Sound = gosxnotifier.Basso
	note.Link = url
	err := note.Push()

	if err != nil {
		log.Println("notification error")
	}
}

func saveNotificationLogo() {
	img, _ := os.Create("/tmp/m1UIjW1.png")
	defer img.Close()

	resp, _ := http.Get("http://www.updemia.com/images/logo.png")
	defer resp.Body.Close()

	_, err := io.Copy(img, resp.Body)
	if err != nil {
		log.Println("Error getting logo")
	}
}
