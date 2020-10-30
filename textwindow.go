package main

import (
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"strings"
)

type TextWindow struct {
	rect       geom.Rect
	tabSize    int
	texts      []gfx.Text
    windowTex  gfx.TexID
}

func MakeTextWindow(w *gfx.Win, rect geom.Rect) TextWindow {
	return TextWindow{
		rect:      rect,
		tabSize:   4,
        windowTex: w.LoadTextureBlank(int(rect.Width()), int(rect.Height())),
	}
}

func (t *TextWindow) Redraw(w *gfx.WinCanvas) {
    texCanvas := w.GetTextureCanvas(t.windowTex)
	for i := range t.texts {
		pos := geom.Vec2{t.rect.Min.X, float32(i * t.texts[i].Size())}
		gfx.DrawText(texCanvas, &t.texts[i], pos)
	}
}

func (t *TextWindow) DrawOn(c gfx.Canvas) {
    gfx.DrawSprite(
        c,
        geom.Ori2{},
        t.rect,
        gfx.White,
        nil,
        &t.windowTex,
    )
}

func (t *TextWindow) AddLine(line string) {
	str := ""
	for _, r := range line {
		switch {
		case r == '\t':
			str += strings.Repeat(" ", t.tabSize)
		default:
			str += string(r)
		}
	}

	t.texts = append(t.texts, gfx.Text{})
	t.texts[len(t.texts)-1].SetString(str)
}
