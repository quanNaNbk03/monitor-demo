package model

import (
	"time"

	"git.ocn.com.vn/ocn/common/ormbase"
)

type Organization struct {
	ID          uint64    `gorm:"column:id;primary_key"`
	Name        string    `gorm:"column:name"`
	Email       string    `gorm:"column:email"`
	ContactName string    `gorm:"column:contact_name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	Zones       []*Zone   `gorm:"foreignKey:OrganizationID"`
}

func (Organization) AvailableSearchFields() []string {
	return []string{"id", "name", "email", "contact_name"}
}

func (Organization) TableName() string {
	return "organizations"
}

type ListOrgParams struct {
	Paging   *ormbase.Paging
	Searches []*ormbase.Search
}
