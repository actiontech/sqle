package model

type File struct {
	Model
	TaskId     uint   `json:"-" gorm:"index"`
	UniqueName string `json:"unique_name" gorm:"type:varchar(255)"`
	FileHost   string `json:"file_host" gorm:"type:varchar(255)"`
	NickName   string `json:"nick_name" gorm:"type:varchar(255)"`
}

const FixFilePath string = "audit_files/"
