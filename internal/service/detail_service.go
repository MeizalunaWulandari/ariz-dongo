package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/MeizalunaWulandari/ariz-dongo/internal/client/fieldsa"
	"github.com/MeizalunaWulandari/ariz-dongo/internal/logger"
	"golang.org/x/sync/errgroup"
)

type DetailService struct {
	fieldsaClient *fieldsa.Client
}

func NewDetailService(
	fieldsaClient *fieldsa.Client,
) *DetailService {
	return &DetailService{
		fieldsaClient: fieldsaClient,
	}
}

type WorkOrderDetailResponse struct {
	WorkOrderNumber       string            `json:"workOrderNumber"`
	ReferenceCode         string            `json:"referenceCode"`
	CustomerName          string            `json:"customerName"`
	CircuitID             string            `json:"circuitId"`
	CID                   string            `json:"cid"`
	EndCustomerName       string            `json:"endCustomerName"`
	EndCustomerAddress    string            `json:"endCustomerAddress"`
	EndCustomerLongtitude string            `json:"endCustomerLongtitude"`
	HomepassIDCoordinate  string            `json:"homepassIdCoordinate"`
	DeviceAllocation      string            `json:"deviceAllocation"`
	WorkOrderStatusItem   string            `json:"workOrderStatusItem"`
	VendorName            string            `json:"vendorName"`
	TechnicianName        string            `json:"technicianName"`
	Equipment             []Equipment       `json:"equipment"`
	Material              []Material        `json:"material"`
	Provisioning          *Provisioning     `json:"provisioning"`
	TechnicalData         *TechnicalData    `json:"technicalData"`
	Timeline              *TimelineResponse `json:"timeline"`
}

type Equipment struct {
	Description       string `json:"description"`
	SerialNumber      string `json:"serialNumber"`
	MaterialStatus    string `json:"materialStatus"`
	MaterialCondition string `json:"materialCondition"`
}

type Material struct {
	Description       string `json:"description"`
	MaterialStatus    string `json:"materialStatus"`
	MaterialCondition string `json:"materialCondition"`
}

type Provisioning struct {
	SerialNumber    string `json:"serialNumber"`
	EquipmentModel  string `json:"equipmentModel"`
	NetworkEntityID string `json:"networkEntityId"`
	Epon            string `json:"epon"`
	VLAN            string `json:"vlan"`
	ReceivedPower   string `json:"receivedPower"`
}

type TechnicalData struct {
	SiteID         string `json:"siteId"`
	SplitterID     string `json:"splitterId"`
	TotalPort      string `json:"totalPort"`
	PortNumber     string `json:"portNumber"`
	AvailablePort  string `json:"availablePort"`
	DPFOCoordinate string `json:"dpfoCoordinate"`
}

type TimelineResponse struct {
	ActionLatitude         string `json:"actionLatitude"`
	ActionLongtitude       string `json:"actionLongtitude"`
	LatitudeLongtitudeLink string `json:"latitudeLongtitudeLink"`
}

