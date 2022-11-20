package zlorderedconwriter

import (
	"bytes"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/rs/zerolog"
)

// writePart appends a formatted part to buf.
func (w OrderedConsoleWriter) writePart(buf *bytes.Buffer, evt *orderedmap.OrderedMap[string, any], p string) {
	var f zerolog.Formatter

	if w.zcw.PartsExclude != nil && len(w.zcw.PartsExclude) > 0 {
		for _, exclude := range w.zcw.PartsExclude {
			if exclude == p {
				return
			}
		}
	}

	switch p {
	case zerolog.LevelFieldName:
		if w.zcw.FormatLevel == nil {
			f = consoleDefaultFormatLevel(w.zcw.NoColor)
		} else {
			f = w.zcw.FormatLevel
		}
	case zerolog.TimestampFieldName:
		if w.zcw.FormatTimestamp == nil {
			f = consoleDefaultFormatTimestamp(w.zcw.TimeFormat, w.zcw.NoColor)
		} else {
			f = w.zcw.FormatTimestamp
		}
	case zerolog.MessageFieldName:
		if w.zcw.FormatMessage == nil {
			f = consoleDefaultFormatMessage
		} else {
			f = w.zcw.FormatMessage
		}
	case zerolog.CallerFieldName:
		if w.zcw.FormatCaller == nil {
			f = consoleDefaultFormatCaller(w.zcw.NoColor)
		} else {
			f = w.zcw.FormatCaller
		}
	default:
		if w.zcw.FormatFieldValue == nil {
			f = consoleDefaultFormatFieldValue
		} else {
			f = w.zcw.FormatFieldValue
		}
	}

	if s := f(evt.GetOrDefault(p, nil)); len(s) > 0 {
		if buf.Len() > 0 {
			buf.WriteByte(' ') // Write space only if not the first part
		}
		buf.WriteString(s)
	}
}
