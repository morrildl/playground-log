/* Copyright © Playground Global, LLC. All rights reserved. */

package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"playground"
)

type LogLevel int

const (
	LEVEL_ERROR LogLevel = iota
	LEVEL_WARNING
	LEVEL_STATUS
	LEVEL_DEBUG
)

var currentLevel LogLevel = LEVEL_STATUS
var quietLog = false
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func SetLogLevel(newLevel LogLevel) {
	_, ok := levelMap[newLevel]
	if !ok {
		Warn("Logger", "someone tried to set invalid log level ", newLevel)
		return
	}
	currentLevel = newLevel
}

func SetQuiet(isQuiet bool) {
	quietLog = isQuiet
	if isQuiet {
		logger = log.New(os.Stdout, "", 0)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}
}

func SetLogFile(fileName string) {
	var err error
	if fileName, err = filepath.Abs(fileName); err != nil {
		msg := "-log value '"+fileName+"' does not resolve"
		Error("log.SetLogFile", msg, err)
		panic(msg)
	}
	if stat, err := os.Stat(fileName); (err != nil && !os.IsNotExist(err)) || (stat != nil && stat.IsDir()) {
    msg := "-log value '"+fileName+"' does not stat or is a directory"
		Error("log.SetLogFile", msg, err)
		panic(msg)
	}
	if f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660); err == nil {
		fmt.Println("Directing log to " + fileName + ".")
		logger = log.New(f, "", log.LstdFlags)
	} else {
		Warn("Logger", "failed to open log file ", fileName)
	}
}

var levelMap map[LogLevel]string = map[LogLevel]string{
	LEVEL_ERROR:   "ERROR",
	LEVEL_WARNING: "WARNING",
	LEVEL_STATUS:  "STATUS",
	LEVEL_DEBUG:   "DEBUG",
}

func doLog(level LogLevel, component string, extras ...interface{}) {
	if level > currentLevel {
		return
	}

	levelString, ok := levelMap[level]
	if !ok {
		levelString = "ERROR"
		Warn("Logger", "called with invalid level ", level)
	}
	var message string
	if _, ok := extras[0].(string); ok {
		if quietLog {
			if level < LEVEL_STATUS {
				message = fmt.Sprintf("%s %s ", levelString, extras[0])
			} else {
				message = fmt.Sprintf("%s ", extras[0])
			}
		} else {
			message = fmt.Sprintf("[%s] (%s) %s ", levelString, component, extras[0])
		}
		extras = extras[1:]
	} else {
		if quietLog {
			if level < LEVEL_STATUS {
				message = fmt.Sprintf("%s ", levelString)
			} else {
				message = fmt.Sprintf(" ")
			}
		} else {
			message = fmt.Sprintf("[%s] (%s) ", levelString, component)
		}
	}
	if len(extras) > 0 {
		extras = append([]interface{}{message}, extras)
	} else {
		extras = []interface{}{message}
	}
	logger.Print(extras...)
}

func Debug(component string, extras ...interface{}) {
	doLog(LEVEL_DEBUG, component, extras...)
}

func Error(component string, extras ...interface{}) {
	doLog(LEVEL_ERROR, component, extras...)
}

func Warn(component string, extras ...interface{}) {
	doLog(LEVEL_WARNING, component, extras...)
}

func Status(component string, extras ...interface{}) {
	doLog(LEVEL_STATUS, component, extras...)
}
