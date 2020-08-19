package db

import "time"

type SelfTest struct {
	//GORM attributes, see: http://gorm.io/docs/conventions.html
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	DeviceWWN string
	Device    Device `json:"-" gorm:"foreignkey:DeviceWWN"` // use DeviceWWN as foreign key

	Date time.Time
}
