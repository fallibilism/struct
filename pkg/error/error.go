package e

import "errors"

var (
	ErrorNotImplemented = errors.New("not implemented")
)

const (
	OnlyAdminAllowed = "only admin can perform this action, please contact your admin"
)
