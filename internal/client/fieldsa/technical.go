package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetWorkOrderTechnicalData(
	ctx context.Context,
	workOrderNumber string,
) (*WorkOrderTechnicalData, error) {

	path := fmt.Sprintf(
		"/api/dispatch/getworkordertechnicaldata/%s",
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

	var result WorkOrderTechnicalData

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order technical data: %w",
			err,
		)
	}

	return &result, nil
}
