package themes

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/mitchellh/mapstructure"
	"github.com/nickalie/go-webpbin"
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

	if themeData.ShowNames {
		fontSize := float64(tileSize) * 0.08
		if err := c.LoadFontFace("src/assets/fonts/Poppins-Regular.ttf", fontSize); err != nil {
			panic(err)
		}
	}

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

			tile := themeData.Tiles[current]

			image, err := media.GetImage(tile.Image)

			if err != nil {
				fmt.Println("Couldn't get image from " + themeData.Tiles[current].Image)
			} else {
				image = imaging.Resize(image, tileSize, tileSize, imaging.Lanczos)
				//ctx.DrawImage(float64(x), height - float64(y) - float64(tileSize), image, 1)
				//ctx.DrawText(float64(x+20), float64(y+20), canvas.NewTextBox(titleFace, string(current), 0, 0, canvas.Left, canvas.Top, 0, 0))
				c.DrawImage(image, x, y)
			}

			if themeData.ShowNames {
				fmt.Println("ALO")
				drawOverlay(c, themeData, tile, float64(x), float64(y), float64(tileSize))
			}

			current++
		}
	}

	if _, err := os.Stat(config.ExportPath); os.IsNotExist(err) {
		fmt.Println("Export path does not exist!")
	}

	fileName := request.ID + ".webp"

	err := webpbin.NewCWebP().
		Quality(config.Grid.Quality).
		InputImage(c.Image()).
		OutputFile(config.ExportPath + fileName).
		Run()

	if err != nil {
		fmt.Println(err.Error())
	}

	duration := time.Since(start)
	fmt.Printf("Duration: %s", duration)
	fmt.Println()

	fmt.Println(config.ExportPath + fileName)

	return duration.Seconds(), fileName

}

func drawOverlay(
	cg *gg.Context,
	themeData GridThemeData,
	tile GridThemeTile,
	x, y, tileSize float64,
) {
	padding := 8.0

	if themeData.Style == "DEFAULT" {
		//font, _ := truetype.Parse(goregular.TTF)

		grad := gg.NewLinearGradient(0, float64(y), 0, float64(tileSize+y))
		grad.AddColorStop(0, color.RGBA{A: 160})
		grad.AddColorStop(0.35, color.RGBA{A: 0})

		cg.SetFillStyle(grad)
		cg.DrawRectangle(x, y, tileSize, tileSize)
		cg.Fill()

		cg.SetRGB(1, 1, 1)

		//cg.SetFontFace(face)

		fmt.Println(tile)

		utils.DrawSizedString(cg, tile.Name, x+padding, y+padding, 0.0, 1.0, tileSize-(padding*2))

		//im := c.Re
		//
		//cg.DrawImage(c, x, y)
	} else if themeData.Style == "CAPTION" {
		cg.SetColor(color.RGBA{A: 150})
		height := tileSize * 0.16

		topPadding := 11.0

		cg.DrawRectangle(x, tileSize-height+y, tileSize, height)
		cg.Fill()

		cg.SetRGB(1, 1, 1)

		utils.DrawSizedString(cg, tile.Name, x+(tileSize/2), y+tileSize-height+topPadding, 0.5, 1.0, tileSize-(padding*2))
	}
}
