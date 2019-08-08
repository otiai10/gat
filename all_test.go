package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/otiai10/mint"
)

func TestMain(m *testing.M) {
	os.Setenv("TERM_PROGRAM", "") // Force use CellGrid
	code := m.Run()
	os.Exit(code)
}

func TestRun(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	run("./samples/red.png", o, e, 2, 3)

	b, err := ioutil.ReadAll(o)
	Expect(t, err).ToBe(nil)
	Expect(t, len(b)).ToBe(15)

	Because(t, "gat can accept image URL via http/https", func(t *testing.T) {
		t.SkipNow()
		err := run("https://raw.githubusercontent.com/otiai10/gat/master/samples/sample.png", o, e, 2, 3)
		Expect(t, err).ToBe(nil)
		When(t, "given URL is not valid URL nor response is not image/*", func(t *testing.T) {
			err = run("foobaa://github.com/otiai10/gat", o, e, 2, 3)
			Expect(t, err).Not().ToBe(nil)
			err = run("https://github.com/otiai10/gat", o, e, 2, 3)
			Expect(t, err).Not().ToBe(nil)
		})
	})
}
