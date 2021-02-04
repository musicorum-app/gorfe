package utils

import (
	"github.com/fogleman/gg"
)

func DrawTextWithEllipsis(ctx *gg.Context, s string, x, y, ax, ay, maxWidth float64) {
	w, _ := ctx.MeasureString(s)

	if w <= maxWidth {
		ctx.DrawStringAnchored(s, x, y, ax, ay)
		return
	}

	for w > maxWidth {
		if len(s) == 2 {
			break
		}

		s = TrimLastChar(s)
		w, _ = ctx.MeasureString(s + "…")
	}

	ctx.DrawStringAnchored(s+"…", x, y, ax, ay)
}

func SetFontFace(ctx *gg.Context, font string, fontSize float64) {
	if err := ctx.LoadFontFace(font, fontSize); err != nil {
		panic(err)
	}
}
