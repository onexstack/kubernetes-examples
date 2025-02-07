├── cmd/ # 存放组件main文件的目录。将应用下的所有组件统一放在一个目录中，便于维护
│   └── kube-scheduler/ # scheduler应用层代码，主要包括应用配置、应用构建和启动代码
│       ├── app/
│       │   ├── config/ 
│       │   │   └── config.go # 保存了应用配置相关代码。配置内容通过命令行选项或者配置文件构建而来
│       │   ├── options/
│       │   │   ├── configfile.go # 配置文件读取或写入
│       │   │   ├── deprecated.go # 包含已弃用命令行 Flag
│       │   │   └── options.go # 用来给应用设置命令行选项，并对这些选项进行补全和校验
│       │   ├── server.go # 包含了创建并启动调度器的代码
│       └── scheduler.go # scheduler main 入口
├── hack/
│   └── local-up-cluster.sh* # 用于本地测试集群的部署脚本，包括 kube-scheduler 的部署
├── Makefile -> build/root/Makefile # 用于编译 kube-scheduler 二进制文件的 Makefile，编译命令为：make all WHAT=cmd/kube-scheduler。 
├── _output/
│   └── bin/kube-sheduler # 编译生成的 kube-scheduler 二进制文件
├── pkg/ # 目录中保存了调度器的核心实现代码
│   ├── features/
│   │   └── kube_features.go # 定义了一些与 kube-scheduler 相关的功能门控
│   └── scheduler/ # 调度器的核心实现
│       ├── apis/ # 调度器组件配置 KubeSchedulerConfiguration 定义
│       │   └── config/
│       │       ├── latest/ 
│       │       │   └── latest.go # 包含了用来获取最新KubeSchedulerConfiguration定义的 Default() 函数
│       │       ├── register.go # 包含了用于注册调度相关联的资源对象的内容版本的代码
│       │       ├── scheme/
│       │       │   └── scheme.go # 用于将资源对象注册到全局的资源注册表中，会同时注册内部版本资源对象和版本化的资源对象
│       │       ├── testing/ # 定义了一些 方法、变量用来参与调度器的单元测试。例如：包含了参与测试的调度器插件的配置。
│       │       │   ├── config.go
│       │       │   └── defaults/
│       │       │       └── defaults.go
│       │       ├── types.go # 包含了 KubeSchedulerConfiguration 内部版本定义
│       │       ├── types_pluginargs.go # 包含了调度器插件配置资源对象的内部版本定义
│       │       └── zz_generated.deepcopy.go # 使用 deepcopy-gen 工具生成的代码（内部版本的深拷贝）
│       │       ├── v1/ # 包含了调度相关资源对象 v1 版本定义（外部版本）
│       │       │   ├── conversion.go # 包含了内部版本和 v1 版本自定义转换规则
│       │       │   ├── default_plugins.go # 包含了调度器插件配置资源对象默认值设置代码
│       │       │   ├── defaults.go # 包含了KubeSchedulerConfiguration资源对象的默认值设置代码
│       │       │   ├── register.go # 包含了用于注册调度相关联的资源对象的 v1 版本的代码
│       │       │   ├── zz_generated.conversion.go # 使用 conversion-gen 工具生成的代码
│       │       │   ├── zz_generated.deepcopy.go # 使用 deepcopy-gen 工具生成的代码（v1 版本的深拷贝）
│       │       │   └── zz_generated.defaults.go # 使用 defaulter-gen 工具生成的代码
│       │       ├── validation/ # 包含了调度器相关资源对象的验证代码
│       │       │   ├── validation.go # 主要验证 KubeSchedulerConfiguration
│       │       │   └── validation_pluginargs.go # 主要验证各类调度器插件资源配置对象
│       ├── eventhandlers.go # 包含了对 Pod 待调度队列和 Node 缓存操作的方法
│       ├── extender.go # 包含了 scheduler extender（调度器扩展器）相关的代码
│       ├── framework/ # 包含了 Scheduling Framework（调度框架） 相关的代码
│       │   ├── autoscaler_contract/
│       │   ├── cycle_state.go
│       │   ├── events.go
│       │   ├── extender.go
│       │   ├── interface.go
│       │   ├── listers.go
│       │   ├── parallelize/
│       │   │   ├── error_channel.go
│       │   │   └── parallelism.go
│       │   ├── plugins/ # 包含了各类内置的调度插件实现
│       │   │   ├── defaultbinder/ # 默认的绑定器，负责将 Pod 绑定到节点上执行
│       │   │   ├── defaultpreemption/ # 在资源不足时，负责预先终止低优先级的 Pod，以腾出资源给高优先级的 Pod
│       │   │   ├── dynamicresources/ # 动态资源调度插件，用于根据实时资源需求动态调整 Pod 的调度
│       │   │   ├── examples/ # 不同类型调度器的实现 Demo
│       │   │   │   ├── multipoint/
│       │   │   │   │   └── multipoint.go # 多点插件示例
│       │   │   │   ├── prebind/
│       │   │   │   │   └── prebind.go # 预绑定插件示例
│       │   │   │   └── stateful/
│       │   │   │       └── stateful.go # 有状态插件示例
│       │   │   ├── feature/ # 定义 Features 类型的结构体，用于存储被调度器插件用到的 Feature Gate 的值。使用该包可以避免直接依赖 Kubernetes 内部的 features 包
│       │   │   ├── helper/ # 包含了一些工具或者帮助类的函数
│       │   │   ├── imagelocality/ # 选择已经存在 Pod 运行所需容器镜像的节点。实现的扩展点：Score
│       │   │   ├── interpodaffinity/ # 根据 Pod 间的亲和性和反亲和性调度 Pod
│       │   │   ├── names/ # 统一定义了内置调度器插件的名字
│       │   │   ├── nodeaffinity/ # 根据节点的属性或标签，调度 Pod 到具有特定属性的节点上
│       │   │   ├── nodename/ # 检查 Pod 指定的节点名称与当前节点是否匹配，也即根据节点名调度
│       │   │   ├── nodeports/ # 根据节点上的端口分配情况，调度 Pod 到端口可用的节点上
│       │   │   ├── noderesources/ # 包含了一些根据节点资源和 Pod 资源请求来调度 Pod 的算法，是调度器中非常核心的调度插件实现
│       │   │   │   ├── balanced_allocation.go # 调度 Pod 时，选择资源使用更为均衡的节点
│       │   │   │   ├── fit.go # 检查节点是否拥有 Pod 请求的所有资源，只有满足资源需求的节点才被允许调度
│       │   │   │   ├── least_allocated.go # 选择资源分配较少的节点
│       │   │   │   ├── most_allocated.go # 选择已分配资源多的节点（常用于降本场景下）
│       │   │   │   ├── requested_to_capacity_ratio.go
│       │   │   │   ├── resource_allocation.go
│       │   │   ├── nodeunschedulable/ # 过滤掉 .spec.unschedulable 值为 true 的节点
│       │   │   ├── nodevolumelimits/ # 根据节点的存储容量限制，调度 Pod 到合适的节点上
│       │   │   ├── podtopologyspread/ # 根据 Pod 的拓扑分布要求，调度 Pod 到不同的区域或节点上
│       │   │   ├── queuesort/ # 根据 Pod 的优先级对待调度的 Pod 进行排序，优先级高的 Pod 会被优先调度
│       │   │   ├── registry.go
│       │   │   ├── schedulinggates/ # 根据调度门控条件，控制 Pod 的调度行为
│       │   │   ├── tainttoleration/ # 根据节点的 Taint 和 Pod 的 Tolerations，以确定 Pod 是否可以调度到节点上
│       │   │   ├── testing/ # 包含了一些用于调度器单元测试的函数
│       │   │   ├── volumebinding/ # 检查节点是否有请求的卷，或是否可以绑定请求的卷，类似的有 VolumeRestrictions/VolumeZone/NodeVolumeLimits/EBSLimits/GCEPDLimits/AzureDiskLimits/CinderVolume
│       │   │   ├── volumerestrictions/ # 根据卷的限制条件，调度 Pod 到符合条件的节点上
│       │   │   └── volumezone/ # 根据卷的存储区域要求，调度 Pod 到具有相应存储区域的节点上
│       │   ├── preemption/ # 包含了实现抢占逻辑的相关代码
│       │   ├── runtime/ # Scheduling Framework 的运行时实现（最核心的代码逻辑）
│       │   │   ├── framework.go
│       │   │   ├── instrumented_plugins.go 
│       │   │   ├── registry.go
│       │   │   └── waiting_pods_map.go
│       │   └── types.go # 包含了 Scheduling Framework 用到的类型定义及方法实现
│       ├── internal/ # 调度器内部包（仅用于调度器实现）
│       │   ├── cache/ # 调度器缓存实现
│       │   │   ├── cache.go
│       │   │   ├── debugger/ # 包含了一些 Debug 方法，用来检查和更新缓存（用于 debug 目的）
│       │   │   ├── interface.go # 定义了实现 Cache 接口的方法列表
│       │   │   ├── node_tree.go # 包含了 nodeTree 类型的结构体定义，于存储每个区域中的节点名称 
│       │   │   └── snapshot.go # 包含了创建和管理缓存快照的代码实现
│       │   ├── heap/ # 包含了一个标准的堆实现，以及对堆进行操作的方法
│       │   └── queue/ # 定义 SchedulingQueue 类型的接口，接口定义了操作待调度 Pod 队列需要实现的方法。目录中也包含了一个具体的 SchedulingQueue 实现：优先级队列。
│       ├── metrics/ # 包含了调度器指标记录相关实现（变量、方法、函数等）
│       ├── profile/ # 包含了根据配置创建 framework.Framework 类型实例的 NewMap 函数
│       ├── schedule_one.go # 处理单个调度循环的实现
│       ├── scheduler.go # 调度器实现的核心逻辑
│       ├── testing/ # 包含了测试 Scheduling Framework 的代码实现
│       └── util/ # 包含了一些工具类的函数
└── staging/
    └── src/
        └── k8s.io/
            └── kube-scheduler/
                ├── config/
                │   └── v1/
                │       ├── register.go # 包含了 v1 版本资源对象注册逻辑
                │       ├── types.go # 包含 v1 版本 KubeSchedulerConfiguration 的定义
                │       ├── types_pluginargs.go # 包含了 v1 版本调度器插件相关资源对象的定义
                │       └── zz_generated.deepcopy.go # 使用 deepcopy-gen 工具生生成的代码用来深拷贝 v1 版本的资源对象
                ├── extender/ # 包含了 scheduler extender 相关的类型定义
                │   └── v1/
                │       ├── types.go # v1 版本的类型定义，例如：ExtenderPreemptionArgs、ExtenderPreemptionResult、ExtenderBindingArgs、ExtenderBindingResult 等
                │       └── zz_generated.deepcopy.go # deepcopy-gen 工具生成的代码
