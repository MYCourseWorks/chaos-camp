package persistance

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrEventCreation comment
const ErrEventCreation = "Can't create a event from entity"

const findEventByIDQuery = `
select e.ID, e.Date, e.EventType from GameEvent as e
where e.IsDeleted = 0 and e.ID = ?
`

// FindEventByID comment
func (dl *DataLayer) FindEventByID(id int64) (*entities.GameEvent, error) {
	row := dl.db.Connection().QueryRow(findEventByIDQuery, id)
	return extractEventFromRow(row, dl.db.GetDateTimeFormat())
}

const findEventQuery = `
select e.ID, e.Date, e.EventType from GameEvent as e
where
	e.IsDeleted = 0 and
    DATE_FORMAT(e.Date,'%Y-%m-%d %H') = DATE_FORMAT(TIMESTAMP(?), '%Y-%m-%d %H')
    and e.EventType = ?
    and e.RelatedGameID = ?
`

// FindEvent comment
func (dl *DataLayer) FindEvent(date *time.Time,
	eventType entities.EventType,
	relatedGameID int64) (*entities.GameEvent, error) {

	row := dl.db.Connection().QueryRow(findEventQuery, date, eventType, relatedGameID)
	return extractEventFromRow(row, dl.db.GetDateTimeFormat())
}

func extractEventFromRow(row *sql.Row, dateFormat string) (*entities.GameEvent, error) {
	var id int64
	var dateStr string
	var eventType int
	err := row.Scan(&id, &dateStr, &eventType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return nil, err
	}

	event := entities.NewGameEvent(id, date, nil, entities.EventType(eventType))

	return event, nil
}

const createEventQuery = `
	insert into GameEvent (Date, EventType, RelatedGameID) values (?, ?, ?)
`

// CreateEvent creates an event and sets the ID.
// NOTE: Assumes that the related game exists!
func (dl *DataLayer) CreateEvent(event *entities.GameEvent) error {
	if event == nil {
		return fmt.Errorf(ErrEventCreation)
	}
	if event.ID > 0 {
		infra.Warn(infra.WarningColor, "GameEvent might already exist in DB")
	}

	result, err := dl.db.Connection().
		Exec(createEventQuery, event.Date, event.EventType, event.RelatedGame.ID)
	if err != nil {
		return err
	}

	event.ID, err = result.LastInsertId()
	return err
}

const deleteEventQuery = `
	update GameEvent
	set IsDeleted = 1
	where ID = ?
`

// DeleteEvent soft delete from database
func (dl *DataLayer) DeleteEvent(id int64) error {
	_, err := dl.db.Connection().Exec(deleteEventQuery, id)
	return err
}
