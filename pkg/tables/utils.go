package tables

import (
	"strings"

	"github.com/google/uuid"
)

/*
UUIDWithPrefix generates a UUID with a prefix
*/
func UUIDWithPrefix(prefix string) string {
	id := uuid.New().String()
	id = prefix + "_" + id
	id = strings.ReplaceAll(id, "-", "")
	return id
}
