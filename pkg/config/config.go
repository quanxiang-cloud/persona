package config

import (
	"git.internal.yunify.com/qxp/persona/pkg/misc/client"
	"io/ioutil"
	"time"

	"git.internal.yunify.com/qxp/persona/pkg/misc/logger"
	"gopkg.in/yaml.v2"
)

// Config 全局配置对象
var Config *Configs

// Configs 总配置结构体
type Configs struct {
	Model          string        `yaml:"model"`
	Port           string        `yaml:"port"`
	HostName       string        `yaml:"hostName"`
	CallbackURL    string        `yaml:"callback"`
	Log            logger.Config `yaml:"log"`
	InternalNet    client.Config `yaml:"internalNet"`
	Etcd           EtcdConfig    `yaml:"etcd"`
	ES             ESConf        `yaml:"elasticsearch"`
	ProcessorNum   int           `yaml:"processorNum"`
	BackendStorage string        `yaml:"backendStorage"`
}

// HTTPServer http服务配置
type HTTPServer struct {
	Port              string        `yaml:"port"`
	ReadHeaderTimeOut time.Duration `yaml:"readHeaderTimeOut"`
	WriteTimeOut      time.Duration `yaml:"writeTimeOut"`
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
}

// Proxy proxy
type Proxy struct {
	Timeout               time.Duration `yaml:"timeout"`
	KeepAlive             time.Duration `yaml:"keepAlive"`
	MaxIdleConns          int           `yaml:"maxIdleConns"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout"`
	TLSHandshakeTimeout   time.Duration `yaml:"tlsHandshakeTimeout"`
	ExpectContinueTimeout time.Duration `yaml:"expectContinueTimeout"`
}

// EtcdConfig db Config
type EtcdConfig struct {
	Addrs    []string
	Username string
	Password string
	Timeout  time.Duration
}

// ESConf Elasticsearch 配置
type ESConf struct {
	Host         []string
	Username     string
	Password     string
	DefaultIndex string
}

// Init 初始化
func Init(configPath string) error {
	if configPath == "" {
		configPath = "../configs/configs.yml"
	}
	//Config = new(Configs)
	err := read(configPath, &Config)
	if err != nil {
		return err
	}
	return nil
}

// read 读取配置文件
func read(yamlPath string, v interface{}) error {
	// Read config file
	buf, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return nil
}
