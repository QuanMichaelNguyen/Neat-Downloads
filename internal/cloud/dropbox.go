package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type DropboxClient struct {
	AppKey       string
	AppSecret    string
	AccessToken  string
	RefreshToken string
}

// Make a new DropBox Client
func NewDropboxClient(appKey, appSecret string) *DropboxClient {
	return &DropboxClient{
		AppKey:    appKey,
		AppSecret: appSecret,
	}
}

// GetAuthURL returns the OAuth URL --> user can authorize the app

func (d *DropboxClient) GetAuthURL() string {
	params := url.Values{}
	params.Add("client_id", d.AppKey)
	params.Add("response_type", "code")
	params.Add("redirect_uri", "http://localhost:8080/auth/callback")

	return "https://www.dropbox.com/oauth2/authorize?" + params.Encode()
}

// ExchangeCodeForToken exchanges the authorization code for an access token

func (d *DropboxClient) ExchangeCodeForToken(code string) error {
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", d.AppKey)
	data.Set("client_secret", d.AppSecret)
	data.Set("redirect_uri", "http://localhost:8080/auth/callback")

	res, err := http.PostForm("https://api.dropboxapi.com/oauth2/token", data)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	if result.Error != "" {
		return fmt.Errorf("dropbox error: %s", result.Error)
	}

	d.AccessToken = result.AccessToken   // User's access token
	d.RefreshToken = result.RefreshToken // User's refresh token
	return nil
}

// ListFiles lists files in a Dropbox directory
func (d *DropboxClient) ListFiles(path string) ([]string, error) {
	url := "https://api.dropboxapi.com/2/files/list_folder"

	requestBody, _ := json.Marshal(map[string]interface{}{
		"path":      path,
		"recursive": false,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+d.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Entries []struct {
			Name string `json:"name"`
			Path string `json:"path_display"`
		} `json:"entries"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range result.Entries {
		files = append(files, entry.Path)
	}

	return files, nil
}

// MoveFile moves a file in Dropbox
func (d *DropboxClient) MoveFile(fromPath, toPath string) error {
	url := "https://api.dropboxapi.com/2/files/move_v2"

	requestBody, _ := json.Marshal(map[string]interface{}{
		"from_path":  fromPath,
		"to_path":    toPath,
		"autorename": true,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Authorization", "Bearer "+d.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to move file: %s", string(body))
	}

	return nil
}

// UploadFile uploads a file to Dropbox
func (d *DropboxClient) UploadFile(localPath, dropboxPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	url := "https://content.dropboxapi.com/2/files/upload"

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return err
	}

	apiArg, _ := json.Marshal(map[string]interface{}{
		"path":       dropboxPath,
		"mode":       "add",
		"autorename": true,
		"mute":       false,
	})

	req.Header.Set("Authorization", "Bearer "+d.AccessToken)
	req.Header.Set("Dropbox-API-Arg", string(apiArg))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.ContentLength = fileInfo.Size()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upload file: %s", string(body))
	}

	return nil
}

func (d *DropboxClient) DownloadFile(dropboxPath, localPath string) error {
	url := "https://content.dropboxapi.com/2/files/download"

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+d.AccessToken)

	apiArg, _ := json.Marshal(map[string]string{
		"path": dropboxPath,
	})
	req.Header.Set("Dropbox-API-Arg", string(apiArg))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to download file: %s", string(body))
	}

	// Create directory if it doesn't exist
	err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	if err != nil {
		return err
	}

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
