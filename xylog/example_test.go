package xylog_test

import (
	"time"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xylog"
)

// Summary:
//      Xylog allows a module to log messages to specified outputs with custom
// and built-in log levels. It has already provided six log levels: INFO, WARN,
// ERROR, CRITICAL, DEBUG, TRACE. But you can use your own log level if needed.
//
//      Xylog is able to turn on or off logs of a module easily, this helps you
//  control which module could write the log in your application.
//

// To create a logger, you must have a module first.
var ExampleModule = xyplatform.NewModule(430000, "ExampleModule")

// Then call xylog.Register with that module. Like xyerror.Register, you should
// call this function in the global scope with a dummy declaration (avoid
// syntax error).
var _ = xylog.Register(ExampleModule)

func ExampleConfig() {
	// The default configuration is set to logger when you register the module
	// with xylog.Register. The default configuration is:
	//     - Allow all log levels to be printed.
	//     - The log string's format is "time -- module [level] message".
	//     - Print to stdandard output.

	// You can customize configurations of a logger by xylog.Config.

	// Allow is the configuration indicating which log level can be printed.
	// Use "ALL" if you want to print all log levels.
	xylog.Config(ExampleModule, xylog.Allow("INFO", "Error"))

	// Format configures the log string format.
	// Xylog provides some macros to format the log string:
	//     $TIME$    - time in format dd-mm-yy hh:mm:ss when logging
	//     $LEVEL$   - the log level
	//     $MODULE$  - module name
	//     $MESSAGE$ - log message
	xylog.Config(ExampleModule, xylog.Format("[$LEVEL$][$MODULE$] - $MESSAGE$"))

	// Writer configures which output the log should print.
	// Xylog provides three types of writer:
	//     Stdout - print logs to standard output
	//     File   - print logs to a specified file
	//     SFile  - similar to File, but it will write to another file if a
	//              stop condition is reached. See ExampleSFile.
	xylog.Config(ExampleModule, xylog.Writer(xylog.Stdout))

	// Note: Config() function allows you to pass a list of configurations
	// instead of one by one.

	// This example doesn't print any output, it only affects on the ExampleLog
	// function below.

	// Output:
}

func ExampleLog() {
	// After you register a module (and configure if needed), you can begin
	// to log anything you want.

	// You can use any custom log level with xylog.Log.
	xylog.Log(ExampleModule, "CUSTOM", "Any level can be logged")

	// Or use built-in log levels.
	xylog.Info(ExampleModule, "Info")
	xylog.Error(ExampleModule, "Something %s", "wrong")

	// The below log level isn't printed because of Allow configuration in the
	// ExampleConfig() function.
	xylog.Debug(ExampleModule, "Debug something")

	// Output:
	// [INFO][ExampleModule] - Info
	// [ERROR][ExampleModule] - Something wrong
}

func ExampleLogger() {
	// Instead of calling functions with the module as a parameter, you can get
	// the logger of that module by calling xylog.Logger.
	ExampleLogger := xylog.Logger(ExampleModule)

	// Functions of logger don't need module as a parameter to call.
	ExampleLogger.Config(
		xylog.Format("$MESSAGE$ of ExampleLogger"),
		xylog.Allow("Critical"),
	)

	ExampleLogger.Critical("Something is critical")

	// Output:
	// Something is critical of ExampleLogger
}

func ExampleFile() {
	// xylog.File allows the log to print to a specified file.
	xylog.Config(ExampleModule, xylog.Writer(xylog.File("example_test.log")))

	// Now, call log functions normally and all log messages will be put in
	// file instead.
	xylog.Config(ExampleModule, xylog.Allow("INFO"))
	xylog.Info(ExampleModule, "Info in file")

	// This example doesn't print any output, see the log file for more details.

	// Output:
}

func ExampleSFile() {
	// Printing to only one file can cause some problems of file size. You can
	// use SFile of Writer configuration to stop printing to a file when a
	// condition is reached, the logger then creates another file to print.
	fnFormat := "example_log_%s.log"
	xylog.Config(ExampleModule, xylog.Writer(xylog.SFile(fnFormat, xylog.LimitSize(50*xylog.Byte))))
	// Or you can split file by day
	// xylog.Config(ExampleModule, xylog.Writer(xylog.SFile(fnFormat, xylog.TimePeriod(xylog.Day))))
	// Or after one hour
	// xylog.Config(ExampleModule, xylog.Writer(xylog.SFile(fnFormat, xylog.TimeAfter(time.Hour))))

	// Now you can log normally and the log message will be printed to many files.
	xylog.Config(ExampleModule, xylog.Allow("Info", "Error"))
	time.Sleep(time.Second)

	xylog.Info(ExampleModule, "Something is logged but too long")
	time.Sleep(time.Second)

	xylog.Error(ExampleModule, "Something is wrong but too long")
	time.Sleep(time.Second)

	xylog.Error(ExampleModule, "Something is wrong but too long")
	time.Sleep(time.Second)

	// This example doesn't print any output, see log files for more details.

	// Output:
}
