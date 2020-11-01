package main

import (
	"bufio"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"github.com/go-gl/glfw/v3.2/glfw"
	"os"
)

var (
	text *TextWindow
)

func setup(w *gfx.Win) error {
    w.GetGlfwWindow().SetCharCallback(charCallback)


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

func charCallback(w *glfw.Window, r rune) {
    text.InsertChar(r)
}

func keyboard(w *gfx.Win, ev gfx.KeyEvent) {
    if ev.Action == glfw.Press || ev.Action == glfw.Repeat {
        switch ev.Key {
        case glfw.KeyEnter:
            text.Enter()
        case glfw.KeyBackspace:
            text.Backspace()
        case glfw.KeyTab:
            text.InsertChar(' ')
            text.InsertChar(' ')
            text.InsertChar(' ')
            text.InsertChar(' ')
        }
    }
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
        KeyFunc:   keyboard,
		Title:     "Text",
		Resizable: true,
	})
}
