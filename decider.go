package paxos

import (
	"math/rand"
	"time"
)

// Decider ...
type Decider struct {
	IDs          []uint32
	round        uint32
	Votes        map[uint32]bool
	VoteCount    int
	Acknowledged map[uint32]bool
	AckCount     int
	VotedID      uint32
	Voted        bool
	Finished     bool
}

// NewDecider ...
func NewDecider() *Decider {
	D := &Decider{}
	D.reset()

	return D
}

// SetIDs ...
func (D *Decider) SetIDs(IDs []uint32) {
	D.IDs = IDs
	D.reset()
}

// NextRound ...
func (D Decider) NextRound() uint32 {
	return D.round + 1
}

func (D *Decider) reset() {
	D.Votes = make(map[uint32]bool)
	D.VoteCount = 0

	D.Acknowledged = make(map[uint32]bool)
	D.AckCount = 0

	D.Finished = false
}

// Check ...
func (D Decider) Check(ID uint32, M Message) uint32 {
	if D.VotedID == ID && D.round == M.Round {
		return OK
	}

	return NotOK
}

// VoteResponse ...
func (D *Decider) VoteResponse(ID uint32, M Message) uint32 {
	if D.round > M.Round {
		return No
	}

	if D.round < M.Round || !D.Voted {
		D.round = M.Round
		D.VotedID = ID

		D.Voted = true
	}

	if D.VotedID == ID {
		return Yes
	}

	return No
}

// IsElected ...
func (D *Decider) IsElected(ID uint32, M Message) bool {
	if D.round < M.Round {
		D.reset()

		return false
	}

	if M.Payload == Yes {
		if !D.Votes[ID] {
			D.VoteCount++
		}
	}

	Half := len(D.IDs) / 2

	return D.VoteCount > Half
}

// ConsensusResult this needs to
// be threadsafe
func (D Decider) ConsensusResult() (uint32, bool) {
	return D.VotedID, D.Finished
}

// Acknowledge ...
func (D *Decider) Acknowledge(ID uint32, M Message) bool {
	if D.round < M.Round {
		D.reset()

		return false
	}

	if M.Payload == OK {
		if !D.Acknowledged[ID] {
			D.AckCount++
		}
	}

	Half := len(D.IDs) / 2

	return D.AckCount > Half
}

// Confirmed ...
func (D *Decider) Confirmed(ID uint32) {
	D.Finished = true
	D.VotedID = ID
}

// SleepDuration ...
func (Decider) SleepDuration() time.Duration {
	Mag := 700 + rand.Intn(600)

	return time.Duration(Mag) * time.Millisecond
}

// Round ...
func (D Decider) Round() uint32 {
	return D.round
}
