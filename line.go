package main

import (
	"github.com/tadeuszjt/gfx"
)

type TextChunk struct {
	Colour gfx.Colour
	Str    string
}

type Line struct {
	Chunks []TextChunk
}
