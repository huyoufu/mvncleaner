package main

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

type MCMetric struct {
	version string
	//开始时间
	start int64
	//结束时间
	end int64
	//花费的毫秒值
	costTime int64
	//检索文件个数
	sumFileNum int64
	//没问题的文件数
	commonFileNum int64
	//错误的文件数
	errFileNum int64
	//索引文件数
	indexFileNum int64
	//jar文件数
	jarFileNum int64
}

type StepInfo struct {
	cmt  string
	time int64
}

type TimerTrack struct {
	topic string
	//是否允许重写时间
	overwrite bool
	//开始时间
	start int64
	//中间步骤时间
	steps []StepInfo
	//步骤数
	stepNum int32
	//步骤名字中缀
	//stepName=cmt+Mid+stepNum
	stepNameMid string
	//结束时间
	end int64
	//花费的时间 如果设置isMs 为true 则得到是毫秒 否则是纳秒
	costTime int64
	//
	isMs bool
}

var DefaultTimerTrack *TimerTrack = &TimerTrack{
	stepNameMid: "step-",
	steps:       []StepInfo{},
	isMs:        true,
}

//步数增长
func (t *TimerTrack) NextStepNum() int32 {
	addStepNum := atomic.AddInt32(&t.stepNum, 1)
	return addStepNum
}

//开始计时
func (t *TimerTrack) Start() {
	if t.start != 0 {
		if t.overwrite {
			//如果不为0 且容许覆盖则 修改为现在时间
			t.start = time.Now().UnixNano()
		} else {
			t.Step4cmt("start overwrite")
		}
	} else {
		//为零则开始计时
		t.start = time.Now().UnixNano()
	}
}
func (t *TimerTrack) Step() {
	t.Step4cmt("")
}
func (t *TimerTrack) Step4cmt(cmt string) {
	stepInfo := StepInfo{
		cmt:  cmt + t.stepNameMid + strconv.Itoa(int(t.NextStepNum())),
		time: time.Now().UnixNano(),
	}
	t.steps = append(t.steps, stepInfo)
}

//结束计时
func (t *TimerTrack) End() {
	if t.end != 0 {
		if t.overwrite {
			//如果不为0 且容许覆盖则 修改为现在时间
			t.end = time.Now().UnixNano()
		} else {
			t.Step4cmt("end overwrite")
		}
	} else {
		//为零则开始计时
		t.end = time.Now().UnixNano()
	}
}

//返回共花费了多久
func (t *TimerTrack) Cost() int64 {
	if t.end == 0 {
		t.End()
	}
	if t.isMs {
		t.costTime = (t.end - t.start) / 1e6
	} else {
		t.costTime = t.end - t.start
	}

	return t.costTime
}

func (t *TimerTrack) PrintBeautiful() {

	//formatter :="2006-01-02 15:04:05.999999999 -0700 MST";
	formatter := "2006-01-02 15:04:05.999 -0700 MST"
	if t.start != 0 {
		//开始了

		fmt.Printf("start is: %s\n", time.Unix(0, t.start).Format(formatter))

		if t.stepNum > 0 {
			for _, stepInfo := range t.steps {
				fmt.Printf("    %s---%s\n", stepInfo.cmt, time.Unix(0, stepInfo.time).Format(formatter))
			}
		}

		fmt.Printf("end is: %s\n", time.Unix(0, t.end).Format(formatter))
		t.Cost()
		fmt.Printf("cost is %d ms\n", t.Cost())
	}

}
