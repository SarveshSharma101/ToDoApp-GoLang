package datamodel

type (
	Project struct {
		Pid  int    `gorm:"primaryKey" json:"pid"`
		Name string `gorm:"project_name; not null" json:"name"`
	}

	Task struct {
		Tid            int    `gorm:"primaryKey"`
		TaskName       string `gorm:"not null"`
		TaskDesc       string `gorm:"not null"`
		Status         string `gorm:"default:Not Completed"`
		ClosureComment string
		ProjectID      int
		Project        Project `gorm:"constraint:OnDelete:CASCADE"`
		UserId         int
		User           User `gorm:"constraint:OnDelete:SET NULL"`
	}

	User struct {
		Uid         int    `gorm:"primaryKey;autoIncrement"`
		Username    string `gorm:"not null"`
		Password    string `gorm:"not null"`
		Type        int
		LoginStatus bool `gorm:"default:false"`
	}

	RedisUser struct {
		Username string
		Role     int
	}
)
