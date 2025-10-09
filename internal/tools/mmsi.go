package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Country represents a country with its code and associated MIDs
type Country struct {
	Code string
	Name string
	MIDs []int
}

// MMSIType represents a specific MMSI type with generation and validation functions
type MMSIType struct {
	Name     string // kebab-case name (e.g., "ship", "sar-aircraft")
	FullName string // full name (e.g., "Ship Station", "SAR Aircraft")
	IsType   func(int) bool
	Generate func(countryCode string) (int, error)
}

// MMSITool implements MMSI number validation and generation
type MMSITool struct {
	countries []Country
	allMIDs   []int
	types     []MMSIType
}

// NewMMSITool creates a new MMSI tool instance
func NewMMSITool() *MMSITool {
	tool := &MMSITool{}
	tool.populateCountries()
	tool.populateAllMIDs()
	tool.populateTypes()
	return tool
}

// Name returns the tool name
func (m *MMSITool) Name() string {
	return "mmsi"
}

// Description returns the tool description
func (m *MMSITool) Description() string {
	return "Generate and validate Maritime Mobile Service Identity (MMSI) numbers. MMSI numbers are 9-digit identifiers used for maritime communication."
}

// Execute runs the MMSI tool
func (m *MMSITool) Execute(params map[string]interface{}) (interface{}, error) {
	operation, _ := params["operation"].(string)
	if operation == "" {
		operation = "validate" // Default to validate
	}

	switch operation {
	case "validate":
		return m.validateMMSI(params)
	case "generate":
		return m.generateMMSI(params)
	default:
		return nil, fmt.Errorf("invalid operation: %s. Must be 'validate' or 'generate'", operation)
	}
}

// validateMMSI validates an MMSI number
func (m *MMSITool) validateMMSI(params map[string]interface{}) (interface{}, error) {
	input, _ := params["input"].(string)
	if input == "" {
		return nil, fmt.Errorf("input parameter is required for validation")
	}

	// Clean the input (remove spaces, dashes, etc.)
	cleanInput := strings.ReplaceAll(strings.ReplaceAll(input, " ", ""), "-", "")

	// Check if it's exactly 9 digits
	if len(cleanInput) != 9 {
		return map[string]interface{}{
			"valid": false,
			"error": "MMSI number must be exactly 9 digits",
			"input": input,
		}, nil
	}

	// Parse as integer and validate
	mmsiNumber, err := strconv.Atoi(cleanInput)
	if err != nil {
		return map[string]interface{}{
			"valid": false,
			"error": "MMSI number must contain only digits",
			"input": input,
		}, nil
	}

	// Basic validation - MMSI should be between 100000000 and 999999999
	if mmsiNumber < 100000000 || mmsiNumber > 999999999 {
		return map[string]interface{}{
			"valid": false,
			"error": "MMSI number must be between 100000000 and 999999999",
			"input": input,
		}, nil
	}

	// Extract MID (first 3 digits)
	mid := mmsiNumber / 1000000

	// Get country name if available
	countryName := m.getCountryName(mid)

	// Determine MMSI type based on format
	mmsiType := m.determineMMSIType(mmsiNumber)

	return map[string]interface{}{
		"valid":        true,
		"mmsi":         fmt.Sprintf("%09d", mmsiNumber),
		"input":        input,
		"mid":          mid,
		"country_name": countryName,
		"type":         mmsiType,
	}, nil
}

// generateMMSI generates MMSI numbers
func (m *MMSITool) generateMMSI(params map[string]interface{}) (interface{}, error) {
	count, _ := params["count"].(int)
	if count <= 0 {
		count = 1
	}

	// Validate count
	if count > 100 {
		return nil, fmt.Errorf("count cannot exceed 100")
	}

	mmsiType, _ := params["type"].(string)
	countryCode, _ := params["country-code"].(string)

	results := make([]int, count)

	for idx := range count {
		var mmsi int
		var err error

		if mmsiType != "" {
			// Generate specific MMSI type
			mmsi, err = m.generateSpecificType(mmsiType, countryCode)
		} else {
			// Generate random MMSI
			mmsi, err = m.generateRandomMMSI(countryCode)
		}

		if err != nil {
			return nil, err
		}

		results[idx] = mmsi
	}

	if count == 1 {
		return fmt.Sprintf("%09d", results[0]), nil
	}

	// Convert to strings for output
	stringResults := make([]string, count)
	for i, mmsi := range results {
		stringResults[i] = fmt.Sprintf("%09d", mmsi)
	}

	return stringResults, nil
}

