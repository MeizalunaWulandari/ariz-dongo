package fieldsa

type WorkOrderDetail struct {
	WorkOrderNumber      string `json:"workOrderNumber"`
	ReferenceCode        string `json:"referenceCode"`
	CustomerName         string `json:"customerName"`
	CircuitID            string `json:"circuitId"`
	CID                  string `json:"cid"`
	EndCustomerName      string `json:"endCustomerName"`
	EndCustomerAddress   string `json:"endCustomerAddress"`
	HomepassIDCoordinate string `json:"homepassIdCoordinate"`
	DeviceAllocation     string `json:"deviceAllocation"`
	WorkOrderStatusItem  string `json:"workOrderStatusItem"`
	VendorName           string `json:"vendorName"`
	TechnicianName       string `json:"technicianName"`
}

type WorkOrderDocument struct {
	WorkOrderNumber  string `json:"workOrderNumber"`
	DocumentTypeID   int    `json:"documentTypeId"`
	DocumentTypeName string `json:"documentTypeName"`
	DocumentFullPath string `json:"documentFullPath"`
	Note             string `json:"note"`
	Created          string `json:"created"`
	CreatedBy        string `json:"createdBy"`
}

type WorkOrderMaterial struct {
	WorkOrderMaterialID  int     `json:"workOrderMaterialId"`
	WorkOrderNumber      string  `json:"workOrderNumber"`
	MaterialNumber       string  `json:"materialNumber"`
	Manufacture          *string `json:"manufacture"`
	EquipmentModel       *string `json:"equipmentModel"`
	Description          string  `json:"description"`
	UOM                  string  `json:"uom"`
	Qty                  float64 `json:"qty"`
	SerialNumber         string  `json:"serialNumber"`
	IsOwnedByFS          bool    `json:"isOwnedByFs"`
	MaterialOwnedBy      string  `json:"materialOwnedBy"`
	MaterialStatusID     int     `json:"materialStatusId"`
	MaterialStatus       string  `json:"materialStatus"`
	IsEquipment          bool    `json:"isEquipment"`
	MaterialConditionID  int     `json:"materialConditionId"`
	MaterialCondition    string  `json:"materialCondition"`
	IsSendToProvisioning bool    `json:"isSendToProvisioning"`
	Created              string  `json:"created"`
	CreatedBy            string  `json:"createdBy"`
	Modified             string  `json:"modified"`
	ModifiedBy           string  `json:"modifiedBy"`
}

type WorkOrderProvisioning struct {
	WorkOrderProvisioningID int    `json:"workOrderProvisioningId"`
	WorkOrderNumber         string `json:"workOrderNumber"`
	ProvisioningTypeID      int    `json:"provisioningTypeId"`
	ProvisioningType        string `json:"provisioningType"`
	SiteID                  string `json:"siteId"`
	SerialNumber            string `json:"serialNumber"`
	EquipmentModel          string `json:"equipmentModel"`
	NetworkEntityID         string `json:"networkEntityId"`
	Epon                    string `json:"epon"`
	VLAN                    string `json:"vlan"`
	Profile                 string `json:"profile"`
	ReceivedPower           string `json:"receivedPower"`
	IsSuccess               *bool  `json:"isSuccess"`
	StatusRequest           string `json:"statusRequest"`
	Note                    string `json:"note"`
	Created                 string `json:"created"`
	CreatedBy               string `json:"createdBy"`
	Modified                string `json:"modified"`
	ModifiedBy              string `json:"modifiedBy"`
}

type WorkOrderTechnicalData struct {
	WorkOrderTechnicalDataID int    `json:"workOrderTechnicalDataId"`
	WorkOrderNumber          string `json:"workOrderNumber"`
	SiteID                   string `json:"siteId"`
	SplitterID               string `json:"splitterId"`
	DPFOID                   string `json:"dpfoId"`
	TotalPort                string `json:"totalPort"`
	PortNumber               string `json:"portNumber"`
	AvailablePort            string `json:"availablePort"`
	DPFOCoordinate           string `json:"dpfoCoordinate"`
	Created                  string `json:"created"`
	CreatedBy                string `json:"createdBy"`
	Modified                 string `json:"modified"`
	ModifiedBy               string `json:"modifiedBy"`
}
type WorkOrderTimeline struct {
	WorkOrderHistoryID     int     `json:"workOrderHistoryId"`
	WorkOrderNumber        string  `json:"workOrderNumber"`
	WorkOrderStatusID      int     `json:"workOrderStatusId"`
	WorkOrderStatusGroup   string  `json:"workOrderStatusGroup"`
	WorkOrderStatusItem    string  `json:"workOrderStatusItem"`
	ActionNote             *string `json:"actionNote"`
	AdditionalNote         string  `json:"additionalNote"`
	AdditionalNoteOriginal string  `json:"additionalNoteOriginal"`
	ActionLatitude         string  `json:"actionLatitude"`
	ActionLongtitude       string  `json:"actionLongtitude"`
	LatitudeLongtitudeLink string  `json:"latitudeLongtitudeLink"`
	VendorCode             string  `json:"vendorCode"`
	Technician             string  `json:"technician"`
	Modified               string  `json:"modified"`
	ModifiedBy             string  `json:"modifiedBy"`
}
