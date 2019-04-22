package render_test

import (
	"image"
	"os"

	"github.com/otiai10/gat/render"
)

func ExampleRenderer_Render() {

	f, err := os.Open("./samples/sample.png")
	if err != nil {
		panic(err)
	}

	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	r := render.GetDefaultRenderer()
	if err := r.Render(os.Stdout, i); err != nil {
		panic(err)
	}

}
