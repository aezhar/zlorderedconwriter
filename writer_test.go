package zlorderedconwriter_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/aezhar/zlorderedconwriter"
)

func TestConsoleWriter(t *testing.T) {
	t.Run("Default field formatter", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{
			Out: buf, NoColor: true, PartsOrder: []string{"foo"}})

		_, err := w.Write([]byte(`{"foo": "DEFAULT"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "DEFAULT foo=DEFAULT\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write colorized", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: false})

		_, err := w.Write([]byte(`{"level": "warn", "message": "Foobar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "\x1b[90m<nil>\x1b[0m \x1b[31mWRN\x1b[0m Foobar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		ts := time.Unix(0, 0)
		d := ts.UTC().Format(time.RFC3339)
		_, err := w.Write([]byte(`{"time": "` + d + `", "level": "debug", "message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := ts.Format(time.RFC3339Nano) + " DBG Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Unix timestamp input format", func(t *testing.T) {
		of := zerolog.TimeFieldFormat
		defer func() {
			zerolog.TimeFieldFormat = of
		}()
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, TimeFormat: time.StampMilli, NoColor: true})

		_, err := w.Write([]byte(`{"time": 1234, "level": "debug", "message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := time.Unix(1234, 0).Format(time.StampMilli) + " DBG Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Unix timestamp ms input format", func(t *testing.T) {
		of := zerolog.TimeFieldFormat
		defer func() {
			zerolog.TimeFieldFormat = of
		}()
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, TimeFormat: time.StampMilli, NoColor: true})

		_, err := w.Write([]byte(`{"time": 1234567, "level": "debug", "message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := time.Unix(1234, 567000000).Format(time.StampMilli) + " DBG Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Unix timestamp us input format", func(t *testing.T) {
		of := zerolog.TimeFieldFormat
		defer func() {
			zerolog.TimeFieldFormat = of
		}()
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro

		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, TimeFormat: time.StampMicro, NoColor: true})

		_, err := w.Write([]byte(`{"time": 1234567891, "level": "debug", "message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := time.Unix(1234, 567891000).Format(time.StampMicro) + " DBG Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("No message field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		_, err := w.Write([]byte(`{"level": "debug", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "<nil> DBG foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("No level field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		_, err := w.Write([]byte(`{"message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "<nil> ??? Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write colorized fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: false})

		_, err := w.Write([]byte(`{"level": "warn", "message": "Foobar", "foo": "bar"}`))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "\x1b[90m<nil>\x1b[0m \x1b[31mWRN\x1b[0m Foobar \x1b[36mfoo=\x1b[0mbar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write error field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		ts := time.Unix(0, 0)
		d := ts.UTC().Format(time.RFC3339)
		evt := `{"time": "` + d + `", "level": "error", "message": "Foobar", "aaa": "bbb", "error": "Error"}`
		// t.Log(evt)

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := ts.Format(time.RFC3339Nano) + " ERR Foobar error=Error aaa=bbb\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write caller field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Cannot get working directory: %s", err)
		}

		ts := time.Unix(0, 0)
		d := ts.UTC().Format(time.RFC3339)
		evt := `{"time": "` + d + `", "level": "debug", "message": "Foobar", "foo": "bar", "caller": "` + cwd + `/foo/bar.go"}`
		// t.Log(evt)

		_, err = w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := ts.Format(time.RFC3339Nano) + " DBG foo/bar.go > Foobar foo=bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write JSON field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		evt := `{"level": "debug", "message": "Foobar", "foo": [1,2,3], "bar": true}`
		// t.Log(evt)

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "<nil> DBG Foobar foo=[1,2,3] bar=true\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write quoted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		evt := `{"level": "debug", "message": "Foobar", "foo": "baa baz"}`

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := `<nil> DBG Foobar foo="baa baz"` + "\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})

	t.Run("Write quoted field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		w := zlorderedconwriter.New(zerolog.ConsoleWriter{Out: buf, NoColor: true})

		evt := `{"level": "debug", "message": "Foo bar"}`

		_, err := w.Write([]byte(evt))
		if err != nil {
			t.Errorf("Unexpected error when writing output: %s", err)
		}

		expectedOutput := "<nil> DBG Foo bar\n"
		actualOutput := buf.String()
		if actualOutput != expectedOutput {
			t.Errorf("Unexpected output %q, want: %q", actualOutput, expectedOutput)
		}
	})
}
