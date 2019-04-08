package gioc

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

func (this *beanTagProcessor) PrepareProcess() {

}

var _ ITagProcessor = (*beanTagProcessor)(nil)

func NewBeanTagProcessor(beanContainer IBeanContainer) ITagProcessor {
	return &beanTagProcessor{beanContainer: beanContainer}
}

func (this *beanTagProcessor) TagProcessorName() string {
	return BeanTagProcessor
}

func (this *beanTagProcessor) TagProcessorPriority() int {
	return PrioritySystem
}

func (this *beanTagProcessor) TagProcess(bean interface{}, fType reflect.StructField, fValue reflect.Value, tags map[string]string) {
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
		fieldTypeBeans := this.beanContainer.GetBeanByType(fValue.Type())
		if len(fieldTypeBeans) == 0 && fValue.Type().Kind() == reflect.Slice {
			elements := this.beanContainer.GetBeanByType(fValue.Type().Elem())
			svalue := reflect.New(fValue.Type()).Elem()
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
					fValue.Type(), len(fieldTypeBeans), reflect.TypeOf(bean))
			}
		}
	}

	if fieldBean != nil {
		fValue.Set(reflect.ValueOf(fieldBean))
	} else if beanRequire {
		util.Panic("required bean %v in instance %v has not found", fValue.Type(), reflect.TypeOf(bean))
	}
}
