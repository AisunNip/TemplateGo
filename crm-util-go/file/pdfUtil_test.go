package file

import (
	"crm-util-go/timeUtil"
	"fmt"
	"testing"
	"time"
)

func TestGeneratePDF(t *testing.T) {
	fmt.Println("######### Test GeneratePDF #########")
	//FontStyle_Bold := "B"
	//FontStyle_Italic := "I"
	//FontStyle_Underline := "U"

	startTime := time.Now()

	pdf := NewPDF("A4")
	defer pdf.Close()

	err := AddFontFilePDF(pdf, "rsu", "../font/rsu.ttf", 14)
	if err != nil {
		t.Errorf("Error AddFont " + err.Error())
		return
	}

	pdf.SetLineWidth(1)
	pdf.SetFillColor(55, 125, 113)

	pdf.Image("../image/lake.jpg", 0, 0, nil)

	pdf.SetXY(200, 40)
	pdf.Text("Link to google.com")
	pdf.AddExternalLink("https://www.google.com/", 195, 25, 100, 15)

	pdf.SetXY(200, 50)
	pdf.Cell(nil, "Hello Lake (font rsu) size 14")

	err = AddFontFilePDF(pdf, "true_light", "../font/true_light.ttf", 14)
	if err != nil {
		t.Errorf("Error AddFont " + err.Error())
		return
	}

	pdf.SetXY(200, 60)
	pdf.Cell(nil, "Hello Lake (font true_light) size 14")

	pdf.SetFont("rsu", "", 20)
	pdf.SetXY(200, 70)
	pdf.Cell(nil, "Hello Lake (font rsu) size 20")

	pdf.WritePdf("D:/TestCreatePDF.pdf")

	endTime := time.Now()

	respTime := timeUtil.DiffTime(endTime, startTime)

	fmt.Println("ResponseTime:", respTime.Milliseconds())
}
