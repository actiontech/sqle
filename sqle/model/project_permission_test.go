package model

// import (
// 	"testing"

// 	sqlmock "github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestStorage_CheckUserHasOpToInstances(t *testing.T) {
// 	// 1. test for common user
// 	query := `
// SELECT
// instances.id
// FROM instances
// LEFT JOIN project_member_roles ON instances.id = project_member_roles.instance_id
// LEFT JOIN users ON project_member_roles.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// LEFT JOIN role_operations ON role_operations.role_id = roles.id
// WHERE
// instances.deleted_at IS NULL
// AND users.id = ?
// AND role_operations.op_code IN (?)
// AND instances.id IN (?, ?)
// GROUP BY instances.id

// UNION
// SELECT
// instances.id
// FROM instances
// LEFT JOIN project_member_group_roles ON instances.id = project_member_group_roles.instance_id
// LEFT JOIN roles ON roles.id = project_member_group_roles.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
// LEFT JOIN role_operations ON role_operations.role_id = roles.id
// LEFT JOIN user_groups ON project_member_group_roles.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL AND user_groups.stat = 0
// JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
// JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
// WHERE
// instances.deleted_at IS NULL
// AND users.id = ?
// AND role_operations.op_code IN (?)
// AND instances.id IN (?, ?)
// GROUP BY instances.id

// UNION
// SELECT instances.id
// FROM instances
// LEFT JOIN projects ON instances.project_id = projects.id
// LEFT JOIN project_manager ON project_manager.project_id = projects.id
// LEFT JOIN users ON project_manager.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// WHERE instances.deleted_at IS NULL
// AND users.id = ?
// AND instances.id IN (?, ?)
// GROUP BY instances.id
// `
// 	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	assert.NoError(t, err)
// 	InitMockStorage(mockDB)
// 	mock.ExpectQuery(query).WithArgs(1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 2).
// 		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

// 	inst1 := &Instance{}
// 	inst1.ID = 1
// 	inst2 := &Instance{}
// 	inst2.ID = 2

// 	user := &User{}
// 	user.ID = 1
// 	exist, err := GetStorage().CheckUserHasOpToInstances(user, []*Instance{inst1, inst2}, []uint{1})
// 	assert.NoError(t, err)
// 	assert.Equal(t, true, exist)
// 	mockDB.Close()
// }

// func Test_GetUserCanOpInstances(t *testing.T) {
// 	query := `
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN project_member_roles ON instances.id = project_member_roles.instance_id
// 		LEFT JOIN users ON project_member_roles.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?, ?, ?)
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN project_member_group_roles ON instances.id = project_member_group_roles.instance_id
// 		LEFT JOIN roles ON roles.id = project_member_group_roles.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		LEFT JOIN user_groups ON project_member_group_roles.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL AND user_groups.stat = 0
// 		JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
// 		JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?, ?, ?)
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_manager ON project_manager.project_id = projects.id
// 		LEFT JOIN users ON project_manager.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		GROUP BY instances.id
// 		`
// 	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	assert.NoError(t, err)
// 	InitMockStorage(mockDB)
// 	mock.ExpectQuery(query).WithArgs(1, 1, 2, 3, 1, 1, 2, 3, 1).
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "inst_1").AddRow(2, "inst_2"))

// 	user := &User{}
// 	user.ID = 1
// 	instances, err := GetStorage().GetUserCanOpInstances(user, []uint{1, 2, 3})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(instances))
// 	assert.Equal(t, uint(1), instances[0].ID)
// 	assert.Equal(t, "inst_1", instances[0].Name)
// 	assert.Equal(t, uint(2), instances[1].ID)
// 	assert.Equal(t, "inst_2", instances[1].Name)
// 	mockDB.Close()
// }

// func Test_GetUserCanOpInstancesFromProject(t *testing.T) {
// 	query := `
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_member_roles ON instances.id = project_member_roles.instance_id
// 		LEFT JOIN users ON project_member_roles.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?)
// 		AND projects.name = ?
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_member_group_roles ON instances.id = project_member_group_roles.instance_id
// 		LEFT JOIN roles ON roles.id = project_member_group_roles.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		LEFT JOIN user_groups ON project_member_group_roles.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL AND user_groups.stat = 0
// 		JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
// 		JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?)
// 		AND projects.name = ?
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_manager ON project_manager.project_id = projects.id
// 		LEFT JOIN users ON project_manager.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND projects.name = ?
// 		GROUP BY instances.id
// 		`
// 	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	assert.NoError(t, err)
// 	InitMockStorage(mockDB)
// 	mock.ExpectQuery(query).WithArgs(1, 1, "project_1", 1, 1, "project_1", 1, "project_1").
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "inst_1").AddRow(2, "inst_2"))

// 	user := &User{}
// 	user.ID = 1
// 	instances, err := GetStorage().GetUserCanOpInstancesFromProject(user, "project_1", []uint{1})
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(instances))
// 	assert.Equal(t, uint(1), instances[0].ID)
// 	assert.Equal(t, "inst_1", instances[0].Name)
// 	assert.Equal(t, uint(2), instances[1].ID)
// 	assert.Equal(t, "inst_2", instances[1].Name)
// 	mockDB.Close()
// }

// func Test_GetInstanceTipsByUserAndOperation(t *testing.T) {
// 	query := `
// 		SELECT
// 		instances.id, instances.name, instances.db_host as host, instances.db_port as port, instances.db_type
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_member_roles ON instances.id = project_member_roles.instance_id
// 		LEFT JOIN users ON project_member_roles.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?, ?)
// 		AND projects.name = ?
// 		AND instances.db_type = ?
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name, instances.db_host as host, instances.db_port as port, instances.db_type
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_member_group_roles ON instances.id = project_member_group_roles.instance_id
// 		LEFT JOIN roles ON roles.id = project_member_group_roles.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
// 		LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 		LEFT JOIN user_groups ON project_member_group_roles.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL AND user_groups.stat = 0
// 		JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
// 		JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND role_operations.op_code IN (?, ?)
// 		AND projects.name = ?
// 		AND instances.db_type = ?
// 		GROUP BY instances.id

// 		UNION
// 		SELECT
// 		instances.id, instances.name, instances.db_host as host, instances.db_port as port, instances.db_type
// 		FROM instances
// 		LEFT JOIN projects ON instances.project_id = projects.id
// 		LEFT JOIN project_manager ON project_manager.project_id = projects.id
// 		LEFT JOIN users ON project_manager.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// 		WHERE
// 		instances.deleted_at IS NULL
// 		AND users.id = ?
// 		AND projects.name = ?
// 		AND instances.db_type = ?
// 		GROUP BY instances.id
// 		`
// 	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 	assert.NoError(t, err)
// 	InitMockStorage(mockDB)
// 	mock.ExpectQuery(query).WithArgs(1, 1, 2, "project_1", "MySQL", 1, 1, 2, "project_1", "MySQL", 1, "project_1", "MySQL").
// 		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "db_type"}).AddRow(1, "inst_1", "MySQL").AddRow(2, "inst_2", "Oracle"))

// 	user := &User{}
// 	user.ID = 1
// 	instances, err := GetStorage().getInstanceTipsByUserAndOperation(user, "MySQL", "project_1", 1, 2)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 2, len(instances))
// 	assert.Equal(t, uint(1), instances[0].ID)
// 	assert.Equal(t, "MySQL", instances[0].DbType)
// 	assert.Equal(t, uint(2), instances[1].ID)
// 	assert.Equal(t, "Oracle", instances[1].DbType)
// 	mockDB.Close()
// }