// populateCountries populates the countries and allMIDs fields
func (m *MMSITool) populateCountries() {
	m.countries = []Country{
		{"US", "United States", []int{366, 367, 368, 369}},
		{"GB", "United Kingdom", []int{232, 233, 234, 235}},
		{"DE", "Germany", []int{211, 218}},
		{"FR", "France", []int{226, 227, 228}},
		{"IT", "Italy", []int{247}},
		{"ES", "Spain", []int{224, 225}},
		{"NL", "Netherlands", []int{244, 245, 246}},
		{"NO", "Norway", []int{257, 258, 259}},
		{"SE", "Sweden", []int{265, 266}},
		{"DK", "Denmark", []int{219, 220}},
		{"FI", "Finland", []int{230, 231}},
		{"PL", "Poland", []int{261, 262}},
		{"RU", "Russia", []int{273, 274, 275, 276}},
		{"JP", "Japan", []int{431, 432}},
		{"CN", "China", []int{412, 413, 414}},
		{"KR", "South Korea", []int{440, 441}},
		{"IN", "India", []int{419, 420}},
		{"AU", "Australia", []int{503, 504}},
		{"NZ", "New Zealand", []int{512}},
		{"CA", "Canada", []int{316, 317}},
		{"BR", "Brazil", []int{710, 711}},
		{"AR", "Argentina", []int{701, 702}},
		{"MX", "Mexico", []int{345, 346}},
		{"ZA", "South Africa", []int{601, 602}},
		{"EG", "Egypt", []int{622, 623}},
		{"NG", "Nigeria", []int{636, 637}},
		{"KE", "Kenya", []int{634, 635}},
		{"MA", "Morocco", []int{242, 243}},
		{"TN", "Tunisia", []int{672, 673}},
		{"DZ", "Algeria", []int{605, 606}},
		{"LY", "Libya", []int{642, 643}},
		{"SD", "Sudan", []int{626, 627}},
		{"ET", "Ethiopia", []int{624, 625}},
		{"GH", "Ghana", []int{620, 621}},
		{"CI", "Ivory Coast", []int{618, 619}},
		{"SN", "Senegal", []int{660, 661}},
		{"ML", "Mali", []int{649, 650}},
		{"BF", "Burkina Faso", []int{633, 634}},
		{"NE", "Niger", []int{656, 657}},
		{"TD", "Chad", []int{670, 671}},
		{"CF", "Central African Republic", []int{612, 613}},
		{"CM", "Cameroon", []int{613, 614}},
		{"GA", "Gabon", []int{626, 627}},
		{"CG", "Congo", []int{676, 677}},
		{"CD", "Democratic Republic of Congo", []int{676, 677}},
		{"AO", "Angola", []int{603, 604}},
		{"ZM", "Zambia", []int{678, 679}},
		{"ZW", "Zimbabwe", []int{679, 680}},
		{"BW", "Botswana", []int{679, 680}},
		{"NA", "Namibia", []int{659, 660}},
		{"SZ", "Swaziland", []int{601, 602}},
		{"LS", "Lesotho", []int{601, 602}},
		{"MG", "Madagascar", []int{647, 648}},
		{"MU", "Mauritius", []int{645, 646}},
		{"SC", "Seychelles", []int{664, 665}},
		{"KM", "Comoros", []int{616, 617}},
		{"DJ", "Djibouti", []int{604, 605}},
		{"SO", "Somalia", []int{666, 667}},
		{"ER", "Eritrea", []int{625, 626}},
		{"SS", "South Sudan", []int{626, 627}},
		{"UG", "Uganda", []int{675, 676}},
		{"RW", "Rwanda", []int{661, 662}},
		{"BI", "Burundi", []int{609, 610}},
		{"TZ", "Tanzania", []int{677, 678}},
		{"MW", "Malawi", []int{655, 656}},
		{"MZ", "Mozambique", []int{650, 651}},
		{"ZM", "Zambia", []int{678, 679}},
		{"ZW", "Zimbabwe", []int{679, 680}},
		{"BW", "Botswana", []int{679, 680}},
		{"NA", "Namibia", []int{659, 660}},
		{"SZ", "Swaziland", []int{601, 602}},
		{"LS", "Lesotho", []int{601, 602}},
		{"MG", "Madagascar", []int{647, 648}},
		{"MU", "Mauritius", []int{645, 646}},
		{"SC", "Seychelles", []int{664, 665}},
		{"KM", "Comoros", []int{616, 617}},
		{"DJ", "Djibouti", []int{604, 605}},
		{"SO", "Somalia", []int{666, 667}},
		{"ER", "Eritrea", []int{625, 626}},
		{"SS", "South Sudan", []int{626, 627}},
		{"UG", "Uganda", []int{675, 676}},
		{"RW", "Rwanda", []int{661, 662}},
		{"BI", "Burundi", []int{609, 610}},
		{"TZ", "Tanzania", []int{677, 678}},
		{"MW", "Malawi", []int{655, 656}},
		{"MZ", "Mozambique", []int{650, 651}},
	}
}

