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
    cursorLine int
    cursorChar int
	lines      []gfx.Text
    windowTex  gfx.TexID
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

func (t *TextWindow) curLine() *gfx.Text {
    return &t.lines[t.cursorLine]
}

func (t *TextWindow) InsertLine(lindIdx int) {
    text := gfx.Text{}
    text.SetSize(t.textSize)

    t.lines = append(t.lines, text)
    if lindIdx == len(t.lines)-1 {
        return
    }

    copy(t.lines[lindIdx+1:], t.lines[lindIdx:])
    t.lines[lindIdx] = text
}

func (t *TextWindow) Enter() {
    line := t.cursorLine
    char := t.cursorChar

    str := t.curLine().GetString()
    t.curLine().SetString(str[:char])

    t.InsertLine(line + 1)
    t.lines[line+1].SetString(str[char:])

    t.cursorLine++
    t.cursorChar = 0
}

func (t *TextWindow) Backspace() {
    if t.cursorChar == 0 {
        return
    }

    str := t.curLine().GetString()
    t.curLine().SetString(str[:t.cursorChar-1] + str[t.cursorChar:])
    t.cursorChar--
}

func (t *TextWindow) InsertChar(r rune) {
    str := t.curLine().GetString()
    t.curLine().SetString(str[:t.cursorChar] + string(r) + str[t.cursorChar:])
    t.cursorChar++
}

func (t *TextWindow) Redraw(w *gfx.Win) {
    texCanvas := w.GetTextureCanvas(t.windowTex)
    texCanvas.Clear(gfx.White)
	for i := range t.lines {
		pos := geom.Vec2{t.rect.Min.X, float32(i * t.textSize)}
		gfx.DrawText(texCanvas, &t.lines[i], pos)
	}

    gfx.DrawSprite(
        texCanvas, 
        geom.Ori2{float32(t.charWidth * float64(t.cursorChar)), float32(t.textSize*t.cursorLine), 0},
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

	t.lines = append(t.lines, gfx.Text{})
	t.lines[len(t.lines)-1].SetSize(t.textSize)
	t.lines[len(t.lines)-1].SetString(str)
    t.charWidth = t.lines[0].CharWidth()
}
