package domain

import "time"

type Device struct {
	ID          string
	UserID      string
	DeviceName  string
	DeviceModel string
	OSVersion   string
	LastSeen    time.Time
	Status      string
	CreatedAt   time.Time
}
