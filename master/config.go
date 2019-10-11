package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type httpConfig struct {
	ApiPort         int `json:"apiPort"`
	ApiReadTimeout  int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
}
type etcdConfig struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}
type mongoDbConfig struct {
	MongodbUri            string `json:"mongodbUri"`
	MongodbConnectTimeout int    `json:"mongodbConnectTimeout"`
}
type LogConfig struct {
	LogDatabase         string `json:"logDatabase"`
	LogCollection       string `json:"logCollection"`
	LogDefaultLimitPage int    `json:"logDefaultLimitPage"`
}
type config struct {
	HttpConfig    httpConfig    `json:"http"`
	EtcdConfig    etcdConfig    `json:"etcd"`
	WebRoot       string        `json:"webroot"`
	MongoDbConfig mongoDbConfig `json:"mongodb"`
	LogConfig     LogConfig     `json:"log"`
}

//设置一个全局单例的config
var (
	G_config *config
)

//加载读取配置文件
func InitConfig(filename string) (err error) {
	var (
		content []byte
		cfg     config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	cfg = config{}
	if err = json.Unmarshal(content, &cfg); err != nil {
		return
	}
	G_config = &cfg
	fmt.Println(G_config)
	return
}
