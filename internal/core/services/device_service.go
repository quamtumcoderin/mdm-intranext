package services

import (
	"context"
	"errors"
	"strings"

	"mdm-intranext/internal/core/domain"
	"mdm-intranext/internal/core/ports"
)

type deviceService struct {
	repo ports.DeviceRepository
}

func NewDeviceService(repo ports.DeviceRepository) ports.DeviceService {
	return &deviceService{repo: repo}
}

func (s *deviceService) RegisterDevice(ctx context.Context, device *domain.Device) error {
	device.UserID = strings.TrimSpace(device.UserID)
	device.DeviceName = strings.TrimSpace(device.DeviceName)

	if device.UserID == "" {
		return errors.New("userid is required")
	}
	if device.DeviceName == "" {
		return errors.New("devicename is required")
	}

	device.Status = "registered"

	return s.repo.Create(ctx, device)
}

func (s *deviceService) GetDevice(ctx context.Context, id string) (*domain.Device, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(ctx, id)
}
