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
	IsType   func(string) bool
	Generate func(countryCode string) (string, error)
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

	// Extract MID (first 3 digits)
	mid := mmsiNumber / 1000000

	// Basic validation - MMSI should be between 100000000 and 999999999
	if mmsiNumber < 100000000 || mmsiNumber > 999999999 {
		return map[string]interface{}{
			"valid": false,
			"error": "MMSI number must be between 100000000 and 999999999",
			"input": input,
		}, nil
	}

	// Get country name if available
	countryName := m.getCountryName(mid)

	// Determine MMSI type based on format
	mmsiType := m.determineMMSIType(cleanInput)

	return map[string]interface{}{
		"valid":        true,
		"mmsi":         cleanInput,
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

	results := make([]string, count)

	for idx := range count {
		var mmsi string
		var err error

		if mmsiType != "" {
			// Generate specific MMSI type
			mmsi, err = m.generateSpecificType(mmsiType, countryCode)
		} else {
			// Generate random MMSI (existing logic)
			mmsi, err = m.generateRandomMMSI(countryCode)
		}

		if err != nil {
			return nil, err
		}

		results[idx] = mmsi
	}

	if count == 1 {
		return results[0], nil
	}

	return results, nil
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
		{"FI", "Finland", []int{230}},
		{"JP", "Japan", []int{431, 432}},
		{"CN", "China", []int{412, 413, 414}},
		{"KR", "South Korea", []int{440, 441}},
		{"AU", "Australia", []int{503}},
		{"CA", "Canada", []int{316}},
		{"BR", "Brazil", []int{710}},
		{"RU", "Russia", []int{273}},
		{"IN", "India", []int{419}},
		{"SG", "Singapore", []int{563, 564, 565, 566}},
		{"AR", "Argentina", []int{701}},
		{"CL", "Chile", []int{725}},
		{"EC", "Ecuador", []int{735}},
		{"PE", "Peru", []int{760}},
		{"UY", "Uruguay", []int{770}},
		{"ZA", "South Africa", []int{601}},
		{"EG", "Egypt", []int{622, 623}},
		{"MA", "Morocco", []int{242}},
		{"TN", "Tunisia", []int{672, 673, 674, 675, 676, 677, 678, 679, 680, 681, 682, 683, 684, 685, 686, 687, 688, 689, 690, 691, 692, 693, 694, 695, 696, 697, 698, 699}},
		{"LY", "Libya", []int{642, 643, 644, 645, 646, 647, 648, 649, 650, 651, 652, 653, 654, 655, 656, 657, 658, 659, 660, 661, 662, 663, 664, 665, 666, 667, 668, 669, 670, 671}},
		{"TR", "Turkey", []int{271}},
		{"GR", "Greece", []int{237, 239, 240, 241}},
		{"PL", "Poland", []int{261}},
		{"RO", "Romania", []int{264}},
		{"BG", "Bulgaria", []int{207}},
		{"HR", "Croatia", []int{238}},
		{"SI", "Slovenia", []int{278}},
		{"HU", "Hungary", []int{243}},
		{"SK", "Slovak Republic", []int{267}},
		{"CZ", "Czech Republic", []int{270}},
		{"AT", "Austria", []int{203}},
		{"CH", "Switzerland", []int{269}},
		{"LI", "Liechtenstein", []int{252}},
		{"MC", "Monaco", []int{254}},
		{"SM", "San Marino", []int{268}},
		{"VA", "Vatican City State", []int{208}},
		{"AD", "Andorra", []int{202}},
		{"MT", "Malta", []int{215, 229, 248, 249, 256}},
		{"CY", "Cyprus", []int{209, 210, 212}},
		{"IS", "Iceland", []int{251}},
		{"IE", "Ireland", []int{250}},
		{"LU", "Luxembourg", []int{253}},
		{"BE", "Belgium", []int{205}},
		{"PT", "Portugal", []int{263}},
		{"EE", "Estonia", []int{276}},
		{"LV", "Latvia", []int{275}},
		{"LT", "Lithuania", []int{277}},
		{"BY", "Belarus", []int{206}},
		{"UA", "Ukraine", []int{272}},
		{"MD", "Moldova", []int{214}},
		{"GE", "Georgia", []int{213}},
		{"AM", "Armenia", []int{216}},
		{"AZ", "Azerbaijan", []int{423}},
		{"KZ", "Kazakhstan", []int{436}},
		{"UZ", "Uzbekistan", []int{437}},
		{"TM", "Turkmenistan", []int{434}},
		{"KG", "Kyrgyz Republic", []int{451}},
		{"TJ", "Tajikistan", []int{472}},
		{"AF", "Afghanistan", []int{401}},
		{"PK", "Pakistan", []int{463}},
		{"BD", "Bangladesh", []int{405}},
		{"LK", "Sri Lanka", []int{417}},
		{"MV", "Maldives", []int{455}},
		{"BT", "Bhutan", []int{410}},
		{"NP", "Nepal", []int{459}},
		{"MM", "Myanmar", []int{506}},
		{"TH", "Thailand", []int{567}},
		{"LA", "Lao People's Democratic Republic", []int{531}},
		{"KH", "Cambodia", []int{514, 515}},
		{"VN", "Viet Nam", []int{574}},
		{"MY", "Malaysia", []int{533}},
		{"BN", "Brunei Darussalam", []int{508}},
		{"ID", "Indonesia", []int{525}},
		{"TL", "East Timor", []int{542}},
		{"PH", "Philippines", []int{548}},
		{"TW", "Taiwan", []int{416}},
		{"HK", "Hong Kong", []int{477}},
		{"MO", "Macao", []int{453}},
		{"MN", "Mongolia", []int{457}},
		{"KP", "Democratic People's Republic of Korea", []int{445}},
		{"IR", "Iran", []int{422}},
		{"IQ", "Iraq", []int{425}},
		{"SY", "Syrian Arab Republic", []int{468}},
		{"LB", "Lebanon", []int{450}},
		{"JO", "Jordan", []int{438}},
		{"IL", "Israel", []int{428}},
		{"PS", "State of Palestine", []int{443}},
		{"SA", "Saudi Arabia", []int{403}},
		{"KW", "Kuwait", []int{447}},
		{"BH", "Bahrain", []int{408}},
		{"QA", "Qatar", []int{466}},
		{"AE", "United Arab Emirates", []int{470}},
		{"OM", "Oman", []int{461}},
		{"YE", "Yemen", []int{473, 475}},
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

// getMID returns the first Maritime Identification Digit for a country (for backward compatibility)
func (m *MMSITool) getMID(countryCode string) int {
	mids := m.getMIDs(countryCode)
	if len(mids) > 0 {
		return mids[0]
	}
	return 0
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

// generateRandomMMSI generates a random MMSI (existing logic)
func (m *MMSITool) generateRandomMMSI(countryCode string) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		// Get MIDs for specific country
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
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

	// Format as 9-digit string with leading zeros if needed
	return fmt.Sprintf("%09d", mmsi), nil
}

// populateTypes initializes the MMSI types slice
func (m *MMSITool) populateTypes() {
	m.types = []MMSIType{
		// Most specific types first
		{
			Name:     "us-coast-guard-ship",
			FullName: "US Coast Guard Group Ship Station",
			IsType: func(mmsi string) bool {
				return mmsi == "036699999"
			},
			Generate: func(countryCode string) (string, error) {
				return "036699999", nil
			},
		},
		{
			Name:     "us-coast-guard-coast",
			FullName: "US Coast Guard Group Coast Station",
			IsType: func(mmsi string) bool {
				return mmsi == "003669999"
			},
			Generate: func(countryCode string) (string, error) {
				return "003669999", nil
			},
		},
		{
			Name:     "us-federal",
			FullName: "US Federal MMSI",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// US Federal: starts with 3669
				return mmsi[:4] == "3669"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateUSFederal()
			},
		},
		{
			Name:     "us-ship-international",
			FullName: "US Ship Station (International/Inmarsat)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// US Ship with Inmarsat: starts with 366 and ends with 000
				return mmsi[:3] == "366" && mmsi[6:] == "000"
			},
			Generate: func(countryCode string) (string, error) {
				// Generate US ship with Inmarsat B/C/M
				remaining := fmt.Sprintf("%03d000", rand.Intn(1000))
				return fmt.Sprintf("366%s", remaining), nil
			},
		},
		{
			Name:     "us-ship-other",
			FullName: "US Ship Station (Other)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// US Ship other: starts with 366, ends with 0, not 000
				return mmsi[:3] == "366" && mmsi[8] == '0' && mmsi[6:] != "000"
			},
			Generate: func(countryCode string) (string, error) {
				// Generate US ship with Inmarsat C
				remaining := fmt.Sprintf("%05d0", rand.Intn(100000))
				return fmt.Sprintf("366%s", remaining), nil
			},
		},
		{
			Name:     "us-ship-regular",
			FullName: "US Ship Station",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// US Ship regular: starts with 366, not ending with 0
				return mmsi[:3] == "366" && mmsi[8] != '0'
			},
			Generate: func(countryCode string) (string, error) {
				// Generate regular US ship
				remaining := fmt.Sprintf("%06d", rand.Intn(1000000))
				return fmt.Sprintf("366%s", remaining), nil
			},
		},
		{
			Name:     "ship-inmarsat-bcm",
			FullName: "Ship Station (Inmarsat B/C/M)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Inmarsat B/C/M: ends with 000, but not coast station (starts with 00)
				return mmsi[6:] == "000" && mmsi[:2] != "00"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateShipStation(countryCode, true, false)
			},
		},
		{
			Name:     "ship-inmarsat-c",
			FullName: "Ship Station (Inmarsat C)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Inmarsat C: ends with 0 but not 000
				return mmsi[8] == '0' && mmsi[6:] != "000"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateShipStation(countryCode, false, true)
			},
		},
		{
			Name:     "group-coast",
			FullName: "Group Coast Station",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Group coast: starts with 00 and ends with 0000 (but not 0000000)
				return mmsi[:2] == "00" && mmsi[5:] == "0000" && mmsi != "000000000"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateCoastStation(countryCode, true)
			},
		},
		{
			Name:     "coast",
			FullName: "Coast Station",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Coast station: starts with 00, not 000, not ending with 0000
				return mmsi[:2] == "00" && mmsi[:3] != "000" && mmsi[5:] != "0000"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateCoastStation(countryCode, false)
			},
		},
		{
			Name:     "group-ship",
			FullName: "Group Ship Station",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Group ship: starts with 0, not 00
				return mmsi[0] == '0' && mmsi[:2] != "00"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateGroupShipStation(countryCode)
			},
		},
		{
			Name:     "sar-aircraft",
			FullName: "SAR Aircraft",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// SAR aircraft: starts with 111
				return mmsi[:3] == "111"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateSARAircraft(countryCode)
			},
		},
		{
			Name:     "handheld-vhf",
			FullName: "Handheld VHF Transceiver",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Handheld VHF: starts with 8
				return mmsi[0] == '8'
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateHandheldVHF()
			},
		},
		{
			Name:     "sar-transponder",
			FullName: "SAR Transponder (AIS-SART)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// SAR transponder: starts with 970
				return mmsi[:3] == "970"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateSARTransponder()
			},
		},
		{
			Name:     "man-overboard",
			FullName: "Man Overboard Device",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Man overboard: starts with 972
				return mmsi[:3] == "972"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateManOverboard()
			},
		},
		{
			Name:     "epirb-ais",
			FullName: "EPIRB-AIS",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// EPIRB-AIS: starts with 974
				return mmsi[:3] == "974"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateEPIRBAIS()
			},
		},
		{
			Name:     "craft-associated",
			FullName: "Craft Associated with Parent Ship",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Craft associated: starts with 98
				return mmsi[:2] == "98"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateCraftAssociated(countryCode)
			},
		},
		{
			Name:     "navigational-aid",
			FullName: "Navigational Aid (AtoN)",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Navigational aid: starts with 99
				return mmsi[:2] == "99"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateNavigationalAid(countryCode)
			},
		},
		{
			Name:     "ship",
			FullName: "Ship Station",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Regular ship: not starting with 0, 1, 8, 9, or 3 (US federal)
				first := mmsi[0]
				return first != '0' && first != '1' && first != '8' && first != '9' && mmsi[:4] != "3669"
			},
			Generate: func(countryCode string) (string, error) {
				return m.generateShipStation(countryCode, false, false)
			},
		},
		{
			Name:     "free-form",
			FullName: "Free-form Device",
			IsType: func(mmsi string) bool {
				if len(mmsi) != 9 {
					return false
				}
				// Free-form: starts with 9, not 97, 98, 99
				return mmsi[0] == '9' && mmsi[:2] != "97" && mmsi[:2] != "98" && mmsi[:2] != "99"
			},
			Generate: func(countryCode string) (string, error) {
				// Generate a free-form device MMSI
				remaining := fmt.Sprintf("%08d", rand.Intn(100000000))
				return fmt.Sprintf("9%s", remaining), nil
			},
		},
	}
}

