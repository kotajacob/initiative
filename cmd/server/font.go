package main

import (
	"embed"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed "font.ttf"
var Files embed.FS

func loadFont() (font.Face, error) {
	b, err := Files.ReadFile("font.ttf")
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              48,
		GlyphCacheEntries: 1,
	}), nil
}
