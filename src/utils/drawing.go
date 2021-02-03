package utils

import (
	"fmt"
	"github.com/fogleman/gg"
)

func DrawSizedString(ctx *gg.Context, s string, x, y, ax, ay, maxWidth float64) {
	w, _ := ctx.MeasureString(s)

	if w <= maxWidth {
		ctx.DrawStringAnchored(s, x, y, ax, ay)
	} else {
		str := s + "..."
		for true {
			if len(str) == 5 {
				break
			} else {
				str = str[:len(str)-4] + "..."
				fmt.Println(str)
				w, _ = ctx.MeasureString(str)
				if w <= maxWidth {
					ctx.DrawStringAnchored(str, x, y, ax, ay)
					break
				}
			}
		}
	}
}
