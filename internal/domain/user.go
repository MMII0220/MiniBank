package domain

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// в domain не должен быть прсистввовать теги, сделать маппинг в каждом слое где они будут использовать
// json можно оставить а вот db надо урать из domain обязателно прочитать про это. Обязателньо это сделай

// Чистая доменная модель без внешних зависимостей
type User struct {
	ID        int
	FullName  string
	Phone     string
	Email     string
	Password  string
	Role      Role
	CreatedAt string
	UpdatedAt string
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
