package utils

import (
	"fmt"
	"github.com/fogleman/gg"
)

func DrawTextWithEllipsis (ctx *gg.Context, s string, x, y, width float64) {
	w, _ := ctx.MeasureString(s)
	for w > width {
		if len(s) == 2 {
			break
		}

		s = TrimLastChar(s)
		w, _ = ctx.MeasureString(s + "…")
	}

	fmt.Println(s, x, y, width)
	ctx.DrawString(s + "…", x, y)
}