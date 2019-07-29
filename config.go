package libsignal

// Config is an interface for libsignal configuration
type Config interface {
	GetTel() (string, error)
	GetServer() (string, error)

	GetHTTPSignalingKey() ([]byte, error)
	SetHTTPSignalingKey([]byte) error

	GetHTTPPassword() (string, error)
	SetHTTPPassword(string) error
}
