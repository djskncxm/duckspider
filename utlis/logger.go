package utlis

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

type Stats struct {
	mu           sync.Mutex
	OverallStats map[string]interface{}
}

func NewStats() *Stats {
	return &Stats{
		OverallStats: make(map[string]interface{}),
	}
}

func (stats *Stats) AddInt(key string, count interface{}) {
	stats.mu.Lock()
	defer stats.mu.Unlock()

	current, ok := stats.OverallStats[key].(int)
	if !ok {
		current = 0
	}

	stats.OverallStats[key] = current + count.(int)
}

func (stats *Stats) AddString(key string, count interface{}) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.OverallStats[key] = count.(string)
}

func (stats *Stats) Clear() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.OverallStats = make(map[string]interface{})
}

func (stats *Stats) OutTableInfo() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"统计项目", "信息"})

	for key, value := range stats.OverallStats {
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case int, int64, float64, float32:
			valueStr = fmt.Sprintf("%v", v)
		default:
			valueStr = "未知类型"
		}

		err := table.Append([]string{key, valueStr})
		if err != nil {
			return
		}
	}
	err := table.Render()
	if err != nil {
		return
	}
}

type Logger struct {
	name  string
	level logrus.Level
	mu    sync.RWMutex
	entry *logrus.Entry
	Stats *Stats
}

func InitLogger(name string, Level string) *Logger {
	baseLogger := logrus.New()
	level := ParseLogLevel(Level)

	baseLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "060102 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		PadLevelText:    true,
	})

	baseLogger.SetOutput(os.Stdout)
	baseLogger.SetLevel(level)

	entry := baseLogger.WithField("name", name)

	return &Logger{
		name:  name,
		level: level,
		entry: entry,
		Stats: NewStats(),
	}
}
func ParseLogLevel(levelStr string) logrus.Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG", "debug":
		return logrus.DebugLevel
	case "INFO", "info":
		return logrus.InfoLevel
	case "WARN", "WARNING", "warn":
		return logrus.WarnLevel
	case "ERROR", "error":
		return logrus.ErrorLevel
	case "FATAL", "fatal":
		return logrus.FatalLevel
	case "PANIC", "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *Logger) Info(msg string, err ...error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Info(msg, err)
}

func (l *Logger) Debug(msg string, err ...error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Debug(msg, err)
}

func (l *Logger) Warn(msg string, err ...error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Warn(msg, err)
}

func (l *Logger) Error(msg string, err ...error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Error(msg, err)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Infof(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Fatalf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.entry.Panicf(format, args...)
}
