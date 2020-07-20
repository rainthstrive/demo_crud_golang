package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func RenameWithMD5(filePath string) string {
	//filePath := "E:/go-work/media_sources/donkey-kong-country-cover.jpg"

	fileName, fileDir, err := getFileName(filePath)
	CheckForError(err)

	hash, err := hashFileMd5(filePath)
	CheckForError(err)

	newFileName := hash + fileName[strings.LastIndex(fileName, "."):]

	renameFileResult := renameFile(filePath, fileDir+newFileName)
	if renameFileResult != nil {
		log.Fatal(renameFileResult)
	}

	return fileDir + newFileName
}

func getFileName(filePath string) (string, string, error) {
	fileStat, err := os.Stat(filePath)
	CheckForError(err)

	return fileStat.Name(), strings.TrimSuffix(filePath, fileStat.Name()), err
}

func renameFile(filePath string, newFilePath string) error {
	err := os.Rename(filePath, newFilePath)
	CheckForError(err)

	return nil
}

func hashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	CheckForError(err)

	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}
