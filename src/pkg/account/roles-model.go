//nolint:unparam
package account

const (
	AccRoleAdmin int = 1
	AccRoleOper  int = 2
)

var RolesInclude = map[int][]int{
	AccRoleAdmin: {AccRoleOper},
}

var Roles = map[int]string{
	AccRoleAdmin: "admin",
	AccRoleOper:  "operator",
}

var RevRoles = map[string]int{
	"admin":    AccRoleAdmin,
	"operator": AccRoleOper,
}

//nolint
type AccountRole struct {
	AccountID int `json:"account_id" gorm:"primary_key"`
	Role      int `json:"role" gorm:"primary_key"`
}

func (AccountRole) TableName() string {
	return "account_role"
}

func (model AccountRole) GetStrRole() string {
	return Roles[model.Role]
}

func GetRole(roleName string) int {
	return RevRoles[roleName]
}

func GetIncludedRoles(role int) []int {
	var roles = RolesInclude[role]
	var allRoles []int
	allRoles = append(allRoles, roles...)
	if len(roles) > 0 {
		for _, r := range roles {
			allRoles = append(allRoles, GetIncludedRoles(r)...)
		}
	}
	return allRoles
}
