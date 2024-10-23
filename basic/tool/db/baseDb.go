package dbtool

import "gorm.io/gorm"

type BaseDb struct {
	*gorm.DB
}
