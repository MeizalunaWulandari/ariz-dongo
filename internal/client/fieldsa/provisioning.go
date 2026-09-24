package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetWorkOrderProvisioning(
	ctx context.Context,
	workOrderNumber string,
) (*WorkOrderProvisioning, error) {

	path := fmt.Sprintf(
		"/api/fieldsa/workorder/getworkorderprovisioning/%s",
		workOrderNumber,
	)

	data, err := c.Get(
		ctx,
		c.baseURL,
		path,
	)

	if err != nil {
		return nil, err
	}

	var result WorkOrderProvisioning

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order provisioning: %w",
			err,
		)
	}

	return &result, nil
}
