package request

import "github.com/actiontech/sqle/sqle/utils"

const (
	KeyCheckUserCanAccess      = "check_user_can_access"
	KeyAvailableInstanceIDList = "available_instance_id_list_str"
	KeyCurrentUserID           = "current_user_id"
)

func GetUserIDFromRequest(input map[string]interface{}) (userID int, found bool) {

	if userIDInt, found := input[KeyCurrentUserID].(int); found {
		return userIDInt, found
	}

	if userIDUint, found := input[KeyCurrentUserID].(uint); found {
		return int(userIDUint), found
	}

	return 0, false
}

func IsNeedToCheckAccessByData(data map[string]interface{}) (need bool) {

	// if user access is required
	b, exist := data[KeyCheckUserCanAccess].(bool)
	if !exist || !b {
		return false
	}

	return true
}

func UpdateDataWithAvailableInstanceIDList(
	data map[string]interface{}, availableInstanceIDList []uint) (
	output map[string]interface{}) {

	data[KeyAvailableInstanceIDList] =
		utils.JoinUintSliceToString(availableInstanceIDList, ", ")

	return data
}
