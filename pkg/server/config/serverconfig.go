package config


type ServerConfig struct {
	//default logger
	DefaultLogPath  string `json:"default_log_path,omitempty"`
	DefaultLogLevel string `json:"default_log_level,omitempty"`

	//go processor number
	Processor int `json:"processor"`

	Listeners []*ListenerConfig `json:"listeners,omitempty"`
}
