package accesscontrolsvc

import (
	"strings"

	"github.com/GDGVIT/configgy-backend/pkg/tables"
)

func StringifyListOfPermissions(permissions []tables.Permission) string {
	return strings.Join(PermissionSliceToStringSlice(permissions), ",")
}

func PermissionSliceToStringSlice(permissions []tables.Permission) []string {
	finalString := []string{}
	for _, permission := range permissions {
		finalString = append(finalString, string(permission))
	}
	return finalString
}
