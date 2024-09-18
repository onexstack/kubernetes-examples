package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Credential struct {
	Token string `json:"token"`
}

func main() {
	// 模拟获取一个令牌  
	token := "your-static-token"

	// 创建凭证响应
	credential := Credential{Token: token}

	// 将响应序列化为 JSON
	response, err := json.Marshal(credential)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		os.Exit(1)
	}

	// 输出 JSON 响应  
	fmt.Println(string(response))
}
