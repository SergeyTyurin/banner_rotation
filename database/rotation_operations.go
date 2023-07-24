package database

import (
	"database/sql"
)

func checkEntityInRotationTx(tx *sql.Tx, bannerId, slotId, groupId int) (bool, error) {
	count := 0
	if err := tx.QueryRow(`SELECT count(*) 
	FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3`,
		slotId, bannerId, groupId).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func increaseDisplay(d *databaseImpl, bannerId, slotId, groupId int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE "Statistic"
	SET click_count = click_count + 1
	WHERE slot_id=$1 AND group_id=$2 AND banner_id=$3`

	_, err = tx.Exec(query, slotId, groupId, bannerId)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *databaseImpl) AddToRotation(bannerId, slotId int) error {
	if err := checkEntityIsExists(d, "Banners", bannerId); err != nil {
		return err
	}
	if err := checkEntityIsExists(d, "Slots", slotId); err != nil {
		return err
	}

	groups, err := d.GetGroups()
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO "Statistic"(banner_id, slot_id, group_id, display_count, click_count)
	VALUES($1, $2, $3, $4, $5)`

	for _, group := range groups {
		inRotation, err := checkEntityInRotationTx(tx, bannerId, slotId, group.Id)
		if err != nil {
			return err
		}
		if inRotation {
			return ErrAlreadyInRotation
		}
		_, err = tx.Exec(query, bannerId, slotId, group.Id, 0, 0)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *databaseImpl) DeleteFromRotation(bannerId, slotId int) error {
	if err := checkEntityIsExists(d, "Banners", bannerId); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Slots", slotId); err != nil {
		return err
	}

	groups, err := d.GetGroups()
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, group := range groups {
		inRotation, err := checkEntityInRotationTx(tx, bannerId, slotId, group.Id)
		if err != nil {
			return err
		}
		if !inRotation {
			return ErrNotInRotation
		}
	}

	query := `DELETE FROM "Statistic" 
	WHERE banner_id=$1 AND slot_id=$2`

	_, err = tx.Exec(query, bannerId, slotId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *databaseImpl) SelectFromRotation(slotId, groupId int) (bannerId int, err error) {
	if err := checkEntityIsExists(d, "Groups", groupId); err != nil {
		return invalidId, err
	}
	if err := checkEntityIsExists(d, "Slots", slotId); err != nil {
		return invalidId, err
	}

	query := `SELECT banner_id, display_count, click_count FROM "Statistic"
	WHERE slot_id=$1 AND group_id=$2`
	rows, err := d.db.Query(query, slotId, groupId)
	if err != nil {
		return invalidId, err
	}
	defer rows.Close()

	banners := make([]int, 0)
	displays := make([]int, 0)
	clicks := make([]int, 0)
	for rows.Next() {
		bannerId := int(0)
		displayCount := int(0)
		clickCount := int(0)
		if err := rows.Scan(&bannerId, &displayCount, &clickCount); err != nil {
			return invalidId, err
		}
		banners = append(banners, bannerId)
		displays = append(displays, displayCount)
		clicks = append(clicks, clickCount)
	}

	if len(banners) == 0 {
		return invalidId, ErrNotInRotation
	}

	bannerIndex, err := banner_selector.SelectBannerIndex(displays, clicks)
	if err != nil {
		return invalidId, err
	}
	if err := increaseDisplay(d, banners[bannerIndex], slotId, groupId); err != nil {
		return invalidId, err
	}

	return banners[bannerIndex], nil
}

func (d *databaseImpl) RegisterTransition(slotId, bannerId, groupId int) error {
	if err := checkEntityIsExists(d, "Banners", bannerId); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Slots", slotId); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Groups", groupId); err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	inRotation, err := checkEntityInRotationTx(tx, bannerId, slotId, groupId)
	if err != nil {
		return err
	}
	if !inRotation {
		return ErrNotInRotation
	}

	query := `UPDATE "Statistic"
	SET click_count = click_count + 1
	WHERE slot_id=$1 AND group_id=$2 AND banner_id=$3`

	_, err = tx.Exec(query, slotId, groupId, bannerId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
