package draw2dsvg

import (
	"encoding/xml"
	_ "errors"
	"io"
	"os"
)

func WriteSvg(w io.Writer, svg *Svg) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "\t")
	return encoder.Encode(svg)
}

func SaveToSvgFile(filePath string, svg *Svg) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteSvg(f, svg)
}
