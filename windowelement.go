package main

import (
	"github.com/tadeuszjt/geom/32"
)

type WindowElement interface {
	GetRect() geom.Rect
	SetRect(geom.Rect)
}
