package acl

type StudentRecognitionSal struct{}

func NewStudentRecognitionSal() *StudentRecognitionSal {
	return &StudentRecognitionSal{}
}

//
// 在传统意义的防腐层实现中，会有一个适配器和一个对应的翻译器，其中适配器的作用是适配对其他上下文的调用，而翻译器则是将调用的结果转换成本地上下文中的元素。
// 在这里，我们为了保持代码的简单，没有特意声明这样两个对象，我们使用了 overpass 来进行rpc调用，overpass 在这里起到了适配器的作用，至于翻译器，我们只是简单的提出了一个方法，在方法名上做了特殊的前缀修饰。
// 最后，NewStudentRecognitionSal 会作为 StudentRecognitionService 的实现类，最终被注入到应用服务中。

// func (s *StudentRecognitionSal) StudentRecognitionFrom(ctx context.Context, paperId int64) (StudentRecognition, error) {
// 	// edu_ark_paper 是overpass自动生成的
// 	resp, err := edu_ark_paper.QueryPaperSheetInfo(ctx, &paper.QueryPaperSheetInfoRequest{
// 		PaperId: paperId,
// 	})
// 	if err != nil || resp.Data == nil {
// 		return StudentRecognition_Unknown, err
// 	}
// 	return transToStudentRecognition(resp.Data), nil
// }
//
// func transToStudentRecognition(info *paper.PaperSheetInfoData) StudentRecognition {
// 	switch info.SheetProp.GetStuNoType() {
// 	case common.StuNoType_JIKENo:
// 		return StudentRecognition_JIKE
// 	case common.StuNoType_StuNo:
// 		return StudentRecognition_STU
// 	}
// 	return StudentRecognition_Unknown
// }
