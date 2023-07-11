package domains

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/domains"
	"go.uber.org/zap"
)

type Client struct {
	apiClient *sdk.APIClient
}

type CreateRequest struct {
	sdk.CreateDomainRequest
}

type UpdateRequest struct {
	sdk.UpdateDomainRequest
	Id int64
}

type DomainResponse interface {
	GetId() int64
	GetName() string
	GetDomainName() string
	GetCnames() []string
	GetCnameAccessOnly() bool
	GetDigitalCertificateId() int64
	GetEdgeApplicationId() int64
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 30 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Get(ctx context.Context, id string) (DomainResponse, error) {
	logger.Debug("Get Domain")
	req := c.apiClient.DomainsApi.GetDomain(ctx, id)
	logger.Debug("request", zap.Any("struct", req))
	res, httpResp, err := req.Execute()
	logger.Debug("response", zap.Any("struct", res), zap.Any("http", httpResp), zap.Error(err))
	if err != nil {
		logger.Error("error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &res.Results, nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (DomainResponse, error) {
	logger.Debug("Create Domain")
	request := c.apiClient.DomainsApi.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)
	logger.Debug("request", zap.Any("struct", request))
	domainsResponse, httpResp, err := request.Execute()
	logger.Debug("response", zap.Any("struct", domainsResponse), zap.Any("http", httpResp), zap.Error(err))
	if err != nil {
		logger.Error("error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &domainsResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (DomainResponse, error) {
	logger.Debug("Update Domain")
	str := strconv.FormatInt(req.Id, 10)
	request := c.apiClient.DomainsApi.UpdateDomain(ctx, str).UpdateDomainRequest(req.UpdateDomainRequest)
	logger.Debug("request", zap.Any("struct", request))

	domainsResponse, httpResp, err := request.Execute()
	logger.Debug("response", zap.Any("struct", domainsResponse), zap.Any("http", httpResp), zap.Error(err))

	if err != nil {
		logger.Error("error", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &domainsResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.DomainResponseWithResults, error) {
	// different from other APIs, domains will return internal server error if order by is empty
	logger.Debug("List Domains")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.DomainsApi.GetDomains(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		Execute()
	logger.Debug("response", zap.Any("struct", resp), zap.Any("http", httpResp), zap.Error(err))

	if err != nil {
		logger.Error("error", zap.Error(err))
		return &sdk.DomainResponseWithResults{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Domain")
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.DomainsApi.DelDomain(ctx, str)
	logger.Debug("request", zap.Any("struct", req))

	httpResp, err := req.Execute()
	logger.Debug("response", zap.Any("http", httpResp), zap.Error(err))
	if err != nil {
		logger.Error("error", zap.Error(err))
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
