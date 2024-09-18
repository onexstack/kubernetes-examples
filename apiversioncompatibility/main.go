package main

import (
	"context"
	"encoding/json"
	"fmt"
)

// 定义结构体 Frobber
type Frobber struct {
	Height int      `json:"height"`
	Width  int      `json:"width"`
	Param  string   `json:"param"`  // 旧版参数
	Params []string `json:"params"` // 新版参数
}

// 假设一个存储，用于示例
var frobberStorage = make(map[string]*Frobber)

// CreateFrobber 创建 Frobber，同时兼容旧版和新版参数
func CreateFrobber(ctx context.Context, frobber *Frobber) error {
	// 检查旧版参数，如果存在则添加到新版参数中
	if frobber.Param != "" {
		frobber.Params = append(frobber.Params, frobber.Param)
	}
	if frobber.Param != "" && len(frobber.Params) > 0 && frobber.Param != frobber.Params[0] {
		return fmt.Errorf("param is not equal to params[0]")
	}
	// 假设用 "frobber1" 作为标识符存储
	frobberStorage["frobber1"] = frobber
}

// GetFrobber 获取 Frobber，确保新旧参数都是可用的
func GetFrobber(ctx context.Context) *Frobber {
	// 从存储中获取 Frobber
	frobber := frobberStorage["frobber1"]
	if frobber == nil {
		return nil // 如果没有找到，返回 nil
	}
	// 如果新版参数为空但旧版参数存在，则转移旧版参数
	if len(frobber.Params) == 0 && frobber.Param != "" {
		frobber.Params = append(frobber.Params, frobber.Param)
	}
	return frobber
}

// UpdateFrobber 更新 Frobber，同时兼容旧版和新版参数
func UpdateFrobber(ctx context.Context, after *Frobber) error {
	// 从存储中获取当前 Frobber
	existingFrobber := frobberStorage["frobber1"]
	if existingFrobber == nil {
		return fmt.Errorf("frobber not found")
	}

	// 检查旧版参数，如果存在则添加到新版参数中
	if newFrobber.Param != "" {
		// 如果 Params 为空，转移旧版参数到新版参数
		if len(existingFrobber.Params) == 0 {
			existingFrobber.Params = append(existingFrobber.Params, newFrobber.Param)
		} else {
			// 如果 Params 不为空且与新 Param 不相等，返回错误
			if newFrobber.Param != existingFrobber.Params[0] {
				return fmt.Errorf("param is not equal to params[0]")
			}
		}
	}

	// 更新现有字段
	existingFrobber.Height = newFrobber.Height
	existingFrobber.Width = newFrobber.Width

	// 更新新版参数
	existingFrobber.Params = newFrobber.Params

	return nil
}

func main() {
	ctx := context.Background()

	// 创建 Frobber 示例
	createFrobber := &Frobber{
		Height: 10,
		Width:  20,
		Param:  "old_param_value",
	}
	CreateFrobber(ctx, createFrobber)

	// 获取 Frobber 示例
	retrievedFrobber := GetFrobber(ctx)
	if retrievedFrobber != nil {
		data, _ := json.Marshal(retrievedFrobber)
		fmt.Println("Retrieved Frobber:", string(data)) // 打印以 JSON 格式输出的 Frobber 数据
	} else {
		fmt.Println("Frobber not found")
	}
}
