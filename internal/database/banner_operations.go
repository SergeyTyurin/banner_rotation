package database

import "github.com/SergeyTyurin/banner_rotation/structures"

func (d *databaseImpl) GetBanners() ([]structures.Banner, error) {
	query := `SELECT id, info FROM "Banners"`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make([]structures.Banner, 0)
	for rows.Next() {
		var id int
		var info string
		rows.Scan(&id, &info)
		banners = append(banners, structures.Banner{Id: id, Info: info})
	}
	return banners, nil
}

func (d *databaseImpl) GetBanner(id int) (structures.Banner, error) {
	query := `SELECT info FROM "Banners" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Banner{Id: invalidId}, ErrNotExist
	}
	return structures.Banner{Id: id, Info: info}, nil
}

func (d *databaseImpl) DeleteBanner(id int) error {
	if err := checkEntityIsExists(d, "Banners", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE banner_id = $1`
	query := `DELETE FROM "Banners" WHERE id = $1`

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

func (d *databaseImpl) CreateBanner(entity structures.Banner) (structures.Banner, error) {
	query := `INSERT INTO "Banners" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Banner{Id: invalidId}, err
	}
	defer tx.Rollback()

	row := tx.QueryRow(query, entity.Info)
	id := invalidId
	if err := row.Scan(&id); err != nil {
		return structures.Banner{Id: invalidId}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Banner{Id: invalidId}, err
	}
	return structures.Banner{Id: id, Info: entity.Info}, nil
}

func (d *databaseImpl) UpdateBanner(entity structures.Banner) error {
	if err := checkEntityIsExists(d, "Banners", entity.Id); err != nil {
		return err
	}
	query := `UPDATE "Banners"
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
