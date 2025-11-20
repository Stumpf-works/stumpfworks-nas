package cli

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table creates a formatted table with StumpfWorks styling
func Table(headers []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(true)
	table.SetHeaderLine(true)
	table.SetRowLine(false)
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}

// KeyValueTable creates a simple key-value table
func KeyValueTable(data map[string]string) {
	rows := make([][]string, 0, len(data))
	for key, value := range data {
		rows = append(rows, []string{key, value})
	}
	Table([]string{"Key", "Value"}, rows)
}
