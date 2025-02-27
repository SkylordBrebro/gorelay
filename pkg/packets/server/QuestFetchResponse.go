package server

import (
	"fmt"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
	"strings"
)

// QuestFetchResponse represents the server packet for quest fetch responses
type QuestFetchResponse struct {
	Quests           []*dataobjects.QuestData
	NextRefreshPrice int16
}

// Type returns the packet type for QuestFetchResponse
func (p *QuestFetchResponse) Type() interfaces.PacketType {
	return interfaces.QuestFetchResponse
}

// Read reads the packet data from the provided reader
func (p *QuestFetchResponse) Read(r interfaces.Reader) error {
	var err error

	// Read Quests array length
	questsCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Quests array
	p.Quests = make([]*dataobjects.QuestData, questsCount)
	for i := 0; i < int(questsCount); i++ {
		p.Quests[i] = dataobjects.NewQuestData()
		err = p.Quests[i].Read(r)
		if err != nil {
			return err
		}
	}

	// Read NextRefreshPrice
	p.NextRefreshPrice, err = r.ReadInt16()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *QuestFetchResponse) Write(w interfaces.Writer) error {
	var err error

	// Write Quests array length
	err = w.WriteInt16(int16(len(p.Quests)))
	if err != nil {
		return err
	}

	// Write Quests array
	for _, quest := range p.Quests {
		err = quest.Write(w)
		if err != nil {
			return err
		}
	}

	// Write NextRefreshPrice
	err = w.WriteInt16(p.NextRefreshPrice)
	if err != nil {
		return err
	}

	return nil
}

// String returns a string representation of the QuestFetchResponse
func (p *QuestFetchResponse) String() string {
	questStrings := make([]string, len(p.Quests))
	for i, quest := range p.Quests {
		questStrings[i] = quest.String()
	}

	return fmt.Sprintf("QuestFetchResponse (NextRefreshPrice = %d, Quests = %s",
		p.NextRefreshPrice, strings.Join(questStrings, "\n"))
}

func (p *QuestFetchResponse) ID() int32 {
	return int32(interfaces.QuestFetchResponse)
}