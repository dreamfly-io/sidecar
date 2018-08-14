package config

type ListenerConfig struct {
	Name          string         `json:"name,omitempty"`
	Address       string         `json:"address,omitempty"`

	LogPath  string `json:"log_path,omitempty"`
	LogLevel string `json:"log_level,omitempty"`
}
