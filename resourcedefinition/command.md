## 生成客户端 SDK

```bash
$ git clone https://github.com/superproj/k8sdemo.git
$ cd resourcedefinition/
$ go mod tidy
$ client-gen -v 10 --go-header-file ./boilerplate.go.txt --output-dir ./generated/clientset --output-pkg=github.com/superproj/k8sdemo/resourcedefinition/generated/clientset --clientset-name=versioned --input-base= --input $PWD/apps/v1beta1
```

## 生成 DeepCopyObject

```bash
$ git clone https://github.com/superproj/k8sdemo.git
$ cd resourcedefinition/
$ go mod tidy
$ deepcopy-gen -v 10 --go-header-file ./boilerplate.go.txt --output-file zz_generated.deepcopy.go ./apps/v1beta1
```
