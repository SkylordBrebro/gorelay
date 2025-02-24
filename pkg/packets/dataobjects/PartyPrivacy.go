package dataobjects

// PartyPrivacy represents the privacy setting of a party
type PartyPrivacy byte

const (
	// Public represents a public party
	Public PartyPrivacy = iota
	// Private represents a private party
	Private
)