// populateAllMIDs populates the allMIDs field from countries
func (m *MMSITool) populateAllMIDs() {
	for _, country := range m.countries {
		m.allMIDs = append(m.allMIDs, country.MIDs...)
	}
}

// getMIDs returns all Maritime Identification Digits for a country
func (m *MMSITool) getMIDs(countryCode string) []int {
	countryCode = strings.ToUpper(countryCode)

	for _, country := range m.countries {
		if country.Code == countryCode {
			return country.MIDs
		}
	}
	return nil
}

// getCountryName returns the country name for a MID
func (m *MMSITool) getCountryName(mid int) string {
	for _, country := range m.countries {
		for _, countryMID := range country.MIDs {
			if countryMID == mid {
				return country.Name
			}
		}
	}
	return fmt.Sprintf("Unknown (MID: %d)", mid)
}

// generateRandomMMSI generates a random MMSI
func (m *MMSITool) generateRandomMMSI(countryCode string) (int, error) {
	var allMIDs []int
	if countryCode != "" {
		// Get MIDs for specific country
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return 0, fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		// Use precomputed all MIDs
		allMIDs = m.allMIDs
	}

	// Select a random MID from all available ones
	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	// Generate 6 random digits for the remaining part
	remainingDigits := rand.Intn(1000000) // 0 to 999999

	// Build the MMSI number
	mmsi := selectedMID*1000000 + remainingDigits

	return mmsi, nil
}

