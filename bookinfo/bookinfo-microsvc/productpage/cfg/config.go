package cfg

type Config struct {
	Port        int               `json:"port,omitempty" mapstructure:"port" yaml:"port"`
	ServerMap   map[string]string `json:"server_map,omitempty" mapstructure:"server_map" yaml:"server_map"`
	ServiceName string            `json:"service_name,omitempty" mapstructure:"service_name" yaml:"service_name"`
	JaegerHost  string            `json:"jaeger_host,omitempty" mapstructure:"jaeger_host" yaml:"jaeger_host"`
}
