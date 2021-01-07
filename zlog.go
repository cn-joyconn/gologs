package log

import (
	"encoding/json"
	"os"
	"strings"

	filetool "github.com/cn-joyconn/joyconn-goutils/filetool"
	filelog "github.com/cn-joyconn/joyconn-gologs/filelog"

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

var loggerMap map[string]zap.Logger
var defaultLogger zap.Logger

func init() {
	initDefaultLogger()
	selfDir := filetool.SelfDir()
	configPath := selfDir + "/conf/log.yml"
	var logconfs logConfs
	if filetool.IsExist(configPath) {
		configBytes, err := filetool.ReadFileToBytes(configPath)
		if err != nil {
			defaultLogger.Error(err.Error())
			return
		}
		err = yaml.Unmarshal(configBytes, &logconfs)
		if err != nil {
			defaultLogger.Error("解析log.yml文件失败")
			return
		}
	} else {
		defaultLogger.Error("未找到log.yml")
		return
	}
	loggerMap = make(map[string]zap.Logger)
	var logconf logConf

	for i := 0; i < len(logconfs.Logs); i++ {
		logconf = logconfs.Logs[i]
		logger, err := newLogger(&logconf)
		if err == nil {
			loggerMap[logconf.Name] = *logger
		}

	}
}

//GetLogger 获取一个Logger
func GetLogger(name string) *zap.Logger {
	logger, ok := loggerMap[name]
	if ok {
		return &logger
	}
	return &defaultLogger
}

func initDefaultLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), os.Stdout, zapcore.DebugLevel)

	defaultLogger = *zap.New(core, zap.AddCaller())
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
		zloger = &defaultLogger
	}

	return zloger, nil
}
