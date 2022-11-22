package log

import "testing"

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "msg", "helloworld", "k", "v")

	logger = DefaultLogger
	logger = With(logger)
	_ = logger.Log(LevelDebug, "msg", "helloworld")
}
