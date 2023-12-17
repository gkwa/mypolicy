package mypolicy

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var replace func([]string, slog.Attr) slog.Attr

// HandlerType enum for handler types
type HandlerType string

const (
	TextHandler HandlerType = "text"
	JsonHandler HandlerType = "json"
)

func Main() int {
	slog.Debug("mypolicy", "test", true)

	var logger *slog.Logger

	// text handler
	logger = MakeLogger(os.Stderr, slog.LevelInfo, TextHandler)
	logger.Info("Hello World!")
	logger.Error("Hello World!")

	slog.Info("Hello World!")
	slog.Error("Hello World!")

	slog.SetDefault(logger)
	slog.Info("Hello World!")
	slog.Error("Hello World!")

	// json handler
	logger = MakeLogger(os.Stderr, slog.LevelInfo, JsonHandler)
	logger.Info("Hello World!")
	logger.Error("Hello World!")

	slog.Info("Hello World!")
	slog.Error("Hello World!")

	slog.SetDefault(logger)
	slog.Info("Hello World!")
	slog.Error("Hello World!")

	// defaults to text handler
	logger = MakeLogger(os.Stderr, slog.LevelInfo)
	logger.Info("Hello World!")
	logger.Error("Hello World!")

	slog.Info("Hello World!")
	slog.Error("Hello World!")

	slog.SetDefault(logger)
	slog.Info("Hello World!")
	slog.Error("Hello World!")

	return 0
}

func setPartialPath(source *slog.Source) {
	fileName := filepath.Base(source.File)
	parentDir := filepath.Base(filepath.Dir(source.File))

	source.File = filepath.Join(parentDir, fileName)
}

func init() {
	replace = func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		if a.Key == slog.SourceKey {
			source, _ := a.Value.Any().(*slog.Source)
			if source != nil {
				setPartialPath(source)
			}
		}
		return a
	}
}

func makeHandlerOptions(level slog.Level) slog.HandlerOptions {
	return slog.HandlerOptions{
		AddSource:   true,
		Level:       &level,
		ReplaceAttr: replace,
	}
}

func MakeLogger(w io.Writer, level slog.Level, handlerType ...HandlerType) *slog.Logger {
	opts := makeHandlerOptions(level)

	var handler slog.Handler
	if len(handlerType) == 0 {
		handler = slog.NewTextHandler(w, &opts)
	} else {
		switch handlerType[0] {
		case "text":
			handler = slog.NewTextHandler(w, &opts)
		case "json":
			handler = slog.NewJSONHandler(w, &opts)
		default:
			handler = slog.NewTextHandler(w, &opts)
		}
	}

	logger := slog.New(handler)
	return logger
}
