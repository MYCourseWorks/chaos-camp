package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// ErrColor comment
const ErrColor = "\033[1;31m%s\033[0m"

// InfoColor comment
const InfoColor = "\033[1;34m%s\033[0m"

func extractPrettyJSON(reader io.ReadCloser) string {
	var body []byte
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return ""
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		fmt.Printf("error = %v\n", err)
		return ""
	}

	return string(prettyJSON.Bytes())
}

func handleResp(resp *http.Response) string {

	defer resp.Body.Close()

	fmt.Printf("METHOD \033[1;33m%s\033[0m ", resp.Request.Method)
	if resp.StatusCode >= 400 {
		fmt.Printf(ErrColor, fmt.Sprintf("StatusCode = %d\n", resp.StatusCode))
	} else {
		fmt.Printf(InfoColor, fmt.Sprintf("StatusCode = %d\n", resp.StatusCode))
	}

	// If response is not json, just return the status
	if resp.Header.Get("content-type") != "application/json; charset=utf-8" {
		return ""
	}

	return extractPrettyJSON(resp.Body)
}

// SendGet comment
func SendGet(url string) string {
	resp, err := http.Get(fmt.Sprintf(url))
	if err != nil {
		return ""
	}

	return handleResp(resp)
}

// SendPost comment
func SendPost(url string, body map[string]string) (string, location string) {

	if body == nil {
		resp, err := http.Post(fmt.Sprintf(url), "application/json; charset=utf-8", nil)
		if err != nil {
			return "", ""
		}
		return handleResp(resp), ""
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", ""
	}

	resp, err2 := http.Post(fmt.Sprintf(url), "application/json; charset=utf-8", bytes.NewBuffer(bodyJSON))
	if err2 != nil {
		return "", ""
	}

	location = resp.Header.Get("Location")
	return handleResp(resp), location
}

// SendPut comment
func SendPut(url string, body interface{}) string {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return ""
	}

	// Create request
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return ""
	}

	// Fetch Request
	client := &http.Client{}
	req.Header.Add("content-type", "application/json; charset=utf-8")
	resp, err := client.Do(req)

	return handleResp(resp)
}

// SendDelete comment
func SendDelete(url string) string {
	// Create request
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return ""
	}

	// Fetch Request
	client := &http.Client{}
	req.Header.Add("content-type", "application/json; charset=utf-8")
	resp, err := client.Do(req)

	return handleResp(resp)
}
