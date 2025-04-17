package cfg

type Config struct {
	Port      int               `json:"port,omitempty" mapstructure:"port" yaml:"port"`
	ServerMap map[string]string `json:"server_map,omitempty" mapstructure:"server_map" yaml:"server_map"`
}
