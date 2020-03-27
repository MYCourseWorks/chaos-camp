package persistance

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// ErrOddCreate comment
const ErrOddCreate = "Can't create an odd from entity"

const findOddByIDQuery = `
select o.ID, o.Source, o.Values from Odd as o
where o.IsDeleted = 0 and o.ID = ?
`

// FindOddByID comment
func (dl *DataLayer) FindOddByID(id int64) (*entities.Odd, error) {
	row := dl.db.Connection().QueryRow(findOddByIDQuery, id)
	return extratOddFromRow(row)
}

const findOddQuery = `
select o.ID, o.Source, o.Values from Odd as o
where o.IsDeleted = 0 and o.Source = ? and o.LineID = ?
`

// FindOdd comment
func (dl *DataLayer) FindOdd(source string, lineID int64) (*entities.Odd, error) {
	row := dl.db.Connection().QueryRow(findOddQuery, source, lineID)
	return extratOddFromRow(row)
}

func extratOddFromRow(row *sql.Row) (*entities.Odd, error) {
	var id int64
	var source string
	var valuesStr string
	err := row.Scan(&id, &source, &valuesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	values, err := parseOdds(valuesStr)
	if err != nil {
		return nil, err
	}

	odd := entities.NewOdd(id, source, values, nil)
	return odd, nil
}

func parseOdds(input string) ([]float32, error) {
	split := strings.Split(input, ":")
	values := make([]float32, len(split))
	for i, s := range split {
		f64, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		values[i] = float32(f64)
	}

	return values, nil
}

const createOddQuery = "insert into Odd (Source, `Values`, LineID) values (?, ?, ?)"

// CreateOdd creates an odd and sets the ID.
// NOTE: assumes that the related line exists!
func (dl *DataLayer) CreateOdd(odd *entities.Odd) error {
	if odd == nil || odd.Source == "" || len(odd.Source) > 254 {
		return fmt.Errorf(ErrOddCreate)
	}
	if odd.ID > 0 {
		infra.Warn(infra.WarningColor, "Odd might already exist in DB")
	}

	var sb strings.Builder
	vLen := len(odd.Values)
	for i := 0; i < vLen-1; i++ {
		o := odd.Values[i]
		sb.WriteString(fmt.Sprintf("%f:", o))
	}
	sb.WriteString(fmt.Sprintf("%f", odd.Values[vLen-1]))

	result, err := dl.db.Connection().Exec(createOddQuery, odd.Source, sb.String(), odd.Line.ID)
	if err != nil {
		return err
	}

	odd.ID, err = result.LastInsertId()
	return err
}

const deleteOddQuery = `
	update Odd
	set IsDeleted = 1
	where ID = ?
`

// DeleteOdd soft delete from database
func (dl *DataLayer) DeleteOdd(id int64) error {
	_, err := dl.db.Connection().Exec(deleteOddQuery, id)
	return err
}
