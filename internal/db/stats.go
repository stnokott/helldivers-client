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
}

func (s tableMergeStats) must(tableName string) {
	if _, ok := s[tableName]; !ok {
		s[tableName] = &mergeStats{}
	}
}

func (s tableMergeStats) IncrUpdate(tableName string) {
	s.must(tableName)
	s[tableName].Updated++
}

func (s tableMergeStats) IncrInsert(tableName string) {
	s.must(tableName)
	s[tableName].Inserted++
}

func (s tableMergeStats) Print(logger *log.Logger) {
	var padTableName, padInserted, padUpdated int
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
	}
	formatString := fmt.Sprintf("** %%-%ds -> %%%dd inserted, %%%dd updated", padTableName, padInserted, padUpdated)
	for tableName, stats := range s {
		logger.Printf(formatString, tableName, stats.Inserted, stats.Updated)
	}
}
