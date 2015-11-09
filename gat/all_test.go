package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/otiai10/mint"
)

func TestRun(t *testing.T) {
	o, e := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	run("../samples/red.png", o, e, 2, 3)

	b, err := ioutil.ReadAll(o)
	Expect(t, err).ToBe(nil)
	Expect(t, len(b)).ToBe(15)
}
