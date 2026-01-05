package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
}

func PostRequest(url string, body []byte) (int, map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if resp.ContentLength != 0 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			// ตรวจสอบว่า Error ไม่ใช่ EOF (ซึ่งมักเกิดเมื่อ Response Body ว่างเปล่า)
			if err.Error() != "EOF" {
				return resp.StatusCode, nil, errors.New("failed to decode response body")
			}
		}
	}

	return resp.StatusCode, result, nil
}
