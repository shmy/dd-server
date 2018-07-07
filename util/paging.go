package util

type Paging struct {
	Page   int // 第几页
	Limit  int // 数量
	Offset int // 计算出的偏移值
}

// 解析 分页参数
func ParsePaging(c *ApiContext) Paging {
	page := c.DefaultQueryInt("page", 1)
	limit := c.DefaultQueryInt("per_page", 20)
	return Paging{
		page,
		limit,
		(page - 1) * limit,
	}
}
