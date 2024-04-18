package models

type DepartmentInfo struct {
	ID                 uint   `json:"id"`
	DepartmentName     string `json:"department_name"`
	StaffQuantity      int    `json:"staff_quantity"`
	DepartmentDirector string `json:"department_director"`
	ModuleId           int    `json:"module_id"`
}

func (DepartmentInfo) TableName() string {
	return "department_info"
}
