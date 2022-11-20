package zlorderedconwriter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/buger/jsonparser"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
)

type OrderedConsoleWriter struct{ zcw zerolog.ConsoleWriter }

func (w OrderedConsoleWriter) Write(p []byte) (n int, err error) {
	// Fix color on Windows
	var buf = consoleBufPool.Get()
	defer consoleBufPool.Put(buf)

	evt := orderedmap.NewOrderedMap[string, any]()
	err = jsonparser.ObjectEach(p, func(k []byte, v []byte, t jsonparser.ValueType, pos int) error {
		switch t {
		case jsonparser.String:
			sv, err := jsonparser.ParseString(v)
			if err != nil {
				return err
			}
			evt.Set(string(k), sv)
		case jsonparser.Number:
			evt.Set(string(k), json.Number(v))
		default:
			evt.Set(string(k), v)
		}
		return nil
	})
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}

	for _, p := range w.zcw.PartsOrder {
		w.writePart(buf, evt, p)
	}

	w.writeFields(evt, buf)

	err = buf.WriteByte('\n')
	if err != nil {
		return n, err
	}

	_, err = buf.WriteTo(w.zcw.Out)
	return len(p), err
}

func New(w zerolog.ConsoleWriter) OrderedConsoleWriter {
	if w.Out == os.Stdout || w.Out == os.Stderr {
		w.Out = colorable.NewColorable(w.Out.(*os.File))
	}

	if w.PartsOrder == nil {
		w.PartsOrder = consoleDefaultPartsOrder()
	}
	return OrderedConsoleWriter{zcw: w}
}
