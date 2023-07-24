package database

import "github.com/SergeyTyurin/banner_rotation/structures"

func (d *databaseImpl) GetSlots() ([]structures.Slot, error) {
	query := `SELECT id, info FROM "Slots"`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := make([]structures.Slot, 0)
	for rows.Next() {
		var id int
		var info string
		rows.Scan(&id, &info)
		slots = append(slots, structures.Slot{Id: id, Info: info})
	}
	return slots, nil
}

func (d *databaseImpl) GetSlot(id int) (structures.Slot, error) {
	query := `SELECT info FROM "Slots" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Slot{Id: invalidId}, ErrNotExist
	}
	return structures.Slot{Id: id, Info: info}, nil
}

func (d *databaseImpl) DeleteSlot(id int) error {
	if err := checkEntityIsExists(d, "Slots", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE Slot_id = $1`
	query := `DELETE FROM "Slots" WHERE id = $1`

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(rotationQuery, id)
	if err != nil {
		return err
	}

	res, err := tx.Exec(query, id)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected < 1 {
		return ErrNotExist
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *databaseImpl) CreateSlot(entity structures.Slot) (structures.Slot, error) {
	query := `INSERT INTO "Slots" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Slot{Id: invalidId}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(query, entity.Info)
	id := invalidId
	if err := row.Scan(&id); err != nil {
		return structures.Slot{Id: invalidId}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Slot{Id: invalidId}, err
	}
	return structures.Slot{Id: id, Info: entity.Info}, nil
}

func (d *databaseImpl) UpdateSlot(entity structures.Slot) error {
	if err := checkEntityIsExists(d, "Slots", entity.Id); err != nil {
		return err
	}
	query := `UPDATE "Slots"
	SET info = $1
	WHERE id = $2`
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(query, entity.Info, entity.Id)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected < 1 {
		return ErrNotExist
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
