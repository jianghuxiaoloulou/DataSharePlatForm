package global

import (
	"WowjoyProject/DataSharePlatForm/pkg/logger"
	"WowjoyProject/DataSharePlatForm/pkg/setting"
)

var (
	GeneralSetting  *setting.GeneralSettingS
	DatabaseSetting *setting.DatabaseSettingS
	ObjectSetting   *setting.ObjectSettingS
	Logger          *logger.Logger
)
