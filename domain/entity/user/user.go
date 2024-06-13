package user

type User struct {
	ID      int64  `db:"id"`                     // 主键id
	Account string `db:"account" json:"account"` // 域账号
	Name    string `db:"name" json:"name"`       // 名字
	Email   string `db:"email" json:"email"`
	State   int32  `db:"state" json:"state"` // 人员状态
}
