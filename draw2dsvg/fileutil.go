package draw2dsvg

import (
	"os"
	"bytes"
	_ "errors"
)

func SaveToSvgFile(filePath string, svg *SVG) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes.NewBuffer((*bytes.Buffer)(svg).Bytes()).WriteTo(f) // clone buffer to make multiple writes possible

	return nil
}