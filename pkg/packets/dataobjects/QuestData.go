package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
	"strings"
)

// QuestData represents quest information
type QuestData struct {
	ID           string
	Name         string
	Description  string
	Expiration   string
	Requirements []int32
	Rewards      []int32
	Completed    bool
	ItemOfChoice bool
	Repeatable   bool
	Category     int32
	Weight       int32
}

// NewQuestData creates a new QuestData instance
func NewQuestData() *QuestData {
	return &QuestData{
		Requirements: make([]int32, 0),
		Rewards:      make([]int32, 0),
	}
}

// Read reads the quest data from a Reader
func (q *QuestData) Read(r interfaces.Reader) error {
	var err error
	q.ID, err = r.ReadString()
	if err != nil {
		return err
	}
	q.Name, err = r.ReadString()
	if err != nil {
		return err
	}
	q.Description, err = r.ReadString()
	if err != nil {
		return err
	}
	q.Expiration, err = r.ReadString()
	if err != nil {
		return err
	}
	q.Weight, err = r.ReadInt32()
	if err != nil {
		return err
	}
	q.Category, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Requirements
	reqCount, err := r.ReadInt16()
	if err != nil {
		return err
	}
	q.Requirements = make([]int32, reqCount)
	for i := 0; i < int(reqCount); i++ {
		q.Requirements[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	// Read Rewards
	rewCount, err := r.ReadInt16()
	if err != nil {
		return err
	}
	q.Rewards = make([]int32, rewCount)
	for i := 0; i < int(rewCount); i++ {
		q.Rewards[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	q.Completed, err = r.ReadBool()
	if err != nil {
		return err
	}
	q.ItemOfChoice, err = r.ReadBool()
	if err != nil {
		return err
	}
	q.Repeatable, err = r.ReadBool()
	return err
}

// Write writes the quest data to a Writer
func (q *QuestData) Write(w interfaces.Writer) error {
	if err := w.WriteString(q.ID); err != nil {
		return err
	}
	if err := w.WriteString(q.Name); err != nil {
		return err
	}
	if err := w.WriteString(q.Description); err != nil {
		return err
	}
	if err := w.WriteString(q.Expiration); err != nil {
		return err
	}
	if err := w.WriteInt32(q.Weight); err != nil {
		return err
	}
	if err := w.WriteInt32(q.Category); err != nil {
		return err
	}

	// Write Requirements
	if err := w.WriteInt16(int16(len(q.Requirements))); err != nil {
		return err
	}
	for _, req := range q.Requirements {
		if err := w.WriteInt32(req); err != nil {
			return err
		}
	}

	// Write Rewards
	if err := w.WriteInt16(int16(len(q.Rewards))); err != nil {
		return err
	}
	for _, rew := range q.Rewards {
		if err := w.WriteInt32(rew); err != nil {
			return err
		}
	}

	if err := w.WriteBool(q.Completed); err != nil {
		return err
	}
	if err := w.WriteBool(q.ItemOfChoice); err != nil {
		return err
	}
	return w.WriteBool(q.Repeatable)
}

// Clone creates a copy of the QuestData
func (q *QuestData) Clone() DataObject {
	requirements := make([]int32, len(q.Requirements))
	copy(requirements, q.Requirements)

	rewards := make([]int32, len(q.Rewards))
	copy(rewards, q.Rewards)

	return &QuestData{
		ID:           q.ID,
		Name:         q.Name,
		Description:  q.Description,
		Expiration:   q.Expiration,
		Requirements: requirements,
		Rewards:      rewards,
		Completed:    q.Completed,
		ItemOfChoice: q.ItemOfChoice,
		Category:     q.Category,
		Repeatable:   q.Repeatable,
		Weight:       q.Weight,
	}
}

// String returns a string representation of the QuestData
func (q *QuestData) String() string {
	reqs := make([]string, len(q.Requirements))
	for i, r := range q.Requirements {
		reqs[i] = fmt.Sprintf("%d", r)
	}

	rews := make([]string, len(q.Rewards))
	for i, r := range q.Rewards {
		rews[i] = fmt.Sprintf("%d", r)
	}

	return fmt.Sprintf("{ Id=%s, Name=%s, Description=%s, Requirements=%s, Rewards=%s, Completed=%v, ItemOfChoice=%v, Category=%d, Repeatable=%v }",
		q.ID, q.Name, q.Description, strings.Join(reqs, " "), strings.Join(rews, " "),
		q.Completed, q.ItemOfChoice, q.Category, q.Repeatable)
}

// ToStringMinified returns a shorter string representation of the QuestData
func (q *QuestData) ToStringMinified() string {
	return fmt.Sprintf("{ Name=%s, Id=%s }", q.Name, q.ID)
}
