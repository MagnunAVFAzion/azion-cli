package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aziontech/azion-cli/pkg/cmd/version"

	"os"

	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/storageapi"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.AddDefaultHeader("Authorization", "Token "+token)
	conf.AddDefaultHeader("Accept", "application/json; version=3")
	conf.AddDefaultHeader("Content-Type", "application/json")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Upload(ctx context.Context, versionID, path string, file *os.File) error {
	c.apiClient.GetConfig().DefaultHeader["Content-Disposition"] = fmt.Sprintf("attachment; filename=\"%s\"", path)
	req := c.apiClient.DefaultApi.StorageVersionIdPost(ctx, versionID).XAzionStaticPath(path).Body(file)
	_, httpResp, err := req.Execute()
	if err != nil {
		fmt.Println(err.Error())
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
