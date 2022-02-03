package themes

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/esimov/stackblur-go"
	"github.com/fogleman/gg"
	"github.com/getsentry/sentry-go"
	"github.com/mitchellh/mapstructure"
	"github.com/oliamb/cutter"
	"gorfe/media"
	"image"

	//"github.com/pixiv/go-libjpeg/jpeg"
	"gorfe/constants"
	"image/color"
	"image/jpeg"
	"log"
	"sync"

	//"image/color"
	"os"

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

type ImageChunk struct {
	Image    image.Image
	Position int
}

var config utils.ConfigFile

func InitializeGridTheme() {
	config = utils.GetConfig()
}

func GenerateGridImage(request structs.GenerateRequest, span *sentry.Span) (float64, string) {
	start := time.Now()
	var themeData GridThemeData

	renderSpan := span.StartChild("rendering")

	mapstructure.Decode(request.Data, &themeData)

	tileSize := themeData.TileSize

	fmt.Println(themeData.TileSize, themeData.Columns)

	width, height := float64(themeData.Columns*tileSize), float64(themeData.Rows*tileSize)

	c := gg.NewContext(int(width), int(height))

	var wg sync.WaitGroup
	wg.Add(themeData.Rows)

	for _i := 0; _i < themeData.Rows; _i++ {

		go func(i int) {
			current := i * themeData.Columns
			rowSpan := renderSpan.StartChild("row")
			rc := gg.NewContext(int(width), int(tileSize))
			defer wg.Done()
			for j := 0; j < themeData.Columns; j++ {
				x := j * tileSize
				y := 0

				if current >= len(themeData.Tiles) {
					continue
				}

				tile := themeData.Tiles[current]

				tileSpan := rowSpan.StartChild("tile.render")
				tileSpan.Description = fmt.Sprintf("Tile: %s", tile.Name)
				tileSpan.Data = map[string]interface{}{
					"name":      tile.Name,
					"secondary": tile.Secondary,
					"image":     tile.Image,
				}

				image, err := media.GetImage(tile.Image)

				if err != nil {
					fmt.Println("Couldn't get image from " + themeData.Tiles[current].Image)
				} else {
					if image.Bounds().Dx() == image.Bounds().Dy() {
						image = imaging.Resize(image, tileSize, tileSize, imaging.Lanczos)
						rc.DrawImage(image, x, y)
					} else {
						fmt.Println(image.Bounds().Dx())

						cropped, _ := cutter.Crop(image, cutter.Config{
							Width:   1,
							Height:  1,
							Mode:    cutter.Centered,
							Options: cutter.Ratio,
						})

						fmt.Println(cropped.Bounds())

						image = imaging.Resize(cropped, tileSize, tileSize, imaging.Lanczos)

						rc.DrawImage(image, x, y)
					}
				}

				if themeData.ShowNames {
					drawOverlay(rc, themeData, tile, float64(x), float64(y), float64(tileSize))
				}
				tileSpan.Finish()

				current++
			}
			compositionSpan := rowSpan.StartChild("composition")
			c.DrawImage(rc.Image(), 0, i*tileSize)
			compositionSpan.Finish()

			rowSpan.Finish()
		}(_i)

	}

	wg.Wait()

	fmt.Println("tiles done")

	if _, err := os.Stat(config.ExportPath); os.IsNotExist(err) {
		fmt.Println("Export path does not exist!")
	}

	fileName := request.ID + ".webp"

	output, err := os.Create(config.ExportPath + fileName)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer output.Close()

	//options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, float32(config.Grid.Quality))
	if err != nil {
		fmt.Println(err.Error())
	}

	renderSpan.Finish()

	compressionSpan := span.StartChild("encoding")

	//if err := webp.Encode(output, c.Image(), options); err != nil {
	//	fmt.Println(err.Error())
	//}

	if err = jpeg.Encode(output, c.Image(), &jpeg.Options{Quality: 60}); err != nil {
		log.Printf("failed to encode: %v", err)
	}

	compressionSpan.Finish()

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
	padding := tileSize * 0.03
	fontSize := tileSize * 0.07
	fontSizeSecondary := tileSize * 0.055

	fontSizeSecondaryGap := tileSize * 0.086

	if themeData.Style == "DEFAULT" {
		//font, _ := truetype.Parse(goregular.TTF)

		grad := gg.NewLinearGradient(0, float64(y), 0, float64(tileSize+y))
		grad.AddColorStop(0, color.RGBA{A: 130})
		grad.AddColorStop(0.26, color.RGBA{A: 30})
		grad.AddColorStop(0.38, color.RGBA{A: 0})

		cg.SetFillStyle(grad)
		cg.DrawRectangle(x, y, tileSize, tileSize)
		cg.Fill()

		cg.SetRGB(1, 1, 1)

		utils.SetFontFace(cg, constants.PoppinsSemiBold, fontSize)

		utils.DrawTextWithEllipsis(cg, tile.Name, x+padding, y+padding, 0.0, 1.0, tileSize-(padding*2))

		if tile.Secondary != nil {
			utils.SetFontFace(cg, constants.PoppinsRegular, fontSizeSecondary)
			utils.DrawTextWithEllipsis(cg, *tile.Secondary, x+padding, y+padding+fontSizeSecondaryGap, 0.0, 1.0, tileSize-(padding*2))
		}

	} else if themeData.Style == "CAPTION" {
		cg.SetColor(color.RGBA{A: 110})
		height := tileSize * 0.18

		cg.DrawRectangle(x, tileSize-height+y, tileSize, height)
		cg.Fill()

		cg.SetRGB(1, 1, 1)

		utils.SetFontFace(cg, constants.PoppinsSemiBold, fontSize)

		if tile.Secondary != nil {
			textsSize := fontSize + fontSizeSecondary

			textsY := y + tileSize - height + ((height - textsSize) / 2)
			textY := textsY
			secondaryY := textsY + ((textsSize / 3) * 2)

			utils.DrawTextWithEllipsis(cg, tile.Name, x+(tileSize/2), textY, 0.5, 1.0, tileSize-(padding*2))

			utils.SetFontFace(cg, constants.PoppinsRegular, fontSizeSecondary)
			utils.DrawTextWithEllipsis(cg, *tile.Secondary, x+(tileSize/2), secondaryY, 0.5, 1.0, tileSize-(padding*2))
		} else {
			utils.DrawTextWithEllipsis(cg, tile.Name, x+(tileSize/2), y+tileSize-(height/2), 0.5, 0.5, tileSize-(padding*2))
		}

	} else if themeData.Style == "SHADOW" {
		intTileSize := int(tileSize)

		tg := gg.NewContext(intTileSize, intTileSize)

		tg.SetRGB(0, 0, 0)

		spread := 1
		spreads := 1
		spreadPoint := float64(spread) / float64(spreads)

		utils.SetFontFace(tg, constants.PoppinsSemiBold, fontSize)
		for i := 0; i <= spreads; i++ {
			for j := 0; j <= spreads; j++ {
				sideX := (padding - float64(spread)) + (spreadPoint * float64(i))
				sideY := (padding - float64(spread)) + (spreadPoint * float64(i))
				utils.DrawTextWithEllipsis(tg, tile.Name, sideX, sideY, 0.0, 1.0, tileSize-(padding*2))
			}
		}

		if tile.Secondary != nil {
			utils.SetFontFace(tg, constants.PoppinsRegular, fontSizeSecondary)
			for i := 0; i <= spreads; i++ {
				for j := 0; j <= spreads; j++ {
					sideX := (padding - float64(spread)) + (spreadPoint * float64(i))
					sideY := (padding + fontSizeSecondaryGap - float64(spread)) + (spreadPoint * float64(i))
					utils.DrawTextWithEllipsis(tg, *tile.Secondary, sideX, sideY, 0.0, 1.0, tileSize-(padding*2))
				}
			}
		}

		radius := 12
		var done = make(chan struct{}, radius)
		im := stackblur.Process(tg.Image(), uint32(intTileSize), uint32(intTileSize), uint32(radius), done)
		<-done

		cg.DrawImage(im, int(x), int(y))

		cg.SetRGB(1, 1, 1)

		utils.SetFontFace(cg, constants.PoppinsSemiBold, fontSize)
		utils.DrawTextWithEllipsis(cg, tile.Name, x+padding, y+padding, 0.0, 1.0, tileSize-(padding*2))

		if tile.Secondary != nil {
			utils.SetFontFace(cg, constants.PoppinsRegular, fontSizeSecondary)
			utils.DrawTextWithEllipsis(cg, *tile.Secondary, x+padding, y+padding+fontSizeSecondaryGap, 0.0, 1.0, tileSize-(padding*2))
		}
	}
}
