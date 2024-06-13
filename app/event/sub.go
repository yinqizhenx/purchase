package event

import (
	"purchase/domain/event"
)

var domainEventApp *DomainEventApp

func init() {
	domainEventApp = &DomainEventApp{}
}

//
// func Handler(ctx context.Context, msgs ...string) error {
// 	for i := range msgs {
// 		e := convertToEvent(msgs[i])
// 		err := domainEventApp.OnPACreated(ctx, e)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func convertToEvent(s string) event.Event {
	return nil
}
