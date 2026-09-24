package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetWorkOrderTimeline(
	ctx context.Context,
	workOrderNumber string,
) ([]WorkOrderTimeline, error) {

	path := fmt.Sprintf(
		"/api/dispatch/getworkordertimeline/%s",
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

	var result []WorkOrderTimeline

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order timeline: %w",
			err,
		)
	}

	return result, nil
}
