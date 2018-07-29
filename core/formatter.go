package core

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

// This helps to generate the data in format of a table to be printed on the screen
func TableOutput(data [][]string, header []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)        // Pass all the header data
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)  // Set Border to false
	table.AppendBulk(data)		   // Pass all the data to be displayed on the table
	table.Render()
}