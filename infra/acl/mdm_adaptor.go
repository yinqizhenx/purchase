package acl

import "purchase/domain/entity/user"

// MDMSearchReq 搜索请求
type MDMSearchReq struct {
	Search    string `json:"search"`
	PageSize  int    `json:"page_size"`
	JobStatus int    `json:"job_status"`
}

func NewMDMSearchReq(account string) MDMSearchReq {
	return MDMSearchReq{
		Search: account,
	}
}

type (
	BaseResp struct {
		RetCode int    `json:"retcode"`
		Message string `json:"message"`
	}

	// MDMSearchUserItem 搜索用户结果项
	MDMSearchUserItem struct {
		NameEn              string `json:"name_en"`
		NameCn              string `json:"name_cn"`
		Nickname            string `json:"nickname"`
		Email               string `json:"email"`
		Avatar              string `json:"avatar"`
		EmpType             string `json:"emp_type"`
		EmpKind             string `json:"emp_kind"`
		IsOutsourcing       bool   `json:"is_outsourcing"`
		CsLoc               string `json:"cs_loc"`
		JobStatus           string `json:"job_status"`
		JobStatusDesc       string `json:"job_status_desc"`
		EnglishAbbreviation string `json:"english_abbreviation"`
		IsAdvisor           bool   `json:"is_advisor"`
	}

	// MDMSearchData 搜索结果data
	MDMSearchData struct {
		List     []*MDMSearchUserItem `json:"list"`
		PageSize string               `json:"page_size"`
	}

	// MDMSearchResp 搜索结果
	MDMSearchResp struct {
		BaseResp
		Data *MDMSearchData `json:"data"`
	}
)

func NewMDMSearchResp() *MDMSearchResp {
	return &MDMSearchResp{}
}

func (res *MDMSearchResp) toUser() *user.User {
	return &user.User{}
}
