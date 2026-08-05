package dto

import "time"

// Defines values for Type.
const (
	None      Type = "None"
	Native    Type = "Native"
	Primary   Type = "Primary"
	Secondary Type = "Secondary"
)

// Defines values for SOAType.
const (
	NONE     SOAType = "NONE"
	DEFAULT  SOAType = "DEFAULT"
	EPOCH    SOAType = "EPOCH"
	INCREASE SOAType = "INCREASE"
	OFF      SOAType = "OFF"
)

type Type string
type SOAType string
type Zone struct {
	Id             uint64    `json:"id"`
	DnsSec         bool      `json:"dnsSec"`
	IsTemplate     bool      `json:"isTemplate,omitempty"`
	Name           string    `json:"name"`
	OrganizationID uint64    `json:"organizationID"`
	SoaType        string    `json:"soaType,omitempty"`
	Type           string    `json:"type,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type UpdateZoneRequest struct {
	DnsSec  *bool   `json:"dnsSec,omitempty"`
	Type    *string `json:"type,omitempty"`
	SoaType *string `json:"soaType,omitempty"`
}

type ListZonesJSONBody struct {
	Paging   PagingRequest   `json:"paging,omitempty"`
	Searches []SearchRequest `json:"searches,omitempty"`
}

// PagingRequest defines model for PagingRequest.
type PagingRequest struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
}

// SearchRequest defines model for SearchRequest.
type SearchRequest struct {
	Field string      `json:"field,omitempty"`
	Op    string      `json:"op,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// ListZoneResponse defines model for ListZoneResponse.
type ListZoneResponse struct {
	Organizations []Zone `json:"organizations,omitempty"`
	Paging        Paging `json:"paging,omitempty"`
}
