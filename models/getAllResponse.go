package models

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type GetAll struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}
