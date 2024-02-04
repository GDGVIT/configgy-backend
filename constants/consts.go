package constants

var Prefix = struct {
	SESSION string
}{
	SESSION: "session",
}

var TokenTypes = struct {
	USER  string
	ADMIN string
}{
	USER:  "user",
	ADMIN: "admin",
}

const DefaultFileStoragePath = "./filestore"
