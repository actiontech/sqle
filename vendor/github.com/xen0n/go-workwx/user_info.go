package workwx

// UserDetail 成员详细信息的公共字段
type UserDetail struct {
	UserID         string   `json:"userid"`
	Name           string   `json:"name,omitempty"`
	DeptIDs        []int64  `json:"department"`
	DeptOrder      []uint32 `json:"order"`
	Position       string   `json:"position"`
	Mobile         string   `json:"mobile,omitempty"`
	Gender         string   `json:"gender,omitempty"`
	Email          string   `json:"email,omitempty"`
	IsLeaderInDept []int    `json:"is_leader_in_dept"`
	AvatarURL      string   `json:"avatar"`
	Telephone      string   `json:"telephone"`
	IsEnabled      int      `json:"enable"`
	Alias          string   `json:"alias"`
	Status         int      `json:"status"`
	QRCodeURL      string   `json:"qr_code"`
	// TODO: extattr external_profile external_position
}

// GetUser 读取成员
func (c *WorkwxApp) GetUser(userid string) (*UserInfo, error) {
	resp, err := c.execUserGet(reqUserGet{
		UserID: userid,
	})
	if err != nil {
		return nil, err
	}

	obj, err := resp.intoUserInfo()
	if err != nil {
		return nil, err
	}

	// TODO: return bare T instead of &T?
	return &obj, nil
}

// UpdateUser 更新成员
func (c *WorkwxApp) UpdateUser(userDetail *UserDetail) error {
	_, err := c.execUserUpdate(reqUserUpdate{
		UserDetail: userDetail,
	})
	if err != nil {
		return err
	}
	return nil
}

// ListUsersByDeptID 获取部门成员详情
func (c *WorkwxApp) ListUsersByDeptID(deptID int64, fetchChild bool) ([]*UserInfo, error) {
	resp, err := c.execUserList(reqUserList{
		DeptID:     deptID,
		FetchChild: fetchChild,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*UserInfo, len(resp.Users))
	for index, user := range resp.Users {
		userInfo, err := user.intoUserInfo()
		if err != nil {
			return nil, err
		}
		users[index] = &userInfo
	}
	return users, nil
}

// ConvertUserIDToOpenID userid转openid
func (c *WorkwxApp) ConvertUserIDToOpenID(userID string) (string, error) {
	resp, err := c.execConvertUserIDToOpenID(reqConvertUserIDToOpenID{
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	return resp.OpenID, nil
}

// ConvertOpenIDToUserID openid转userid
func (c *WorkwxApp) ConvertOpenIDToUserID(openID string) (string, error) {
	resp, err := c.execConvertOpenIDToUserID(reqConvertOpenIDToUserID{
		OpenID: openID,
	})
	if err != nil {
		return "", err
	}
	return resp.UserID, nil
}

// GetUserJoinQrcode 获取加入企业二维码
func (c *WorkwxApp) GetUserJoinQrcode(sizeType SizeType) (string, error) {
	resp, err := c.execUserJoinQrcode(reqUserJoinQrcode{
		SizeType: sizeType,
	})
	if err != nil {
		return "", err
	}
	return resp.JoinQrcode, nil
}

// GetUserIDByMobile 通过手机号获取 userid
func (c *WorkwxApp) GetUserIDByMobile(mobile string) (string, error) {
	resp, err := c.execUserIDByMobile(reqUserIDByMobile{
		Mobile: mobile,
	})
	if err != nil {
		return "", err
	}
	return resp.UserID, nil
}

// GetUserIDByEmail 通过邮箱获取 userid
func (c *WorkwxApp) GetUserIDByEmail(email string, emailType EmailType) (string, error) {
	if emailType == 0 {
		emailType = EmailTypeCorporate
	}
	resp, err := c.execUserIDByEmail(reqUserIDByEmail{
		Email:     email,
		EmailType: emailType,
	})
	if err != nil {
		return "", err
	}
	return resp.UserID, nil
}

// GetUserInfoByCode 获取访问用户身份，根据code获取成员信息
func (c *WorkwxApp) GetUserInfoByCode(code string) (*UserIdentityInfo, error) {
	resp, err := c.execUserInfoGet(reqUserInfoGet{
		Code: code,
	})
	if err != nil {
		return nil, err
	}
	return &resp.UserIdentityInfo, nil
}
