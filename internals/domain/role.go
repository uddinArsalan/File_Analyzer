package domain

type RoleName string

const (
	USER  RoleName = "user"
	ADMIN RoleName = "admin"
)

type Role struct {
	ID   int64    `json:"id"`
	Name RoleName `json:"role_name"`
}
