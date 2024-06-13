package department

type Department struct {
	Code          string   `json:"code"`
	I18nName      string   `json:"i18n_name"`
	ParentCode    string   `json:"parent_code"`
	PathCodes     []string `json:"path_codes"`
	I18nPathNames string   `json:"i18n_path_names"`
	OrgID         int64    `json:"org_id"`     // mdm 组织ID
	Leader        string   `json:"leader"`     // 部门leader
	ProjOwner     string   `json:"proj_owner"` // 项目负责人
	Disabled      bool     `json:"disabled"`   // 是否已经失效
	IsEmptyCache  bool     `json:"is_empty_cache"`
}
