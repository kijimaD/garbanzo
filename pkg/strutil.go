package garbanzo

import (
	"fmt"
	"os"
)

func buildHomeMD() (string, error) {
	homedir, _ := os.UserHomeDir()
	conf := NewConfig(homedir)

	tmpl, err := buildTemplateMD()
	if err != nil {
		return "", err
	}
	feedTable, err := buildFeedStatus(conf)
	if err != nil {
		return "", err
	}
	tokenStatus, err := buildTokenStatus(conf)
	if err != nil {
		return "", err
	}

	md := tmpl + feedTable + tokenStatus
	return md, nil
}

func buildTemplateMD() (string, error) {
	data, err := fss.ReadFile("static/home.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func buildFeedStatus(c *Config) (string, error) {
	header := "## RSS Feed lists\n"
	filepath := fmt.Sprintf("`%s`\n\n", c.feedFilePath())

	b, _ := os.ReadFile(c.feedFilePath())
	fss := c.loadFeedSources(b)

	result := header + filepath + fss.dumpFeedsTable()
	return result, nil
}

func buildTokenStatus(c *Config) (string, error) {
	tokenHeader := "## GitHub Token\n"
	tokenPath := "`~/.garbanzo/token`\n\n"

	f, err := os.Open(c.tokenFilePath())
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 1024)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}

	var tokenStatus string
	if len(string(buf)) > 0 {
		tokenStatus = "🟢 ok"
	} else {
		tokenStatus = "🔴 not set"
	}

	return tokenHeader + tokenPath + tokenStatus, nil
}
