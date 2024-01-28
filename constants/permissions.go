package constants

import "github.com/GDGVIT/configgy-backend/pkg/tables"

var CreatePermissions = []tables.Permission{
	tables.AdminPermission,
	tables.OwnerPermission,
}

var UpdatePermissions = []tables.Permission{
	tables.AdminPermission,
	tables.OwnerPermission,
	tables.EditPermission,
}

var DeletePermissions = []tables.Permission{
	tables.AdminPermission,
	tables.OwnerPermission,
}

var ReadPermissions = []tables.Permission{
	tables.AdminPermission,
	tables.OwnerPermission,
	tables.EditPermission,
	tables.ViewPermission,
}
