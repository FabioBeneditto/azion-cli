package edgefunction

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
	"go.uber.org/zap"
)

const javascript = "javascript"

func (c *Client) Get(ctx context.Context, id int64) (EdgeFunctionResponse, error) {
	logger.Debug("Get Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while getting an edge function", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)

	httpResp, err := request.Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting an edge function", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeFunctionResponse, error) {
	// Although there's only one option, the API requires the `language` field.
	// Hard-coding javascript for now
	logger.Debug("Create Edge Function")
	req.CreateEdgeFunctionRequest.SetLanguage(javascript)

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).CreateEdgeFunctionRequest(req.CreateEdgeFunctionRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an edge function", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (EdgeFunctionResponse, error) {
	logger.Debug("Update Edge Function")
	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdPatch(ctx, req.Id).PatchEdgeFunctionRequest(req.PatchEdgeFunctionRequest)

	edgeFuncResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating an edge function", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return edgeFuncResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.ListEdgeFunctionResponse, error) {
	logger.Debug("List Edge Functions")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.EdgeFunctionsApi.EdgeFunctionsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing the edge functions", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
