package system_handler

import (
	"fmt"
	"im_go/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// 系统状态信息
func systemHandle(c *gin.Context) {
	mem, _ := mem.VirtualMemory()
	cpuNum, _ := cpu.Counts(true)
	//cpuInfo, _ := cpu.Percent(10*time.Microsecond, true)
	info, _ := cpu.Info()
	data := make(map[string]interface{})
	//data["im.conn"] = len(ClientMaps)
	data["mem.total"] = fmt.Sprintf("%vMB", mem.Total/1024/1024)
	data["mem.free"] = fmt.Sprintf("%vMB", mem.Free/1024/1024)
	data["mem.used_percent"] = fmt.Sprintf("%s%%", strconv.FormatFloat(mem.UsedPercent, 'f', 2, 64))
	data["cpu.num"] = cpuNum
	data["cpu.info"] = info

	c.IndentedJSON(http.StatusOK, model.NewIMResponseData(data, ""))
}
