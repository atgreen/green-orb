package main

import "testing"

func TestLogger(t *testing.T) {
	var logger Logger

	var a = []any{1, 2, 3}

	logger.Debug(a)
	logger.Debugf("%v", a)
	logger.Debugln(a)

	logger.Info(a)
	logger.Infof("%v", a)
	logger.Infoln(a)

	logger.Print(a)
	logger.Printf("%v", a)
	logger.Println(a)

	logger.Warn(a)
	logger.Warnf("%v", a)
	logger.Warnln(a)

	logger.Warning(a)
	logger.Warningf("%v", a)
	logger.Warningln(a)

	logger.Error(a)
	logger.Errorf("%v", a)
	logger.Errorln(a)

	logger.Fatal(a)
	logger.Fatalf("%v", a)
	logger.Fatalln(a)

	logger.Panic(a)
	logger.Panicf("%v", a)
	logger.Panicln(a)
}