// generateSpecificType generates a specific type of MMSI using the types slice
func (m *MMSITool) generateSpecificType(mmsiType, countryCode string) (string, error) {
	for _, t := range m.types {
		if t.Name == mmsiType {
			return t.Generate(countryCode)
		}
	}
	return "", fmt.Errorf("unsupported MMSI type: %s", mmsiType)
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
func (m *MMSITool) generateShipStation(countryCode string, inmarsatBCM, inmarsatC bool) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	var suffix string
	if inmarsatBCM {
		suffix = "000" // Inmarsat B/C/M
	} else if inmarsatC {
		suffix = "0" // Inmarsat C
	} else {
		// Regular ship - generate random 6 digits
		suffix = fmt.Sprintf("%06d", rand.Intn(1000000))
	}

	return fmt.Sprintf("%03d%s", selectedMID, suffix), nil
}

// generateGroupShipStation generates a group ship station MMSI
func (m *MMSITool) generateGroupShipStation(countryCode string) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	selectedMID := allMIDs[rand.Intn(len(allMIDs))]
	remaining := fmt.Sprintf("%05d", rand.Intn(100000))

	return fmt.Sprintf("0%03d%s", selectedMID, remaining), nil
}

// generateCoastStation generates a coast station MMSI
func (m *MMSITool) generateCoastStation(countryCode string, group bool) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	selectedMID := allMIDs[rand.Intn(len(allMIDs))]

	if group {
		return fmt.Sprintf("00%03d0000", selectedMID), nil
	} else {
		remaining := fmt.Sprintf("%04d", rand.Intn(10000))
		return fmt.Sprintf("00%03d%s", selectedMID, remaining), nil
	}
}

