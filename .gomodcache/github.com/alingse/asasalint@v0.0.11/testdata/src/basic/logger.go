package main

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
}

func fromLogger(logger Logger) {
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
