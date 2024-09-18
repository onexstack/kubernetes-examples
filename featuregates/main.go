package main

import (
	"fmt"

	"github.com/superproj/k8sdemo/featuregates/feature"
)

func main() {
	if feature.DefaultFeatureGate.Enabled(feature.MyNewFeature) {
		// 启用新特性时的逻辑
		// 例如，执行某些新功能的代码
		fmt.Println("Feature Gates MyNewFeature is opend")
	} else {
		// 新特性未启用时的逻辑
		// 例如，执行旧版功能的代码
		fmt.Println("Feature Gates MyNewFeature is closed")
	}
}
