package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

// 定义新的功能门控  
const MyNewFeature featuregate.Feature = "MyNewFeature"

func init() {
	// runtime.Must(utilfeature.DefaultMutableFeatureGate.Add(defaultFeatureGates))
	runtime.Must(DefaultMutableFeatureGate.Add(defaultFeatureGates))
}

// defaultFeatureGates consists of all known specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	MyNewFeature: {Default: false, PreRelease: featuregate.Alpha},
}
