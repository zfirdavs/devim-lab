package cast

import (
	"fmt"
	"strconv"
)

func ToFloat64(name, value string) (float64, error) {
	if len(value) == 0 {
		return 0, fmt.Errorf("the %v argument must not be empty", name)
	}

	casted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to cast to float64: %w", err)
	}
	return casted, nil
}
