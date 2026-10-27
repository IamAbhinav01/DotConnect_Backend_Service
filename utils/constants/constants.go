package constants

const ROLE_ADMIN = "admin"
const ROLE_OWNER = "owner"
const ROLE_MEMBER = "member"

func Is_Manager(role string) bool {
	return role == ROLE_ADMIN || role == ROLE_OWNER
}
