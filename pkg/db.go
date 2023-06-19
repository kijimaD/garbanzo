package garbanzo

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
)

const APPDIR = ".garbanzo"
const SAVEFILE = "mark.csv"

// 既読ファイルに書き込む
func (c *Config) markToFile(url string) {
	if _, err := os.Stat(c.saveFilePath()); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(c.saveFilePath())
		defer f.Close()
		if err != nil {
			log.Println(err)
		}

		writer := csv.NewWriter(f)
		writer.Write([]string{"url"}) // ヘッダー
		writer.Flush()                // Writeだと内部バッファに書き込まれるだけ
	}

	file, err := os.OpenFile(c.saveFilePath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// CSV 形式でデータを書き込む
	writer := csv.NewWriter(file)
	defer writer.Flush() // Writeだと内部バッファに書き込まれるだけ
	writer.Write([]string{url})
}

// ファイルにURLが存在すれば既読状態としてtrueを返す
func (c *Config) isMarked(url string) bool {
	file, err := os.Open(c.saveFilePath())
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	var exists bool
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			exists = false
			break
		}
		if record[0] == url {
			exists = true
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return exists
}
