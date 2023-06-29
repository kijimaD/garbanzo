package garbanzo

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const FeedFile = "feeds.yml"
const TokenFile = "token"

type Config struct {
	baseDir string
}

func NewConfig(baseDir string) *Config {
	return &Config{
		baseDir: baseDir,
	}
}

func (c *Config) appDirPath() string {
	return filepath.Join(c.baseDir, AppDir)
}

func (c *Config) saveFilePath() string {
	return filepath.Join(c.baseDir, AppDir, SaveFile)
}

func (c *Config) feedFilePath() string {
	return filepath.Join(c.baseDir, AppDir, FeedFile)
}

func (c *Config) tokenFilePath() string {
	return filepath.Join(c.baseDir, AppDir, TokenFile)
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
	const saveFileContent = `# marked list
`
	const feedFileContent = `# feed list
- desc: RFC
  url: https://www.rfc-editor.org/rfcrss.xml
- desc: Hacker News
  url: https://news.ycombinator.com/rss
`
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
			if _, err = f.Write([]byte(saveFileContent)); err != nil {
				log.Println(err)
			}
			if err != nil {
				log.Println(err)
			}
		}
		{
			f, err := os.Create(c.feedFilePath())
			defer f.Close()
			if _, err = f.Write([]byte(feedFileContent)); err != nil {
				log.Println(err)
			}
			if err != nil {
				log.Println(err)
			}
		}
	}
}
