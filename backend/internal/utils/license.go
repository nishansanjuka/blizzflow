package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const (
	segmentLength = 8
	saltKey       = 0x5A // Simple XOR salt
	segmentCount  = 4
	separator     = "-"
)

type License struct {
	Key         string
	Username    string
	ExpiresAt   time.Time
	Fingerprint string
}

func padSegment(input string, length int) string {
	if len(input) > length {
		return input[:length]
	}
	for len(input) < length {
		input = input + "0"
	}
	return input
}

func encodeUsername(username string) string {
	usernameBytes := []byte(username)
	if len(usernameBytes) > 4 {
		usernameBytes = usernameBytes[:4]
	}
	for len(usernameBytes) < 4 {
		usernameBytes = append(usernameBytes, byte('_'))
	}

	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)

	result := make([]byte, 4)
	for i := 0; i < 4; i++ {
		result[i] = (usernameBytes[i] ^ saltKey) & randomBytes[i]
	}

	encoded := hex.EncodeToString(result)
	// Ensure even length and pad to segmentLength
	if len(encoded)%2 != 0 {
		encoded = "0" + encoded
	}
	return padSegment(encoded, segmentLength)
}

func encodeFingerprint(fingerprint string) string {
	fpBytes := []byte(fingerprint)
	if len(fpBytes) > 4 {
		fpBytes = fpBytes[:4]
	}
	result := make([]byte, 4)
	for i := range fpBytes {
		result[i] = fpBytes[i] ^ saltKey
	}
	return padSegment(hex.EncodeToString(result), segmentLength)
}

func encodeExpiration(expiresAt time.Time) string {
	timeHex := fmt.Sprintf("%x", expiresAt.Unix())
	return padSegment(timeHex, segmentLength)
}

func generateTail() string {
	tail := make([]byte, 4)
	rand.Read(tail)
	return padSegment(hex.EncodeToString(tail), segmentLength)
}

func GenerateLicenseKey(username string, expiresAt time.Time) (*License, error) {

	fingerprint, _ := GenerateFingerprint()
	segment1 := encodeUsername(username)
	segment2 := encodeFingerprint(fingerprint)
	segment3 := encodeExpiration(expiresAt)
	segment4 := generateTail()

	licenseKey := fmt.Sprintf("%s-%s-%s-%s",
		segment1, segment2, segment3, segment4)

	return &License{
		Key:         licenseKey,
		Username:    username,
		ExpiresAt:   expiresAt,
		Fingerprint: fingerprint,
	}, nil
}

func DecodeLicense(licenseKey string) (*License, error) {
	segments := strings.Split(licenseKey, separator)

	if len(segments) != segmentCount {
		return nil, fmt.Errorf("invalid segment count: got %d, want %d", len(segments), segmentCount)
	}

	// Validate segment lengths
	for i, segment := range segments {
		if len(segment) != segmentLength {
			return nil, fmt.Errorf("segment %d invalid length: got %d, want %d", i, len(segment), segmentLength)
		}
	}

	// Clean and validate username segment
	usernameHex := strings.TrimRight(segments[0], "0")
	if len(usernameHex)%2 != 0 {
		usernameHex = usernameHex + "0"
	}

	// Decode username
	usernameBytes, err := hex.DecodeString(usernameHex)
	if err != nil {
		return nil, fmt.Errorf("invalid username encoding: %v", err)
	}

	// Clean username
	username := strings.TrimRight(string(usernameBytes), "_")

	// Decode fingerprint from segment 1
	fingerprint := segments[1]

	// Decode expiration from segment 2
	var expTimestamp int64
	_, err = fmt.Sscanf(segments[2], "%x", &expTimestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration format: %v", err)
	}
	expiresAt := time.Unix(expTimestamp, 0)

	return &License{
		Key:         licenseKey,
		Username:    username,
		ExpiresAt:   expiresAt,
		Fingerprint: fingerprint,
	}, nil
}

func ValidateLicenseKey(licenseKey string) (bool, error) {
	segments := strings.Split(licenseKey, separator)

	if len(segments) != segmentCount {
		return false, fmt.Errorf("invalid segment count: got %d, want %d", len(segments), segmentCount)
	}

	// Check each segment is exactly 8 characters
	for i, segment := range segments {
		if len(segment) != 8 {
			return false, fmt.Errorf("segment %d invalid length: got %d, want 8", i, len(segment))
		}
	}

	// Decode expiration time from segment 3
	expHex := segments[2]
	var expTimestamp int64
	_, err := fmt.Sscanf(expHex, "%x", &expTimestamp)
	if err != nil {
		return false, fmt.Errorf("invalid expiration format")
	}

	// Check if license is expired
	expiresAt := time.Unix(expTimestamp, 0)
	if time.Now().After(expiresAt) {
		return false, fmt.Errorf("license expired")
	}

	fingerprint, _ := GenerateFingerprint()
	// Verify fingerprint
	encodedFingerprint := encodeFingerprint(fingerprint)
	if segments[1] != encodedFingerprint {
		return false, fmt.Errorf("invalid fingerprint")
	}

	return true, nil
}
