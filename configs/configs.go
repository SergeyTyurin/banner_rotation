package configs

import (
	"errors"
)

var errInputIsNil = errors.New("input is nil")

// func (c *dbConnectionImpl) DSN() string {
// 	format := `host=%s port=%d user=%s password=%s dbname=%s`
// 	return fmt.Sprintf(format, c.HostDB, c.PortDB, c.UserDB, c.PasswordDB, c.NameDB)
// }
