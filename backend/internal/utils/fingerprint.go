package utils

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

// GenerateFingerprint generates a hardware fingerprint based on the processor ID and disk serial number.
func GenerateFingerprint() (string, error) {
	cpuID, err := GetProcessorID()
	if err != nil {
		return "", err
	}

	diskID, err := GetDiskID()
	if err != nil {
		return "", err
	}

	// Combine hardware IDs and create SHA-256 hash
	combined := cpuID + diskID
	hash := sha256.Sum256([]byte(combined))
	return fmt.Sprintf("%x", hash), nil
}

func GetProcessorID() (string, error) {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorId")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		return strings.TrimSpace(lines[1]), nil
	}
	return "", fmt.Errorf("processor ID not found")
}

func GetDiskID() (string, error) {
	cmd := exec.Command("wmic", "diskdrive", "get", "SerialNumber")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		return strings.TrimSpace(lines[1]), nil
	}
	return "", fmt.Errorf("disk serial not found")
}
