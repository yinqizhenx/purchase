package service

// 在 domain.service 中定义领域服务的接口
// type SpamChecker interface {
// 	Check(content string) (Spam, error)
// }
//
// type PostCommentService struct {
// 	// repository.CommentRepository 和 ContentSal 在领域层都是以接口的形式存在
// 	// 因为在领域层不关心具体的实现
// 	commentRepo repository.CommentRepository
// 	spamChecker SpamChecker
// }
//
// func NewPostCommentService(repo repository.CommentRepository, checker SpamChecker) *PostCommentService {
// 	return &PostCommentService{
// 		commentRepo: repo,
// 		spamChecker: checker,
// 	}
// }
//
// func (s *PostCommentService) AddComment(ctx context.Context, post *Post, user *User, content string) (*Comment, error) {
// 	...
// }
