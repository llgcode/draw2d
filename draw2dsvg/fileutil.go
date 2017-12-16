package draw2dsvg

import (
	"os"
	"bytes"
	"errors"
	svgo "github.com/ajstarks/svgo/float"
)

func SaveToSvgFile(filePath string, svg *svgo.SVG) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if b, ok := svg.Writer.(*bytes.Buffer); ok {
		bytes.NewBuffer(b.Bytes()).WriteTo(f) // clone buffer to make multiple writes possible
	} else {
		return errors.New("Svg has not been not created from with NewSvg (dow not have byte.Buffer as its Writer)")
	}

	return nil
}