// populateTypes initializes the MMSI types slice
func (m *MMSITool) populateTypes() {
	m.types = []MMSIType{
		// Most specific types first
		{
			Name:     "us-coast-guard-ship",
			FullName: "US Coast Guard Group Ship Station",
			IsType: func(mmsi int) bool {
				return mmsi == 36699999
			},
			Generate: func(countryCode string) (int, error) {
				return 36699999, nil
			},
		},
		{
			Name:     "us-coast-guard-coast",
			FullName: "US Coast Guard Group Coast Station",
			IsType: func(mmsi int) bool {
				return mmsi == 3669999
			},
			Generate: func(countryCode string) (int, error) {
				return 3669999, nil
			},
		},
		{
			Name:     "us-federal",
			FullName: "US Federal MMSI",
			IsType: func(mmsi int) bool {
				// US Federal: starts with 3669
				return mmsi >= 366900000 && mmsi <= 366999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(100000) // 0 to 99999
				return 366900000 + remaining, nil
			},
		},
		{
			Name:     "us-ship-international",
			FullName: "US Ship Station (International/Inmarsat)",
			IsType: func(mmsi int) bool {
				// US Ship with Inmarsat: starts with 366 and ends with 000
				return mmsi >= 366000000 && mmsi <= 366999999 && mmsi%1000 == 0
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(1000) * 1000 // 000, 1000, 2000, ..., 999000
				return 366000000 + remaining, nil
			},
		},
		{
			Name:     "us-ship-other",
			FullName: "US Ship Station (Other)",
			IsType: func(mmsi int) bool {
				// US Ship other: starts with 366, not ending with 000, not federal
				return mmsi >= 366000000 && mmsi <= 366999999 && mmsi%1000 != 0 && mmsi < 366900000
			},
			Generate: func(countryCode string) (int, error) {
				// Generate random US ship (not Inmarsat, not federal)
				remaining := rand.Intn(900000) + 100000 // 100000 to 999999
				return 366000000 + remaining, nil
			},
		},
		{
			Name:     "us-ship-regular",
			FullName: "US Ship Station (Regular)",
			IsType: func(mmsi int) bool {
				// US Ship regular: starts with 367, 368, 369
				return (mmsi >= 367000000 && mmsi <= 367999999) ||
					(mmsi >= 368000000 && mmsi <= 368999999) ||
					(mmsi >= 369000000 && mmsi <= 369999999)
			},
			Generate: func(countryCode string) (int, error) {
				// Generate random US ship with other MIDs
				mids := []int{367, 368, 369}
				selectedMID := mids[rand.Intn(len(mids))]
				remaining := rand.Intn(1000000)
				return selectedMID*1000000 + remaining, nil
			},
		},
		{
			Name:     "sar-aircraft",
			FullName: "SAR Aircraft",
			IsType: func(mmsi int) bool {
				// SAR aircraft: 111xxxxxx
				return mmsi >= 111000000 && mmsi <= 111999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(1000000)
				return 111000000 + remaining, nil
			},
		},
		{
			Name:     "ais-sart",
			FullName: "AIS-SART",
			IsType: func(mmsi int) bool {
				// AIS-SART: 970xxxxxx
				return mmsi >= 970000000 && mmsi <= 970999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(1000000)
				return 970000000 + remaining, nil
			},
		},
		{
			Name:     "handheld-vhf",
			FullName: "Handheld VHF",
			IsType: func(mmsi int) bool {
				// Handheld VHF: 8xxxxxxx
				return mmsi >= 800000000 && mmsi <= 899999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(100000000)
				return 800000000 + remaining, nil
			},
		},
		{
			Name:     "man-overboard",
			FullName: "Man Overboard Device",
			IsType: func(mmsi int) bool {
				// Man overboard: 972xxxxxx
				return mmsi >= 972000000 && mmsi <= 972999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(1000000)
				return 972000000 + remaining, nil
			},
		},
		{
			Name:     "epirb-ais",
			FullName: "EPIRB-AIS",
			IsType: func(mmsi int) bool {
				// EPIRB-AIS: 974xxxxxx
				return mmsi >= 974000000 && mmsi <= 974999999
			},
			Generate: func(countryCode string) (int, error) {
				remaining := rand.Intn(1000000)
				return 974000000 + remaining, nil
			},
		},
		{
			Name:     "ship",
			FullName: "Ship Station",
			IsType: func(mmsi int) bool {
				// Regular ship: 100000000 to 799999999, not special ranges
				return mmsi >= 100000000 && mmsi <= 799999999 &&
					!(mmsi >= 111000000 && mmsi <= 111999999) && // Not SAR aircraft
					!(mmsi >= 800000000 && mmsi <= 899999999) && // Not handheld VHF
					!(mmsi >= 970000000 && mmsi <= 970999999) && // Not AIS-SART
					!(mmsi >= 972000000 && mmsi <= 972999999) && // Not man overboard
					!(mmsi >= 974000000 && mmsi <= 974999999) // Not EPIRB-AIS
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateShipStation(countryCode, false, false)
			},
		},
		{
			Name:     "group-ship",
			FullName: "Group Ship Station",
			IsType: func(mmsi int) bool {
				// Group ship: 0xxxxxxx (first digit is 0)
				return mmsi >= 10000000 && mmsi <= 99999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateGroupShipStation(countryCode)
			},
		},
		{
			Name:     "coast-station",
			FullName: "Coast Station",
			IsType: func(mmsi int) bool {
				// Coast station: 00xxxxxxx (first two digits are 00)
				return mmsi >= 1000000 && mmsi <= 9999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateCoastStation(countryCode, false)
			},
		},
		{
			Name:     "group-coast-station",
			FullName: "Group Coast Station",
			IsType: func(mmsi int) bool {
				// Group coast station: 000xxxxxx (first three digits are 000)
				return mmsi >= 100000 && mmsi <= 999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateCoastStation(countryCode, true)
			},
		},
		{
			Name:     "craft-associated",
			FullName: "Craft Associated with Parent Ship",
			IsType: func(mmsi int) bool {
				// Craft associated: 98xxxxxxx
				return mmsi >= 980000000 && mmsi <= 989999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateCraftAssociated()
			},
		},
		{
			Name:     "navigational-aid",
			FullName: "Navigational Aid",
			IsType: func(mmsi int) bool {
				// Navigational aid: 99xxxxxxx
				return mmsi >= 990000000 && mmsi <= 999999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateNavigationalAid()
			},
		},
		{
			Name:     "free-form",
			FullName: "Free-form Device",
			IsType: func(mmsi int) bool {
				// Free-form: any other valid MMSI
				return mmsi >= 100000000 && mmsi <= 999999999
			},
			Generate: func(countryCode string) (int, error) {
				return m.generateRandomMMSI(countryCode)
			},
		},
	}
}

