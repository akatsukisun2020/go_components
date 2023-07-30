package serverinit

import (
	"github.com/akatsukisun2020/go_components/config"
	_ "github.com/akatsukisun2020/go_components/config" // 配置
	"github.com/akatsukisun2020/go_components/logger"
)

// 服务初始化总体框架
// 包含：配置、日志等。

// ServerInit 服务初始化
func ServerInit() {
	// 获取系统配置
	config := config.GetSystemConfig()

	// 初始化日志
	logger.InitLogger(&logger.UdfLoggerConfig{
		FileName:     config.LogConfig.FileName,
		MaxFileSize:  config.LogConfig.MaxFileSize,
		MaxBackups:   config.LogConfig.MaxBackups,
		LevelEnabler: config.LogConfig.LevelEnabler,
	})
}
