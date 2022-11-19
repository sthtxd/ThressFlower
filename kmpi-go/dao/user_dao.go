package dao

import "time"

const (
	SUPER_ADMIN_AUTHORITY = 100
	ADMIN_AUTHORITY       = 99
	UPDATE_AUTHORITY      = 98
	SELECT_AUTHORITY      = 97
)

type UserDao struct {
	ID            int       `gorm:"primary_key;auto_increment" json:"id"`
	UserId        string    `gorm:"type:varchar(50);not null;uniqueIndex:user_id_idx;comment:'用户id'" json:"userId"`
	UserAccount   string    `gorm:"type:varchar(50);not null;uniqueIndex:user_account_id_idx;comment:'用户账户'" json:"userAccount"`
	UserPassword  string    `gorm:"type:varchar(50);not null;comment:'用户密码'" json:"userPassword"`
	UserAuthority int       `gorm:"type:int(10);not null;default:0;comment:'100=admin,99=update,98=select'" json:"userAuthority"`
	Timestamp     time.Time `gorm:"type:timestamp(3) on update current_timestamp(3);omitempty;default:current_timestamp(3);comment:'时间戳'" json:"timestamp"`
}

func (p UserDao) TableName() string {

	return "user"
}
