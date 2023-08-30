package v1

import (
	"fmt"
	"strings"
)

// login config
var (
	JwtSigningKey = []byte("actiontech dms secret")
)

// router
var (
	SessionRouterGroup            = "/dms/sessions"
	UserRouterGroup               = "/dms/users"
	DBServiceRouterGroup          = "/dms/db_services"
	ProxyRouterGroup              = "/dms/proxys"
	PluginRouterGroup             = "/dms/plugins"
	MemberRouterGroup             = "/dms/members"
	NamespaceRouterGroup          = "/dms/namespaces"
	NotificationRouterGroup       = "/dms/notifications"
	WebHookRouterGroup            = "/dms/webhooks"
	MemberForInternalRouterSuffix = "/internal"
	InternalDBServiceRouterGroup  = "/internal/db_services"
)

// api group
var (
	GroupV1             = "/v1"
	CurrentGroupVersion = GroupV1
)

func GetUserOpPermissionRouter(userUid string) string {
	return fmt.Sprintf("%s%s/%s/op_permission", CurrentGroupVersion, UserRouterGroup, userUid)
}

func GetUserOpPermissionRouterWithoutPrefix(userUid string) string {
	router := GetUserOpPermissionRouter(userUid)
	return strings.TrimPrefix(strings.TrimPrefix(router, CurrentGroupVersion), UserRouterGroup)
}

func GetDBServiceRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, DBServiceRouterGroup)
}

func GetUserRouter(userUid string) string {
	return fmt.Sprintf("%s%s/%s", CurrentGroupVersion, UserRouterGroup, userUid)
}

func GetUsersRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, UserRouterGroup)
}

func GetListMembersForInternalRouter() string {
	return fmt.Sprintf("%s%s%s", CurrentGroupVersion, MemberRouterGroup, MemberForInternalRouterSuffix)
}

func GetProxyRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, ProxyRouterGroup)
}

func GetPluginRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, PluginRouterGroup)
}

func GetNamespacesRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, NamespaceRouterGroup)
}

func GetNotificationRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, NotificationRouterGroup)
}

func GetWebHooksRouter() string {
	return fmt.Sprintf("%s%s", CurrentGroupVersion, WebHookRouterGroup)
}

func GetDBConnectionAbleRouter() string {
	return fmt.Sprintf("%s%s/connection", CurrentGroupVersion, InternalDBServiceRouterGroup)
}