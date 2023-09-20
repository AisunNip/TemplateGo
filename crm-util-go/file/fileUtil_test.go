package file

import (
	"fmt"
	"testing"
	"time"
)

func TestExecCommand(t *testing.T) {
	fmt.Println("######### Test Cmd Line wait to completed #########")
	timeout := time.Duration(20) * time.Second

	output, err := ExecCommand(timeout, "cmd", "/C", "dir")
	// output, err := file.ExecCommand(isWait, timeoutSec, "C:\\Program Files\\internet explorer\\iexplore.exe", "http://www.google.com")

	if err != nil {
		t.Errorf("ExecCommand Error %s", err.Error())
	} else {
		fmt.Println("ExecCommand success")
		fmt.Println(output)
	}
}

func TestExecCommandNoWait(t *testing.T) {
	fmt.Println("######### Test Cmd Line NoWait #########")
	err := ExecCommandNoWait("C:\\Program Files\\internet explorer\\iexplore.exe", "http://www.google.com")

	if err != nil {
		t.Errorf("ExecCommandNoWait Error %s", err.Error())
	} else {
		fmt.Println("ExecCommandNoWait success")
	}
}

func TestGetFileSize(t *testing.T) {
	fmt.Println("######### Test GetFileSize #########")
	fileSizeBytes, err := GetFileSize("D:/pui.pdf")
	if err != nil {
		t.Errorf("GetFileSize Error %s", err.Error())
	} else {
		fmt.Println("GetFileSize success. fileSizeBytes=", fileSizeBytes)
	}
}

func TestGzipAndUnGzip(t *testing.T) {
	fmt.Println("######### Test Gzip #########")
	srcFileName := "D:/Test/source/abc.docx"
	gzipDir := "D:/Test/gzip/"
	err := Gzip(srcFileName, gzipDir)

	if err != nil {
		t.Errorf("Gzip Error %s", err.Error())
		return
	} else {
		fmt.Println("Gzip success")
	}

	fmt.Println("######### Test UnGzip #########")
	gzipFile := gzipDir + "abc.docx.gz"
	ungzipDir := "D:/Test/ungzip/"

	err = UnGzip(gzipFile, ungzipDir)
	if err != nil {
		t.Errorf("UnGzip Error %s", err.Error())
	} else {
		fmt.Println("UnGzip success")
	}
}

func TestListFile(t *testing.T) {
	fmt.Println("######### Test ListFile #########")
	fileInfoList, err := ListFile("D:/")
	if err != nil {
		t.Errorf("ListFile Error: " + err.Error())
	} else {
		for _, file := range fileInfoList {
			fmt.Println(fmt.Sprintf("File Name: %s, Size: %d bytes, ModifyTime: %v",
				file.Name(), file.Size(), file.ModTime()))
		}
	}
}

func TestListDirectory(t *testing.T) {
	fmt.Println("######### Test ListDirectory #########")
	fileInfoList, err := ListDirectory("D:/")
	if err != nil {
		t.Errorf("ListDirectory Error: " + err.Error())
	} else {
		for _, file := range fileInfoList {
			fmt.Println(fmt.Sprintf("Dir Name: %s, Size: %d bytes, ModifyTime: %v",
				file.Name(), file.Size(), file.ModTime()))
		}
	}
}

/*
	fmt.Println("######### Test ReadFile + WriteFile #########")
	bytes, err := ReadFile(source)
	if err != nil {
		fmt.Println("ReadFile Error: " + err.Error())
	} else {
		fmt.Println("ReadFile success")

		err = WriteFile("D:/Test/abcWrite.docx", false, bytes)
		if err != nil {
			fmt.Println("WriteFile Error: " + err.Error())
		} else {
			fmt.Println("WriteFile success")
		}
	}

	fmt.Println("######### Test DeleteFile #########")
	err = DeleteFile(source + ".gz")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("DeleteFile success")
	}

	fmt.Println("######### Test DeleteDirectory #########")
	dirName := "D:/TestAA"
	err = DeleteDirectory(dirName)
	if err != nil {
		fmt.Println("DeleteDirectory Error: " + err.Error())
	} else {
		fmt.Println("DeleteDirectory success")
	}
*/
