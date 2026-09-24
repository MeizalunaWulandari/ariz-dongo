package service

import (
	"context"
	"fmt"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/client/fieldsa"
)

type PhotoService struct {
	fieldsaClient *fieldsa.Client
}

func NewPhotoService(
	fieldsaClient *fieldsa.Client,
) *PhotoService {
	return &PhotoService{
		fieldsaClient: fieldsaClient,
	}
}

type WorkOrderPhoto struct {
	WorkOrderNumber  string `json:"workOrderNumber"`
	DocumentTypeName string `json:"documentTypeName"`
	DocumentFullPath string `json:"documentFullPath"`
	Note             string `json:"note"`
	Created          string `json:"created"`
	CreatedBy        string `json:"createdBy"`
	DocumentTypeID   int    `json:"documentTypeId"`
}

func (s *PhotoService) GetPhotos(
	ctx context.Context,
	workOrderNumber string,
) ([]WorkOrderPhoto, error) {

	if workOrderNumber == "" {
		return nil, fmt.Errorf(
			"workOrderNumber wajib diisi",
		)
	}

	documents, err :=
		s.fieldsaClient.GetWorkOrderDocuments(
			ctx,
			workOrderNumber,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"get work order documents: %w",
			err,
		)
	}

	photos := make(
		[]WorkOrderPhoto,
		0,
		len(documents),
	)

	for _, document := range documents {

		// Hanya ambil dokumen yang mempunyai URL.
		if document.DocumentFullPath == "" {
			continue
		}

		photos = append(
			photos,
			WorkOrderPhoto{
				WorkOrderNumber:  document.WorkOrderNumber,
				DocumentTypeName: document.DocumentTypeName,
				DocumentFullPath: document.DocumentFullPath,
				Note:             document.Note,
				Created:          document.Created,
				CreatedBy:        document.CreatedBy,
				DocumentTypeID:   document.DocumentTypeID,
			},
		)
	}

	return photos, nil
}