func (s *DetailService) GetDetail(
	ctx context.Context,
	workOrderNumber string,
) (*WorkOrderDetailResponse, error) {

	if workOrderNumber == "" {
		return nil, fmt.Errorf(
			"workOrderNumber wajib diisi",
		)
	}

	/*
		==================================================
		RESULT VARIABLES
		==================================================

		Semua hasil disimpan di variabel masing-masing.

		Jadi walaupun goroutine selesai dengan urutan:
		Material → Timeline → Detail → Equipment → ...

		Response tetap akan kita susun dengan urutan
		yang sudah ditentukan di bawah.
	*/

	var (
		detail        *fieldsa.WorkOrderDetail
		equipment     []fieldsa.WorkOrderMaterial
		material      []fieldsa.WorkOrderMaterial
		provisioning  *fieldsa.WorkOrderProvisioning
		technicalData *fieldsa.WorkOrderTechnicalData
		timelineData  []fieldsa.WorkOrderTimeline
	)

	/*
		==================================================
		CONCURRENT REQUEST
		==================================================

		6 endpoint Fieldsa dijalankan bersamaan.

		Endpoint optional tidak membuat seluruh request
		gagal. Error akan dicatat sebagai warning.
	*/

	var group errgroup.Group

	// 1. Work Order Detail
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderDetail(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get detail gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		detail = result

		return nil
	})

	// 2. Equipment
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderEquipment(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get equipment gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		equipment = result

		return nil
	})

	// 3. Material
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderMaterial(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get material gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		material = result

		return nil
	})

	// 4. Provisioning
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderProvisioning(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get provisioning gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		provisioning = result

		return nil
	})

	// 5. Technical Data
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderTechnicalData(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get technical data gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		technicalData = result

		return nil
	})

	// 6. Timeline
	group.Go(func() error {

		result, err :=
			s.fieldsaClient.GetWorkOrderTimeline(
				ctx,
				workOrderNumber,
			)

		if err != nil {
			logger.Warn(
				"Get timeline gagal untuk %s: %v",
				workOrderNumber,
				err,
			)

			return nil
		}

		timelineData = result

		return nil
	})

	// Tunggu semua goroutine selesai.
	if err := group.Wait(); err != nil {
		return nil, err
	}

	/*
		==================================================
		VALIDATE MAIN DETAIL
		==================================================

		Detail adalah sumber utama.

		Kalau endpoint utama gagal, baru request dianggap
		gagal. Endpoint lainnya bersifat optional.
	*/

	if detail == nil {
		return nil, fmt.Errorf(
			"work order detail tidak tersedia",
		)
	}

	/*
		==================================================
		EQUIPMENT
		==================================================
	*/

	equipmentResponse :=
		make([]Equipment, 0, len(equipment))

	for _, item := range equipment {

		equipmentResponse = append(
			equipmentResponse,
			Equipment{
				Description:       item.Description,
				SerialNumber:      item.SerialNumber,
				MaterialStatus:    item.MaterialStatus,
				MaterialCondition: item.MaterialCondition,
			},
		)
	}

	/*
		==================================================
		MATERIAL
		==================================================
	*/

	materialResponse :=
		make([]Material, 0, len(material))

	for _, item := range material {

		materialResponse = append(
			materialResponse,
			Material{
				Description:       item.Description,
				MaterialStatus:    item.MaterialStatus,
				MaterialCondition: item.MaterialCondition,
			},
		)
	}

	/*
		==================================================
		PROVISIONING
		==================================================
	*/

	var provisioningResponse *Provisioning

	if provisioning != nil {

		provisioningResponse = &Provisioning{
			SerialNumber:    provisioning.SerialNumber,
			EquipmentModel:  provisioning.EquipmentModel,
			NetworkEntityID: provisioning.NetworkEntityID,
			Epon:            provisioning.Epon,
			VLAN:            provisioning.VLAN,
			ReceivedPower:   provisioning.ReceivedPower,
		}
	}

	/*
		==================================================
		TECHNICAL DATA
		==================================================
	*/

	var technicalResponse *TechnicalData

	if technicalData != nil {

		technicalResponse = &TechnicalData{
			SiteID:         technicalData.SiteID,
			SplitterID:     technicalData.SplitterID,
			TotalPort:      technicalData.TotalPort,
			PortNumber:     technicalData.PortNumber,
			AvailablePort:  technicalData.AvailablePort,
			DPFOCoordinate: technicalData.DPFOCoordinate,
		}
	}

	/*
		==================================================
		TIMELINE
		==================================================
	*/

	var timelineResponse *TimelineResponse

	// Prioritas 1:
	// Provisioning In Progress
	for _, item := range timelineData {

		if item.WorkOrderStatusGroup !=
			"Provisioning In Progress" {
			continue
		}

		latitude := strings.ReplaceAll(
			item.ActionLatitude,
			",",
			".",
		)

		longitude := strings.ReplaceAll(
			item.ActionLongtitude,
			",",
			".",
		)

		timelineResponse = &TimelineResponse{
			ActionLatitude:         latitude,
			ActionLongtitude:       longitude,
			LatitudeLongtitudeLink: item.LatitudeLongtitudeLink,
		}

		break
	}

	// Prioritas 2:
	// On Progress
	if timelineResponse == nil {

		for _, item := range timelineData {

			if item.WorkOrderStatusGroup !=
				"On Progress" {
				continue
			}

			latitude := strings.ReplaceAll(
				item.ActionLatitude,
				",",
				".",
			)

			longitude := strings.ReplaceAll(
				item.ActionLongtitude,
				",",
				".",
			)

			timelineResponse = &TimelineResponse{
				ActionLatitude:         latitude,
				ActionLongtitude:       longitude,
				LatitudeLongtitudeLink: item.LatitudeLongtitudeLink,
			}

			break
		}
	}

	/*
		==================================================
		FINAL RESPONSE
		==================================================

		URUTAN DI SINI TETAP.

		Tidak peduli goroutine mana yang selesai lebih dulu.
	*/

	return &WorkOrderDetailResponse{

		// 1. Detail
		WorkOrderNumber:      detail.WorkOrderNumber,
		ReferenceCode:        detail.ReferenceCode,
		CustomerName:         detail.CustomerName,
		CircuitID:            detail.CircuitID,
		CID:                  detail.CID,
		EndCustomerName:      detail.EndCustomerName,
		EndCustomerAddress:   detail.EndCustomerAddress,
		HomepassIDCoordinate: detail.HomepassIDCoordinate,
		DeviceAllocation:     detail.DeviceAllocation,
		WorkOrderStatusItem:  detail.WorkOrderStatusItem,
		VendorName:           detail.VendorName,
		TechnicianName:       detail.TechnicianName,

		// 2. Equipment
		Equipment: equipmentResponse,

		// 3. Material
		Material: materialResponse,

		// 4. Provisioning
		Provisioning: provisioningResponse,

		// 5. Technical Data
		TechnicalData: technicalResponse,

		// 6. Timeline
		Timeline: timelineResponse,
	}, nil
}
