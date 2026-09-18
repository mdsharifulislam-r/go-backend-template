package email

import (
	"fmt"
	"time"
)

const primary = "#009A54"

func baseTemplate(content string) string {
	return fmt.Sprintf(`
<body style="margin:0;padding:0;background:#f4f6f8;font-family:Arial,Helvetica,sans-serif;">
  <table width="100%%" cellspacing="0" cellpadding="0" border="0" style="background:#f4f6f8;padding:20px 0;">
    <tr>
      <td align="center">
        <table width="600" cellspacing="0" cellpadding="0" border="0"
          style="background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 4px 12px rgba(0,0,0,0.05);">
          <tr>
            <td align="center" style="padding:30px 20px;border-bottom:1px solid #eee;">
              <h2 style="margin:0;color:%s;">Go Backend</h2>
            </td>
          </tr>
          <tr>
            <td style="padding:30px 25px;">%s</td>
          </tr>
          <tr>
            <td align="center" style="padding:20px;color:#999;font-size:12px;border-top:1px solid #eee;">
              © %d Go Backend Template. All rights reserved.
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>`, primary, content, time.Now().Year())
}

func otpBox(otp int) string {
	return fmt.Sprintf(`
<div style="background:%s;color:#fff;font-size:28px;letter-spacing:4px;padding:14px 24px;border-radius:8px;display:inline-block;margin:20px 0;font-weight:bold;">
%d
</div>`, primary, otp)
}

func CreateAccount(name, emailAddr string, otp int) Message {
	content := fmt.Sprintf(`
    <h2 style="margin:0 0 10px;color:#111;">Verify your account</h2>
    <p style="color:#555;font-size:15px;line-height:1.6;">
      Hi %s,<br/><br/>
      Welcome! Use the verification code below to activate your account.
    </p>
    <div style="text-align:center;">%s</div>
    <p style="color:#666;font-size:14px;">This code is valid for <b>3 minutes</b>.</p>
  `, name, otpBox(otp))

	return Message{
		To:      emailAddr,
		Subject: "Verify your account",
		HTML:    baseTemplate(content),
	}
}

func ResetPassword(emailAddr string, otp int) Message {
	content := fmt.Sprintf(`
    <h2 style="margin:0 0 10px;color:#111;">Reset your password</h2>
    <p style="color:#555;font-size:15px;line-height:1.6;">
      We received a request to reset your password. Use the code below to continue.
    </p>
    <div style="text-align:center;">%s</div>
    <p style="color:#666;font-size:14px;">This code is valid for <b>3 minutes</b>.</p>
  `, otpBox(otp))

	return Message{
		To:      emailAddr,
		Subject: "Reset your password",
		HTML:    baseTemplate(content),
	}
}
