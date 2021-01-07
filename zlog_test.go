package log

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestZlog(t *testing.T) {
	go func() {
		i := 0
		logger := GetLogger("xadmin").With(zap.String("tag", "aaa"))
		for {
			i++
			logger.Info(fmt.Sprintf("我是日志%d", i))
		}
	}()
	go func() {
		i := 0
		logger := GetLogger("xadmin").With(zap.String("tag", "bbb"))
		for {
			i++
			logger.Info(fmt.Sprintf("我是日志%d", i))
		}
	}()

	go func() {
		i := 0
		logger := GetLogger("xadmin2")
		for {
			i++
			logger.Info(fmt.Sprintf("我是日志%d", i))
		}
	}()
	// go func() {
	// 	i := 0
	// 	logger := log.GetLogger("xadmin3")
	// 	for {
	// 		i++
	// 		logger.Info(fmt.Sprintf("PPPPPP我是日志PPPPPP%d",i))
	// 	}
	// }()
	time.Sleep(time.Duration(1) * time.Second)
}
