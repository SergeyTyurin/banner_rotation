package database

import "github.com/SergeyTyurin/banner-rotation/structures"

func (d *databaseImpl) DatabaseGetSlots() ([]structures.Slot, error) {
	query := `SELECT id, info FROM "Slots"`
	rows, err := d.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	slots := make([]structures.Slot, 0)
	for rows.Next() {
		var id int
		var info string
		if err := rows.Scan(&id, &info); err != nil {
			return nil, err
		}
		slots = append(slots, structures.Slot{ID: id, Info: info})
	}
	return slots, nil
}

func (d *databaseImpl) DatabaseGetSlot(id int) (structures.Slot, error) {
	query := `SELECT info FROM "Slots" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Slot{ID: invalidID}, ErrNotExist
	}
	return structures.Slot{ID: id, Info: info}, nil
}

func (d *databaseImpl) DatabaseDeleteSlot(id int) error {
	if err := checkEntityIsExists(d, "Slots", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE Slot_id = $1`
	query := `DELETE FROM "Slots" WHERE id = $1`

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

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

	return tx.Commit()
}

func (d *databaseImpl) DatabaseCreateSlot(entity structures.Slot) (structures.Slot, error) {
	query := `INSERT INTO "Slots" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Slot{ID: invalidID}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	row := tx.QueryRow(query, entity.Info)
	id := invalidID
	if err := row.Scan(&id); err != nil {
		return structures.Slot{ID: invalidID}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Slot{ID: invalidID}, err
	}
	return structures.Slot{ID: id, Info: entity.Info}, nil
}

func (d *databaseImpl) DatabaseUpdateSlot(entity structures.Slot) error {
	if err := checkEntityIsExists(d, "Slots", entity.ID); err != nil {
		return err
	}
	query := `UPDATE "Slots"
	SET info = $1
	WHERE id = $2`
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.Exec(query, entity.Info, entity.ID)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected < 1 {
		return ErrNotExist
	}

	return tx.Commit()
}
