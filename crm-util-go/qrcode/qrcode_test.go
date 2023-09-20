package qrcode

import (
	"crm-util-go/file"
	"fmt"
	"testing"
)

func TestCreateBarcode128(t *testing.T) {
	binary, err := CreateBarcode128("Hello Barcode128", 250, 50)

	if err != nil {
		t.Errorf("TestCreateBarcode128 Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/Barcode128.png", false, binary)
		fmt.Println("CreateBarcode128 success")
	}
}

func TestCreateBarcode93(t *testing.T) {
	binary, err := CreateBarcode93("Hello Barcode93", 250, 50)

	if err != nil {
		t.Errorf("TestCreateBarcode93 Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/Barcode93.png", false, binary)
		fmt.Println("CreateBarcode93 success")
	}
}

func TestCreateBarcode39(t *testing.T) {
	binary, err := CreateBarcode39("Hello Barcode39", 250, 50)

	if err != nil {
		t.Errorf("TestCreateBarcode39 Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/Barcode39.png", false, binary)
		fmt.Println("CreateBarcode39 success")
	}
}

func TestCreateBarcodeCodabar(t *testing.T) {
	binary, err := CreateBarcodeCodabar("1234567890", 250, 50)

	if err != nil {
		t.Errorf("TestCreateBarcodeCodabar Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/BarcodeCodabar.png", false, binary)
		fmt.Println("CreateBarcodeCodabar success")
	}
}

func TestCreateBarcodeITF(t *testing.T) {
	binary, err := CreateBarcodeITF("1234567890", 250, 50)

	if err != nil {
		t.Errorf("TestCreateBarcodeITF Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/BarcodeITF.png", false, binary)
		fmt.Println("CreateBarcodeITF success")
	}
}

func TestCreateQRCode(t *testing.T) {
	binary, err := CreateQRCode("1234567890 Paravit Tunvichian", 250)

	if err != nil {
		t.Errorf("TestCreateQRCode Error %s", err.Error())
	} else {
		err = file.WriteFile("D:/QRCode.png", false, binary)
		fmt.Println("CreateQRCode success")
	}
}