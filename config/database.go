package config

import "fmt"

type Database struct {
	Host     string
	Name     string
	User     string
	Password string
	Sslmode  string
}

func (d *Database) DSN() (dsn string) {
	dsn = fmt.Sprintf("user=%s dbname=%s sslmode=%s", d.User, d.Name, d.Sslmode)
	if len(d.Host) > 0 {
		dsn = dsn + " host=" + d.Host
	}
	if len(d.Password) > 0 {
		dsn = dsn + " password=" + d.Password
	}
	return
}

func (d *Database) Defaults() {
	d.Sslmode = "disable"
}
