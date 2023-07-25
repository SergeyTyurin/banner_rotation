package database

import (
	"github.com/SergeyTyurin/banner-rotation/structures"
)

func (d *databaseImpl) DatabaseGetGroups() ([]structures.Group, error) {
	query := `SELECT id, info FROM "Groups"`
	rows, err := d.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]structures.Group, 0)
	for rows.Next() {
		var id int
		var info string
		if err := rows.Scan(&id, &info); err != nil {
			return nil, err
		}
		groups = append(groups, structures.Group{ID: id, Info: info})
	}
	return groups, nil
}

func (d *databaseImpl) DatabaseGetGroup(id int) (structures.Group, error) {
	query := `SELECT info FROM "Groups" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Group{ID: invalidID}, ErrNotExist
	}
	return structures.Group{ID: id, Info: info}, nil
}

func (d *databaseImpl) DatabaseDeleteGroup(id int) error {
	if err := checkEntityIsExists(d, "Groups", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE Group_id = $1`
	query := `DELETE FROM "Groups" WHERE id = $1`

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

func (d *databaseImpl) DatabaseCreateGroup(entity structures.Group) (structures.Group, error) {
	query := `INSERT INTO "Groups" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Group{ID: invalidID}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	row := tx.QueryRow(query, entity.Info)
	id := invalidID
	if err := row.Scan(&id); err != nil {
		return structures.Group{ID: invalidID}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Group{ID: invalidID}, err
	}
	return structures.Group{ID: id, Info: entity.Info}, nil
}

func (d *databaseImpl) DatabaseUpdateGroup(entity structures.Group) error {
	if err := checkEntityIsExists(d, "Groups", entity.ID); err != nil {
		return err
	}
	query := `UPDATE "Groups"
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
