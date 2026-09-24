package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) GetWorkOrderDocuments(
	ctx context.Context,
	workOrderNumber string,
) ([]WorkOrderDocument, error) {

	path := fmt.Sprintf(
		"/api/dispatch/getallworkorderdocument/%s",
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

	var documents []WorkOrderDocument

	if err := json.Unmarshal(
		data,
		&documents,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order documents: %w",
			err,
		)
	}

	return documents, nil
}
