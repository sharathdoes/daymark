package quiz

type ProgressEvent struct {
	Stage   string      `json:"stage"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}