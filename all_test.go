package gat

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	_ "image/png"

	. "github.com/otiai10/mint"
)

func TestFprint(t *testing.T) {
	out := bytes.NewBuffer(nil)
	Fprint(out, 12, "foo")
	b, err := ioutil.ReadAll(out)
	Expect(t, err).ToBe(nil)
	Expect(t, b).ToBe([]byte("\x1b[48;5;12mfoo\x1b[m"))
}

func TestNewClient(t *testing.T) {
	client := NewClient(Rect{100, 200})
	Expect(t, client).TypeOf("*gat.Client")
	Expect(t, client.Canvas).TypeOf("gat.Rect")
	Expect(t, client.Canvas.Row).ToBe(uint16(100))
	Expect(t, client.Canvas.Col).ToBe(uint16(200))
	Expect(t, client.Border).TypeOf("gat.DefaultBorder")
	Expect(t, client.Out).TypeOf("*os.File")
}

func TestClient_Set(t *testing.T) {
	client := NewClient(Rect{123, 456})
	Expect(t, client.Canvas.Row).ToBe(uint16(123))
	Expect(t, client.Canvas.Col).ToBe(uint16(456))
	Expect(t, client.Border).ToBe(DefaultBorder{})

	client.Set(Rect{17, 19})
	Expect(t, client.Canvas.Row).ToBe(uint16(17))
	Expect(t, client.Canvas.Col).ToBe(uint16(19))

	client.Set(SimpleBorder{})
	Expect(t, client.Border).ToBe(SimpleBorder{})
}

func TestClient_PrintImage(t *testing.T) {
	client := NewClient(Rect{Row: 2, Col: 3})
	out := bytes.NewBuffer(nil)
	client.Out = out
	Expect(t, out.Len()).ToBe(0)

	img := getImage("red.png")
	err := client.PrintImage(img)
	Expect(t, err).ToBe(nil)
	Expect(t, out.Len()).Not().ToBe(0)
}

func getImage(filename string) image.Image {
	f, err := os.Open(filepath.Join("samples", filename))
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return img
}
