package utils

import "strings"

func TrimStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

// IDGetter is a constraint for types that have an ID field
type IDGetter interface {
	GetID() int64
}

// ConvertArrToMapID converts a slice of items to a map indexed by ID
// If duplicate IDs exist, only the first occurrence is stored
// T must implement IDGetter interface (have a GetID() method)
func ConvertArrToMapID[T IDGetter](ls []T) map[int64]T {
	mapper := make(map[int64]T)
	for _, record := range ls {
		id := record.GetID()
		if _, exist := mapper[id]; !exist {
			mapper[id] = record
		}
	}
	return mapper
}

// ConvertArrToMapIDSlice converts a slice of pointer items to a map of value slices indexed by ID
// Groups all items with the same ID together and dereferences pointers
// getID is a function that extracts the ID from each record
func ConvertArrToMapIDSlice[T any](ls []*T, getID func(*T) int64) map[int64][]T {
	mapper := make(map[int64][]T)
	for _, record := range ls {
		if record != nil {
			id := getID(record)
			mapper[id] = append(mapper[id], *record)
		}
	}
	return mapper
}

func DerefIntOrDefault(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	}
	return *value
}

func DerefInt64OrDefault(value *int64, defaultValue int64) int64 {
	if value == nil {
		return defaultValue
	}
	return *value
}

func DerefInt16OrDefault(value *int16, defaultValue int16) int16 {
	if value == nil {
		return defaultValue
	}
	return *value
}

func DerefFloat64OrDefault(value *float64, defaultValue float64) float64 {
	if value == nil {
		return defaultValue
	}
	return *value
}
