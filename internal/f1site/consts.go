package f1site

import "time"

// time duration
const (
	LOGIN_EVENT_TIMEOUT_DURATION = 7 * time.Second
	LOGOUT_TIMEOUT_DURATION      = 10 * time.Second
)

// endpoints
const (
	HOMEPAGE_URL = "https://www.formula1.com/"
	LOGIN_URL    = "https://account.formula1.com/#/en/login?redirect=https%3A%2F%2Fwww.formula1.com%2Fen&lead_source=web_f1core"
	LOGOUT_URL   = "https://account.formula1.com/#/en/logout?redirect=https%3A%2F%2Fwww.formula1.com%2Fen"

	BY_PASSWORD_URL = "https://api.formula1.com/v2/account/subscriber/authenticate/by-password"
	REJECT_ALL_URL  = "https://consent.formula1.com/wrapper/v2/choice/reject-all*"
	GDPR_URL        = "https://consent.formula1.com/wrapper/v2/choice/gdpr*"
)

// selectors
var COOKIE_BANNER_SELECTORS = struct {
	I_FRAME            string
	ESSENTIAL_ONLY_BTN string
}{
	I_FRAME:            "#sp_message_iframe_1406947",
	ESSENTIAL_ONLY_BTN: "#notice > div.message-component.message-row.unstack > button.message-component.message-button.no-children.focusable.button.button-hover.sp_choice_type_13",
}

var LOGIN_FORM_SELECTORS = struct {
	EMAIL_INPUT    string
	PASSWORD_INPUT string
	SUBMIT_BTN     string
}{
	EMAIL_INPUT:    "#loginform > div:nth-child(2) > input",
	PASSWORD_INPUT: "#loginform > div.field.password > input",
	SUBMIT_BTN:     "#loginform > div.actions > button",
}

// javascript
const (
	WAIT_COOKIE_EXISTS = `() => document.cookie.includes("user-metadata")`
	WAIT_COOKIE_GONE   = `() => !document.cookie.includes("user-metadata")`
)
