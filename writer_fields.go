package zlorderedconwriter

import (
	"bytes"
	"sort"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/rs/zerolog"
)

// writeFields appends formatted key-value pairs to buf.
func (w OrderedConsoleWriter) writeFields(evt *orderedmap.OrderedMap[string, any], buf *bytes.Buffer) {
	var fields = make([]string, 0, evt.Len())

	for el := evt.Front(); el != nil; el = el.Next() {
		field := el.Key

		var isExcluded bool
		for _, excluded := range w.FieldsExclude {
			if field == excluded {
				isExcluded = true
				break
			}
		}
		if isExcluded {
			continue
		}

		switch field {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName, zerolog.CallerFieldName:
			continue
		}
		fields = append(fields, field)
	}

	// Write space only if something has already been written to the buffer, and if there are fields.
	if buf.Len() > 0 && len(fields) > 0 {
		buf.WriteByte(' ')
	}

	// Move the "error" field to the front
	ei := sort.Search(len(fields), func(i int) bool { return fields[i] >= zerolog.ErrorFieldName })
	if ei < len(fields) && fields[ei] == zerolog.ErrorFieldName {
		fields[ei] = ""
		fields = append([]string{zerolog.ErrorFieldName}, fields...)
		var xfields = make([]string, 0, len(fields))
		for _, field := range fields {
			if field == "" { // Skip empty fields
				continue
			}
			xfields = append(xfields, field)
		}
		fields = xfields
	}

	for i, field := range fields {
		var fn zerolog.Formatter
		var fv zerolog.Formatter

		if field == zerolog.ErrorFieldName {
			if w.FormatErrFieldName == nil {
				fn = consoleDefaultFormatErrFieldName(w.NoColor)
			} else {
				fn = w.FormatErrFieldName
			}

			if w.FormatErrFieldValue == nil {
				fv = consoleDefaultFormatErrFieldValue(w.NoColor)
			} else {
				fv = w.FormatErrFieldValue
			}
		} else {
			if w.FormatFieldName == nil {
				fn = consoleDefaultFormatFieldName(w.NoColor)
			} else {
				fn = w.FormatFieldName
			}

			if w.FormatFieldValue == nil {
				fv = consoleDefaultFormatFieldValue
			} else {
				fv = w.FormatFieldValue
			}
		}

		buf.WriteString(fn(field))

		buf.WriteString(fv(evt.GetOrDefault(field, nil)))

		if i < len(fields)-1 { // Skip space for last field
			buf.WriteByte(' ')
		}
	}
}
