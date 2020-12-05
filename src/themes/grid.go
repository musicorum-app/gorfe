package themes

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/mitchellh/mapstructure"
	"image/color"

	//"github.com/nickalie/go-webpbin"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"gorfe/constants"
	"gorfe/media"
	"gorfe/structs"
	"gorfe/utils"
	//"image/color"
	"time"
)

type GridThemeTile struct {
	Image     string  `json:"image"`
	Name      string  `json:"name"`
	Secondary *string `json:"secondary"`
}

type GridThemeData struct {
	Tiles         []GridThemeTile `json:"tiles"`
	Rows          int             `json:"rows"`
	Columns       int             `json:"columns"`
	ShowNames     bool            `mapstructure:"show_names"`
	ShowPlaycount string          `mapstructure:"show_playcount"`
	Style         string          `json:"style"`
}

var config utils.ConfigFile

func InitializeGridTheme() {
	config = utils.GetConfig()
}

func GenerateGridImage(request structs.GenerateRequest) {
	start := time.Now()
	var themeData GridThemeData

	mapstructure.Decode(request.Data, &themeData)

	tileSize := config.Grid.TileSize

	width, height := float64(themeData.Columns*tileSize), float64(themeData.Rows*tileSize)

	c := canvas.New(width, height)
	ctx := canvas.NewContext(c)

	ctx.SetFillColor(color.RGBA{255, 100, 0, 100})
	//ctx.DrawPath(0, 0, canvas.Rectangle(width, height))
	ctx.SetFillColor(canvas.White)

	titleFace := constants.Poppins.Face(300, color.White, canvas.FontSemibold, canvas.FontNormal)



	current := 0
	for i := 0; i + 300 < themeData.Columns; i++ {
		y := i * tileSize
		for j := 0; j < themeData.Rows; j++ {
			x := j * tileSize
			tile := themeData.Tiles[current]

			image, err := media.GetImage(tile.Image)

			if err != nil {
				fmt.Println("Couldn't get image from " + themeData.Tiles[current].Image)
			} else {
				image = imaging.Resize(image, tileSize, tileSize, imaging.Lanczos)
				ctx.DrawImage(float64(x), float64(y), image, 1)
				ctx.DrawText(float64(x+20), float64(y+20), canvas.NewTextBox(titleFace, string(current), 0, 0, canvas.Left, canvas.Top, 0, 0))
				fmt.Println(x, y)
			}

			if themeData.ShowNames {
				//drawOverlay(cctx, themeData.Style, tile, float64(x), float64(y), float64(tileSize))
			}

			current++
		}

		//if themeData.ShowNames && themeData.Style == "DEFAULT" {
		//	grad := gg.NewLinearGradient(0, float64(y), 0, float64(tileSize+y))
		//	grad.AddColorStop(0, color.RGBA{A: 255 * 0.8})
		//	grad.AddColorStop(0.5, color.RGBA{A: 255 * 0})
		//
		//	ctx.SetFillStyle(grad)
		//	ctx.DrawRectangle(0, float64(y), float64(themeData.Columns*tileSize), float64(tileSize))
		//	ctx.Fill()
		//}
	}

	//webpbin.NewCWebP().
	//	Quality(88).
	//	InputImage(c.I)).
	//	OutputFile(config.ExportPath + request.ID + ".webp").
	//	Run()

	text := canvas.NewTextBox(titleFace, "alo alo bom dia", 400, 300, canvas.Left, canvas.Top, 10, 0)

	ctx.DrawText(2, 2, text)

	c.WriteFile(config.ExportPath + request.ID + ".jpg", rasterizer.JPGWriter(0.8, nil))

	duration := time.Since(start)
	fmt.Printf("Duration: %s", duration)
	fmt.Println()

}

func drawOverlay(
	ctx *gg.Context,
	style string,
	tile GridThemeTile,
	x, y, tileSize float64,
) {
	if style == "DEFAULT" {
		ctx.SetHexColor("#FFFFFF")
		//ctx.SetFontFace(constants.PoppinsSemiBold)
		utils.DrawTextWithEllipsis(ctx, tile.Name, x, y, tileSize)
	}
}
