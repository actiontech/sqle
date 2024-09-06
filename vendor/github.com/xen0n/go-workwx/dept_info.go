package workwx

// CreateDept 创建部门
func (c *WorkwxApp) CreateDept(deptInfo *DeptInfo) (deptID int64, err error) {
	resp, err := c.execDeptCreate(reqDeptCreate{
		DeptInfo: deptInfo,
	})
	if err != nil {
		return 0, err
	}
	return resp.ID, nil
}

// ListAllDepts 获取全量组织架构。
func (c *WorkwxApp) ListAllDepts() ([]*DeptInfo, error) {
	resp, err := c.execDeptList(reqDeptList{
		HaveID: false,
		ID:     0,
	})
	if err != nil {
		return nil, err
	}

	return resp.Department, nil
}

// ListDepts 获取指定部门及其下的子部门。
func (c *WorkwxApp) ListDepts(id int64) ([]*DeptInfo, error) {
	resp, err := c.execDeptList(reqDeptList{
		HaveID: true,
		ID:     id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Department, nil
}

// SimpleListAllDepts 获取全量组织架构（简易）。
func (c *WorkwxApp) SimpleListAllDepts() ([]*DeptInfo, error) {
	resp, err := c.execDeptSimpleList(reqDeptSimpleList{
		HaveID: false,
		ID:     0,
	})
	if err != nil {
		return nil, err
	}

	return resp.DepartmentIDs, nil
}

// SimpleListDepts 获取指定部门及其下的子部门（简易）。
func (c *WorkwxApp) SimpleListDepts(id int64) ([]*DeptInfo, error) {
	resp, err := c.execDeptSimpleList(reqDeptSimpleList{
		HaveID: true,
		ID:     id,
	})
	if err != nil {
		return nil, err
	}

	return resp.DepartmentIDs, nil
}
