package fieldsa

import (
	"context"
	"encoding/json"
	"fmt"
)

type MaterialPagingRequest struct {
	First        int            `json:"first"`
	Rows         int            `json:"rows"`
	SortOrder    int            `json:"sortOrder"`
	Filters      map[string]any `json:"filters"`
	GlobalFilter any            `json:"globalFilter"`
}

func newMaterialPagingRequest() MaterialPagingRequest {
	return MaterialPagingRequest{
		First:        0,
		Rows:         10,
		SortOrder:    1,
		Filters:      map[string]any{},
		GlobalFilter: nil,
	}
}

func (c *Client) GetWorkOrderEquipment(
	ctx context.Context,
	workOrderNumber string,
) ([]WorkOrderMaterial, error) {

	path := fmt.Sprintf(
		"/api/material/getworkorderequipmentpaging/%s",
		workOrderNumber,
	)

	payload := newMaterialPagingRequest()

	data, err := c.Post(
		ctx,
		c.dispatchBaseURL,
		path,
		payload,
	)

	if err != nil {
		return nil, err
	}

	var result []WorkOrderMaterial

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order equipment: %w",
			err,
		)
	}

	return result, nil
}

func (c *Client) GetWorkOrderMaterial(
	ctx context.Context,
	workOrderNumber string,
) ([]WorkOrderMaterial, error) {

	path := fmt.Sprintf(
		"/api/material/getworkordermaterialpaging/%s",
		workOrderNumber,
	)

	payload := newMaterialPagingRequest()

	data, err := c.Post(
		ctx,
		c.dispatchBaseURL,
		path,
		payload,
	)

	if err != nil {
		return nil, err
	}

	var result []WorkOrderMaterial

	if err := json.Unmarshal(
		data,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"decode work order material: %w",
			err,
		)
	}

	return result, nil
}
