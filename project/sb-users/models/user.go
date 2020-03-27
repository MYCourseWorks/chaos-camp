package models

// RoleFlags comment
type RoleFlags int

const (
	// Anonymous comment
	Anonymous RoleFlags = 0
	// SiteUser comment
	SiteUser RoleFlags = 1 << iota
	// Operator comment
	Operator RoleFlags = 1 << iota
	// Admin comment
	Admin RoleFlags = 1 << iota
)

// User models a REST API user
type User struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"-"`
	Roles    RoleFlags `json:"roles"`
}
