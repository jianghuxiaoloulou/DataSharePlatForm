package setting

type GeneralSettingS struct {
	LogSavePath string
	LogFileName string
	LogFileExt  string
	LogMaxSize  int
	LogMaxAge   int
	MaxThreads  int
	CronSpec    string
}

type DatabaseSettingS struct {
	DBConn       string
	DBType       string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  int
}

type ObjectSettingS struct {
	Object_InteractionId        string
	Object_OrganizationName     string
	Object_OrganizationCode     string
	Object_RetrieveAETitle      string
	Object_OrganizationAK       string
	Object_PIX                  string
	Object_CentralizedStorageID string
	Object_CentralizedStorage   string
	Object_MaxTasks             int
	Object_Count                int
	Start_Time                  string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
