package domain

// Is not used right now
type Role string

// Is not used right now
const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	FullName  string `json:"full_name,omitempty" db:"full_name"`
	Phone     string `json:"phone" db:"phone"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"`
	Role      Role   `json:"role" db:"role"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" db:"updated_at"`
}

// Is not used right now
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
