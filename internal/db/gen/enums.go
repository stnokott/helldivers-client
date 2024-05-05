// Code generated by sqlc-gen-enum. DO NOT EDIT.
package gen

//go:generate go run golang.org/x/tools/cmd/stringer@v0.20.0 -type=Table -linecomment
type Table int

const (
	TableWars                Table = iota + 1 // Wars
	TableCampaigns                            // Campaigns
	TableEvents                               // Events
	TableBiomes                               // Biomes
	TableHazards                              // Hazards
	TablePlanets                              // Planets
	TableAssignmentTasks                      // Assignment Tasks
	TableAssignments                          // Assignments
	TableDispatches                           // Dispatches
	TableWarSnapshots                         // War Snapshots
	TableEventSnapshots                       // Event Snapshots
	TableAssignmentSnapshots                  // Assignment Snapshots
	TableSnapshotStatistics                   // Snapshot Statistics
	TablePlanetSnapshots                      // Planet Snapshots
	TableSnapshots                            // Snapshots
)

var AllTables = []Table{
	TableWars,
	TableCampaigns,
	TableEvents,
	TableBiomes,
	TableHazards,
	TablePlanets,
	TableAssignmentTasks,
	TableAssignments,
	TableDispatches,
	TableWarSnapshots,
	TableEventSnapshots,
	TableAssignmentSnapshots,
	TableSnapshotStatistics,
	TablePlanetSnapshots,
	TableSnapshots,
}
