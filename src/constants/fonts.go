package constants

import (
	"github.com/tdewolff/canvas"
	"gorfe/utils"
)

var Poppins *canvas.FontFamily
var PoppinsSemiBold *canvas.FontFamily

func LoadFonts() {
	var err error

	Poppins = canvas.NewFontFamily("PoppinsRegular")
	Poppins.Use(canvas.CommonLigatures)

	err = Poppins.LoadFontFile("./src/assets/fonts/Poppins-Regular.ttf", canvas.FontRegular)
	utils.FailOnError(err)

	err = Poppins.LoadFontFile("./src/assets/fonts/Poppins-SemiBold.ttf", canvas.FontSemibold)
	utils.FailOnError(err)
}
