package draw2d


import (
	"freetype-go.googlecode.com/hg/freetype"
	"freetype-go.googlecode.com/hg/freetype/truetype"
	"path"
	"log"
	"io/ioutil"
)


var (
	fontFolder = "../../fonts/"
	fonts = make(map[string] *truetype.Font)
)


type FontStyle byte

const (
	FontStyleNormal FontStyle = iota
	FontStyleBold
	FontStyleItalic
)

type FontFamily byte

const (
	FontFamilySans FontFamily = iota
	FontFamilySerif
	FontFamilyMono
)


type FontData struct {
	Name string
	Family FontFamily
	Style FontStyle
}



func GetFont(fontData FontData) (*truetype.Font) {
	fontFileName := fontData.Name
	switch fontData.Family {
		case FontFamilySans: fontFileName += "s"
		case FontFamilySerif: fontFileName += "r"
		case FontFamilyMono: fontFileName += "m"
	}
	if(fontData.Style & FontStyleBold != 0) {
		fontFileName += "b"
	} else {
		fontFileName += "r"
	} 
	
	if(fontData.Style & FontStyleItalic != 0) {
		fontFileName += "i"
	}
	fontFileName += ".ttf"
	font := fonts[fontFileName]
	if(font != nil) {
		return font
	}
	fonts[fontFileName] = loadFont(fontFileName)
	return fonts[fontFileName]
}

func GetFontFolder() string {
	return fontFolder
}

func SetFontFolder(folder string) {
	fontFolder = folder
}

func loadFont(fontFileName string) (*truetype.Font) {
	fontBytes, err := ioutil.ReadFile(path.Join(fontFolder, fontFileName))
	if err != nil {
		log.Println(err)
		return nil
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}