// generateSpecificType generates a specific type of MMSI using the types slice
func (m *MMSITool) generateSpecificType(mmsiType, countryCode string) (int, error) {
	for _, t := range m.types {
		if t.Name == mmsiType {
			return t.Generate(countryCode)
		}
	}
	return 0, fmt.Errorf("unsupported MMSI type: %s", mmsiType)
}

// getSupportedTypes returns a slice of supported MMSI type names
func (m *MMSITool) getSupportedTypes() []string {
	types := make([]string, len(m.types))
	for i, t := range m.types {
		types[i] = t.Name
	}
	return types
}

// generateShipStation generates a ship station MMSI
func (m *MMSITool) generateShipStation(countryCode string, inmarsatBCM, inmarsatC bool) (int, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return 0, fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	// Select a random MID from all available ones
	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	// Generate remaining digits
	var remaining int
	if inmarsatBCM || inmarsatC {
		// Inmarsat B/C/M: ends with 000
		remaining = rand.Intn(1000) * 1000 // 000, 1000, 2000, ..., 999000
	} else {
		// Regular ship
		remaining = rand.Intn(1000000) // 0 to 999999
	}

	return selectedMID*1000000 + remaining, nil
}

// generateGroupShipStation generates a group ship station MMSI
func (m *MMSITool) generateGroupShipStation(countryCode string) (int, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return 0, fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	// Select a random MID from all available ones
	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	// Generate remaining digits (6 digits, but first digit must be 0 for group)
	remaining := rand.Intn(100000) // 0 to 99999

	return selectedMID*1000000 + remaining, nil
}

// generateCoastStation generates a coast station MMSI
func (m *MMSITool) generateCoastStation(countryCode string, group bool) (int, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return 0, fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	// Select a random MID from all available ones
	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	// Generate remaining digits
	var remaining int
	if group {
		// Group coast station: 000xxxxxx
		remaining = rand.Intn(100000) // 0 to 99999
	} else {
		// Regular coast station: 00xxxxxxx
		remaining = rand.Intn(1000000) // 0 to 999999
	}

	return selectedMID*1000000 + remaining, nil
}

