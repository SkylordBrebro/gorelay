package models

import "math"

// MoveRecord represents a single movement record
type MoveRecord struct {
	Time int64   `json:"time"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}

// MoveRecords manages a collection of movement records
type MoveRecords struct {
	LastClearTime int64
	Records       []MoveRecord
}

// NewMoveRecords creates a new move records manager
func NewMoveRecords() *MoveRecords {
	return &MoveRecords{
		LastClearTime: -1,
		Records:       make([]MoveRecord, 0),
	}
}

// AddRecord adds a new movement record
func (mr *MoveRecords) AddRecord(time int64, x, y float32) {
	if mr.LastClearTime < 0 {
		return
	}

	id := mr.getID(time)
	if id < 1 || id > 10 {
		return
	}

	if len(mr.Records) == 0 {
		mr.Records = append(mr.Records, MoveRecord{
			Time: time,
			X:    x,
			Y:    y,
		})
		return
	}

	currentRecord := &mr.Records[len(mr.Records)-1]
	currentID := mr.getID(currentRecord.Time)

	if id != currentID {
		mr.Records = append(mr.Records, MoveRecord{
			Time: time,
			X:    x,
			Y:    y,
		})
		return
	}

	score := mr.getScore(id, time)
	currentScore := mr.getScore(currentID, currentRecord.Time)

	if score < currentScore {
		currentRecord.Time = time
		currentRecord.X = x
		currentRecord.Y = y
	}
}

// Clear clears all records and sets a new clear time
func (mr *MoveRecords) Clear(time int64) {
	mr.Records = mr.Records[:0]
	mr.LastClearTime = time
}

// getID calculates the ID for a given time
func (mr *MoveRecords) getID(time int64) int64 {
	return int64(math.Round(float64(time-mr.LastClearTime+50) / 100))
}

// getScore calculates the score for a given ID and time
func (mr *MoveRecords) getScore(id, time int64) int64 {
	return int64(math.Round(math.Abs(float64(time - mr.LastClearTime - id*100))))
}
