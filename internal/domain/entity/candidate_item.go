package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

// CandidateItem represents a product in the candidate pool
type CandidateItem struct {
	ID               int64          `json:"id" db:"id"`
	GroupKey         string         `json:"groupKey" db:"group_key"`
	Region           string         `json:"region" db:"region"`
	TitleVotes       map[string]int `json:"titleVotes" db:"title_votes"`
	TotalOccurrences int            `json:"totalOccurrences" db:"total_occurrences"`
	LastPrice        float64        `json:"lastPrice" db:"last_price"`
	LastStatus       int            `json:"lastStatus" db:"last_status"`
	FirstSeenTime    time.Time      `json:"firstSeenTime" db:"first_seen_time"`
	LastSeenTime     time.Time      `json:"lastSeenTime" db:"last_seen_time"`
	CreateTime       time.Time      `json:"createTime" db:"create_time"`
	UpdateTime       time.Time      `json:"updateTime" db:"update_time"`
}

// TitleVotesDB is the database representation of title_votes (JSONB)
type TitleVotesDB map[string]int

// Scan implements sql.Scanner for TitleVotesDB
func (t *TitleVotesDB) Scan(value interface{}) error {
	if value == nil {
		*t = make(map[string]int)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into TitleVotesDB", value)
	}
	if err := json.Unmarshal(bytes, t); err != nil {
		return fmt.Errorf("failed to unmarshal TitleVotesDB: %w", err)
	}
	return nil
}

// MarshalJSON implements json.Marshaler for TitleVotesDB
func (t TitleVotesDB) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string]int(t))
}

// AddTitleVote adds a vote for a title
func (c *CandidateItem) AddTitleVote(title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if c.TitleVotes == nil {
		c.TitleVotes = make(map[string]int)
	}
	c.TitleVotes[title]++
	c.TotalOccurrences++
	return nil
}

// UpdateLastSeen updates the last seen time
func (c *CandidateItem) UpdateLastSeen(price float64, status int) {
	c.LastPrice = price
	c.LastStatus = status
	c.LastSeenTime = time.Now()
}

// ShouldPromote checks if the candidate should be promoted to master
func (c *CandidateItem) ShouldPromote(threshold int) bool {
	return c.TotalOccurrences >= threshold
}

// ElectStandardTitle returns the title with the most votes
func (c *CandidateItem) ElectStandardTitle() string {
	if len(c.TitleVotes) == 0 {
		return ""
	}

	var winner string
	maxVotes := 0

	for title, votes := range c.TitleVotes {
		// Skip empty titles and invalid votes
		if title == "" || votes <= 0 {
			continue
		}
		if votes > maxVotes {
			maxVotes = votes
			winner = title
		} else if votes == maxVotes && len(title) > len(winner) {
			// Tie-breaker: prefer longer title
			winner = title
		}
	}

	return winner
}
