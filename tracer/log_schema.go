package tracer

import _ "github.com/jinzhu/gorm"

type LogEntry struct {
  Id uint64
  TerminalId uint64 `sql:"not null"`
  Timestamp string `sql:"not null"`
  Message string `sql:"not null"`
}
