package dbfs

import (
	"gorm.io/gorm"
	"time"
)

const (
	LogServerFatal   = int8(1)
	LogServerError   = int8(2)
	LogServerWarning = int8(3)
	LogServerInfo    = int8(4)
	LogServerDebug   = int8(5)
	LogServerTrace   = int8(6)
)

type Log struct {
	LogID      int    `gorm:"primaryKey;autoIncrement"`
	LogType    int8   `gorm:"not null"`
	Message    string `gorm:"not null"`
	ServerName string
	TimeStamp  time.Time `gorm:"not null; autoUpdateTime"`
}

type Alert struct {
	LogID      int    `gorm:"primaryKey;"`
	LogType    int8   `gorm:"not null"`
	Message    string `gorm:"not null"`
	ServerName string
	TimeStamp  time.Time `gorm:"not null;"`
}

type Logger struct {
	db     *gorm.DB
	server string
}

func NewLogger(db *gorm.DB, server string) *Logger {
	return &Logger{
		db:     db,
		server: server,
	}
}

func (l *Logger) LogFatal(message string) {
	log := Log{
		LogType:    LogServerFatal,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
	alert := Alert{
		LogID:      log.LogID,
		LogType:    log.LogType,
		Message:    log.Message,
		ServerName: log.ServerName,
		TimeStamp:  log.TimeStamp,
	}
	l.db.Create(&alert)
}

func (l *Logger) LogError(message string) {
	log := Log{
		LogType:    LogServerError,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
	alert := Alert{
		LogID:      log.LogID,
		LogType:    log.LogType,
		Message:    log.Message,
		ServerName: log.ServerName,
		TimeStamp:  log.TimeStamp,
	}
	l.db.Create(&alert)
}

func (l *Logger) LogWarning(message string) {
	log := Log{
		LogType:    LogServerWarning,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
	alert := Alert{
		LogID:      log.LogID,
		LogType:    log.LogType,
		Message:    log.Message,
		ServerName: log.ServerName,
		TimeStamp:  log.TimeStamp,
	}
	l.db.Create(&alert)
}

func (l *Logger) LogInfo(message string) {
	log := Log{
		LogType:    LogServerInfo,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
}

func (l *Logger) LogDebug(message string) {
	log := Log{
		LogType:    LogServerDebug,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
}

func (l *Logger) LogTrace(message string) {
	log := Log{
		LogType:    LogServerTrace,
		Message:    message,
		ServerName: l.server,
		TimeStamp:  time.Now(),
	}
	l.db.Create(&log)
}
