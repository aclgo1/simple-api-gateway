package user

type ErrInvalidEmail struct {
	Message string
}

func (e ErrInvalidEmail) Error() string {
	e.Message = "invalid email"

	return e.Message
}

type ErrEmailCadastred struct {
	Message string
}

func (e ErrEmailCadastred) Error() string {
	e.Message = "email cadastred"
	return e.Message
}

type ErrEmailNotCadastred struct {
	Message string
}

func (e ErrEmailNotCadastred) Error() string {
	e.Message = "email not cadastred"
	return e.Message
}

type ErrEmailSentCheckInbox struct {
	Message string
}

func (e ErrEmailSentCheckInbox) Error() string {
	e.Message = "email sent check this inbox"
	return e.Message
}

type ErrInvalidCode struct {
	Message string
}

func (e ErrInvalidCode) Error() string {
	e.Message = "invalid code"
	return e.Message
}

type ErrEmptyToken struct {
}

func (e ErrEmptyToken) Error() string {
	return "empty token"
}

type ErrUserNotVerified struct {
}

func (e ErrUserNotVerified) Error() string {
	return "user not verified"
}

type ErrFailedVerifyCaptcha struct {
	Message string
}

func (e ErrFailedVerifyCaptcha) Error() string {
	e.Message = "failed to verify captcha"
	return e.Message
}