// generateCraftAssociated generates a craft associated MMSI
func (m *MMSITool) generateCraftAssociated() (int, error) {
	// Craft associated: 98xxxxxxx (fixed format, doesn't depend on country)
	remaining := rand.Intn(10000000) // 0 to 9999999

	return 980000000 + remaining, nil
}

// generateNavigationalAid generates a navigational aid MMSI
func (m *MMSITool) generateNavigationalAid() (int, error) {
	// Navigational aid: 99xxxxxxx (fixed format, doesn't depend on country)
	remaining := rand.Intn(10000000) // 0 to 9999999

	return 990000000 + remaining, nil
}

// determineMMSIType determines the type of MMSI based on its format using the types slice
func (m *MMSITool) determineMMSIType(mmsi int) string {
	// Check each type in order of specificity (most specific first)
	for _, t := range m.types {
		if t.IsType(mmsi) {
			return t.FullName
		}
	}

	return "Unknown"
}

// ValidateParams validates the input parameters
func (m *MMSITool) ValidateParams(params map[string]interface{}) error {
	// Validate operation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok {
			if operationStr != "validate" && operationStr != "generate" {
				return fmt.Errorf("operation must be 'validate' or 'generate'")
			}
		} else {
			return fmt.Errorf("operation must be a string")
		}
	}

	// Validate count for generation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok && operationStr == "generate" {
			if count, ok := params["count"]; ok {
				if countInt, ok := count.(int); ok {
					if countInt < 1 {
						return fmt.Errorf("count must be at least 1")
					}
					if countInt > 100 {
						return fmt.Errorf("count cannot exceed 100")
					}
				} else {
					return fmt.Errorf("count must be an integer")
				}
			}
		}
	}

	// Validate input for validation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok && operationStr == "validate" {
			if input, ok := params["input"]; ok {
				if _, ok := input.(string); !ok {
					return fmt.Errorf("input must be a string")
				}
			} else {
				return fmt.Errorf("input parameter is required for validation")
			}
		}
	}

	// Validate type parameter
	if mmsiType, ok := params["type"]; ok {
		if typeStr, ok := mmsiType.(string); ok {
			supportedTypes := m.getSupportedTypes()
			found := false
			for _, t := range supportedTypes {
				if t == typeStr {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("unsupported MMSI type: %s. Supported types: %v", typeStr, supportedTypes)
			}
		} else {
			return fmt.Errorf("type must be a string")
		}
	}

	// Validate country code
	if countryCode, ok := params["country-code"]; ok {
		if countryStr, ok := countryCode.(string); ok {
			if countryStr != "" {
				mids := m.getMIDs(countryStr)
				if len(mids) == 0 {
					return fmt.Errorf("invalid country code: %s", countryStr)
				}
			}
		} else {
			return fmt.Errorf("country-code must be a string")
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for the tool's input parameters
func (m *MMSITool) GetInputSchema() map[string]interface{} {
	return CreateJSONSchema([]ParameterDefinition{
		{
			Name:        "operation",
			Type:        "string",
			Description: "Operation to perform: 'validate' or 'generate'",
			Required:    false,
			Enum:        []string{"validate", "generate"},
		},
		{
			Name:        "input",
			Type:        "string",
			Description: "MMSI number to validate (required for validation operation)",
			Required:    false,
		},
		{
			Name:        "count",
			Type:        "integer",
			Description: "Number of MMSI numbers to generate (1-100, default: 1)",
			Required:    false,
		},
		{
			Name:        "type",
			Type:        "string",
			Description: "Specific MMSI type to generate",
			Required:    false,
			Enum:        m.getSupportedTypes(),
		},
		{
			Name:        "country-code",
			Type:        "string",
			Description: "Country code for MMSI generation (e.g., 'US', 'GB', 'DE')",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for the tool's output
func (m *MMSITool) GetOutputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"result": map[string]interface{}{
				"description": "Generated MMSI number(s) or validation result",
			},
		},
	}
}
