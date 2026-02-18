package user

import (
	"time"
)

var (
	BASE_API_URL          string
	DefaultVerifiedNo     = "no"
	DefaultTimeSendEmails = time.Hour
	DefaultFromSendMail   = "marcellosanttos2014@gmail.com"

	DefaultServiceName         = "gmail"
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

	DefaultSubjectResetPass  = "Reset Pass"
	DefaultBodyResetPass     = ""
	DefaultTemplateResetPass = `
        <div style="font-family: sans-serif; max-width: 600px; margin: 0 auto;">
            <h2>Reset You Pass</h2>
            <p>We received a request to change your account password.</p>
            <p>Click the link below to create a new password:</p>
            <p><a href="%s">Reset Password Now</a></p>
            <p style="color: #d9534f;"><strong>Attention:</strong> This link expires in 1 hours for security reasons.</p>
            <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
            <p style="font-size: 12px; color: #666;">
                If you didn't request this, we recommend changing your current password or contacting support.
            </p>
        </div>`
)

func SetBaseAPiUrl(apiurl string) {
	BASE_API_URL = apiurl
	DefaulfBodySendConfirm = BASE_API_URL + "/api/user/confirm/%s"
	DefaultBodyResetPass = BASE_API_URL + "/newpass?code=%s"
}
