// Package stats provides utilities for collecting and printing statistics
package stats

import (
	"log"
	"strconv"
	"strings"

	"github.com/stnokott/helldivers-client/internal/db/gen"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Collector collects per-table C(R)U(D)-statistics.
type Collector map[gen.Table]*tblMergeStats

type tblMergeStats struct {
	Inserted int64
	Updated  int64
	Noop     int64
}

// NewCollector creates a new mergeStats instance.
func NewCollector() Collector {
	s := Collector{}
	for _, tbl := range gen.AllTables {
		s[tbl] = &tblMergeStats{}
	}
	return s
}

// Inserted adds `n` inserts to the statistics for `table`.
func (c Collector) Inserted(table gen.Table, n int64) {
	c[table].Inserted += n
}

// Updated adds `n` inserts to the statistics for `table`.
func (c Collector) Updated(table gen.Table, n int64) {
	c[table].Updated += n
}

// Noop adds `n` inserts to the statistics for `table`.
func (c Collector) Noop(table gen.Table, n int64) {
	c[table].Noop += n
}

func statOrNan(x int64) string {
	if x > 0 {
		return strconv.FormatInt(x, 10)
	}
	return "-"
}

func (c Collector) renderTable() string {
	w := table.NewWriter()

	w.AppendHeader(table.Row{"Table Name", "Inserted", "Updated", "Unchanged"})

	total := tblMergeStats{}
	for tableName, stats := range c {
		w.AppendRow(table.Row{
			tableName,
			statOrNan(stats.Inserted),
			statOrNan(stats.Updated),
			statOrNan(stats.Noop),
		})
		total.Inserted += stats.Inserted
		total.Updated += stats.Updated
		total.Noop += stats.Noop
	}
	w.AppendSeparator()
	w.AppendFooter(table.Row{"Total", total.Inserted, total.Updated, total.Noop})
	w.SetStyle(table.StyleLight)
	w.SortBy([]table.SortBy{
		{Name: "Inserted", Mode: table.DscNumericAlpha},
		{Name: "Updated", Mode: table.DscNumericAlpha},
		{Name: "Unchanged", Mode: table.DscNumericAlpha},
	})
	return w.Render()
}

// Print prints the collected statistics to `logger` per line, retaining
// potential logging prefixes.
func (c Collector) Print(logger *log.Logger) {
	rendered := c.renderTable()
	lines := strings.Split(rendered, "\n")
	for _, line := range lines {
		logger.Println(line)
	}
}
