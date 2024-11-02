package logger

import (
	"log"
	"os"

	"github.com/TinyWAF/TinyWAF/internal"
)

var info *log.Logger
var debug *log.Logger
var warn *log.Logger
var block *log.Logger
var error *log.Logger
var loadedCfg *internal.MainConfig

func Init(cfg *internal.MainConfig) {
	loadedCfg = cfg

	outFile := os.Stdout
	errOutFile := os.Stderr

	if cfg.Log.File != "" {
		f, err := os.OpenFile(cfg.Log.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Error creating or opening log file '%v': %v", cfg.Log.File, err.Error())
			return
		}
		defer f.Close()

		outFile = f
		errOutFile = outFile
	}

	info = log.New(outFile, "[ SYSTEM  ] ", log.LstdFlags|log.Lmsgprefix)
	debug = log.New(outFile, "[ DEBUG   ] ", log.LstdFlags|log.Lmsgprefix)
	warn = log.New(outFile, "[ WARN    ] ", log.LstdFlags|log.Lmsgprefix)
	block = log.New(outFile, "[ BLOCKED ] ", log.LstdFlags|log.Lmsgprefix)
	error = log.New(errOutFile, "[ ERROR   ] ", log.LstdFlags|log.Lmsgprefix)
}

func Info(format string, v ...any) {
	info.Printf(format, v...)
}

func Debug(format string, v ...any) {
	if loadedCfg.Log.Levels.Debug {
		debug.Printf(format, v...)
	}
}

func Warn(format string, v ...any) {
	if loadedCfg.Log.Levels.Warn {
		warn.Printf(format, v...)
	}
}

func Block(format string, v ...any) {
	if loadedCfg.Log.Levels.Block {
		block.Printf(format, v...)
	}
}

func Error(format string, v ...any) {
	error.Printf(format, v...)
}

func Fatal(format string, v ...any) {
	error.Fatalf(format, v...)
}
