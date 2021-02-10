package client

import (
	"context"
	"fmt"

	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/permission"
	"github.com/tidepool-org/platform/platform"
	"github.com/tidepool-org/platform/request"
)

type Client struct {
	client                 *platform.Client
	permissionClientConfig httpClientConfig
}

type httpClientConfig struct {
	clientType string
	urlPrefix  string
	httpMethod string
}

type CoastguardRequestBody struct {
	Service       string `json:"service"`
	RequestUserID string `json:"requestUserId"`
	TargetUserID  string `json:"targetUserId"`
}
type CoastguardResponseBody struct {
	Authorized bool   `json:"authorized"`
	Route      string `json:"route"`
}

var (
	permissionClientTypes = map[string]httpClientConfig{
		"gatekeeper": {
			clientType: "gatekeeper",
			urlPrefix:  "access",
			httpMethod: "GET",
		},
		"coastguard": {
			clientType: "coastguard",
			urlPrefix:  "v1/data/backloops/platform",
			httpMethod: "POST",
		},
	}
)

func New(config *platform.Config, authorizeAs platform.AuthorizeAs, permissionType string) (*Client, error) {
	clnt, err := platform.NewClient(config, authorizeAs)
	if err != nil {
		return nil, err
	}
	permsClientConfig, ok := permissionClientTypes[permissionType]
	if !ok {
		return nil, fmt.Errorf("unknown permission client type: %s", permissionType)
	}
	return &Client{
		client:                 clnt,
		permissionClientConfig: permsClientConfig,
	}, nil
}

func (c *Client) GetUserPermissions(ctx context.Context, requestUserID string, targetUserID string) (permission.Permissions, error) {

	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if requestUserID == "" {
		return nil, errors.New("request user id is missing")
	}
	if targetUserID == "" {
		return nil, errors.New("target user id is missing")
	}
	result := permission.Permissions{}

	if requestUserID == targetUserID {
		result[permission.Owner] = permission.Permission{}
		return permission.FixOwnerPermissions(result), nil
	}

	authConfig := c.permissionClientConfig
	if authConfig.clientType == "gatekeeper" {
		url := c.client.ConstructURL(authConfig.urlPrefix, targetUserID, requestUserID)
		result := permission.Permissions{}
		if err := c.client.RequestData(ctx, authConfig.httpMethod, url, nil, nil, &result); err != nil {
			if request.IsErrorResourceNotFound(err) {
				return nil, request.ErrorUnauthorized()
			}
			return nil, err
		}
		return permission.FixOwnerPermissions(result), nil
	}

	if c.permissionClientConfig.clientType == "coastguard" {
		url := c.client.ConstructURL(authConfig.urlPrefix)
		coastguardResponse := CoastguardResponseBody{}
		requestBody := CoastguardRequestBody{
			Service:       "platform",
			RequestUserID: requestUserID,
			TargetUserID:  targetUserID,
		}
		if err := c.client.RequestData(ctx, authConfig.httpMethod, url, nil, &requestBody, &coastguardResponse); err != nil {
			return nil, err
		}
		if coastguardResponse.Authorized {
			result[permission.Read] = permission.Permission{}
		}
		return permission.FixOwnerPermissions(result), nil
	}
	return permission.FixOwnerPermissions(result), nil
}
