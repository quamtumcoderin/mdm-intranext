package ports

import (
	"context"
	"mdm-intranext/internal/core/domain"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *domain.Device) error
	GetByID(ctx context.Context, id string) (*domain.Device, error)
}
