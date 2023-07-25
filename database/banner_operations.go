package database

import (
	"github.com/SergeyTyurin/banner-rotation/structures"
)

func (d *databaseImpl) DatabaseGetBanners() ([]structures.Banner, error) {
	query := `SELECT id, info FROM "Banners"`
	rows, err := d.db.Query(query)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make([]structures.Banner, 0)
	for rows.Next() {
		var id int
		var info string
		err := rows.Scan(&id, &info)
		if err != nil {
			return nil, err
		}
		banners = append(banners, structures.Banner{ID: id, Info: info})
	}
	return banners, nil
}

func (d *databaseImpl) DatabaseGetBanner(id int) (structures.Banner, error) {
	query := `SELECT info FROM "Banners" WHERE id = $1`
	row := d.db.QueryRow(query, id)

	var info string
	if err := row.Scan(&info); err != nil {
		return structures.Banner{ID: invalidID}, ErrNotExist
	}
	return structures.Banner{ID: id, Info: info}, nil
}

func (d *databaseImpl) DatabaseDeleteBanner(id int) error {
	if err := checkEntityIsExists(d, "Banners", id); err != nil {
		return err
	}
	rotationQuery := `DELETE FROM "Statistic" WHERE banner_id = $1`
	query := `DELETE FROM "Banners" WHERE id = $1`

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

func (d *databaseImpl) DatabaseCreateBanner(entity structures.Banner) (structures.Banner, error) {
	query := `INSERT INTO "Banners" (info) VALUES($1)
	RETURNING id`
	tx, err := d.db.Begin()
	if err != nil {
		return structures.Banner{ID: invalidID}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	row := tx.QueryRow(query, entity.Info)
	id := invalidID
	if err := row.Scan(&id); err != nil {
		return structures.Banner{ID: invalidID}, err
	}

	if err := tx.Commit(); err != nil {
		return structures.Banner{ID: invalidID}, err
	}
	return structures.Banner{ID: id, Info: entity.Info}, nil
}

func (d *databaseImpl) DatabaseUpdateBanner(entity structures.Banner) error {
	if err := checkEntityIsExists(d, "Banners", entity.ID); err != nil {
		return err
	}
	query := `UPDATE "Banners"
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
