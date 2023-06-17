package garbanzo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const APPDIR = ".garbanzo"
const SAVEFILE = "mark.csv"

// 初期化する
func (c *Config) putConfDir() {
	fileInfo, err := os.Lstat(c.baseDir)
	if err != nil {
		fmt.Println(err)
	}
	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm

	if _, err := os.Stat(c.appDirPath()); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(c.appDirPath(), unixPerms); err != nil {
			log.Fatal(err)
		}
	}
}

func (c *Config) markToFile(url string) {
	if _, err := os.Stat(c.saveFilePath()); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(c.saveFilePath())
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}

		writer := csv.NewWriter(f)
		writer.Write([]string{"url"})
		writer.Flush() // Writeだと内部バッファに書き込まれるだけ
	}

	file, err := os.OpenFile(c.saveFilePath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// CSV 形式でデータを書き込む
	writer := csv.NewWriter(file)
	defer writer.Flush() // Writeだと内部バッファに書き込まれるだけ
	writer.Write([]string{url})
}

func (c *Config) isMarked(url string) bool {
	file, err := os.Open(c.saveFilePath())
	if err != nil {
		log.Fatal(err)
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

// homedir, _ := os.UserHomeDir()
