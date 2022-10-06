package yaml

import (
	"github.com/yunduansing/gocommon/logger"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//解析配置文件
func Resolver(filename string, out interface{}) error {
	summary := "读取配置文件"
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatal(summary, zap.String("errinfo", "读取配置文件错误"), zap.Error(err))
		return err
	}

	err = yaml.Unmarshal(yamlFile, out)
	if err != nil {
		logger.Fatal(summary, zap.String("errinfo", "配置文件解码错误"), zap.Error(err))
		return err
	}
	return nil
}
