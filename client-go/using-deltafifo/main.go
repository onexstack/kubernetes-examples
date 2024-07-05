package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func main() {

	// 创建一个DeltaFIFO对象
	fifo := cache.NewDeltaFIFO(cache.MetaNamespaceKeyFunc, nil)

	dep1 := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "dep1", Namespace: metav1.NamespaceDefault}}
	dep2 := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "dep2", Namespace: metav1.NamespaceDefault}}
	// 1. 将对象添加事件放入 DeltaFIFO 中
	_ = fifo.Add(dep1)
	_ = fifo.Add(dep2)

	// 2. 将对象变更事件放入 DeltaFIFO 中
	dep1.Name = "dep1-modified"
	_ = fifo.Update(dep1)

	// 3. 以列表形式返回所有 Key
	fmt.Println(fifo.ListKeys())

	// 4. 将对象删除事件放入 DeltaFIFO 中
	_ = fifo.Delete(dep1)

	// 5. "不断"从 DeltaFIFO 中 Pop 资源对象。
	// 当中有个回调函数，作用是分别不同事件所有做的不同回调方法，这里只打印了一条信息
	for {
		_, _ = fifo.Pop(func(obj interface{}, isInInitialList bool) error {
			for _, delta := range obj.(cache.Deltas) {
				deploy := delta.Object.(*appsv1.Deployment)

				// 这里进行回调，区分不同事件，可以执行业务逻辑 ex: 统计次数 加入本地缓存等操作。
				switch delta.Type {
				case cache.Added:
					fmt.Printf("Added: %s/%s\n", deploy.Namespace, deploy.Name)
				case cache.Updated:
					fmt.Printf("Updated: %s/%s\n", deploy.Namespace, deploy.Name)
				case cache.Deleted:
					fmt.Printf("Deleted: %s/%s\n", deploy.Namespace, deploy.Name)
				}
			}

			return nil
		})
	}
}
