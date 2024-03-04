package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAddPage(t *testing.T) {
	Data := `{
        "Mainchara": {
            "Chara": "主角A",
            "Cast": "演员A",
            "Info": "关于主角A的信息"
        },
        "Relatives": {
            "Chara": "亲戚A",
            "Cast": "演员B",
            "Info": "关于亲戚A的信息"
        },
        "Team": {
            "Chara": "小组A",
            "Grade": "A",
            "Group": "1"
        }
    }`

	// 解析 JSON 数据到 Page 结构体实例
	var page Page
	if err := json.Unmarshal([]byte(Data), &page); err != nil {
		fmt.Println("解析 JSON 数据出错:", err)
		return
	}
	// 调用 Create 方法插入数据
	AddPage(page)
}
