// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dong-tran/gotrade/helper"
)

// FileSystemRepository stores and retrieves asset snapshots using
// the local file system.
type FileSystemWeekRepository struct {
	// base is the root directory where asset snapshots are stored.
	base string
}

// NewFileSystemRepository initializes a file system repository with
// the given base directory.
func NewFileSystemWeekRepository(base string) Repository {
	return &FileSystemWeekRepository{
		base: base,
	}
}

// Assets returns the names of all assets in the repository.
func (r *FileSystemWeekRepository) Assets() ([]string, error) {
	files, err := os.ReadDir(r.base)
	if err != nil {
		return nil, err
	}

	var assets []string

	suffix := ".csv"

	for _, file := range files {
		name := file.Name()

		if strings.HasSuffix(name, suffix) {
			assets = append(assets, strings.TrimSuffix(name, suffix))
		}
	}

	return assets, nil
}

// Get attempts to return a channel of snapshots for the asset with the given name.
func (r *FileSystemWeekRepository) Get(name string) (<-chan *Snapshot, error) {
	snapshotChan, err := helper.ReadFromCsvFile[Snapshot](r.getCsvFileName(name), true)
	if err != nil {
		return nil, err
	}
	snapshotSlice := helper.ChanToSlice(snapshotChan)
	snapshotWeek := r.calculateWeeklyAggregates(snapshotSlice)
	return helper.SliceToChan(snapshotWeek), nil
}

// GetSince attempts to return a channel of snapshots for the asset with the given name since the given date.
func (r *FileSystemWeekRepository) GetSince(name string, date time.Time) (<-chan *Snapshot, error) {
	snapshots, err := r.Get(name)
	if err != nil {
		return nil, err
	}

	snapshots = helper.Filter(snapshots, func(s *Snapshot) bool {
		return s.Date.Equal(date) || s.Date.After(date)
	})

	return snapshots, nil
}

// LastDate returns the date of the last snapshot for the asset with the given name.
func (r *FileSystemWeekRepository) LastDate(name string) (time.Time, error) {
	var last time.Time

	snapshots, err := r.Get(name)
	if err != nil {
		return last, err
	}

	snapshot, ok := <-helper.Last(snapshots, 1)
	if !ok {
		return last, errors.New("empty asset")
	}

	return snapshot.Date, nil
}

// Append adds the given snapshows to the asset with the given name.
func (r *FileSystemWeekRepository) Append(name string, snapshots <-chan *Snapshot) error {
	return helper.AppendOrWriteToCsvFile(r.getCsvFileName(name), true, snapshots)
}

// getCsvFileName gets the CSV file name for the given asset name.
func (r *FileSystemWeekRepository) getCsvFileName(name string) string {
	return filepath.Join(r.base, fmt.Sprintf("%s.csv", name))
}

// calculateWeeklyAggregates groups snapshots by week and calculates weekly aggregates.
func (r *FileSystemWeekRepository) calculateWeeklyAggregates(snapshots []*Snapshot) []*Snapshot {
	if len(snapshots) == 0 {
		return nil
	}

	var weeklyData []*Snapshot
	var currentWeek []*Snapshot
	var weekStart time.Time

	for _, snapshot := range snapshots {
		// Determine the start of the week for the snapshot's date.
		currentWeekStart := snapshot.Date.Truncate(time.Hour * 24 * 7)

		if len(currentWeek) == 0 || currentWeekStart == weekStart {
			// Add snapshot to the current week.
			currentWeek = append(currentWeek, snapshot)
			weekStart = currentWeekStart
		} else {
			// Process the previous week and start a new one.
			weeklyData = append(weeklyData, r.processWeek(currentWeek))
			currentWeek = []*Snapshot{snapshot}
			weekStart = currentWeekStart
		}
	}

	// Process the last week.
	if len(currentWeek) > 0 {
		weeklyData = append(weeklyData, r.processWeek(currentWeek))
	}

	return weeklyData
}

// processWeek calculates the weekly aggregates from a slice of snapshot pointers.
func (r *FileSystemWeekRepository) processWeek(week []*Snapshot) *Snapshot {
	high := week[0].High
	low := week[0].Low
	open := week[0].Open
	close := week[len(week)-1].Close
	volume := 0.0

	for _, snapshot := range week {
		if snapshot.High > high {
			high = snapshot.High
		}
		if snapshot.Low < low {
			low = snapshot.Low
		}
		volume += snapshot.Volume
	}

	return &Snapshot{
		Date:   week[0].Date,
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: volume,
	}
}
