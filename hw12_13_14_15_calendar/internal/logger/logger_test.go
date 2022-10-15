package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("debug level", func(t *testing.T) {
		logg := New("debug", 0)
		logg.debugPrefix = "D"
		logg.infoPrefix = "I"
		logg.warnPrefix = "W"
		logg.errorPrefix = "E"

		var buf bytes.Buffer
		logg.SetOutput(&buf)

		logg.Debug("debug")
		logg.Info("info")
		logg.Warn("warn")
		logg.Error("error")

		require.Equal(t, "D debug\nI info\nW warn\nE error\n", buf.String())
	})

	t.Run("info level", func(t *testing.T) {
		logg := New("info", 0)
		logg.debugPrefix = "D"
		logg.infoPrefix = "I"
		logg.warnPrefix = "W"
		logg.errorPrefix = "E"

		var buf bytes.Buffer
		logg.SetOutput(&buf)

		logg.Debug("debug")
		logg.Info("info")
		logg.Warn("warn")
		logg.Error("error")

		require.Equal(t, "I info\nW warn\nE error\n", buf.String())
	})

	t.Run("warn level", func(t *testing.T) {
		logg := New("warn", 0)
		logg.debugPrefix = "D"
		logg.infoPrefix = "I"
		logg.warnPrefix = "W"
		logg.errorPrefix = "E"

		var buf bytes.Buffer
		logg.SetOutput(&buf)

		logg.Debug("debug")
		logg.Info("info")
		logg.Warn("warn")
		logg.Error("error")

		require.Equal(t, "W warn\nE error\n", buf.String())
	})

	t.Run("error level", func(t *testing.T) {
		logg := New("error", 0)
		logg.debugPrefix = "D"
		logg.infoPrefix = "I"
		logg.warnPrefix = "W"
		logg.errorPrefix = "E"

		var buf bytes.Buffer
		logg.SetOutput(&buf)

		logg.Debug("debug")
		logg.Info("info")
		logg.Warn("warn")
		logg.Error("error")

		require.Equal(t, "E error\n", buf.String())
	})

	t.Run("off level", func(t *testing.T) {
		logg := New("off", 0)
		logg.debugPrefix = "D"
		logg.infoPrefix = "I"
		logg.warnPrefix = "W"
		logg.errorPrefix = "E"

		var buf bytes.Buffer
		logg.SetOutput(&buf)

		logg.Debug("debug")
		logg.Info("info")
		logg.Warn("warn")
		logg.Error("error")

		require.Equal(t, "", buf.String())
	})
}
