package util

import (
	"os"
	"net/http"
	"bytes"
	"fmt"
	"io"
)

func HttpExists(uri, username, password string) (bool, error) {
	req, err := http.NewRequest("HEAD", uri, http.NoBody)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return false, nil
	}

	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("HTP error downloading file [%v]", resp.StatusCode)
	}

	return true, nil
}

func DownloadFile(uri, username, password, filepath string) error {
	req, err := http.NewRequest("GET", uri, http.NoBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writer the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil  {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTP error downloading file [%v]", resp.StatusCode)
	}

	return nil
}


func UploadFile(uri, username, password, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", uri, file)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := &bytes.Buffer{}

	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error [%v]: %v", resp.StatusCode, body)
	}

	return nil
}

