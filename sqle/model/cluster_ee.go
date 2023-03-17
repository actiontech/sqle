package model

import "time"

func init() {
	autoMigrateList = append(autoMigrateList, &Leader{})
}

/*
使用基于 MySQL 表进行集群选主，方案参考自：https://blog.csdn.net/weichi7549/article/details/108136118
*/

const leaderTableAnchor = 1

type Leader struct {
	Anchor         int       `gorm:"primary_key"` // 常量值，保证改表仅有一行不重复记录。无其他意义。
	ServerId       string    `gorm:"not null"`
	LastSeenActive time.Time `gorm:"not null"`
}

func (a Leader) TableName() string {
	return "cluster_leader"
}

var GetLeader = "SELECT server_id AS leader FROM cluster_leader WHERE anchor=1 LIMIT 1"

func (s *Storage) GetClusterLeader() (string, error) {
	var leader = &Leader{}
	err := s.db.Select("server_id").Where("anchor = ?", leaderTableAnchor).First(leader).Error
	if err != nil {
		return "", err
	}
	return leader.ServerId, nil
}

var AttemptLeadership = `
INSERT ignore INTO cluster_leader (anchor, server_id, last_seen_active) VALUES (?, ?, now()) 
ON DUPLICATE KEY UPDATE 
server_id = IF(last_seen_active < now() - interval 30 second, VALUES(server_id), server_id), 
last_seen_active = IF(server_id = VALUES(server_id), VALUES(last_seen_active), last_seen_active)
`

func (s *Storage) AttemptClusterLeadership(serverId string) error {
	return s.db.Exec(AttemptLeadership, leaderTableAnchor, serverId).Error
}
