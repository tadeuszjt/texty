package main

import (
	"bufio"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"os"
)

var (
	text *TextWindow
    cmd  *TextWindow
    currentWindow *TextWindow
)

func setup(w *gfx.Win) error {
	w.GetGlfwWindow().SetCharCallback(charCallback)

	f, err := os.Open("text")
	if err != nil {
		return err
	}

	text = NewTextWindow(w, geom.RectOrigin(400, 400))
    currentWindow = text

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text.AddLine(scanner.Text())
	}

	return nil
}

func charCallback(w *glfw.Window, r rune) {
    switch r {
    case ';':
    default:
        TextBufferInsertChar(currentWindow, r)
    }
}

func keyboard(w *gfx.Win, ev gfx.KeyEvent) {
	if ev.Action == glfw.Press || ev.Action == glfw.Repeat {
		switch ev.Key {
		case glfw.KeyEnter:
            TextBufferEnter(currentWindow)
		case glfw.KeyBackspace:
            TextBufferBackspace(currentWindow)
		case glfw.KeyTab:
            TextBufferInsertChar(currentWindow, ' ')
            TextBufferInsertChar(currentWindow, ' ')
            TextBufferInsertChar(currentWindow, ' ')
            TextBufferInsertChar(currentWindow, ' ')
		case glfw.KeyUp:
            TextBufferMoveCursor(currentWindow, -1, 0)
		case glfw.KeyDown:
            TextBufferMoveCursor(currentWindow, 1, 0)
		case glfw.KeyLeft:
            TextBufferMoveCursor(currentWindow, 0, -1)
		case glfw.KeyRight:
            TextBufferMoveCursor(currentWindow, 0, 1)
        case glfw.KeySemicolon:
            rect := w.GetFrameRect()
            if cmd == nil {
                cmd = NewTextWindow(w, geom.MakeRect(0, rect.Height() - 100., rect.Width(), 100.))
                resize(w, int(rect.Width()), int(rect.Height()))
                currentWindow = cmd
            } else {
                cmd.Free(w)
                cmd = nil
                resize(w, int(rect.Width()), int(rect.Height()))
                currentWindow = text
            }
		}
	}
}

func resize(w *gfx.Win, width, height int) {
	text.Resize(w, geom.RectOrigin(float32(width), float32(height)))

    if cmd != nil {
        text.Resize(w, geom.MakeRect(
            0,
            0,
            float32(width),
            float32(height) - 100.,
        ))

        cmd.Resize(w, geom.MakeRect(
            0,
            float32(height) - 100.,
            float32(width),
            100.,
        ))

        return
    }

    text.Resize(w, geom.MakeRect(
        0,
        0,
        float32(width),
        float32(height),
    ))
}

func draw(w *gfx.Win, c gfx.Canvas) {
    text.Redraw(w)
    text.DrawOn(c)

    if cmd != nil {
        cmd.Redraw(w)
        cmd.DrawOn(c)
    }
}

func main() {
	gfx.RunWindow(gfx.WinConfig{
		SetupFunc:  setup,
		DrawFunc:   draw,
		ResizeFunc: resize,
		KeyFunc:    keyboard,
		Title:      "Text",
		Resizable:  true,
	})
}
