package httprs

import (
	"fmt"
	"strconv"
	"strings"
)

// ContentRange represents a content range header.
type ContentRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Total int64 `json:"total"`
}

func ParseContentRange(header string) (*ContentRange, error) {
	cr := &ContentRange{-1, -1, -1}

	// Clean up the header.
	clean := strings.TrimSpace(header)
	if !strings.HasPrefix(clean, "bytes ") {
		return nil, fmt.Errorf("invalid content range header: %s", header)
	}
	clean = strings.TrimPrefix(clean, "bytes ")

	// Split into parts
	parts := strings.Split(clean, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid content range header: %s", header)
	}

	// Parse the start and end.
	if startEnd := strings.Split(parts[0], "-"); len(startEnd) == 2 {
		start, err := strconv.ParseInt(startEnd[0], 10, 64)
		if err != nil {
			start = -1
		}
		end, err := strconv.ParseInt(startEnd[1], 10, 64)
		if err != nil {
			end = -1
		}
		cr.Start = start
		cr.End = end
	}

	// Parse the total.
	total, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		total = -1
	}
	cr.Total = total

	return cr, nil
}
