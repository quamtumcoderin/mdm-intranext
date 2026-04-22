package postgres

import (
	"context"
	"fmt"

	"mdm-intranext/internal/core/domain"
	"mdm-intranext/internal/core/ports"

	"github.com/jackc/pgx/v5/pgxpool"
)

type deviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) ports.DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *domain.Device) error {
	query := `INSERT INTO devices (user_id, device_name, device_model, os_version, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, device.UserID, device.DeviceName, device.DeviceModel, device.OSVersion, device.Status).Scan(&device.ID, &device.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	return nil
}

func (r *deviceRepository) GetByID(ctx context.Context, id string) (*domain.Device, error) {
	query := `SELECT id, user_id, device_name, device_model, os_version, last_seen, status, created_at FROM devices WHERE id = $1`

	device := &domain.Device{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&device.ID, &device.UserID, &device.DeviceName, &device.DeviceModel,
		&device.OSVersion, &device.LastSeen, &device.Status, &device.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	return device, nil
}
