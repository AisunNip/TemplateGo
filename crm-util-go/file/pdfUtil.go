package file

import "github.com/signintech/gopdf"

/*
pageSize is A0,A1,A2,A3,A4,A5,B4,B5
*/
func NewPDF(pageSize string) *gopdf.GoPdf {
	pdf := new(gopdf.GoPdf)

	switch pageSize {
	case "A0":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA0})
	case "A1":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA1})
	case "A2":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA2})
	case "A3":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA3})
	case "A4":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	case "A5":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA5})
	case "B4":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeB4})
	case "B5":
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeB5})
	default:
		pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	}

	pdf.SetCompressLevel(9)
	pdf.AddPage()

	return pdf
}

func AddFontFilePDF(pdf *gopdf.GoPdf, family string, ttfPath string, size interface{}) error {
	err := pdf.AddTTFFont(family, ttfPath)
	if err != nil {
		return err
	}

	err = pdf.SetFont(family, "", size)
	return err
}
