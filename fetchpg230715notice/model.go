package fetchpg230715notice

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Row struct {
	URL                string
	County             string
	ScreenshotFilename string
	Content            string
}

type Data struct {
	Rows     []Row
	filename string
}

func NewData() *Data {
	return &Data{
		Rows: make([]Row, 0),
	}
}

func (d *Data) Add(url, county, ssFilename, content string) {

	if d.Rows == nil {
		d.Rows = make([]Row, 0)
	}

	d.Rows = append(d.Rows, Row{
		URL:                url,
		County:             county,
		ScreenshotFilename: ssFilename,
		Content:            content,
	})
}
func (d *Data) SaveFile(fn string) {
	// change file name if already exist
	if d.filename == "" {
		for _, err := os.Stat(fn); err == nil; _, err = os.Stat(fn) {
			fnArr := strings.Split(fn, ".")
			fn = strings.Join(fnArr[:len(fnArr)-1], ".") + "_copy." + fnArr[len(fnArr)-1]
		}
		d.filename = fn
	}
	fn = d.filename

	// setting up file
	file, err := os.Create(fn)
	if err != nil {
		fmt.Println("ERROR Create:", err)
		panic(err)
	}
	defer file.Close()

	// setting up writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write header
	err = writer.Write([]string{
		"URL",
		"County",
		"Screenshot Filename",
		"Content",
	})
	if err != nil {
		fmt.Println("ERROR Write:", err)
		panic(err)
	}

	// write data
	for _, row := range d.Rows {
		err = writer.Write([]string{
			row.URL,
			row.County,
			row.ScreenshotFilename,
			row.Content,
		})
		if err != nil {
			fmt.Println("ERROR Write:", err)
			panic(err)
		}
	}
}
