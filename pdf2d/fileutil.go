package pdf2d

import "github.com/stanim/gofpdf"

// SaveToPdfFile create and save a pdf document to a file
func SaveToPdfFile(filePath string, pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose(filePath)
}
