package v2

import (
	"fmt"
	"strings"
)

// router
var (
	DBServiceRouterGroup = "/dms/projects/:project_uid/db_services"
	ProjectRouterGroup   = "/dms/projects"
)

// api group
var (
	GroupV2             = "/v2"
	CurrentGroupVersion = GroupV2
)

func GetDBServiceRouter(projectUid string) string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, strings.Replace(DBServiceRouterGroup, ":project_uid", projectUid, 1))
}

func GetProjectsRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, ProjectRouterGroup)
}
