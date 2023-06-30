package garbanzo

import (
	"fmt"
	"io"
	"os"
)

func buildHomeMD(c *Config) (string, error) {
	tmpl, err := buildTemplateMD()
	if err != nil {
		return "", err
	}
	feedTable, err := buildFeedStatus(c)
	if err != nil {
		return "", err
	}
	tokenStatus, err := buildTokenStatus(c)
	if err != nil {
		return "", err
	}

	md := tmpl + feedTable + tokenStatus
	return md, nil
}

const templateMDPath = "static/home.md"

func buildTemplateMD() (string, error) {
	data, err := fss.ReadFile(templateMDPath)
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

	bs, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	var tokenStatus string
	if len(string(bs)) > 0 {
		tokenStatus = "🟢 ok"
	} else {
		tokenStatus = "🔴 not set"
	}

	return tokenHeader + tokenPath + tokenStatus, nil
}
