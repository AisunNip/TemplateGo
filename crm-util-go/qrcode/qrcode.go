package qrcode

import (
	"bytes"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
	"image/png"
)

const (
	barcodeMargin int = 4
	qrcodeMargin  int = 1
)

func createBarcode(writer gozxing.Writer, format gozxing.BarcodeFormat,
	content string, width int, height int, margin int) ([]byte, error) {

	var binary []byte

	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = margin

	// Generate a barcode image (*BitMatrix)
	// *BitMatrix implements the image.Image interface
	img, err := writer.Encode(content, format, width, height, hints)

	if err != nil {
		return binary, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)

	if err == nil {
		binary = buf.Bytes()
	}

	return binary, err
}

func CreateBarcode128(content string, width int, height int) ([]byte, error) {
	return createBarcode(oned.NewCode128Writer(), gozxing.BarcodeFormat_CODE_128,
		content, width, height, barcodeMargin)
}

func CreateBarcode93(content string, width int, height int) ([]byte, error) {
	return createBarcode(oned.NewCode93Writer(), gozxing.BarcodeFormat_CODE_93,
		content, width, height, barcodeMargin)
}

func CreateBarcode39(content string, width int, height int) ([]byte, error) {
	return createBarcode(oned.NewCode39Writer(), gozxing.BarcodeFormat_CODE_39,
		content, width, height, barcodeMargin)
}

func CreateBarcodeCodabar(content string, width int, height int) ([]byte, error) {
	return createBarcode(oned.NewCodaBarWriter(), gozxing.BarcodeFormat_CODABAR,
		content, width, height, barcodeMargin)
}

func CreateBarcodeITF(content string, width int, height int) ([]byte, error) {
	return createBarcode(oned.NewITFWriter(), gozxing.BarcodeFormat_ITF,
		content, width, height, barcodeMargin)
}

func CreateQRCode(content string, size int) ([]byte, error) {
	return createBarcode(qrcode.NewQRCodeWriter(), gozxing.BarcodeFormat_QR_CODE,
		content, size, size, qrcodeMargin)
}
