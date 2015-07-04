package pdf2d

import "github.com/stanim/gofpdf"

// SaveToPdfFile creates and saves a pdf document to a file
func SaveToPdfFile(filePath string, pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose(filePath)
}
