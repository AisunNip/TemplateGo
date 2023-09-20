package file

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func ExecCommandNoWait(cmdPath string, cmdArgs ...string) error {
	cmd := exec.Command(cmdPath, cmdArgs...)
	return cmd.Start()
}

func ExecCommand(timeout time.Duration, cmdPath string, cmdArgs ...string) (output string, err error) {
	var cmd *exec.Cmd
	var zeroTimeout time.Duration

	if timeout != zeroTimeout {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd = exec.CommandContext(ctx, cmdPath, cmdArgs...)
	} else {
		cmd = exec.Command(cmdPath, cmdArgs...)
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return output, err
	}

	defer stdout.Close()

	err = cmd.Start()

	if err != nil {
		return output, err
	}

	binary, _ := io.ReadAll(stdout)
	output = string(binary)

	err = cmd.Wait()

	return output, err
}

func GetFileSize(fileName string) (int64, error) {
	var fileSizeBytes int64

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return fileSizeBytes, err
	}

	fileSizeBytes = fileInfo.Size()
	return fileSizeBytes, err
}

func getFileInfoList(dirName string) ([]os.FileInfo, error) {
	file, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file.Readdir(-1)
}

func ListFile(dirName string) ([]os.FileInfo, error) {
	fileInfoList, err := getFileInfoList(dirName)
	if err != nil {
		return nil, err
	}

	var resultFileInfo []os.FileInfo

	for _, fileInfo := range fileInfoList {
		if !fileInfo.IsDir() {
			resultFileInfo = append(resultFileInfo, fileInfo)
		}
	}

	fileInfoList = nil

	return resultFileInfo, nil
}

func ListDirectory(dirName string) ([]os.FileInfo, error) {
	fileInfoList, err := getFileInfoList(dirName)
	if err != nil {
		return nil, err
	}

	var resultFileInfo []os.FileInfo

	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			resultFileInfo = append(resultFileInfo, fileInfo)
		}
	}

	fileInfoList = nil

	return resultFileInfo, nil
}

func ReadFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func WriteFile(fileName string, isAppend bool, binary []byte) error {
	var permFileMode os.FileMode = 0766
	var flag int
	if isAppend {
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	file, err := os.OpenFile(fileName, flag, permFileMode)

	if err != nil {
		return errors.New(fmt.Sprintf("Opening file error: %v", err))
	}

	defer file.Close()

	_, err = file.Write(binary)

	return err
}

func DeleteFile(fileName string) error {
	fileInfo, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	if fileInfo.IsDir() {
		err = errors.New("DeleteFile error because it is a directory: " + fileName)
	} else {
		err = os.Remove(fileName)
	}

	return err
}

func DeleteDirectory(dirName string) error {
	return os.RemoveAll(dirName)
}

func Gzip(srcFileName string, gzipDir string) error {
	reader, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer reader.Close()

	fileName := filepath.Base(srcFileName)
	gzipDir = filepath.Join(gzipDir, fmt.Sprintf("%s.gz", fileName))

	writer, err := os.Create(gzipDir)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver, _ := gzip.NewWriterLevel(writer, gzip.BestCompression)
	archiver.Name = fileName
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

func UnGzip(gzFileName string, ungzipDir string) error {
	reader, err := os.Open(gzFileName)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	ungzipDir = filepath.Join(ungzipDir, archive.Name)
	writer, err := os.Create(ungzipDir)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}
