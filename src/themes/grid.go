package themes

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/mitchellh/mapstructure"
	"github.com/nickalie/go-webpbin"
	"github.com/tdewolff/canvas"
	"gorfe/constants"
	"image/color"
	"sync"

	//"image/color"
	"os"

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
	TileSize      int             `mapstructure:"tile_size"`
	ShowNames     bool            `mapstructure:"show_names"`
	ShowPlaycount bool            `mapstructure:"show_playcount"`
	Style         string          `json:"style"`
}

var config utils.ConfigFile

func InitializeGridTheme() {
	config = utils.GetConfig()
}

func GenerateGridImage(request structs.GenerateRequest) (float64, string) {
	start := time.Now()
	var themeData GridThemeData

	mapstructure.Decode(request.Data, &themeData)

	tileSize := themeData.TileSize

	fmt.Println(themeData.TileSize, themeData.Columns)

	width, height := float64(themeData.Columns*tileSize), float64(themeData.Rows*tileSize)
	//
	//c := canvas.New(width, height)
	//ctx := canvas.NewContext(c)
	//
	//ctx.SetFillColor(color.RGBA{255, 100, 0, 100})
	//ctx.DrawPath(0, 0, canvas.Rectangle(width, height))
	//ctx.SetFillColor(canvas.White)

	c := gg.NewContext(int(width), int(height))

	//titleFace := constants.Poppins.Face(300, color.White, canvas.FontSemibold, canvas.FontNormal)

	var wg sync.WaitGroup
	wg.Add(len(themeData.Tiles))

	current := 0
	for i := 0; i < themeData.Rows; i++ {
		for j := 0; j < themeData.Columns; j++ {
			x := j * tileSize
			y := i * tileSize

			if current >= len(themeData.Tiles) {
				continue
			}

			go func(curr int) {
				defer wg.Done()

				tile := themeData.Tiles[curr]

				image, err := media.GetImage(tile.Image)

				if err != nil {
					fmt.Println("Couldn't get image from " + themeData.Tiles[curr].Image)
				} else {
					image = imaging.Resize(image, tileSize, tileSize, imaging.Lanczos)
					//ctx.DrawImage(float64(x), height - float64(y) - float64(tileSize), image, 1)
					//ctx.DrawText(float64(x+20), float64(y+20), canvas.NewTextBox(titleFace, string(current), 0, 0, canvas.Left, canvas.Top, 0, 0))
					c.DrawImage(image, x, y)
				}

				if themeData.ShowNames {
					drawOverlay(c, themeData, tile, float64(x), float64(y), float64(tileSize))
				}
			}(current)

			current++
		}
	}

	wg.Wait()

	if themeData.ShowNames && themeData.Style == "DEFAULT" {
		for i := 0; i < themeData.Rows; i++ {
			y := i * tileSize
			grad := gg.NewLinearGradient(0, float64(y), 0, float64(tileSize+y))
			grad.AddColorStop(0, color.RGBA{A: 147})
			grad.AddColorStop(0.5, color.RGBA{A: 0})

			c.SetFillStyle(grad)
			c.DrawRectangle(0, float64(y), float64(themeData.Columns*tileSize), float64(tileSize))
			c.Fill()
		}
	}

	if _, err := os.Stat(config.ExportPath); os.IsNotExist(err) {
		fmt.Println("Export path does not exist!")
	}

	fileName := request.ID + ".webp"

	webpbin.NewCWebP().
		Quality(config.Grid.Quality).
		InputImage(c.Image()).
		OutputFile(config.ExportPath + fileName).
		Run()

	duration := time.Since(start)
	fmt.Printf("Duration: %s", duration)
	fmt.Println()

	return duration.Seconds(), fileName

}

func drawOverlay(
	cg *gg.Context,
	themeData GridThemeData,
	tile GridThemeTile,
	x, y, tileSize float64,
) {
	if themeData.Style == "DEFAULT" {
		c := canvas.New(float64(themeData.TileSize), float64(themeData.TileSize))
		ctx := canvas.NewContext(c)

		ctx.SetFillColor(color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		})
		titleFace := constants.Poppins.Face(300, color.White, canvas.FontSemibold, canvas.FontNormal)

		text := canvas.NewTextBox(titleFace, "alo alo bom dia", float64(themeData.TileSize), float64(themeData.TileSize), canvas.Left, canvas.Top, 10, 0)

		ctx.DrawText(x, y, text)

		//im := c.Re
		//
		//cg.DrawImage(c, x, y)
	}
}
