package zlorderedconwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oxtoacart/bpool"
	"github.com/rs/zerolog"
)

var (
	consoleBufPool = bpool.NewBufferPool(10)
)

const (
	consoleDefaultTimeFormat = time.RFC3339Nano
)

func consoleDefaultPartsOrder() []string {
	return []string{
		zerolog.TimestampFieldName,
		zerolog.LevelFieldName,
		zerolog.CallerFieldName,
		zerolog.MessageFieldName,
	}
}

func consoleDefaultFormatTimestamp(timeFormat string, noColor bool) zerolog.Formatter {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	return func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(zerolog.TimeFieldFormat, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Local().Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixMs:
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				case zerolog.TimeFormatUnixMicro:
					nsec = int64(time.Duration(i) * time.Microsecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec)
				t = ts.Format(timeFormat)
			}
		}
		return colorize(t, colorDarkGray, noColor)
	}
}

func consoleDefaultFormatLevel(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case zerolog.LevelTraceValue:
				l = colorize("TRC", colorMagenta, noColor)
			case zerolog.LevelDebugValue:
				l = colorize("DBG", colorYellow, noColor)
			case zerolog.LevelInfoValue:
				l = colorize("INF", colorGreen, noColor)
			case zerolog.LevelWarnValue:
				l = colorize("WRN", colorRed, noColor)
			case zerolog.LevelErrorValue:
				l = colorize(colorize("ERR", colorRed, noColor), colorBold, noColor)
			case zerolog.LevelFatalValue:
				l = colorize(colorize("FTL", colorRed, noColor), colorBold, noColor)
			case zerolog.LevelPanicValue:
				l = colorize(colorize("PNC", colorRed, noColor), colorBold, noColor)
			default:
				l = colorize("???", colorBold, noColor)
			}
		} else {
			if i == nil {
				l = colorize("???", colorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func consoleDefaultFormatCaller(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
			c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
		}
		return c
	}
}

func consoleDefaultFormatMessage(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatFieldName(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorCyan, noColor)
	}
}

func consoleDefaultFormatFieldValue(i interface{}) string {
	return fmt.Sprintf("%s", i)
}

func consoleDefaultFormatErrFieldName(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s=", i), colorCyan, noColor)
	}
}

func consoleDefaultFormatErrFieldValue(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		return colorize(fmt.Sprintf("%s", i), colorRed, noColor)
	}
}
