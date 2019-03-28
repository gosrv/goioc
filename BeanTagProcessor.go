package goioc

import (
	"github.com/gosrv/goioc/util"
	"reflect"
)

/**
bean相关tag处理
*/
const (
	BeanTagProcessor = "bean"
	// 根据类型装配
	BeanTag = "bean"
	// 根据名字装配
	BeanNameTag = "bean.name"
	// 必须存在，如果装配时不存在会报错，默认是true
	BeanRequiredTag = "bean.required"
)

type beanTagProcessor struct {
	beanContainer IBeanContainer
}

func NewBeanTagProcessor(beanContainer IBeanContainer) ITagProcessor {
	return &beanTagProcessor{beanContainer: beanContainer}
}

func (this *beanTagProcessor) TagProcessorName() string {
	return BeanTagProcessor
}

func (this *beanTagProcessor) GetPriority() int {
	return PrioritySystem
}

func (this *beanTagProcessor) TagProcess(bean interface{}, field reflect.Value, tags map[string]string) {
	_, beanOk := tags[BeanTag]
	beanName, beanNameOk := tags[BeanNameTag]
	beanRequiredVal, beanRequireOk := tags[BeanRequiredTag]
	beanRequire := true
	if beanRequireOk {
		beanRequire = beanRequiredVal == "true"
	}
	if beanNameOk || beanRequireOk {
		beanOk = true
	}
	if !beanOk {
		return
	}

	var fieldBean interface{} = nil
	if beanNameOk {
		fieldBean = this.beanContainer.GetBeanByName(beanName)
	} else {
		fieldTypeBeans := this.beanContainer.GetBeanByType(field.Type())
		if len(fieldTypeBeans) == 0 && field.Type().Kind() == reflect.Slice {
			elements := this.beanContainer.GetBeanByType(field.Type().Elem())
			svalue := reflect.New(field.Type()).Elem()
			for _, ab := range elements {
				svalue = reflect.Append(svalue, reflect.ValueOf(ab))
			}
			fieldBean = svalue.Interface()
		} else {
			switch len(fieldTypeBeans) {
			case 0:
			case 1:
				fieldBean = fieldTypeBeans[0]
			default:
				util.Panic("ambiguous bean type [%v] injct except 1 but find %v in bean %v",
					field.Type(), len(fieldTypeBeans), reflect.TypeOf(bean))
			}
		}
	}

	if fieldBean != nil {
		field.Set(reflect.ValueOf(fieldBean))
	} else if beanRequire {
		util.Panic("required bean %v in instance %v has not found", field.Type(), reflect.TypeOf(bean))
	}
}
