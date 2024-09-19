package main

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/superproj/k8sdemo/featuregates/feature"
)

func main() {
	fs := pflag.NewFlagSet("feature", pflag.ExitOnError)
	feature.DefaultMutableFeatureGate.AddFlag(fs)
	pflag.Parse()

	// 设置功能门控的值
	//_ = feature.DefaultMutableFeatureGate.SetFromMap(*featureGates)

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
