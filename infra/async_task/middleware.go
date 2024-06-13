package async_task

type Middleware func(Handler) Handler

// func TaskExecuteTime(Handler) Handler {
// 	s := time
// }
