package main

import (
	"github.com/tadeuszjt/geom/32"
	"github.com/tadeuszjt/gfx"
	"golang.org/x/image/math/fixed"
)

type TextWindow struct {
	rect       geom.Rect
	tabSize    int
	textSize   int
	cursorLine int
	cursorChar int
	lines      []gfx.Text
	windowTex  gfx.TexID
}


func NewTextWindow(w *gfx.Win, rect geom.Rect) *TextWindow {
	return &TextWindow{
		rect:      rect,
		tabSize:   4,
		textSize:  13,
		windowTex: w.LoadTextureBlank(int(rect.Width()), int(rect.Height())),
	}
}

func (t *TextWindow) curLine() *gfx.Text {
	return &t.lines[t.cursorLine]
}

func (t *TextWindow) GetLine(idx int) string {
    return t.lines[idx].GetString()
}

func (t *TextWindow) SetLine(idx int, str string) {
    t.lines[idx].SetString(str)
    
}

func (t *TextWindow) GetCursor() (line int, char int) {
    return t.cursorLine, t.cursorChar
}

func (t *TextWindow) SetCursor(line, char int) {
    t.cursorChar = char
    t.cursorLine = line
}

func (t *TextWindow) InsertLine(lindIdx int, str string) {
	text := gfx.MakeText()
	text.SetSize(t.textSize)
	text.SetString(str)
	text.SetColour(Solarized.Base0)

	t.lines = append(t.lines, text)
	if lindIdx == len(t.lines)-1 {
		return
	}

	copy(t.lines[lindIdx+1:], t.lines[lindIdx:])
	t.lines[lindIdx] = text
}

func (t *TextWindow) NumLines() int {
    return len(t.lines)
}

func (t *TextWindow) RemoveLine(lindIdx int) {
	if lindIdx < 0 || lindIdx >= len(t.lines) {
		panic("")
	}
	t.lines = append(t.lines[:lindIdx], t.lines[lindIdx+1:]...)
}


func (t *TextWindow) Redraw(w *gfx.Win) {
	texCanvas := w.GetTextureCanvas(t.windowTex)
	texCanvas.Clear(gfx.Colour{0, 43. / 255., 54. / 255., 1})

	for i := range t.lines {
		line := &t.lines[i]
		pos := geom.Vec2{t.rect.Min.X, float32(i * line.Height())}
		gfx.DrawText(texCanvas, line, pos)
	}

	charWidth := t.curLine().CharWidth()
	charHeight := t.curLine().Height()
	cursorX := float32(charWidth.Mul(fixed.I(t.cursorChar)).Ceil())
	cursorY := float32(charHeight * t.cursorLine)

	gfx.DrawSprite(
		texCanvas,
		geom.Ori2{cursorX, cursorY, 0},
		geom.RectOrigin(float32(charWidth.Ceil()), float32(charHeight)),
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
	gfx.DrawRect(c, &t.windowTex, t.rect, geom.RectOrigin(1, 1))
}

func (t *TextWindow) AddLine(line string) {
	t.InsertLine(len(t.lines), line)
}
