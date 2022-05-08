package log

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	filelog "github.com/cn-joyconn/gologs/filelog"
	filetool "github.com/cn-joyconn/goutils/filetool"

	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"
)

//配置文件中字母要小写，结构体属性首字母要大写
type logConf struct {
	Name      string
	Adapter   string
	Formatter string
	Conf      map[string]interface{}
}
type logConfs struct {
	Logs []logConf
}

var loggerMap map[string]*zap.Logger
var defaultLogger *zap.Logger

func init() {
	if loggerMap == nil {
		selfDir := filetool.SelfDir()
		configPath := selfDir + "/conf/log.yml"
		initLoggers(configPath)
	}
}

//加载日志配置
func LoadConfig(configPath string) {
	if strings.HasPrefix(configPath, "./") {
		configPath = filetool.SelfDir() + configPath[1:]
	}
	for {
		if !strings.Contains(configPath, "//") {
			break
		}
		configPath = strings.ReplaceAll(configPath, "//", "/")
	}
	initLoggers(configPath)
}
func initLoggers(configPath string) {

	var logconfs logConfs
	if filetool.IsExist(configPath) {
		configBytes, err := filetool.ReadFileToBytes(configPath)
		if err != nil {
			fmt.Println(err.Error())
			// defaultLogger.Error(err.Error())
			return
		}
		err = yaml.Unmarshal(configBytes, &logconfs)
		if err != nil {
			fmt.Println("解析log.yml文件失败")
			// defaultLogger.Error("解析log.yml文件失败")
			return
		}
	} else {
		fmt.Println("未找到log.yml")
		// defaultLogger.Error("未找到log.yml")
		return
	}
	_loggerMap := make(map[string]*zap.Logger)
	var logconf logConf
	initDefalutLoggered := false
	for i := 0; i < len(logconfs.Logs); i++ {
		logconf = logconfs.Logs[i]
		logger, err := newLogger(&logconf)
		if err == nil {
			_loggerMap[logconf.Name] = logger
		}
		if logconf.Name == "default" {
			initDefalutLoggered = true
			defaultLogger = logger
		}

	}
	if !initDefalutLoggered {
		initDefaultLogger()
	}
	loggerMap = _loggerMap
}

//GetLogger 获取一个Logger
func GetLogger(name string) *zap.Logger {
	logger, ok := loggerMap[name]
	if ok {
		return logger
	}
	return defaultLogger
}
func initDefaultLogger() {
	defaultLogger = creatConsoleLogger(zapcore.DebugLevel)
}
func creatConsoleLogger(level zapcore.Level) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, level)

	resultLogger := *zap.New(core, zap.AddCaller())
	return &resultLogger
	// defaultLogger.
}
func newLogger(lc *logConf) (*zap.Logger, error) {
	var zloger = new(zap.Logger)
	if strings.ToLower(lc.Adapter) == "file" {
		if lc.Conf != nil {
			enc, err := json.Marshal(lc.Conf)
			if err == nil {
				logger, lerr := filelog.NewFileLogger(string(enc))
				if lerr != nil {
					return nil, lerr
				}
				zloger = logger
			} else {
				return nil, err
			}
		}
	} else if strings.ToLower(lc.Adapter) == "console" {
		level := zapcore.DebugLevel
		_level, ok := lc.Conf["level"]
		if !ok {
			level = zapcore.DebugLevel
		} else {
			switch _level.(type) {
			case int:
				level = _level.(zapcore.Level)
			case string:
				switch strings.ToLower(_level.(string)) {
				case "debug":
					level = zapcore.DebugLevel
				case "info":
					level = zapcore.InfoLevel
				case "error":
					level = zapcore.ErrorLevel
				case "warn":
					level = zapcore.WarnLevel
				case "fetal":
					level = zapcore.FatalLevel
				case "dpanic":
					level = zapcore.DPanicLevel
				case "panic":
					level = zapcore.FatalLevel
				}
			}
		}
		zloger = creatConsoleLogger(level)
	}

	return zloger, nil
}
