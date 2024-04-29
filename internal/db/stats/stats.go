// Package stats provides utilities for collecting and printing statistics
package stats

import (
	"fmt"
	"log"
	"strconv"
	"unicode/utf8"

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
	var padTableName, padInserted, padUpdated, padNoop int
	for tableName, stats := range s {
		if x := utf8.RuneCountInString(tableName.String()); x > padTableName {
			padTableName = x
		}
		if x := len(strconv.FormatInt(stats.Inserted, 10)); x > padInserted {
			padInserted = x
		}
		if x := len(strconv.FormatInt(stats.Updated, 10)); x > padUpdated {
			padUpdated = x
		}
		if x := len(strconv.FormatInt(stats.Noop, 10)); x > padNoop {
			padNoop = x
		}
	}
	formatString := fmt.Sprintf("** %%-%ds -> %%%dd inserted, %%%dd updated, %%%dd unchanged", padTableName, padInserted, padUpdated, padNoop)
	for tableName, stats := range s {
		logger.Printf(formatString, tableName, stats.Inserted, stats.Updated, stats.Noop)
	}
}
