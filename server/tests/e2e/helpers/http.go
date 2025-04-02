package helpers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type TestClient struct {
    BaseURL string
    Token   string
    Client  *http.Client
}

func NewTestClient(baseURL string) *TestClient {
    return &TestClient{
        BaseURL: baseURL,
        Client:  &http.Client{},
    }
}

func (c *TestClient) DoRequest(method, path string, body interface{}, target interface{}) error {
    var bodyReader io.Reader
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return fmt.Errorf("failed to marshal request body: %v", err)
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
    if err != nil {
        return fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Authorization", "Bearer "+c.Token)
    }

    resp, err := c.Client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to execute request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
    }

    if target != nil {
        if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
            return fmt.Errorf("failed to decode response: %v", err)
        }
    }

    return nil
}

func (c *TestClient) Register(username, password string) error {
    body := map[string]string{
        "username": username,
        "password": password,
    }
    return c.DoRequest("POST", "/api/register", body, nil)
}

func (c *TestClient) Login(username, password string) error {
    body := map[string]string{
        "username": username,
        "password": password,
    }
    var response struct {
        Token string `json:"token"`
    }
    if err := c.DoRequest("POST", "/api/login", body, &response); err != nil {
        return err
    }
    c.Token = response.Token
    return nil
}