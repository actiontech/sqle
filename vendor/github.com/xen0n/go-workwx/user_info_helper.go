package workwx

import (
	"fmt"
	"strconv"
)

func reshapeDeptInfo(
	ids []int64,
	orders []uint32,
	leaderStatuses []int,
) ([]UserDeptInfo, error) {
	if len(ids) != len(orders) {
		return nil, fmt.Errorf(
			"server API breakage: len(DeptIDs) (%d) != len(DeptOrder) (%d)",
			len(ids),
			len(orders),
		)
	}
	// sometimes leaderStatuses could be empty, but if not, the arrays should
	// have equal length
	if len(ids) != len(leaderStatuses) && len(leaderStatuses) != 0 {
		return nil, fmt.Errorf(
			"server API breakage: len(DeptIDs) (%d) != len(IsLeaderInDept) (%d)",
			len(ids),
			len(leaderStatuses),
		)
	}

	result := make([]UserDeptInfo, len(ids))
	for i := range ids {
		result[i].DeptID = ids[i]
		result[i].Order = orders[i]
		if i < len(leaderStatuses) {
			// apparently leaderStatuses could sometimes be empty, don't set
			// anybody as leader in that case
			// see https://github.com/xen0n/go-workwx/pull/78
			result[i].IsLeader = leaderStatuses[i] != 0
		}
	}

	return result, nil
}

func userGenderFromGenderStr(x string) (UserGender, error) {
	if x == "" {
		return UserGenderUnspecified, nil
	}
	n, err := strconv.Atoi(x)
	if err != nil {
		return UserGenderUnspecified, fmt.Errorf("gender string parse failed: %+v", err)
	}

	return UserGender(n), nil
}

func (x UserDetail) intoUserInfo() (UserInfo, error) {
	deptInfo, err := reshapeDeptInfo(x.DeptIDs, x.DeptOrder, x.IsLeaderInDept)
	if err != nil {
		return UserInfo{}, err
	}

	gender, err := userGenderFromGenderStr(x.Gender)
	if err != nil {
		return UserInfo{}, err
	}

	return UserInfo{
		UserID:      x.UserID,
		Name:        x.Name,
		Position:    x.Position,
		Departments: deptInfo,
		Mobile:      x.Mobile,
		Gender:      gender,
		Email:       x.Email,
		AvatarURL:   x.AvatarURL,
		Telephone:   x.Telephone,
		IsEnabled:   x.IsEnabled != 0,
		Alias:       x.Alias,
		Status:      UserStatus(x.Status),
		QRCodeURL:   x.QRCodeURL,
	}, nil
}
