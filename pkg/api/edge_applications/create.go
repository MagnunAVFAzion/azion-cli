package edge_applications

import (
	"context"
	"io"

	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type ResponseCreate interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetApplicationAcceleration() bool
	GetCaching() bool
	GetDeliveryProtocol() string
	GetDeviceDetection() bool
	GetEdgeFirewall() bool
	GetEdgeFunctions() bool
	GetHttpPort() interface{}
	GetHttpsPort() interface{}
	GetImageOptimization() bool
	GetL2Caching() bool
	GetLoadBalancer() bool
	GetMinimumTlsVersion() string
	GetRawLogs() bool
	GetWebApplicationFirewall() bool
}

type EdgeApplicationsResponse interface {
	GetId() int64
	GetName() string
}

type CreateRequest struct {
	sdk.CreateApplicationRequest
}

func (c *Client) Create(ctx context.Context, req *CreateRequest,
) (EdgeApplicationsResponse, error) {
	logger.Debug("Create Edge Application")
	request := c.apiClient.EdgeApplicationsMainSettingsApi.
		EdgeApplicationsPost(ctx).CreateApplicationRequest(req.CreateApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an edge application", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(err))
				return nil, utils.ErrorPerStatusCode(httpResp, err)
			}
			logger.Debug("", zap.Any("Body", string(body)))
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}