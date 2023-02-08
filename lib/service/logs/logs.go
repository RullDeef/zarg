package logs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const LOGS_DIR = "./logs"

type Logger struct {
	file      *os.File
	AndStdout bool
	lock      sync.Mutex
}

func New() *Logger {
	filename := fmt.Sprintf("%s/%s.log", LOGS_DIR, time.Now().Format("2006-01-02_15-04-05"))
	log.Print(filename)

	if err := os.MkdirAll(LOGS_DIR, os.ModePerm); err != nil {
		log.Panic(err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}

	return &Logger{
		file:      file,
		AndStdout: true,
		lock:      sync.Mutex{},
	}
}

func (l *Logger) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.file.Close()
}

func (l *Logger) Printf(format string, args ...any) {
	l.lock.Lock()
	defer l.lock.Unlock()

	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.file, "%s: %s\n", time.Now().Format(time.StampMilli), strings.TrimRight(msg, "\n"))
	if l.AndStdout {
		log.Print(msg)
	}
}

func (l *Logger) Panicf(format string, args ...any) {
	l.lock.Lock()
	defer l.lock.Unlock()

	msg := fmt.Sprintf(format, args...)
	fmt.Print(l.file, msg)
	if l.AndStdout {
		log.Print(msg)
	}
	panic(msg)
}
