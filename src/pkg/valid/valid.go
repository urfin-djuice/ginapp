package valid

import (
	"oko/pkg/account"
	"oko/pkg/cfg"
	"oko/pkg/db"
	"oko/pkg/repost"
	"oko/pkg/util"
	"oko/srv/proxy/entity"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type Item struct {
	Key     string
	Handler validator.Func
	Message string
}

var Validators = []Item{
	{"StrongPass", StrongPass, "Password must contain one lowercase letter, one uppercase letter and one number and must be at least 8 characters"}, //nolint
	{"UniqueEmail", UniqueEmail, "Email already exist"},
	{"ExistsEmail", ExistsEmail, "Not registered e-mail"},
	{"NotEmpty", NotEmpty, "Field must not be empty"},
	{"UniqueList", UniqueList, "Duplicate items in list"},
	{"SignUpToken", SignUpToken, "Invalid sign up token"},
	{"RecoverToken", RecoverToken, "Invalid recover token"},
	{"RecoverTokenNotExp", RecoverTokenNotExp, "Recover token expired"},
	{"RecoverTokenNotUsed", RecoverTokenNotUsed, "Recover token already used"},
	{"UniqueRepostRequest", UniqueRepostRequest, "Repost request already exist"},
	{"ExistsRepostRequest", ExistsRepostRequest, "Repost request not found"},
	{"CheckRuleStatus", CheckRuleStatus, "Unknown rule status"},
	{"email", nil, "Is not a valid e-mail"},
	{"url", nil, "Is not a valid URL"},
	{"eqfield", nil, "Don\"t match"},
	{"required", nil, "Field is required"},
}

func StrongPass(fl validator.FieldLevel) bool {
	if pass, ok := fl.Field().Interface().(string); ok {
		if len(pass) < cfg.App.MinPassLen ||
			!strings.ContainsAny(pass, "1234567890") ||
			strings.ToLower(pass) == pass ||
			strings.ToUpper(pass) == pass {
			return false
		}
	}
	return true
}

func UniqueEmail(fl validator.FieldLevel) bool {
	if email, ok := fl.Field().Interface().(string); ok {
		var acc account.Account
		db.GetDB().Where("email = ?", email).First(&acc)
		if acc.ID != 0 {
			return false
		}
	}
	return true
}

func ExistsEmail(fl validator.FieldLevel) bool {
	if email, ok := fl.Field().Interface().(string); ok {
		var acc account.Account
		db.GetDB().Where("email = ?", email).First(&acc)
		if acc.ID == 0 {
			return false
		}
	}
	return true
}

func NotEmpty(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().([]string); ok && len(val) == 0 {
		return false
	}
	return true
}

func UniqueList(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().([]string); ok && len(val) > 0 {
		bl := make(map[string]bool)
		for _, item := range val {
			if bl[item] {
				return false
			}
			bl[item] = true
		}
	}
	return true
}

func SignUpToken(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		sut, err := account.GetBySignUpToken(val)
		if err != nil {
			return false
		}
		if sut.Account.ID == 0 {
			return false
		}
	}
	return true
}

func RecoverToken(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		rt, err := account.GetByRecoverToken(val)
		if err != nil {
			return false
		}
		if rt.Account.ID == 0 {
			return false
		}
	}
	return true
}

func RecoverTokenNotExp(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		rt, err := account.GetRecoverToken(val)
		if err != nil {
			return false
		}
		if rt.ExpireAt.Before(time.Now()) {
			return false
		}
	}
	return true
}
func RecoverTokenNotUsed(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		rt, err := account.GetRecoverToken(val)
		if err != nil {
			return false
		}
		if rt.IsUsed {
			return false
		}
	}
	return true
}

func UniqueRepostRequest(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		var m repost.Request
		if db.GetDB().Where("(url = ? or url = ?)", val, util.URLEncoded(val)).First(&m).RecordNotFound() {
			return true
		}
	}
	return false
}

func ExistsRepostRequest(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(string); ok {
		var m repost.Request
		if db.GetDB().Where("(url = ? or url = ?)", val, util.URLEncoded(val)).First(&m).RecordNotFound() {
			return false
		}
	}
	return true
}

func CheckRuleStatus(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(int); ok {
		if val != entity.RuleStatusEnabled && val != entity.RuleStatusDisabled {
			return false
		}
	}
	return true
}
