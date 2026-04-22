package ports

import (
	"context"
	"mdm-intranext/internal/core/domain"
)

type DeviceService interface {
	RegisterDevice(ctx context.Context, device *domain.Device) error
	GetDevice(ctx context.Context, id string) (*domain.Device, error)
}
