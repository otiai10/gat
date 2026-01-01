package main

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"testing"

	"github.com/otiai10/gat/render"

	. "github.com/otiai10/mint"
)

func TestMain(m *testing.M) {
	os.Setenv("TERM_PROGRAM", "") // Force use CellGrid
	code := m.Run()
	os.Exit(code)
}

func TestRun(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	run([]string{"./samples/red.png"}, o, e, 2, 3)

	b, err := ioutil.ReadAll(o)
	Expect(t, err).ToBe(nil)
	Expect(t, len(b)).ToBe(15)

	Because(t, "gat can accept image URL via http/https", func(t *testing.T) {
		t.SkipNow()
		err := run([]string{"https://raw.githubusercontent.com/otiai10/gat/master/samples/sample.png"}, o, e, 2, 3)
		Expect(t, err).ToBe(nil)
		When(t, "given URL is not valid URL nor response is not image/*", func(t *testing.T) {
			err = run([]string{"foobaa://github.com/otiai10/gat"}, o, e, 2, 3)
			Expect(t, err).Not().ToBe(nil)
			err = run([]string{"https://github.com/otiai10/gat"}, o, e, 2, 3)
			Expect(t, err).Not().ToBe(nil)
		})
	})
}

func TestRun_NoFiles(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	err := run([]string{}, o, e, 2, 3)
	Expect(t, err).ToBe(nil)
	// Should print "No files" message to stderr
	Expect(t, e.String()).ToBe("No files \n")
}

func TestRun_DebugMode(t *testing.T) {
	// Save and restore debug flag
	origDebug := debug
	defer func() { debug = origDebug }()

	debug = true
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	err := run([]string{}, o, e, 2, 3)
	Expect(t, err).ToBe(nil)
	// Debug mode with no files should dump colors to stderr
	Expect(t, len(e.String()) > 0).ToBe(true)
}

func TestRun_MultipleFiles(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	err := run([]string{"./samples/red.png", "./samples/red.png"}, o, e, 2, 3)
	Expect(t, err).ToBe(nil)
	// Should have output for both files
	Expect(t, len(o.String()) > 15).ToBe(true)
}

func TestRun_InvalidFile(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	err := run([]string{"./nonexistent.png"}, o, e, 2, 3)
	Expect(t, err).Not().ToBe(nil)
}

func TestRun_NotAnImage(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	err := run([]string{"./main.go"}, o, e, 2, 3)
	Expect(t, err).Not().ToBe(nil) // Should fail to decode as image
}

func TestGetInputReader_LocalFile(t *testing.T) {
	rc, err := getInputReader("./samples/red.png")
	Expect(t, err).ToBe(nil)
	Expect(t, rc).Not().ToBe(nil)
	rc.Close()
}

func TestGetInputReader_NonexistentFile(t *testing.T) {
	rc, err := getInputReader("./nonexistent.png")
	Expect(t, err).Not().ToBe(nil)
	Expect(t, rc).ToBe(nil)
}

func TestGetRenderer_UseCellTrue(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	r := getRenderer(true, 10, 20, "  ", 1.0, img)
	_, ok := r.(*render.CellGrid)
	Expect(t, ok).ToBe(true)
}

func TestGetRenderer_UseCellFalse_NoITerm(t *testing.T) {
	// TERM_PROGRAM is already set to "" in TestMain
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	r := getRenderer(false, 10, 20, "  ", 1.0, img)
	_, ok := r.(*render.CellGrid)
	Expect(t, ok).ToBe(true)
}

func TestGetRenderer_WithITerm(t *testing.T) {
	origTerm := os.Getenv("TERM_PROGRAM")
	defer os.Setenv("TERM_PROGRAM", origTerm)

	os.Setenv("TERM_PROGRAM", "iTerm.app")
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	r := getRenderer(false, 10, 20, "  ", 2.0, img)
	_, ok := r.(*render.ITerm)
	Expect(t, ok).ToBe(true)
}

func TestClearTerminal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	clearTerminal(buf)
	Expect(t, buf.String()).ToBe("\033c")
}
