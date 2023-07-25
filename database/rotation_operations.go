package database

import (
	"database/sql"

	"github.com/SergeyTyurin/banner-rotation/bannerselector"
)

func checkEntityInRotationTx(tx *sql.Tx, bannerID, slotID, groupID int) (bool, error) {
	count := 0
	if err := tx.QueryRow(`SELECT count(*) 
	FROM "Statistic"
	WHERE slot_id=$1 AND banner_id=$2 AND group_id=$3`,
		slotID, bannerID, groupID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func increaseDisplay(d *databaseImpl, bannerID, slotID, groupID int) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `UPDATE "Statistic"
	SET click_count = click_count + 1
	WHERE slot_id=$1 AND group_id=$2 AND banner_id=$3`

	_, err = tx.Exec(query, slotID, groupID, bannerID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *databaseImpl) DatabaseAddToRotation(bannerID, slotID int) error { //nolint:stylecheck
	if err := checkEntityIsExists(d, "Banners", bannerID); err != nil {
		return err
	}
	if err := checkEntityIsExists(d, "Slots", slotID); err != nil {
		return err
	}

	groups, err := d.DatabaseGetGroups()
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO "Statistic"(banner_id, slot_id, group_id, display_count, click_count)
	VALUES($1, $2, $3, $4, $5)`

	for _, group := range groups {
		inRotation, err := checkEntityInRotationTx(tx, bannerID, slotID, group.ID)
		if err != nil {
			return err
		}
		if inRotation {
			return ErrAlreadyInRotation
		}
		_, err = tx.Exec(query, bannerID, slotID, group.ID, 0, 0)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *databaseImpl) DatabaseDeleteFromRotation(bannerID, slotID int) error {
	if err := checkEntityIsExists(d, "Banners", bannerID); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Slots", slotID); err != nil {
		return err
	}

	groups, err := d.DatabaseGetGroups()
	if err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, group := range groups {
		inRotation, err := checkEntityInRotationTx(tx, bannerID, slotID, group.ID)
		if err != nil {
			return err
		}
		if !inRotation {
			return ErrNotInRotation
		}
	}

	query := `DELETE FROM "Statistic" 
	WHERE banner_id=$1 AND slot_id=$2`

	_, err = tx.Exec(query, bannerID, slotID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *databaseImpl) DatabaseSelectFromRotation(slotID, groupID int) (bannerID int, err error) {
	if err := checkEntityIsExists(d, "Groups", groupID); err != nil {
		return invalidID, err
	}
	if err := checkEntityIsExists(d, "Slots", slotID); err != nil {
		return invalidID, err
	}

	query := `SELECT banner_id, display_count, click_count FROM "Statistic"
	WHERE slot_id=$1 AND group_id=$2`
	rows, err := d.db.Query(query, slotID, groupID)
	if err != nil || rows.Err() != nil {
		return invalidID, err
	}
	defer rows.Close()

	banners := make([]int, 0)
	displays := make([]int, 0)
	clicks := make([]int, 0)
	for rows.Next() {
		bannerID := int(0)
		displayCount := int(0)
		clickCount := int(0)
		if err := rows.Scan(&bannerID, &displayCount, &clickCount); err != nil {
			return invalidID, err
		}
		banners = append(banners, bannerID)
		displays = append(displays, displayCount)
		clicks = append(clicks, clickCount)
	}

	if len(banners) == 0 {
		return invalidID, ErrNotInRotation
	}

	bannerIndex, err := bannerselector.SelectBannerIndex(displays, clicks)
	if err != nil {
		return invalidID, err
	}
	if err := increaseDisplay(d, banners[bannerIndex], slotID, groupID); err != nil {
		return invalidID, err
	}

	return banners[bannerIndex], nil
}

func (d *databaseImpl) DatabaseRegisterTransition(slotID, bannerID, groupID int) error {
	if err := checkEntityIsExists(d, "Banners", bannerID); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Slots", slotID); err != nil {
		return err
	}

	if err := checkEntityIsExists(d, "Groups", groupID); err != nil {
		return err
	}

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	inRotation, err := checkEntityInRotationTx(tx, bannerID, slotID, groupID)
	if err != nil {
		return err
	}
	if !inRotation {
		return ErrNotInRotation
	}

	query := `UPDATE "Statistic"
	SET click_count = click_count + 1
	WHERE slot_id=$1 AND group_id=$2 AND banner_id=$3`

	_, err = tx.Exec(query, slotID, groupID, bannerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
