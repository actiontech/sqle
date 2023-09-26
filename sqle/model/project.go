package model

// import (
// 	"database/sql"
// 	"fmt"

// 	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
// 	"github.com/actiontech/sqle/sqle/errors"
// 	"github.com/actiontech/sqle/sqle/utils"

// 	"github.com/jinzhu/gorm"
// )

const ProjectIdForGlobalRuleTemplate = "0"

type ProjectUID string

const (
	ProjectStatusArchived = "archived"
	ProjectStatusActive   = "active"
)

