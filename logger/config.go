package logger

type Config struct {
	Level       string `json:",default=Info,options=Debug|Info|Warn|Error|Panic|fatal"` //日志级别，默认为info
	FilePath    string `json:",default=/log,optional"`                                  //日志文件路径
	LogType     string `json:",default=zap,options=logrus|zap,optional"`                //日志类型，默认zap，目前支持zap和logrus
	ServiceName string `json:",optional"`                                               //所属服务
	MaxSize     int    `json:",default=10,optional"`                                    //日志文件最大数量
	MaxAge      int    `json:",default=30,optional"`                                    //最大保留天数
	BackupNum   int    `json:",default=100,optional"`                                   //最大保留日志文件数量
	Compress    bool   `json:",default=false,optional"`                                 //是否压缩
}

type KeyPair struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}
