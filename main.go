package main

import (
	"WowjoyProject/DataSharePlatForm/global"
	"WowjoyProject/DataSharePlatForm/internal/model"
	"WowjoyProject/DataSharePlatForm/pkg/object"
	"WowjoyProject/DataSharePlatForm/pkg/workpattern"

	"github.com/robfig/cron"
)

// @title 浙江省检验检查结果互认共享平台服务
// @version 1.0.0.1
// @description 与莱达对接检验检查结果
// @termsOfService
func main() {
	global.Logger.Info("*******开始运行检验检查结果互认共享平台服务********")

	global.ObjectDataChan = make(chan global.ObjectData)

	// 注册工作池，传入任务
	// 参数1 初始化worker(工人)设置最大线程数
	wokerPool := workpattern.NewWorkerPool(global.GeneralSetting.MaxThreads)
	// 有任务就去做，没有就阻塞，任务做不过来也阻塞
	wokerPool.Run()
	// 处理任务
	go func() {
		for {
			select {
			case data := <-global.ObjectDataChan:
				sc := &Dosomething{key: data}
				wokerPool.JobQueue <- sc
			}
		}
	}()
	// 后台获取需要处理的任务
	global.RunStatus = false
	// object.WritePatientRegistryAddXML_MingTian()

	run()
}

func run() {
	// 方式一：
	// for {
	// 	// 获取需要注册的任务（检查报告经过审核后）
	// 	model.AutoGetObjectData()
	// 	time.Sleep(time.Second * 15)
	// }
	// 方式二：获取任务(定时任务)
	MyCron := cron.New()
	MyCron.AddFunc(global.GeneralSetting.CronSpec, func() {
		global.Logger.Info("开始执行定时任务")
		model.AutoGetObjectData()
	})
	MyCron.Start()
	defer MyCron.Stop()
	select {}
}

type Dosomething struct {
	key global.ObjectData
}

func (d *Dosomething) Do() {
	global.Logger.Info("正在处理的数据是：", d.key)
	//处理获取的任务
	object.Work(&d.key)
}
