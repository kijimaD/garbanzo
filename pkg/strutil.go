package garbanzo

import (
	"fmt"
	"io"
	"os"
)

// stringを組み立てる関数群

const templateMDPath = "templates/home.md"

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

func buildTemplateMD() (string, error) {
	data, err := fst.ReadFile(templateMDPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func buildFeedStatus(c *Config) (string, error) {
	const header = "## RSS Feed lists\n"
	filePath := fmt.Sprintf("`%s`\n\n", c.feedFilePath())
	const example = "[example config](https://github.com/kijimaD/dotfiles/blob/main/.garbanzo/feeds.yml)\n"

	b, err := os.ReadFile(c.feedFilePath())
	if err != nil {
		return "", err
	}
	feedSource := c.loadFeedSources(b)

	result := header + filePath + feedSource.dumpFeedsTable() + example
	return result, nil
}

func buildTokenStatus(c *Config) (string, error) {
	const header = "## GitHub Token\n"
	tokenPath := fmt.Sprintf("`%s`\n\n", c.tokenFilePath())

	f, err := os.Open(c.tokenFilePath())
	if err != nil {
		return "", err
	}
	defer f.Close()

	bs, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	var status string
	if len(string(bs)) > 0 {
		status = "🟢 ok"
	} else {
		status = "🔴 not set"
	}

	return header + tokenPath + status, nil
}
