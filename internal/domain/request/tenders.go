package request

type GetTender struct {
	Limit       int32    `query:"limit" validate:"gte=0"`
	Offset      int32    `query:"offset" validate:"gte=0"`
	ServiceType []string `query:"service_type" validate:"dive,oneof=Construction Delivery Manufacture"`
}
type CreateTender struct {
	Name            string `json:"name" validate:"required,max=100"`
	Description     string `json:"description" validate:"required,max=500"`
	ServiceType     string `json:"serviceType" validate:"required,oneof=Construction Delivery Manufacture"`
	OrganizationId  string `json:"organizationId" validate:"required,uuid4"`
	CreatorUsername string `json:"creatorUsername" validate:"required,max=50"`
}
type GetTenderByUser struct {
	Limit    int32  `query:"limit" validate:"gte=0"`
	Offset   int32  `query:"offset" validate:"gte=0"`
	Username string `query:"username" validate:"required,max=50"`
}
type TenderStatus struct {
	TenderId string `param:"tenderId" validate:"required,uuid4"`
	Username string `query:"username" validate:"max=50"`
}
type UpdateTenderStatus struct {
	TenderId string `param:"tenderId" validate:"required,uuid4"`
	Status   string `query:"status" validate:"required,oneof=Created Published Closed"`
	Username string `query:"username" validate:"required,max=50"`
}
type EditTender struct {
	TenderId    string `param:"tenderId" validate:"required"`
	Username    string `query:"username" validate:"required"`
	Name        string `json:"name" validate:"max=100"`
	Description string `json:"description" validate:"max=500"`
	ServiceType string `json:"serviceType" validate:"omitempty,oneof=Construction Delivery Manufacture"`
}
type RollbackTender struct {
	TenderId string `param:"tenderId" validate:"required,uuid4"`
	Version  int32  `param:"version" validate:"required,gt=0"`
	Username string `query:"username" validate:"required,max=50"`
}
