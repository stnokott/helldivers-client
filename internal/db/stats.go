package db

import (
	"fmt"
	"log"
	"strconv"
)

type tableMergeStats map[string]*mergeStats

type mergeStats struct {
	Inserted int
	Updated  int
	Noop     int
}

func (s tableMergeStats) must(tableName string) {
	if _, ok := s[tableName]; !ok {
		s[tableName] = &mergeStats{}
	}
}

func (s tableMergeStats) Incr(tableName string, exists bool, affectedRows int64) {
	s.must(tableName)
	if !exists {
		s[tableName].Inserted++
	} else {
		if affectedRows == 0 {
			s[tableName].Noop++
		} else {
			s[tableName].Updated++
		}
	}
}

func (s tableMergeStats) Print(logger *log.Logger) {
	var padTableName, padInserted, padUpdated, padNoop int
	for tableName, stats := range s {
		if x := len(tableName); x > padTableName {
			padTableName = x
		}
		if x := len(strconv.Itoa(stats.Inserted)); x > padInserted {
			padInserted = x
		}
		if x := len(strconv.Itoa(stats.Updated)); x > padUpdated {
			padUpdated = x
		}
		if x := len(strconv.Itoa(stats.Noop)); x > padNoop {
			padNoop = x
		}
	}
	formatString := fmt.Sprintf("** %%-%ds -> %%%dd inserted, %%%dd updated, %%%dd unchanged", padTableName, padInserted, padUpdated, padNoop)
	for tableName, stats := range s {
		logger.Printf(formatString, tableName, stats.Inserted, stats.Updated, stats.Noop)
	}
}
