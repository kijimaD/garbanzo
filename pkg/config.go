package garbanzo

import "path/filepath"

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
