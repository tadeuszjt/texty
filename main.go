package main

import (
	"bufio"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"os"
)

var (
	text TextWindow
)

func setup(w *gfx.Win) error {
	f, err := os.Open("text")
	if err != nil {
		return err
	}

    text = MakeTextWindow(w, geom.RectOrigin(400, 400))

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text.AddLine(scanner.Text())
	}

	return nil
}

func draw(w *gfx.WinCanvas) {
    text.Redraw(w)
	text.DrawOn(w)
}

func main() {
	gfx.RunWindow(gfx.WinConfig{
		SetupFunc: setup,
		DrawFunc:  draw,
		Title:     "Text",
		Resizable: true,
	})
}
