package database

import "github.com/SergeyTyurin/banner_rotation/structures"

func (d *databaseImpl) GetGroups() ([]structures.Group, error) {
	query := `SELECT id, info FROM "Groups"`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]structures.Group, 0)
	for rows.Next() {
		var id int
		var info string
		rows.Scan(&id, &info)
		groups = append(groups, structures.Group{Id: id, Info: info})
	}
	return groups, nil
}

func (d *databaseImpl) GetGroup(id int) (structures.Group, error) {
	query := `SELECT info FROM "Groups" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Group{Id: invalidId}, ErrNotExist
	}
	return structures.Group{Id: id, Info: info}, nil
}

func (d *databaseImpl) DeleteGroup(id int) error {
	if err := checkEntityIsExists(d, "Groups", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE Group_id = $1`
	query := `DELETE FROM "Groups" WHERE id = $1`

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

func (d *databaseImpl) CreateGroup(entity structures.Group) (structures.Group, error) {
	query := `INSERT INTO "Groups" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Group{Id: invalidId}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(query, entity.Info)
	id := invalidId
	if err := row.Scan(&id); err != nil {
		return structures.Group{Id: invalidId}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Group{Id: invalidId}, err
	}
	return structures.Group{Id: id, Info: entity.Info}, nil
}

func (d *databaseImpl) UpdateGroup(entity structures.Group) error {
	if err := checkEntityIsExists(d, "Groups", entity.Id); err != nil {
		return err
	}
	query := `UPDATE "Groups"
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
