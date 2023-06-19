package garbanzo

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const FEEDFILE = "feeds.yml"

type Config struct {
	baseDir string
}

func NewConfig(baseDir string) *Config {
	return &Config{
		baseDir: baseDir,
	}
}

func (c *Config) appDirPath() string {
	return filepath.Join(c.baseDir, APPDIR)
}

func (c *Config) saveFilePath() string {
	return filepath.Join(c.baseDir, APPDIR, SAVEFILE)
}

func (c *Config) feedFilePath() string {
	return filepath.Join(c.baseDir, APPDIR, FEEDFILE)
}

func (c *Config) loadFeedFile() feedSources {
	b, err := os.ReadFile(c.feedFilePath())
	if err != nil {
		log.Println(err)
	}
	f := c.loadFeedSources(b)
	return f
}

func (c *Config) loadFeedSources(b []byte) feedSources {
	feeds := feedSources{}
	yaml.Unmarshal(b, &feeds)
	return feeds
}

// 設定ディレクトリを初期化する。すでにあれば何もしない
func (c *Config) PutConfDir() {
	fileInfo, err := os.Lstat(c.baseDir)
	if err != nil {
		fmt.Println(err)
	}
	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm

	// 設定ファイルを初期化する
	if _, err := os.Stat(c.appDirPath()); errors.Is(err, os.ErrNotExist) {
		{
			if err := os.Mkdir(c.appDirPath(), unixPerms); err != nil {
				log.Fatal(err)
			}
		}
		{
			f, err := os.Create(c.saveFilePath())
			defer f.Close()
			if err != nil {
				log.Println(err)
			}
		}
		{
			f, err := os.Create(c.feedFilePath())
			defer f.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
