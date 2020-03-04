package account

import (
	"fmt"
	"oko/pkg/cfg"
	"oko/pkg/db"
	"oko/pkg/e"
	"oko/pkg/redis"
	"oko/pkg/rndstr"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GenerateToken(accID int) string {
	return fmt.Sprintf("%d-%s", accID, rndstr.RandString(cfg.App.AuthTokenLength))
}

func GetUserID(token string) (int, error) {
	var intID int
	strID, err := redis.Get(token)
	if err != nil || len(strID) == 0 {
		return 0, err
	}
	intID, err = strconv.Atoi(string(strID))
	if err != nil {
		return 0, err
	}
	return intID, nil
}

func SetToken(accID int) (string, error) {
	token := GenerateToken(accID)
	err := redis.SetEx(token, []byte(strconv.Itoa(accID)), int32(cfg.App.AuthTokenLifetime))
	if err != nil {
		return "", err
	}
	return token, nil
}

func RefreshToken(accID int, token string) bool {
	err := redis.SetEx(token, []byte(strconv.Itoa(accID)), int32(cfg.App.AuthTokenLifetime))
	return err != nil
}

func DropToken(token string) error {
	err := redis.Delete(token)
	return err
}

func SetContextAccount(c *gin.Context, acc Account) {
	c.Set("account_id", acc.ID)
	c.Set("account_model", acc)
}

//nolint
func GetContextAcc(c *gin.Context) Account {
	if acc, ok := c.Get("account_model"); !ok {
		return Account{}
	} else {
		return acc.(Account)
	}
}

func Auth(required bool, roles []int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var acc Account

		code = e.Success

		token := c.GetHeader(cfg.App.AuthTokenKey)
		if token == "" {
			code = e.ErrorAuthToken
		} else {
			accID, err := GetUserID(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else {
				db.GetDB().Preload("Roles").First(&acc, accID)
				if acc.ID != 0 {
					switch acc.Status {
					case AccStatusUnconfirmed:
						code = e.ErrorAuthUnconfirmed
					case AccStatusBaned:
						code = e.ErrorAuthBanned
					case AccStatusActive:
						SetContextAccount(c, acc)
					}
				}
				if _, ok := c.Get("account_id"); !ok && code == e.Success {
					code = e.ErrorAuth
				}
			}
		}

		if len(roles) > 0 && acc.ID != 0 {
			if !CheckRole(acc, roles) {
				code = e.ErrorAuthRole
			}
		}

		if code != e.Success {
			if required {
				e.ErrorResponse(c, code, "Authentication is required to perform this action")
				c.Abort()
				return
			}
			SetContextAccount(c, acc)
		} else {
			accID := c.GetInt("account_id")
			RefreshToken(accID, token)
		}

		c.Next()
	}
}

func GetToken(c *gin.Context) string {
	return c.GetHeader(cfg.App.AuthTokenKey)
}

func DropCurrentToken(c *gin.Context) error {
	token := GetToken(c)
	if token != "" {
		err := DropToken(GetToken(c))
		if err != nil {
			return err
		}
		c.Header(cfg.App.AuthTokenKey, "")
	}
	return nil
}

func CheckRole(acc Account, roles []int) bool {
	allAccRoles := make([]int, 0, len(roles))
	accRoles := acc.Roles
	roleAccessed := false
	for _, role := range accRoles {
		allAccRoles = append(allAccRoles, role.Role)
		allAccRoles = append(allAccRoles, GetIncludedRoles(role.Role)...)
	}
	for _, needRole := range roles {
		for _, existsRole := range allAccRoles {
			if needRole == existsRole {
				roleAccessed = true
				break
			}
		}
		if roleAccessed {
			break
		}
	}
	return roleAccessed
}
