package persistance

import (
	"database/sql"
	"fmt"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrLineCreate comment
const ErrLineCreate = "Can't create a line from entity"

const findLineByIDQuery = `
select l.ID, l.OddFormat, l.LineType, l.Description from Line as l
where l.IsDeleted = 0 and l.ID = ?
`

// FindLineByID comment
func (dl *DataLayer) FindLineByID(id int64) (*entities.Line, error) {
	row := dl.db.Connection().QueryRow(findLineByIDQuery, id)
	return extratLineFromRow(row)
}

const findLineQuery = `
select l.ID, l.OddFormat, l.LineType, l.Description from Line as l
where l.IsDeleted = 0 and l.Description = ? and l.EventID = ?
`

// FindLine comment
func (dl *DataLayer) FindLine(description string, eventID int64) (*entities.Line, error) {
	row := dl.db.Connection().QueryRow(findLineQuery, description, eventID)
	return extratLineFromRow(row)
}

func extratLineFromRow(row *sql.Row) (*entities.Line, error) {
	var id int64
	var oddFormat int // NOTE: we ignore and use the default for now.
	var lineType int
	var description string
	err := row.Scan(&id, &oddFormat, &lineType, &description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	line := entities.NewLine(id, entities.LineType(lineType), description, nil)

	return line, nil
}

const createLineQuery = `
	insert into Line (OddFormat, LineType, Description, EventID) values (?, ?, ?, ?)
`

// CreateLine creates a line and sets the ID.
// NOTE: Assumes that the related event exists!
func (dl *DataLayer) CreateLine(line *entities.Line) error {
	if line == nil || line.Description == "" || len(line.Description) > 254 {
		return fmt.Errorf(ErrLineCreate)
	}
	if line.ID > 0 {
		infra.Warn(infra.WarningColor, "Line might already exist in DB")
	}

	result, err := dl.db.Connection().
		Exec(createLineQuery, line.OddFormat, line.LineType, line.Description, line.Event.ID)

	if err != nil {
		return err
	}

	line.ID, err = result.LastInsertId()
	return err
}

const deleteLineQuery = `
	update Line
	set IsDeleted = 1
	where ID = ?
`

// DeleteLine soft delete from database
func (dl *DataLayer) DeleteLine(id int64) error {
	_, err := dl.db.Connection().Exec(deleteLineQuery, id)
	return err
}
