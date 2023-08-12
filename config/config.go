package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// GrpcGoConfig grpc_go的静态配置信息
type GrpcGoConfig struct {
	SystemConfig SystemConfig `yaml:"SystemConfig"` // 系统配置
}

// ServerConfig 服务配置
type ServerConfig struct {
	IP       string `yaml:"IP"`       // 服务ip地址
	TrpcPort int32  `yaml:"TrpcPort"` // 服务trpc协议端口
	HttpPort int32  `yaml:"HttpPort"` // 服务http协议端口
}

// LogConfig 日志配置
type LogConfig struct {
	FileName     string `yaml:"FileName"`     // 日志文件
	MaxFileSize  int64  `yaml:"MaxFileSize"`  // 最大文件长度，单位MB
	MaxBackups   int32  `yaml:"MaxBackups"`   // 最大保留时长,单位天
	LevelEnabler int32  `yaml:"LevelEnabler"` // 日志级别,控制输出的日志的最低级别, -1是debug日志
}

// SystemConfig 系统配置
type SystemConfig struct {
	ServerConfig ServerConfig `yaml:"ServerConfig"` // 服务配置
	LogConfig    LogConfig    `yaml:"LogConfig"`    // 日志配置
}

// UserConfig 用户配置
type UserConfig struct {
}

// gGrpcGoConfig 服务静态配置
var gGrpcGoConfig *GrpcGoConfig

func init() {
	file, err := ioutil.ReadFile("grpcgo_formal.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var grpcGoConfig GrpcGoConfig
	err = yaml.Unmarshal(file, &grpcGoConfig)
	if err != nil {
		log.Fatal(err)
	}

	gGrpcGoConfig = &grpcGoConfig
	log.Printf("InitGrpcGoConfig success, grpcGoConfig:%v", grpcGoConfig)
}

// GetSystemConfig 获取系统配置信息
func GetSystemConfig() *SystemConfig {
	if gGrpcGoConfig != nil {
		return &gGrpcGoConfig.SystemConfig
	}
	return &SystemConfig{}
}
