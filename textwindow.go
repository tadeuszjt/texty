package main

import (
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"strings"
)

type TextWindow struct {
	rect       geom.Rect
	tabSize    int
    textSize   int
    charWidth  float64
	texts      []gfx.Text
    windowTex  gfx.TexID
    cursor     struct{ line, char int }
}

func NewTextWindow(w *gfx.Win, rect geom.Rect) *TextWindow {
	return &TextWindow{
		rect:      rect,
		tabSize:   4,
        textSize:  14,
        charWidth: 10.,
        windowTex: w.LoadTextureBlank(int(rect.Width()), int(rect.Height())),
	}
}

func (t *TextWindow) Redraw(w *gfx.Win) {
    texCanvas := w.GetTextureCanvas(t.windowTex)
	for i := range t.texts {
		pos := geom.Vec2{t.rect.Min.X, float32(i * t.textSize)}
		gfx.DrawText(texCanvas, &t.texts[i], pos)
	}

    gfx.DrawSprite(
        texCanvas, 
        geom.Ori2{float32(t.charWidth * float64(t.cursor.char)), float32(t.textSize*t.cursor.line), 0},
        geom.RectOrigin(float32(t.charWidth), float32(t.textSize)),
        gfx.Colour{0, 1, 0, 0.4},
        nil,
        nil,
    )
}

func (t *TextWindow) Resize(w *gfx.Win, rect geom.Rect) {
    t.rect = rect
    w.FreeTexture(t.windowTex)
    t.windowTex = w.LoadTextureBlank(int(rect.Width()), int(rect.Height()))
}

func (t *TextWindow) DrawOn(c gfx.Canvas, pos geom.Vec2) {
    gfx.DrawSprite(c, pos.Ori2(), t.rect, gfx.White, nil, &t.windowTex)
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
	t.texts[len(t.texts)-1].SetSize(t.textSize)
	t.texts[len(t.texts)-1].SetString(str)
    t.charWidth = t.texts[0].CharWidth()
}
