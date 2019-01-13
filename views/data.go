package views

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
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	d.Alert = &Alert{Level: AlertLevelError,
		Message: err.Error(),
	}
}
