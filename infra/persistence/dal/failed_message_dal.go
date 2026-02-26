package dal

import (
	"context"
	"encoding/json"

	"purchase/infra/mq"
	"purchase/infra/persistence/dal/db/ent"
)

type FailedMessageDal struct {
	db *ent.Client
}

func NewFailedMessageDal(cli *ent.Client) mq.FailedMessageStore {
	return &FailedMessageDal{db: cli}
}

func (dal *FailedMessageDal) Save(ctx context.Context, m *mq.Message, errMsg string) error {
	headerJSON, _ := json.Marshal(m.Header())
	return dal.db.FailedMessage.Create().
		SetMessageID(m.ID).
		SetTopic(m.EventName()).
		SetBizCode(m.BizCode()).
		SetBody(string(m.Body)).
		SetHeader(string(headerJSON)).
		SetErrorMsg(errMsg).
		SetState("pending").
		Exec(ctx)
}
