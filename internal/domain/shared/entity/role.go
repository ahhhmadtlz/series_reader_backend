package entity


type Role uint8


const (
	UserRole Role = iota + 1
	ManagerRole 
	AdminRole
)

const (
	UserRoleStr    = "user"
	ManagerRoleStr = "manager"
	AdminRoleStr   = "admin"
)


func (r Role) String() string {
	switch r {
	case UserRole:
		return UserRoleStr
	case ManagerRole:
		return  ManagerRoleStr
	case AdminRole:
		return AdminRoleStr
	}
	return ""
}

func MapToRoleEntity(roleStr string) (Role, bool) {
	switch roleStr {
	case UserRoleStr:
		return UserRole, true
	case ManagerRoleStr:
		return ManagerRole, true
	case AdminRoleStr:
		return AdminRole, true
	}
	return Role(0), false
}