package dto

import "time"

type Organization struct {
	Contact     string    `json:"contact"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	Description string    `json:"description"`
	Email       string    `json:"email"`
	Id          uint64    `json:"id"`
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
