package garbanzo

import (
	"bytes"

	"github.com/olekukonko/tablewriter"
)

type feedSource struct {
	Desc string
	URL  string
}

type feedSources []feedSource

func (ss *feedSources) dumpFeedsTable() string {
	d := [][]string{}
	for _, f := range *ss {
		d = append(d, []string{f.Desc, f.URL})
	}
	var output bytes.Buffer
	table := tablewriter.NewWriter(&output)
	table.SetHeader([]string{"Description", "Feed"})
	for _, dd := range d {
		table.Append(dd)
	}
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.Render()
	return output.String()
}
