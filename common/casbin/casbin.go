package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:256;uniqueIndex:unique_index"`
	V0    string `gorm:"size:256;uniqueIndex:unique_index"`
	V1    string `gorm:"size:256;uniqueIndex:unique_index"`
	V2    string `gorm:"size:256;uniqueIndex:unique_index"`
	V3    string `gorm:"size:256;uniqueIndex:unique_index"`
	V4    string `gorm:"size:256;uniqueIndex:unique_index"`
	V5    string `gorm:"size:256;uniqueIndex:unique_index"`
}

func New(db *gorm.DB, conf string) *casbin.Enforcer {
	a, _ := gormadapter.NewAdapterByDBWithCustomTable(db, &CasbinRule{})
	e, _ := casbin.NewEnforcer(conf, a)
	return e
}
