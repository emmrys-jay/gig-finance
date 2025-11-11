package utils

import (
	"fmt"
	"strconv"
	"strings"
)

const GIGPrefix = "GIG"

// ParseCustomerID removes the GIG prefix from customer_id if present and returns the numeric ID
func ParseCustomerID(customerID string) (int64, error) {
	// Remove GIG prefix if present
	idStr := strings.TrimPrefix(customerID, GIGPrefix)

	// Parse to int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid customer_id format: %s", customerID)
	}

	return id, nil
}

// FormatCustomerID adds the GIG prefix to customer ID with appropriate padding
// Pads to 5 digits for IDs up to 99999 (total length 8: GIG + 5 digits)
// For IDs exceeding 99999, uses the actual number of digits
func FormatCustomerID(id int64) string {
	idStr := strconv.FormatInt(id, 10)

	// If id is 99999 or less, pad to 5 digits (total length will be 8: GIG + 5 digits)
	if id <= 99999 {
		padding := 5 - len(idStr)
		idStr = strings.Repeat("0", padding) + idStr
	}
	// For IDs greater than 99999, use as-is (e.g., GIG100000)

	return GIGPrefix + idStr
}
