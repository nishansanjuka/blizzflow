package services

import (
	"blizzflow/backend/domain/model"
	repository "blizzflow/backend/domain/repositories"
	license "blizzflow/backend/internal/utils"
	"fmt"
	"time"
)

// Custom errors
var (
	ErrInvalidLicenseKey = fmt.Errorf("invalid license key format")
	ErrExpiredLicense    = fmt.Errorf("license has expired")
	ErrLicenseNotFound   = fmt.Errorf("license not found")
	ErrDatabaseOperation = fmt.Errorf("database operation failed")
	ErrInvalidUsername   = fmt.Errorf("invalid username provided")
)

type LicenseService struct {
	licenseRepo *repository.LicenseRepository
}

func NewLicenseService(licenseRepo *repository.LicenseRepository) *LicenseService {
	return &LicenseService{licenseRepo: licenseRepo}
}

func (s *LicenseService) GenerateLicense(username string, expiresAt time.Time) (*model.License, error) {
	if username == "" {
		return nil, ErrInvalidUsername
	}

	lic, err := license.GenerateLicenseKey(username, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate license: %w", err)
	}

	licenseModel := &model.License{
		Key:         lic.Key,
		Username:    lic.Username,
		ExpiresAt:   lic.ExpiresAt,
		Fingerprint: lic.Fingerprint,
	}

	if err := s.licenseRepo.Create(licenseModel); err != nil {
		return nil, fmt.Errorf("failed to store license: %w", ErrDatabaseOperation)
	}

	return licenseModel, nil
}

func (s *LicenseService) ValidateLicense(key string) (bool, error) {
	if key == "" {
		return false, ErrInvalidLicenseKey
	}

	valid, err := license.ValidateLicenseKey(key)
	if err != nil {
		return false, ErrInvalidLicenseKey
	}

	return valid, nil
}

func (s *LicenseService) DecodeLicense(key string) (*model.License, error) {
	if key == "" {
		return nil, ErrInvalidLicenseKey
	}

	lic, err := license.DecodeLicense(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode license: %w", err)
	}

	if time.Now().After(lic.ExpiresAt) {
		return nil, ErrExpiredLicense
	}

	return &model.License{
		Key:         lic.Key,
		Username:    lic.Username,
		ExpiresAt:   lic.ExpiresAt,
		Fingerprint: lic.Fingerprint,
	}, nil
}
