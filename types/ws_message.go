package types

const (
	Idle = iota
	UploadStarted
	UploadFinished
	UploadError
	UploadCancelled
	NotFound
)

type WSMessage struct {
	Filename     string  `json:"fileName"`
	Percentage   float64 `json:"percent"`
	Status       int     `json:"status"`
	ErrorMessage string  `json:"errorMessage, omitempty"`
}

type Persist func()

var persistStatus Persist

func NewWSMessage(fileName string, persist Persist) *WSMessage {
	persistStatus = persist
	return &WSMessage{
		Filename:     fileName,
		Percentage:   0,
		Status:       Idle,
		ErrorMessage: "",
	}
}

func (up *WSMessage) SetPercentage(percentage float64) {
	up.Percentage = percentage
	persistStatus()
}

func (up *WSMessage) SetUploadStarted() {
	up.Status = UploadStarted
	persistStatus()
}

func (up *WSMessage) SetUploadFinished() {
	up.Status = UploadFinished
	up.Percentage = 100
	persistStatus()
}

func (up *WSMessage) SetUploadError(errorMessage string) {
	up.Status = UploadError
	up.ErrorMessage = errorMessage
	persistStatus()
}

func (up *WSMessage) SetUploadCancelled() {
	up.Status = UploadCancelled
	persistStatus()
}
