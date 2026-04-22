package http

import (
	"mdm-intranext/internal/core/domain"
	"time"
)

type RegisterDeviceRequest struct {
	UserID      string `json:"user_id"`
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
	OSVersion   string `json:"os_version"`
}

type DeviceResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DeviceName  string    `json:"device_name"`
	DeviceModel string    `json:"device_model"`
	OSVersion   string    `json:"os_version"`
	LastSeen    time.Time `json:"last_seen"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *RegisterDeviceRequest) mapToDomain() *domain.Device {
	return &domain.Device{
		UserID:      r.UserID,
		DeviceName:  r.DeviceName,
		DeviceModel: r.DeviceModel,
		OSVersion:   r.OSVersion,
	}
}

func mapToResponse(d *domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:          d.ID,
		UserID:      d.UserID,
		DeviceName:  d.DeviceName,
		DeviceModel: d.DeviceModel,
		OSVersion:   d.OSVersion,
		LastSeen:    d.LastSeen,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt,
	}
}
