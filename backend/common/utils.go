package common

import "github.com/google/uuid"

func UUIDtoString(ids []uuid.UUID) []string {
	res := make([]string, len(ids))
	for idx, id := range ids {
		res[idx] = id.String()
	}
	return res
}
