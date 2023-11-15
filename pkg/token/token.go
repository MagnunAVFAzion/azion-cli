package token

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/constants"
)

const settingsFilename = "settings.toml"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	endpoint string
	client   HTTPClient
	filePath string
	valid    bool
	out      io.Writer
}

type Settings struct {
	Token string
	UUID  string
}

type Config struct {
	Client HTTPClient
	Out    io.Writer
}

type Response struct {
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

func New(c *Config) (*Token, error) {
	dir, err := config.Dir()
	if err != nil {
		return nil, err
	}

	return &Token{
		client:   c.Client,
		endpoint: constants.AuthURL,
		filePath: filepath.Join(dir, settingsFilename),
		out:      c.Out,
	}, nil
}

func (t *Token) Validate(token *string) (bool, error) {
	logger.Debug("validated token", zap.Any("Token", *token))

	req, err := http.NewRequest("GET", utils.Concat(t.endpoint, "/user/me"), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Accept", "application/json; version=3")
	req.Header.Add("Authorization", "token "+*token)

	resp, err := t.client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	t.valid = true

	return true, nil
}

func (t *Token) Save(b []byte) (string, error) {
	logger.Debug("save token", zap.Any("byte", string(b)))
	filePath, err := config.Dir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(t.filePath, b, 0777)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}

	return t.filePath, nil
}

func ReadFromDisk() (string, error) {
	dir, err := config.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get token dir: %w", err)
	}

	fileData, err := os.ReadFile(filepath.Join(dir, settingsFilename))
	if err != nil {
		return "", err
	}

	var settings Settings
	err = toml.Unmarshal(fileData, &settings)
	if err != nil {
		return "", fmt.Errorf("failed parse byte to struct settings: %w", err)
	}

	return settings.Token, nil
}

func (t *Token) Create(b64 string) (*Response, error) {
	logger.Debug("create token", zap.Any("base64", b64))
	req, err := http.NewRequest(http.MethodPost, utils.Concat(t.endpoint, "/tokens"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json; version=3")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", utils.Concat("Basic ", b64))

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
