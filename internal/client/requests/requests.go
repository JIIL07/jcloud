package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/JIIL07/jcloud/internal/client/app"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/requests/jreq"
	"github.com/JIIL07/jcloud/pkg/params"
)

type UserData struct {
	UserID   int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Protocol string `db:"hashprotocol" json:"hashprotocol"`
	Admin    int    `db:"admin" json:"admin"`
}

var URL = "https://jcloud.up.railway.app"

func Login(u *UserData) (*http.Response, error) {
	jsonData, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	p := params.NewParams()
	p.Set("type", "POST")
	p.Set("url", URL+"/api/v1/login")
	p.Set("body", bytes.NewBuffer(jsonData))

	resp, err := jreq.Post(&p)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}

func UploadFile(a *app.ClientContext, f *[]models.File) (*http.Response, error) {
	jsonData, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data: %w", err)
	}

	req, err := http.NewRequest("POST", URL+"/api/v1/files/upload", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	cookieString, err := ReadFromFile(a.Paths.P.Jcookie.Name())
	if err != nil {
		return nil, fmt.Errorf("error reading cookies: %w", err)
	}
	c, err := Deserialize(cookieString)
	if err != nil {
		return nil, fmt.Errorf("error deserializing cookies: %w", err)
	}
	for i := range c {
		req.AddCookie(c[i])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}

	return resp, nil
}

func DeleteFile(f *models.File) error {
	baseURL := URL + "/api/v1/files/delete"
	p := url.Values{}
	p.Add("filename", f.Meta.Name)
	fullURL := fmt.Sprintf("%s?%s", baseURL, p.Encode())

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func GetFiles() error {
	p := params.NewParams()
	p.Set("type", "GET")
	p.Set("url", URL+"/api/v1/files/get")

	resp, err := jreq.Get(&p)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	response := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(response)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	fmt.Println(string(response))
	return nil
}

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
	return os.WriteFile(filename, []byte(data), 0600)
}

func ReadFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
