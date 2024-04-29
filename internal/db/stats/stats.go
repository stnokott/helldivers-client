// Package stats provides utilities for collecting and printing statistics
package stats

import (
	"fmt"
	"log"
	"text/tabwriter"

	"github.com/stnokott/helldivers-client/internal/db/gen"
)

type Collector interface {
	Inserted(table gen.Table, n int64)
	Updated(table gen.Table, n int64)
	Noop(table gen.Table, n int64)
	Print(logger *log.Logger)
}

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

// TODO: print table instead.
func (s mergeStats) Print(logger *log.Logger) {
	w := tabwriter.NewWriter(logger.Writer(), 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "\t TABLE NAME \t INSERTED \t UPDATED \t UNCHANGED \t")
	for tableName, stats := range s {
		fmt.Fprintf(w, "\t %s \t %d \t %d \t %d \t\n", tableName, stats.Inserted, stats.Updated, stats.Noop)
	}
	w.Flush()
}
