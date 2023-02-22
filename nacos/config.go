package nacos

// Config nacos配置项
type Config struct {
	ServerIp    string            `json:",optional"`
	ServerPort  uint64            `json:",optional"`
	Servers     []ServerConfig    `json:",optional"`
	ClientIp    string            `json:",optional"`
	ClientPort  uint64            `json:",optional"`
	Metadata    map[string]string `json:",optional"`
	NamespaceId string            `json:",optional"`
	GroupName   string            `json:",optional"`
	ServiceName string            `json:",optional"`
	Username    string            `json:",optional"`
	Password    string            `json:",optional"`
}

type ServerConfig struct {
	ServerIp   string `json:",optional"`
	ServerPort uint64 `json:",optional"`
}

type ConfigMultiServer struct {
	Servers     []ServerConfig
	ClientIp    string            `json:",optional"`
	ClientPort  uint64            `json:",optional"`
	Metadata    map[string]string `json:",optional"`
	NamespaceId string            `json:",optional"`
	GroupName   string            `json:",optional"`
	ServiceName string            `json:",optional"`
	Username    string            `json:",optional"`
	Password    string            `json:",optional"`
}