// generateSARAircraft generates a SAR aircraft MMSI
func (m *MMSITool) generateSARAircraft(countryCode string) (string, error) {
	// SAR aircraft: 111xxxxxx
	remaining := fmt.Sprintf("%06d", rand.Intn(1000000))
	return fmt.Sprintf("111%s", remaining), nil
}

// generateHandheldVHF generates a handheld VHF MMSI
func (m *MMSITool) generateHandheldVHF() (string, error) {
	// Handheld VHF: 8xxxxxxx
	remaining := fmt.Sprintf("%08d", rand.Intn(100000000))
	return fmt.Sprintf("8%s", remaining), nil
}

// generateSARTransponder generates a SAR transponder MMSI
func (m *MMSITool) generateSARTransponder() (string, error) {
	// SAR transponder: 970xxxxxx
	remaining := fmt.Sprintf("%06d", rand.Intn(1000000))
	return fmt.Sprintf("970%s", remaining), nil
}

// generateManOverboard generates a man overboard device MMSI
func (m *MMSITool) generateManOverboard() (string, error) {
	// Man overboard: 972xxxxxx
	remaining := fmt.Sprintf("%06d", rand.Intn(1000000))
	return fmt.Sprintf("972%s", remaining), nil
}

// generateEPIRBAIS generates an EPIRB-AIS MMSI
func (m *MMSITool) generateEPIRBAIS() (string, error) {
	// EPIRB-AIS: 974xxxxxx
	remaining := fmt.Sprintf("%06d", rand.Intn(1000000))
	return fmt.Sprintf("974%s", remaining), nil
}

