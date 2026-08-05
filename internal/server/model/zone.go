package model

import (
	"strings"
	"time"

	"git.ocn.com.vn/ocn/common/ormbase"
)

type Zone struct {
	ID             uint64        `gorm:"column:id;primary_key"`
	OrganizationID uint64        `gorm:"column:organization_id"`
	Name           string        `gorm:"column:name"`
	IsTemplate     bool          `gorm:"column:is_template"`
	DNSSEC         bool          `gorm:"column:dns_sec"`
	Type           Type          `gorm:"column:type"`
	SOAType        SOAType       `gorm:"column:soa_type"`
	CreatedAt      time.Time     `gorm:"column:created_at"`
	UpdatedAt      time.Time     `gorm:"column:updated_at"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID"`
}

func (Zone) TableName() string {
	return "zones"
}

func (Zone) AvailableSearchFields() []string {
	return []string{"id", "organization_id", "name", "type", "soa_type", "is_template"}
}

type ListZoneParams struct {
	Paging   *ormbase.Paging
	Searches []*ormbase.Search
}
type Type int

const (
	None      Type = 0
	Native    Type = 1
	Primary   Type = 2
	Secondary Type = 3
)

func (s Type) String() string {
	return [...]string{
		"None", "Native", "Primary", "Secondary",
	}[s]
}

type SOAType int

const (
	DEFAULT  SOAType = 0
	EPOCH    SOAType = 1
	INCREASE SOAType = 2
	OFF      SOAType = 3
)

func (s SOAType) String() string {
	return [...]string{
		"DEFAULT", "EPOCH", "INCREASE", "OFF",
	}[s]
}

func ParseType(s string) Type {
	switch strings.ToLower(s) {
	case "none":
		return None
	case "native":
		return Native
	case "primary":
		return Primary
	case "secondary":
		return Secondary
	default:
		return None
	}
}

func ParseSOAType(s string) SOAType {
	switch strings.ToLower(s) {
	case "default":
		return DEFAULT
	case "epoch":
		return EPOCH
	case "increase":
		return INCREASE
	case "off":
		return OFF
	default:
		return DEFAULT
	}
}
