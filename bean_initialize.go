package gioc

import (
	"math"
	"sort"
)

type IBeanInit interface {
	BeanInit()
	BeanUninit()
}

type IBeanStart interface {
	BeanStart()
	BeanStop()
}

type BeanInitDriver struct {
	sorted     bool
	BeansInit  []IBeanInit  `bean:""`
	BeansStart []IBeanStart `bean:""`
}

func NewBeanInitDriver() *BeanInitDriver {
	return &BeanInitDriver{
		sorted: false,
	}
}

func (this *BeanInitDriver) sort() {
	if this.sorted {
		return
	}
	this.sorted = true

	sort.Slice(this.BeansInit, func(i, j int) bool {
		ip := math.MaxInt32
		jp := math.MaxInt32
		if iPriority, ok := this.BeansInit[i].(IPriority); ok {
			ip = iPriority.GetPriority()
		}
		if jPriority, ok := this.BeansInit[j].(IPriority); ok {
			jp = jPriority.GetPriority()
		}
		return ip < jp
	})
	sort.Slice(this.BeansStart, func(i, j int) bool {
		ip := math.MaxInt32
		jp := math.MaxInt32
		if iPriority, ok := this.BeansStart[i].(IPriority); ok {
			ip = iPriority.GetPriority()
		}
		if jPriority, ok := this.BeansStart[j].(IPriority); ok {
			jp = jPriority.GetPriority()
		}
		return ip < jp
	})
}

func (this *BeanInitDriver) CallInit() {
	this.sort()
	for _, bean := range this.BeansInit {
		bean.BeanInit()
	}
}

func (this *BeanInitDriver) CallUnInit() {
	for i := len(this.BeansInit) - 1; i >= 0; i-- {
		this.BeansInit[i].BeanUninit()
	}
}

func (this *BeanInitDriver) CallStart() {
	this.sort()
	for _, bean := range this.BeansStart {
		bean.BeanStart()
	}
}

func (this *BeanInitDriver) CallStop() {
	for i := len(this.BeansStart) - 1; i >= 0; i-- {
		this.BeansStart[i].BeanStop()
	}
}
