// Package stats provides utilities for collecting and printing statistics
package stats

import (
	"log"
	"strconv"
	"strings"

	"github.com/stnokott/helldivers-client/internal/db/gen"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Collector interface {
	Inserted(table gen.Table, n int64)
	Updated(table gen.Table, n int64)
	Noop(table gen.Table, n int64)
	Print(logger *log.Logger)
}

// compile-time interface implementation assertion.
var _ Collector = mergeStats{}

type mergeStats map[gen.Table]*tblMergeStats

type tblMergeStats struct {
	Inserted int64
	Updated  int64
	Noop     int64
}

func NewCollector() mergeStats {
	s := mergeStats{}
	for _, tbl := range gen.AllTables {
		s[tbl] = &tblMergeStats{}
	}
	return s
}

func (s mergeStats) Inserted(table gen.Table, n int64) {
	s[table].Inserted += n
}

func (s mergeStats) Updated(table gen.Table, n int64) {
	s[table].Updated += n
}

func (s mergeStats) Noop(table gen.Table, n int64) {
	s[table].Noop += n
}

func statOrNan(x int64) string {
	if x > 0 {
		return strconv.FormatInt(x, 10)
	}
	return "-"
}

func (s mergeStats) renderTable() string {
	w := table.NewWriter()

	w.AppendHeader(table.Row{"Table Name", "Inserted", "Updated", "Unchanged"})

	total := tblMergeStats{}
	for tableName, stats := range s {
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

func (s mergeStats) Print(logger *log.Logger) {
	rendered := s.renderTable()
	lines := strings.Split(rendered, "\n")
	for _, line := range lines {
		logger.Println(line)
	}
}
