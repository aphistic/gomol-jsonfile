package jsonstream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aphistic/gomol"
)

type Logger struct {
	base        *gomol.Base
	cfg         *Config
	initialized bool

	wc  io.WriteCloser
	buf *bufio.Writer
}

func NewLogger(cfg *Config) (*Logger, error) {
	return &Logger{
		cfg: cfg,
	}, nil
}

func (l *Logger) SetBase(base *gomol.Base) {
	l.base = base
}

func (l *Logger) InitLogger() error {
	if l.initialized {
		return nil
	}

	if l.cfg.w == nil {
		return fmt.Errorf("a valid io.Writer must be supplied")
	}

	l.buf = bufio.NewWriter(l.cfg.w)

	if len(l.cfg.headers) > 0 {
		// After initializing the buffered writer, write out the headers
		// and flush them.
		headers := &outputMsg{
			Type: "header",
			Data: &HeaderMsg{
				Timestamp: l.cfg.clock.Now(),
				Headers:   l.cfg.headers,
			},
		}

		data, err := json.Marshal(headers)
		if err != nil {
			return err
		}
		data = append(data, '\n')

		_, err = l.buf.Write(data)
		if err != nil {
			return err
		}
		l.buf.Flush()
	}

	l.initialized = true

	return nil
}

func (l *Logger) ShutdownLogger() error {
	if l.wc != nil {
		// Make sure everything is flushed before we shut down
		if l.buf != nil {
			err := l.buf.Flush()
			if err != nil {
				return err
			}
			l.buf = nil
		}

		err := l.wc.Close()
		if err != nil {
			return err
		}
		l.wc = nil
	}

	l.initialized = false

	return nil
}

func (l *Logger) IsInitialized() bool {
	return l.initialized
}

func (l *Logger) Flush() error {
	return l.buf.Flush()
}

func (l *Logger) Logm(timestamp time.Time, level gomol.LogLevel, attrs map[string]interface{}, msg string) error {
	if !l.IsInitialized() {
		return fmt.Errorf("logger is uninitialized")
	}

	log := &outputMsg{
		Type: "log",
		Data: &LogMsg{
			Timestamp: timestamp,
			Level:     level,
			Msg:       msg,
			Attrs:     attrs,
		},
	}

	data, err := json.Marshal(log)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	_, err = l.buf.Write(data)
	if err != nil {
		return err
	}

	return nil
}