// generateCraftAssociated generates a craft associated MMSI
func (m *MMSITool) generateCraftAssociated(countryCode string) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	selectedMID := allMIDs[rand.Intn(len(allMIDs))]
	remaining := fmt.Sprintf("%05d", rand.Intn(100000))

	// Craft associated: 98xxxxxxx
	return fmt.Sprintf("98%03d%s", selectedMID, remaining), nil
}

// generateNavigationalAid generates a navigational aid MMSI
func (m *MMSITool) generateNavigationalAid(countryCode string) (string, error) {
	var allMIDs []int
	if countryCode != "" {
		allMIDs = m.getMIDs(countryCode)
		if len(allMIDs) == 0 {
			return "", fmt.Errorf("invalid country code: %s", countryCode)
		}
	} else {
		allMIDs = m.allMIDs
	}

	selectedMID := allMIDs[rand.Intn(len(allMIDs))]
	remaining := fmt.Sprintf("%05d", rand.Intn(100000))

	// Navigational aid: 99xxxxxxx
	return fmt.Sprintf("99%03d%s", selectedMID, remaining), nil
}

// generateUSFederal generates a US Federal MMSI
func (m *MMSITool) generateUSFederal() (string, error) {
	// US Federal MMSI: 3669xxxxx
	remaining := fmt.Sprintf("%05d", rand.Intn(100000))
	return fmt.Sprintf("3669%s", remaining), nil
}

