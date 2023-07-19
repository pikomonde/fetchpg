package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	dir01 := "/home/zenga/project/personal/pikomonde/fetchpg-230715-notice/fetchpg230715notice/output/part_01/"
	dir02 := "/home/zenga/project/personal/pikomonde/fetchpg-230715-notice/fetchpg230715notice/output/part_02/"
	dirff := "/home/zenga/project/personal/pikomonde/fetchpg-230715-notice/fetchpg230715notice/output/final/"
	files01 := fetchFilenames(dir01)
	files02 := fetchFilenames(dir02)

	for filename, _ := range files01 {
		if _, ok := files02[filename]; !ok {
			continue
		}
		// fmt.Println("---->", filename)

		file01, _ := os.Open(dir01 + filename)
		defer file01.Close()
		file01Reader := csv.NewReader(file01)
		file01Content, _ := file01Reader.ReadAll()
		// file01Content = file01Content[1:] // remove header

		file02, _ := os.Open(dir02 + filename)
		defer file02.Close()
		file02Reader := csv.NewReader(file02)
		file02Content, _ := file02Reader.ReadAll()
		file02Content = file02Content[1:] // remove header

		fileff, _ := os.Create(dirff + filename)
		defer fileff.Close()
		fileffWriter := csv.NewWriter(fileff)
		defer fileffWriter.Flush()
		fileffContent := append(file01Content, file02Content...)
		for _, fileffRow := range fileffContent {
			fileffWriter.Write(fileffRow)
		}

	}

}

func fetchFilenames(dir string) map[string]string {
	csvFiles := make(map[string]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".csv") {
			// fmt.Println("---->", file.Name())
			csvFiles[file.Name()] = file.Name()
		}
	}

	return csvFiles
}
