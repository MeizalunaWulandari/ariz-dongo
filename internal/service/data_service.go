package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/client/fieldsa"
)

type DataService struct {
	fieldsaClient *fieldsa.Client
}

func NewDataService(
	fieldsaClient *fieldsa.Client,
) *DataService {
	return &DataService{
		fieldsaClient: fieldsaClient,
	}
}

func (s *DataService) GetData(
	ctx context.Context,
	date string,
	pagination string,
	region string,
) ([]byte, error) {

	// Default tanggal hari ini.
	if date == "" {
		date = time.Now().Format(
			"2006-01-02T00:00:00",
		)
	}

	// Default pagination Fieldsa.
	if pagination == "" {
		pagination = "MSwyLDMsNSw2LDc="
	}

	// Default region.
	if region == "" {
		region = "KALIMANTAN"
	}

	data, err := s.fieldsaClient.GetWorkOrderDispatch(
		ctx,
		date,
		pagination,
		region,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get work order dispatch: %w",
			err,
		)
	}

	return data, nil
}
