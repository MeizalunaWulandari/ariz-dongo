package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetWorkOrderDispatch(
	ctx context.Context,
	date string,
	pagination string,
	region string,
) ([]byte, error) {

	path := fmt.Sprintf(
		"/api/dispatch/getworkorderdispatch/%s/%s/%s",
		date,
		pagination,
		region,
	)

	return c.Get(
		ctx,
		c.dispatchBaseURL,
		path,
	)
}

func (c *Client) GetWorkOrderDetail(
	ctx context.Context,
	workOrderNumber string,
) (*WorkOrderDetail, error) {

	path := fmt.Sprintf(
		"/api/dispatch/getworkorderdetail/%s",
		workOrderNumber,
	)

	data, err := c.Get(
		ctx,
		c.dispatchBaseURL,
		path,
	)

	if err != nil {
		return nil, err
	}

	var result WorkOrderDetail

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order detail: %w",
			err,
		)
	}

	return &result, nil
}
