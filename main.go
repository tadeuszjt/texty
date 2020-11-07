package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
    "os"
    "fmt"
    "bufio"
    "strings"
)

var (
	text *TextWindow
    cmd  *TextWindow
    currentWindow *TextWindow
)

func setup(w *gfx.Win) error {
	w.GetGlfwWindow().SetCharCallback(charCallback)
	text = NewTextWindow(w, w.GetFrameRect())
    currentWindow = text
	return nil
}

func charCallback(w *glfw.Window, r rune) {
    switch r {
    case ';':
    default:
        TextBufferInsertChar(currentWindow, r)
    }
}

func cmdStr(str string) error {
    words := strings.Split(str, " ")
    if len(words) <= 0 {
        return fmt.Errorf("invalid cmd")
    }

    switch words[0] {
    case "save": {
            if len(words) > 2 {
                return fmt.Errorf("invalid cmd")
            }

            f, err := os.Create(words[1])
            if err != nil {
                return err
            }

            _, err = f.WriteString(TextBufferString(text))
            if err != nil {
                return err
            }
        }
    case "open": {
            if len(words) > 2 {
                return fmt.Errorf("invalid cmd")
            }
            f, err := os.Open(words[1])
            if err != nil {
                return err
            }

            TextBufferClear(text)
            scanner := bufio.NewScanner(f)
            for scanner.Scan() {
                TextBufferAppend(text, scanner.Text())
            }
        }
    
    }


    return nil
}

func keyboard(w *gfx.Win, ev gfx.KeyEvent) {
	if ev.Action == glfw.Press || ev.Action == glfw.Repeat {
		switch ev.Key {
		case glfw.KeyEnter:
            if currentWindow == cmd {
                str := cmd.GetLine(0)
                err := cmdStr(str)
                if err != nil {
                    fmt.Println(err)
                }

                cmd.Free(w)
                cmd = nil 
                currentWindow = text
                resize(w)
            } else {
                TextBufferEnter(currentWindow)
            }
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
                resize(w)
                currentWindow = cmd
            } else {
                cmd.Free(w)
                cmd = nil
                resize(w)
                currentWindow = text
            }
		}
	}
}

func resize(w *gfx.Win) {
    rect := w.GetFrameRect()

    if cmd != nil {
        text.Resize(w, geom.MakeRect(0, 0, rect.Width(), rect.Height() - 100.))
        cmd.Resize(w, geom.MakeRect(0, rect.Height() - 100., rect.Width(), 100.))
        return
    }

    text.Resize(w, rect)
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
