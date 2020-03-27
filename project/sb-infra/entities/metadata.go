package entities

import "time"

// Metadata is common data for all objects
type Metadata struct {
	CDate     time.Time `json:"cdate"` // Creation
	MDate     time.Time `json:"mdate"` // Modification
	IsDeleted bool      `json:"isdeleted"`
}
