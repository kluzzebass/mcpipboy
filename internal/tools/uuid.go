package tools

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUIDTool implements comprehensive UUID generation and validation
type UUIDTool struct{}

// NewUUIDTool creates a new UUID tool instance
func NewUUIDTool() *UUIDTool {
	return &UUIDTool{}
}

// Name returns the tool's name
func (u *UUIDTool) Name() string {
	return "uuid"
}

// Description returns the tool's description
func (u *UUIDTool) Description() string {
	return "Generate and validate UUIDs with various versions (v1, v4, v5, v7)"
}

// Execute runs the UUID tool
func (u *UUIDTool) Execute(params map[string]interface{}) (interface{}, error) {
	// Parse parameters
	version, _ := params["version"].(string)
	if version == "" {
		version = "v4"
	}

	count, countProvided := params["count"].(float64)
	if !countProvided {
		count = 1
	}

	// Validate count if provided
	if countProvided && (count < 1 || count > 1000) {
		return nil, fmt.Errorf("count must be between 1 and 1000")
	}

	// Execute based on version
	switch version {
	case "v1":
		return u.generateV1(int(count))
	case "v4":
		return u.generateV4(int(count))
	case "v5":
		return u.generateV5(params, int(count))
	case "v7":
		return u.generateV7(int(count))
	case "validate":
		return u.validateUUID(params)
	default:
		return nil, fmt.Errorf("invalid version: %s, must be one of: v1, v4, v5, v7, validate", version)
	}
}

