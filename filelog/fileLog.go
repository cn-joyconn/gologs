package filelog

import (
	"encoding/json"
	"errors"

	strtool "github.com/cn-joyconn/goutils/strtool"

	lumberjack "github.com/natefinch/lumberjack"
	zap "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

type fileLogConf struct {
	Filename   string `json:"filename" yaml:"filename"`
	Level      string `json:"level" yaml:"level"`
	Maxsize    int    `json:"maxsize" yaml:"maxsize"`    //文件大小限制,单位MB
	Maxbackups int    `json:"maxbackups" yaml:"maxbackups"` //最大保留日志文件数量
	Maxage     int    `json:"maxage" yaml:"maxage"`     //日志文件保留天数
	Compress   bool   `json:"compress" yaml:"compress"`  //是否压缩处理
	LocalTime  bool `json:"localtime" yaml:"localtime"`
}

//NewFileLogger 获取file Logger
// filename:         //日志文件存放目录
// maxsize:          //文件大小限制,单位MB
// maxbackups:       //最大保留日志文件数量
// maxage:           //日志文件保留天数
// level:           //日志等级
// compress:         //是否压缩处理
func NewFileLogger(conf string) (*zap.Logger, error) {
	var flc fileLogConf
	err := json.Unmarshal([]byte(conf), &flc)
	if err != nil {
		return nil, err
	}
	if strtool.IsBlank(flc.Filename) {
		return nil, errors.New("filename is blank")
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// ext := path.Ext(filename)
	// filename = filename[0:len(filename)-len(ext)]
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   flc.Filename,   //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    flc.Maxsize,    //文件大小限制,单位MB
		MaxBackups: flc.Maxbackups, //最大保留日志文件数量
		MaxAge:     flc.Maxage,     //日志文件保留天数
		Compress:   flc.Compress,   //是否压缩处理
	})
	recordLevel := zapcore.Level(-1)
	if !strtool.IsBlank(flc.Level) {
		err := recordLevel.Set(flc.Level)
		if err != nil {
			recordLevel = zapcore.Level(-1)
		}
	}
	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= recordLevel
	})
	// defaultFileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	defaultFileCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(infoFileWriteSyncer), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr := make([]zapcore.Core, 0)
	coreArr = append(coreArr, defaultFileCore)
	log := zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller()为显示文件名和行号，可省略
	// log := zap.New(zapcore.NewTee(coreArr...)) //zap.AddCaller()为显示文件名和行号，可省略

	return log, nil
}