// determineMMSIType determines the type of MMSI based on its format using the types slice
func (m *MMSITool) determineMMSIType(mmsi string) string {
	if len(mmsi) != 9 {
		return "Invalid"
	}

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

	// Validate country code for generation
	if countryCode, ok := params["country-code"]; ok {
		if countryCodeStr, ok := countryCode.(string); ok {
			if m.getMID(countryCodeStr) == 0 {
				return fmt.Errorf("invalid country code: %s", countryCodeStr)
			}
		} else {
			return fmt.Errorf("country-code must be a string")
		}
	}

	// Validate input for validation
	if operation, ok := params["operation"]; ok {
		if operationStr, ok := operation.(string); ok && operationStr == "validate" {
			if input, ok := params["input"]; !ok {
				return fmt.Errorf("input parameter is required for validation")
			} else if _, ok := input.(string); !ok {
				return fmt.Errorf("input must be a string")
			}
		}
	}

	return nil
}

// GetInputSchema returns the JSON schema for tool input parameters
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
			Name:        "type",
			Type:        "string",
			Description: "MMSI type to generate (optional for generation)",
			Required:    false,
			Enum:        m.getSupportedTypes(),
		},
		{
			Name:        "country-code",
			Type:        "string",
			Description: "Country code for generation (e.g., US, GB, DE, FR, etc.)",
			Required:    false,
		},
		{
			Name:        "count",
			Type:        "integer",
			Description: "Number of MMSI numbers to generate (max: 100)",
			Required:    false,
		},
	})
}

// GetOutputSchema returns the JSON schema for tool output
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
