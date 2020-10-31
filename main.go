package main

import (
	"bufio"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"os"
)

var (
	text *TextWindow
)

func setup(w *gfx.Win) error {
	f, err := os.Open("text")
	if err != nil {
		return err
	}

    text = NewTextWindow(w, geom.RectOrigin(400, 400))
    text.Resize(w, geom.RectOrigin(600, 600))

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text.AddLine(scanner.Text())
	}

	return nil
}

func resize(w *gfx.Win, width, height int) {
    text.Resize(w, geom.RectOrigin(float32(width), float32(height)))
}

func draw(w *gfx.Win, c gfx.Canvas) {
    text.Redraw(w)
	text.DrawOn(c, geom.Vec2{})
}

func main() {
	gfx.RunWindow(gfx.WinConfig{
		SetupFunc: setup,
		DrawFunc:  draw,
        ResizeFunc: resize,
		Title:     "Text",
		Resizable: true,
	})
}
