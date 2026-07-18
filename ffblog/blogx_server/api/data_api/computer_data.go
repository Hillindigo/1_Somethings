package data_api

import (
	"blogx_server/common/res"
	"blogx_server/utils/computer"
	"github.com/gin-gonic/gin"
)

type ComputerDataResponse struct {
	CpuPercent  float64 `json:"cpuPercent"`
	MemPercent  float64 `json:"memPercent"`
	DiskPercent float64 `json:"diskPercent"`
}

func (DataApi) ComputerDataView(c *gin.Context) {
	var data = ComputerDataResponse{
		CpuPercent:  computer.GetCpuPercent(),
		MemPercent:  computer.GetMemPercent(),
		DiskPercent: computer.GetDiskPercent(),
	}
	res.OkWithData(data, c)
}