// generateV1 generates UUID v1 (time-based)
func (u *UUIDTool) generateV1(count int) (interface{}, error) {
	var results []string
	for range count {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v1: %v", err)
		}
		results = append(results, id.String())
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateV4 generates UUID v4 (random)
func (u *UUIDTool) generateV4(count int) (interface{}, error) {
	var results []string
	for range count {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v4: %v", err)
		}
		results = append(results, id.String())
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateV5 generates UUID v5 (name-based SHA-1)
func (u *UUIDTool) generateV5(params map[string]interface{}, count int) (interface{}, error) {
	namespace, _ := params["namespace"].(string)
	name, _ := params["name"].(string)

	if namespace == "" {
		return nil, fmt.Errorf("namespace parameter is required for UUID v5")
	}
	if name == "" {
		return nil, fmt.Errorf("name parameter is required for UUID v5")
	}

	// Parse namespace UUID
	namespaceUUID, err := uuid.Parse(namespace)
	if err != nil {
		return nil, fmt.Errorf("invalid namespace UUID: %v", err)
	}

	var results []string
	for i := range count {
		// For multiple generations, append index to make them unique
		uniqueName := name
		if count > 1 {
			uniqueName = fmt.Sprintf("%s-%d", name, i)
		}

		id := uuid.NewSHA1(namespaceUUID, []byte(uniqueName))
		results = append(results, id.String())
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateV7 generates UUID v7 (time-ordered)
func (u *UUIDTool) generateV7(count int) (interface{}, error) {
	var results []string
	for range count {
		// Generate UUID v7 using time-based approach
		// Note: Go's uuid package doesn't have v7 yet, so we'll implement a basic version
		id, err := u.generateV7Impl()
		if err != nil {
			return nil, fmt.Errorf("failed to generate UUID v7: %v", err)
		}
		results = append(results, id)
	}

	// Return single value if count is 1, otherwise return array
	if count == 1 {
		return results[0], nil
	}
	return results, nil
}

// generateV7Impl implements a basic UUID v7 generation
func (u *UUIDTool) generateV7Impl() (string, error) {
	// UUID v7: 48-bit timestamp + 12-bit random + 4-bit version + 2-bit variant + 62-bit random
	now := time.Now()
	timestamp := now.UnixMilli()

	// Generate random bytes
	randomBytes := make([]byte, 10)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Construct UUID v7
	// First 6 bytes: timestamp (48 bits)
	uuidBytes := make([]byte, 16)
	uuidBytes[0] = byte(timestamp >> 40)
	uuidBytes[1] = byte(timestamp >> 32)
	uuidBytes[2] = byte(timestamp >> 24)
	uuidBytes[3] = byte(timestamp >> 16)
	uuidBytes[4] = byte(timestamp >> 8)
	uuidBytes[5] = byte(timestamp)

	// Next 2 bytes: random (12 bits) + version (4 bits)
	uuidBytes[6] = randomBytes[0] & 0x0F // Clear top 4 bits
	uuidBytes[6] |= 0x70                 // Set version to 7
	uuidBytes[7] = randomBytes[1]

	// Next byte: variant (2 bits) + random (6 bits)
	uuidBytes[8] = randomBytes[2] & 0x3F // Clear top 2 bits
	uuidBytes[8] |= 0x80                 // Set variant bits

	// Remaining 7 bytes: random
	copy(uuidBytes[9:], randomBytes[3:])

	// Convert to UUID string
	id := uuid.UUID(uuidBytes)
	return id.String(), nil
}

// validateUUID validates a UUID string
func (u *UUIDTool) validateUUID(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Parse the UUID
	id, err := uuid.Parse(input)
	if err != nil {
		return map[string]interface{}{
			"valid":   false,
			"error":   err.Error(),
			"version": nil,
		}, nil
	}

	// Get version
	version := int(id.Version())
	result := map[string]interface{}{
		"valid":   true,
		"version": version,
		"uuid":    id.String(),
	}

	// Extract additional information based on version
	switch version {
	case 1:
		// UUID v1: Extract timestamp and MAC address
		timestamp := u.extractV1Timestamp(id)
		result["timestamp"] = timestamp
		result["mac_address"] = u.extractV1MAC(id)
	case 5:
		// UUID v5: Extract namespace (if possible to determine)
		namespace := u.extractV5Namespace()
		if namespace != "" {
			result["namespace"] = namespace
		}
	case 7:
		// UUID v7: Extract timestamp
		timestamp := u.extractV7Timestamp(id)
		result["timestamp"] = timestamp
	}

	return result, nil
}

// extractV1Timestamp extracts the timestamp from a UUID v1
func (u *UUIDTool) extractV1Timestamp(id uuid.UUID) string {
	// UUID v1 timestamp is in the first 8 bytes (60 bits)
	// Convert to Unix timestamp (seconds since epoch)
	bytes := id[:8]

	// The timestamp in UUID v1 is in 100-nanosecond intervals since 1582-10-15 00:00:00 UTC
	// We need to convert this to Unix timestamp
	// UUID v1 epoch: 1582-10-15 00:00:00 UTC
	// Unix epoch: 1970-01-01 00:00:00 UTC
	// Difference: 122192928000000000 nanoseconds

	// Read the timestamp as big-endian
	timestamp := int64(bytes[0])<<40 | int64(bytes[1])<<32 | int64(bytes[2])<<24 |
		int64(bytes[3])<<16 | int64(bytes[4])<<8 | int64(bytes[5])

	// Convert from 100-nanosecond intervals to seconds
	// and adjust for the epoch difference
	unixTimestamp := (timestamp - 122192928000000000) / 10000000

	return time.Unix(unixTimestamp, 0).UTC().Format(time.RFC3339)
}

// extractV1MAC extracts the MAC address from a UUID v1
func (u *UUIDTool) extractV1MAC(id uuid.UUID) string {
	// MAC address is in the last 6 bytes
	bytes := id[10:]
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5])
}

// extractV5Namespace attempts to determine the namespace from a UUID v5
func (u *UUIDTool) extractV5Namespace() string {
	// For UUID v5, we can't easily reverse-engineer the namespace
	// without knowing the original name that was hashed
	// This is a placeholder - in practice, you'd need the original name too
	return "unknown (requires original name to determine)"
}

// extractV7Timestamp extracts the timestamp from a UUID v7
func (u *UUIDTool) extractV7Timestamp(id uuid.UUID) string {
	// UUID v7 timestamp is in the first 6 bytes (48 bits)
	// It's a Unix timestamp in milliseconds
	bytes := id[:6]

	// Read the timestamp as big-endian (48 bits)
	timestamp := int64(bytes[0])<<40 | int64(bytes[1])<<32 | int64(bytes[2])<<24 |
		int64(bytes[3])<<16 | int64(bytes[4])<<8 | int64(bytes[5])

	// Convert from milliseconds to seconds
	unixTimestamp := timestamp / 1000
	nanoSeconds := (timestamp % 1000) * 1000000

	return time.Unix(unixTimestamp, nanoSeconds).UTC().Format(time.RFC3339)
}

// ValidateParams validates the input parameters
func (u *UUIDTool) ValidateParams(params map[string]interface{}) error {
	// Validate version
	if version, ok := params["version"]; ok {
		if versionStr, ok := version.(string); ok {
			validVersions := []string{"v1", "v4", "v5", "v7", "validate"}
			if !contains(validVersions, versionStr) {
				return fmt.Errorf("invalid version: %s, must be one of: %s", versionStr, strings.Join(validVersions, ", "))
			}
		} else {
			return fmt.Errorf("version parameter must be a string")
		}
	}

	// Validate count
	if count, ok := params["count"]; ok {
		if countFloat, ok := count.(float64); ok {
			if countFloat < 1 || countFloat > 1000 {
				return fmt.Errorf("count must be between 1 and 1000")
			}
		} else {
			return fmt.Errorf("count parameter must be a number")
		}
	}

	// Validate namespace for v5
	if version, ok := params["version"]; ok {
		if versionStr, ok := version.(string); ok && versionStr == "v5" {
			if namespace, ok := params["namespace"]; ok {
				if namespaceStr, ok := namespace.(string); ok {
					if _, err := uuid.Parse(namespaceStr); err != nil {
						return fmt.Errorf("invalid namespace UUID: %v", err)
					}
				} else {
					return fmt.Errorf("namespace parameter must be a string")
				}
			}
		}
	}

	// Validate name for v5
	if version, ok := params["version"]; ok {
		if versionStr, ok := version.(string); ok && versionStr == "v5" {
			if name, ok := params["name"]; ok {
				if _, ok := name.(string); !ok {
					return fmt.Errorf("name parameter must be a string")
				}
			}
		}
	}

	// Validate input for validation
	if version, ok := params["version"]; ok {
		if versionStr, ok := version.(string); ok && versionStr == "validate" {
			if input, ok := params["input"]; ok {
				if _, ok := input.(string); !ok {
					return fmt.Errorf("input parameter must be a string")
				}
			}
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
func (u *UUIDTool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "version",
			Type:        "string",
			Description: "UUID version: v1 (time-based), v4 (random), v5 (name-based SHA-1), v7 (time-ordered), validate",
			Required:    false,
		},
		{
			Name:        "count",
			Type:        "number",
			Description: "Number of UUIDs to generate (1-1000)",
			Required:    false,
		},
		{
			Name:        "namespace",
			Type:        "string",
			Description: "Namespace UUID for v5 generation (required for v5)",
			Required:    false,
		},
		{
			Name:        "name",
			Type:        "string",
			Description: "Name for v5 generation (required for v5)",
			Required:    false,
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "UUID string to validate (required for validate)",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
func (u *UUIDTool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"type":        "array",
				"description": "Array of UUIDs (or single UUID if count=1), or validation result for validate",
				"items": map[string]interface{}{
					"oneOf": []map[string]interface{}{
						{"type": "string"},
						{"type": "object"},
					},
				},
			},
		},
	}
}
