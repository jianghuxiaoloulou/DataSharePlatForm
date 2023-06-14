package init

import (
	"WowjoyProject/DataSharePlatForm/global"
	"WowjoyProject/DataSharePlatForm/internal/model"
	"WowjoyProject/DataSharePlatForm/pkg/logger"
	"WowjoyProject/DataSharePlatForm/pkg/setting"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitSetup() {
	ReadSetup()
}

func SetupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("General", &global.GeneralSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Object", &global.ObjectSetting)
	if err != nil {
		return err
	}
	return nil
}

func SetupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.GeneralSetting.LogSavePath + "/" + global.GeneralSetting.LogFileName + global.GeneralSetting.LogFileExt,
		MaxSize:   global.GeneralSetting.LogMaxSize,
		MaxAge:    global.GeneralSetting.LogMaxAge,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func SetupReadDBEngine() error {
	var err error
	global.ReadDBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func SetupWriteDBEngine() error {
	var err error
	global.WriteDBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func ReadSetup() {
	err := SetupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	err = SetupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
	err = SetupReadDBEngine()
	if err != nil {
		log.Fatalf("init.setupReadDBEngine err: %v", err)
	}
	err = SetupWriteDBEngine()
	if err != nil {
		log.Fatalf("init.setupWriteDBEngine err: %v", err)
	}
}
