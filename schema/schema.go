package schema

// StatusText 定义状态文本
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

const (
	OKStatus    StatusText = "OK"
	ErrorStatus StatusText = "ERROR"
	FailStatus  StatusText = "FAIL"
)

type StatusResult struct {
	Status StatusText `json:"status"` // 状态(OK)
}

type ErrorResult struct {
	Error ErrorItem `json:"error"` // 错误项
}

type ErrorItem struct {
	Code    int    `json:"code"`    // 错误码
	Message string `json:"message"` // 错误信息
}

type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type PaginationResult struct {
	Total    int  `json:"total"`
	Current  uint `json:"current"`
	PageSize uint `json:"pageSize"`
}

type PaginationParam struct {
	Pagination bool `form:"-"`                                     // 是否使用分页查询
	OnlyCount  bool `form:"-"`                                     // 是否仅查询count
	Current    uint `form:"current,default=1"`                     // 当前页
	PageSize   uint `form:"pageSize,default=10" binding:"max=100"` // 页大小
}

func (a PaginationParam) GetCurrent() uint {
	return a.Current
}

func (a PaginationParam) GetPageSize() uint {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 100
	}
	return pageSize
}

type OrderDirection int

const (
	// OrderByASC 升序排序
	OrderByASC OrderDirection = 1
	// OrderByDESC 降序排序
	OrderByDESC OrderDirection = 2
)

func NewOrderFieldWithKeys(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

type OrderField struct {
	Key       string         // 字段名(字段名约束为小写蛇形)
	Direction OrderDirection // 排序方向
}

func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

type IDResult struct {
	ID string `json:"id"`
}
