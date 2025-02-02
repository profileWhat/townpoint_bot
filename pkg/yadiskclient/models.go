package yadiskclient

type OperationResponse struct {
	OperationID string `json:"operation_id"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Templated   bool   `json:"false"`
}

type PublishInfoResponse struct {
	PublicUrl string `json:"public_url"`
}

type OperationStatus struct {
	Status string `json:"status"`
}

func (s OperationStatus) IsSuccess() bool {
	return s.Status == "success"
}

func (s OperationStatus) IsFailed() bool {
	return s.Status == "failed"
}
