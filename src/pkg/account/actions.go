package account

import (
	"net/http"
	"oko/pkg/cfg"
	"oko/pkg/db"
	"oko/pkg/e"
	"oko/pkg/ginapp/types"
	"oko/pkg/log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SignUp godoc
// @Summary Sign up
// @Description Sing up
// @ID post-account-sign-up
// @Tags Account
// @Accept json
// @Produce json
// @Param object body account.SignUpForm true "Sign up fields"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /account/sign-up [post]
func SignUp(c *gin.Context) {
	var form SignUpForm

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	dbt := db.GetDB().Begin()

	acc := Account{
		Email:  strings.ToLower(form.Email),
		Name:   form.Name,
		Status: AccStatusUnconfirmed,
	}
	if err := acc.SetPassword(form.Password); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	if err := dbt.Create(&acc).Error; err != nil {
		panic(err)
	}

	token := GenerateSignUpToken(acc.ID)
	sut := SignUpToken{
		AccountID: acc.ID,
		Token:     token,
		ExpireAt:  time.Now().Add(time.Second * time.Duration(cfg.App.SignUpTokenLifetime)),
	}
	if err := dbt.Create(&sut).Error; err != nil {
		dbt.Rollback()
		panic(err)
	}

	if err := acc.SendSignUpConfirmed(); err != nil {
		log.Println("Fail to send registration confirmation", err)
	}

	dbt.Commit()

	c.JSON(
		http.StatusCreated,
		types.StdResponse{
			Status:  e.Created,
			Message: e.GetMsg(e.Created),
		},
	)
}

// SignIn godoc
// @Summary Sign in
// @Description Sing in
// @ID post-account-sign-in
// @Tags Account
// @Accept json
// @Produce json
// @Param object body account.SignInForm true "Sign up fields"
// @Success 200 {object} types.ResponseSignIn
// @Failure 401 {object} types.ResponseErrorSwg
// @Router /account/sign-in [post]
func SignIn(c *gin.Context) {
	var form SignInForm

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	acc, err := FindAccount(Account{Email: form.Email})
	if err != nil {
		loginErrorResponse(c, form.Email)
		return
	}
	if acc.ID == 0 {
		loginErrorResponse(c, form.Email)
		return
	}
	err = acc.CheckPassword(form.Password)
	if err != nil {
		loginErrorResponse(c, form.Email)
		return
	}
	token, err := SetToken(acc.ID)
	if err != nil {
		panic(err)
	}

	c.JSON(
		http.StatusOK,
		types.ResponseSignIn{
			StdResponse: types.StdResponse{
				Status:  e.Success,
				Message: e.GetMsg(e.Success),
			},
			Data: types.ResponseAccessToken{AccessToken: token},
		},
	)
}

// SignOut godoc
// @Summary Sign out
// @Description Sign out
// @ID get-account-sign-out
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} types.StdResponse
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /account/sign-out [post]
// @Security ApiKeyAuth
func SignOut(c *gin.Context) {
	err := DropCurrentToken(c)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, types.StdResponse{Status: http.StatusOK, Message: e.GetMsg(http.StatusOK)})
}

// Profile godoc
// @Summary Get account profile
// @Description Get account profile
// @ID get-account-profile
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} types.ResponseProfile
// @Failure 400 {object} types.ResponseErrorSwg
// @Router /account/profile [get]
// @Security ApiKeyAuth
func Profile(c *gin.Context) {
	acc := GetContextAcc(c)
	if acc.ID == 0 {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}
	profile := types.Profile{
		ID:        acc.ID,
		Email:     acc.Email,
		Name:      acc.Name,
		Status:    acc.Status,
		StatusStr: acc.GetStrStatus(),
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": e.GetMsg(http.StatusOK),
			"data":    profile,
		})
}
