package account

const (
	SignUpConfSubject = "Sign Up confirmation"
	SignUpConfTmpl    = "sign-up-conf"
	SignUpLink        = "/confirm-sign-up"

	SignUpConfirmedSubject = "Sign Up confirmed"
	SignUpConfirmedTmpl    = "sign-up-confirmed"

	RecoverLink           = "/recover"
	RecoverSubject        = "Password recover"
	ChangePasswordSubject = "Change password"
	RecoverTmpl           = "recover"

	AccStatusUnconfirmed uint = 1
	AccStatusActive      uint = 2
	AccStatusBaned       uint = 3
	AccStatusService     uint = 4
)
