package v1beta1

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_XXX sets defaults for XXX.
func SetDefaults_XXX(obj *XXX) {
	// XXX name prefix is fixed to `hello-`
	if obj.ObjectMeta.GenerateName == "" {
		obj.ObjectMeta.GenerateName = "hello-"
	}

	SetDefaults_XXXSpec(&obj.Spec)
}

// SetDefaults_XXXSpec sets defaults for XXX spec.
func SetDefaults_XXXSpec(obj *XXXSpec) {
	if obj.DisplayName == "" {
		obj.DisplayName = "xxxdefaulter"
	}
}
