package common

import "github.com/google/uuid"

// UUIDtoString convert a slice of `uuid.UUID` to a slice of strings
func UUIDtoString(ids []uuid.UUID) []string {
	res := make([]string, len(ids))
	for idx, id := range ids {
		res[idx] = id.String()
	}
	return res
}
