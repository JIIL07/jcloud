package cookies

import (
	"encoding/json"
	"net/http"
	"os"
)

func Serialize(cookies []*http.Cookie) (string, error) {
	data, err := json.Marshal(cookies)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Deserialize(data string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	err := json.Unmarshal([]byte(data), &cookies)
	if err != nil {
		return nil, err
	}
	return cookies, nil
}

func WriteToFile(filename, data string) error {
	return os.WriteFile(filename, []byte(data), 0644)
}

func ReadFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
