package enums

type RoleEnum string

const (
	Free    RoleEnum = "FREE"
	Premium RoleEnum = "PREMIUM"
)

func (r RoleEnum) String() string {
	return string(r)
}

func (r RoleEnum) IsValid() bool {
	switch r {
	case Free, Premium:
		return true
	default:
		return false
	}
}

func GetAllRoles() []RoleEnum {
	return []RoleEnum{Free, Premium}
}
