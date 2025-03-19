package repository

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model
	Key         string `gorm:"index"`
	Name        string
	City        string
	Country     string
	Alias       []string  `gorm:"type:text;serializer:json"`
	Regions     []string  `gorm:"type:text;serializer:json"`
	Coordinates []float64 `gorm:"type:text;serializer:json"`
	Province    string
	Timezone    string
	Unlocs      []string `gorm:"type:text;serializer:json"`
	Code        string
}
