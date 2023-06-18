package garbanzo

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const FEEDFILE = "feeds.yml"

type Config struct {
	baseDir string
}

type feedSources []feedSource

type feedSource struct {
	Title string
	URL   string
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
		log.Println(c.feedFilePath(), "is not exists")
	}
	f := c.loadFeedSources(b)
	return f
}

func (c *Config) loadFeedSources(b []byte) feedSources {
	feeds := feedSources{}
	yaml.Unmarshal(b, &feeds)
	return feeds
}
