package views

import "github.com/mihirkelkar/lenslocked.com/models"

const (
	AlertLevelError   = "danger"
	AlertLevelWarn    = "warning"
	AlertLevelInfo    = "info"
	AlertLevelSuccess = "success"
)

type Alert struct {
	Level   string
	Message string
}

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

//a new publicerror interface that implements the error interface
//and an extra public function.
type PublicError interface {
	error
	Public() string
}

func (d *Data) SetAlert(err error) {
	var msg string
	//check if this is a public error
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		msg = err.Error()
	}
	d.Alert = &Alert{
		Level:   AlertLevelError,
		Message: msg,
	}
}
