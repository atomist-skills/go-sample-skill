package handler

type Skill struct {
	Id        string `json:"id"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
}

type Status struct {
	Code       int8   `json:"code"`
	Reason     string `json:"reason"`
	Visibility string `json:"visibility,omitempty"`
}

type Team struct {
	Id string `json:"id"`
}

type StatusHandlerResponse struct {
	ApiVersion    string `json:"api_version"`
	CorrelationId string `json:"correlation_id"`
	Team          Team   `json:"team"`
	Status        Status `json:"status"`
	Skill         Skill  `json:"skill"`
}
