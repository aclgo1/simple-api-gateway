package admin

import (
	"context"
	"errors"
	"regexp"
	"time"
)

type AdminUC interface {
	Create(ctx context.Context, params *ParamsCreateAdmin) (*Admin, error)
	Search(ctx context.Context, params *ParamsSearch) ([]*Admin, error)
	Delete(ctx context.Context, params *ParamsDeleteUser) (string, error)
}

type ParamsCreateAdmin struct {
	Name          string  `json:"name"`
	Lastname      string  `json:"lastname"`
	Password      string  `json:"password"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Verified      string  `json:"verified"`
	Balance       float64 `json:"balance"`
	CaptchaId     string  `json:"captcha_id"`
	CaptchaAwnser string  `json:"captcha_awnser"`
}

func (p *ParamsCreateAdmin) Validate() error {
	if p.Name == "" {
		return errors.New("name empty")
	}
	if p.Lastname == "" {
		return errors.New("lastname empty")
	}
	if p.Password == "" {
		return errors.New("password empty")
	}
	if p.Email == "" {
		return errors.New("email empty")
	}

	regexp := regexp.MustCompile(`^[a-zA-Z0-9._%-+]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexp.MatchString(p.Email) {
		return errors.New("email invalid")
	}

	var finalVerified string
	switch p.Verified {
	case DefaultVerfiedNo:
		finalVerified = DefaultVerfiedNo
	case DefaultVerfiedYes:
		finalVerified = DefaultVerfiedYes
	default:
		finalVerified = DefaultVerfiedNo
	}

	p.Verified = finalVerified

	// if p.CaptchaId == "" {
	// 	return errors.New("captcha id empty")
	// }
	// if p.CaptchaAwnser == "" {
	// 	return errors.New("captcha awnser empty")
	// }

	return nil
}

type Admin struct {
	UserID    string    `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Verified  string    `json:"verified"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParamsSearch struct {
	Query  string `json:"query"`
	Role   string `json:"role"`
	Page   int    `json:"page"`
	OffSet int    `json:"offset"`
	Limit  int    `json:"limit"`
}

func (p *ParamsSearch) Validate() error { return nil }

type ParamsDeleteUser struct {
	UserId string
}

func (p *ParamsDeleteUser) Validate() error {
	return nil
}

type ErrEmaiCadastred struct {
	Message string
}

func (e ErrEmaiCadastred) Error() string {
	return "email cadastred"
}

type ErrEmailSentCheckInbox struct {
	Message string
}

func (e ErrEmailSentCheckInbox) Error() string {
	e.Message = "email sent check this inbox"
	return e.Message
}

type ErrFailedVerifyCaptcha struct {
	Message string
}

func (e ErrFailedVerifyCaptcha) Error() string {
	e.Message = "failed to verify captcha"
	return e.Message
}

var (
	DefaultVerfiedNo      = "no"
	DefaultVerfiedYes     = "yes"
	BASE_API_URL          = ""
	DefaultTimeSendEmails = time.Hour
	DefaultFromSendMail   = ""

	DefaultServiceName = "gmail"

	DefaultSubjectSendConfirm  = "Confirm signup"
	DefaulfBodySendConfirm     = ""
	DefaulfTemplateSendConfirm = `
        <div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
            <h2>Welcome!</h2>
            <p>Thank you for registering. Click the button below to verify your account:</p>
            <a href="%s" 
               style="background-color: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">
               Confirm You Account
            </a>
            <p style="margin-top: 20px; font-size: 12px; color: #666;">
               If you did not request this email, you can safely ignore it.
            </p>
        </div>`
	DefaultServiceSendName = "gmail"
)

func SetConfigUserPackage(apiurl string, defaultFromSendMail string, defaultTimeSendEmails time.Duration, defaultServiceSendEmailName string) {
	BASE_API_URL = apiurl
	DefaulfBodySendConfirm = BASE_API_URL + "/api/user/confirm/%s"

	DefaultFromSendMail = defaultFromSendMail
	DefaultTimeSendEmails = defaultTimeSendEmails
	DefaultServiceName = defaultServiceSendEmailName
}
