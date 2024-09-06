//go:build enterprise
// +build enterprise

package model

import (
	"time"

	"github.com/actiontech/sqle/sqle/errors"
)

func init() {
	autoMigrateList = append(autoMigrateList, &Leader{})
	autoMigrateList = append(autoMigrateList, &Node{})
}

const leaderTableAnchor = 1

type Leader struct {
	Anchor       int       `gorm:"primary_key"` // 常量值，保证该表仅有一行不重复记录。无其他意义。
	ServerId     string    `gorm:"not null;size:255"`
	LastSeenTime time.Time `gorm:"not null"`
}

func (a Leader) TableName() string {
	return "cluster_leader"
}

func (s *Storage) GetClusterLeader() (string, error) {
	var leader = &Leader{}
	err := s.db.Select("server_id").Where("anchor = ?", leaderTableAnchor).First(leader).Error
	if err != nil {
		return "", err
	}
	return leader.ServerId, nil
}

var MaintainLeader = `
INSERT ignore INTO cluster_leader (anchor, server_id, last_seen_time) VALUES (?, ?, now()) 
ON DUPLICATE KEY UPDATE 
server_id = IF(last_seen_time < now() - interval 30 second, VALUES(server_id), server_id), 
last_seen_time = IF(server_id = VALUES(server_id), VALUES(last_seen_time), last_seen_time)
`

func (s *Storage) MaintainClusterLeader(serverId string) error {
	return s.db.Exec(MaintainLeader, leaderTableAnchor, serverId).Error
}

type Node struct {
	Model
	ServerId     string `json:"server_id" gorm:"size:255;index:unique"`

	HardwareSign string `json:"hardware_sign" gorm:"type:varchar(3000)"`
}

func (l *Node) TableName() string {
	return "cluster_node_info"
}

var RegisterNode = `
INSERT INTO cluster_node_info (server_id, hardware_sign) VALUES (?,?) 
ON DUPLICATE KEY UPDATE hardware_sign = VALUES(hardware_sign)
`

func (s *Storage) RegisterClusterNode(serverId, HardwareSign string) error {
	return errors.New(errors.ConnectStorageError, s.db.Exec(RegisterNode, serverId, HardwareSign).Error)
}

func (s *Storage) GetClusterNodes() ([]*Node, error) {
	var nodes []*Node
	err := s.db.Model(Node{}).Find(&nodes).Error
	return nodes, errors.New(errors.ConnectStorageError, err)
}
