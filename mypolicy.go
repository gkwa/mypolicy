package mypolicy

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func Main() int {
	slog.Debug("mypolicy", "test", true)

	logger := makeLog(os.Stderr, slog.LevelInfo)

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

func makeLog(w io.Writer, level slog.Level) *slog.Logger {
	var replace func([]string, slog.Attr) slog.Attr

	replace = func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "source" {
			src := a.Value.Any().(*slog.Source)
			return slog.String("source", src.File)
		}
		if a.Key == "time" {
			t := a.Value.Time()
			return slog.String("time", t.Format(time.Kitchen))
		}
		return a
	}

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

	opts := slog.HandlerOptions{
		AddSource:   true,
		Level:       &level,
		ReplaceAttr: replace,
	}

	handler := slog.NewTextHandler(os.Stderr, &opts)
	logger := slog.New(handler)
	return logger
}
