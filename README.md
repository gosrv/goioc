# goioc
go ioc framework

# star it && QQ群：869428810

### bean
放入容器中的对象就是一个bean，容器可以针对bean的成员变量做一些注入操作

### bean容器builder
用来构建bean容器，并完成bean的注入，bean的注入操作由ITagProcessor完成  
构建过程
1. 收集满足条件的bean
2. 获取所有的ITagProcessor，并按优先级排序
3. 按ITagProcessor的优先级先后处理所有满足条件的bean

### bean的有效条件
通过内嵌接口IBeanCondition来实现条件判断
1. NewConditionAnd() 只有所有的条件都满足才能通过
2. NewConditionOr() 有一个条件满足就能通过
3. NewConditionOnValue() 当配置文件中存在配置项时才能通过
4. NewConditionOnBeanType() 当存在某个类型的bean时才能通过
5. NewConditionOnBeanName() 当存在某个名字的bean时才能通过

### 系统级tag
标注在bean上，可以执行一些特定意义的注入，目前支持两种系统级tag，bean和cfg

### tag bean
1. "bean"
通过类型注入，会注入其它同类型的bean，如果不是slice，要求容器内有且仅有一个次类型的bean，如果是slice，则可以有任意个
2. "bean.name"
通过bean名字注入，必须存在并且只能有一个
3. "bean.required"
如果没找到会不会报错，默认是true

### tag cfg
1. "cfg"
注入配置文件中的配置项
2. "cfg.d"
注入特定域的配置项，与IConfigBase配合使用,调整配置注入的根节点
3. "cfg.default"
如果配置文件中不存，则使用的默认配置

### ITagProcessor
可以通过实现接口ITagProcessor来自定义tag注入器，你还可以为它指定一个优先级，如果它没有任何依赖，则可以设置为系统级tag

### ITagParser
解析tag标签，容器中只能存在一个，而且必须有一个

### exaple
请查看example1.go/example2.go/example3.go/example4.go