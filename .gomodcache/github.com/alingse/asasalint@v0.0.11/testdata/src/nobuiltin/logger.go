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

	logger.Debug(a)        // want `pass \[\]any as any to func logger.Debug func\(args \.\.\.interface{}\)`
	logger.Debugf("%v", a) // want `pass \[\]any as any to func logger.Debugf func\(format string, args \.\.\.interface{}\)`
	logger.Debugln(a)      // want `pass \[\]any as any to func logger.Debugln func\(args \.\.\.interface{}\)`

	logger.Info(a)        // want `pass \[\]any as any to func logger.Info func\(args \.\.\.interface{}\)`
	logger.Infof("%v", a) // want `pass \[\]any as any to func logger.Infof func\(format string, args \.\.\.interface{}\)`
	logger.Infoln(a)      // want `pass \[\]any as any to func logger.Infoln func\(args \.\.\.interface{}\)`

	logger.Print(a)        // want `pass \[\]any as any to func logger.Print func\(args \.\.\.interface{}\)`
	logger.Printf("%v", a) // want `pass \[\]any as any to func logger.Printf func\(format string, args \.\.\.interface{}\)`
	logger.Println(a)      // want `pass \[\]any as any to func logger.Println func\(args \.\.\.interface{}\)`

	logger.Warn(a)        // want `pass \[\]any as any to func logger.Warn func\(args \.\.\.interface{}\)`
	logger.Warnf("%v", a) // want `pass \[\]any as any to func logger.Warnf func\(format string, args \.\.\.interface{}\)`
	logger.Warnln(a)      // want `pass \[\]any as any to func logger.Warnln func\(args \.\.\.interface{}\)`

	logger.Warning(a)        // want `pass \[\]any as any to func logger.Warning func\(args \.\.\.interface{}\)`
	logger.Warningf("%v", a) // want `pass \[\]any as any to func logger.Warningf func\(format string, args \.\.\.interface{}\)`
	logger.Warningln(a)      // want `pass \[\]any as any to func logger.Warningln func\(args \.\.\.interface{}\)`

	logger.Error(a)        // want `pass \[\]any as any to func logger.Error func\(args \.\.\.interface{}\)`
	logger.Errorf("%v", a) // want `pass \[\]any as any to func logger.Errorf func\(format string, args \.\.\.interface{}\)`
	logger.Errorln(a)      // want `pass \[\]any as any to func logger.Errorln func\(args \.\.\.interface{}\)`

	logger.Fatal(a)        // want `pass \[\]any as any to func logger.Fatal func\(args \.\.\.interface{}\)`
	logger.Fatalf("%v", a) // want `pass \[\]any as any to func logger.Fatalf func\(format string, args \.\.\.interface{}\)`
	logger.Fatalln(a)      // want `pass \[\]any as any to func logger.Fatalln func\(args \.\.\.interface{}\)`

	logger.Panic(a)        // want `pass \[\]any as any to func logger.Panic func\(args \.\.\.interface{}\)`
	logger.Panicf("%v", a) // want `pass \[\]any as any to func logger.Panicf func\(format string, args \.\.\.interface{}\)`
	logger.Panicln(a)      // want `pass \[\]any as any to func logger.Panicln func\(args \.\.\.interface{}\)`
}
