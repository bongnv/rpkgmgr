package model

import (
	"time"
)

type Package struct {
	ID              int64 `gorm:"primary_key"`
	Name            string
	Version         string
	PublicationDate *time.Time `gorm:"column:publication_date"`
	Title           string
	Description     string
	Authors         string
	Maintainers     string
}
