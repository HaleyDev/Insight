package form

type EditPermission struct {
	Id       uint   `form:"id" json:"id" binding:"omitempty"`                                                  // id
	Name     string `form:"name" json:"name" binding:"required,max=60"`                                        // 权限名称
	Desc     string `form:"desc" json:"desc" binding:"omitempty"`                                              // 权限描述
	Method   string `form:"method" json:"method" binding:"required,oneof=GET POST PUT DELETE" label:"接口请求方法""` // 请求方法
	Route    string `form:"route" json:"route" binding:"required"`                                             // 请求路由
	Func     string `form:"func" json:"func" binding:"required"`                                               // 权限功能
	FuncPath string `form:"func_path" json:"func_path" binding:"required"`                                     // 功能路径
	IsAuth   int8   `form:"is_auth" json:"is_auth" binding:"required"`                                         // 是否需要认证
	Sort     int32  `form:"sort" json:"sort" binding:"required"`                                               // 排序
}

func NewEditPermissionForm() *EditPermission {
	return &EditPermission{}
}

type ListPermission struct {
	Paginate
	Name   string `form:"name" json:"name" binding:"omitempty,max=60"`                                                          // 权限名称
	Method string `form:"method" json:"method" binding:"omitempty,oneof=GET POST PUT DELETE OPTIONS HEAD PATCH" label:"接口请求方法"` // 接口请求方法
	Route  string `form:"route" json:"route" binding:"omitempty"`                                                               // 接口路由
	IsAuth int8   `form:"is_auth" json:"is_auth" binding:"omitempty"`                                                           // 是否需要认证
}

func NewListPermissionQuery() *ListPermission {
	return &ListPermission{}
}
