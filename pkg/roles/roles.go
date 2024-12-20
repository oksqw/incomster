package roles

const (
	Admin = "admin"
	User  = "user"
)

func IsValid(role string) bool {
	switch role {
	case Admin, User:
		return true
	default:
		return false
	}
}
