package global

const (
	PlatFormLaiDa    int = iota // 莱达互认平台
	PlatFormMingTian            // 明天医网互认平台
)

type ObjectData struct {
	Uid_Enc            string // 患者唯一ID
	Report_update_time string
	Count              int // 执行次数
}

// 报告相关信息
type PatientReportInfo struct {
	AccessionNumber     string // 检查号
	StudyInstanceUid    string // 影像检查唯一ID
	Clinic_Id           string // 门( 急) 诊号 住院号
	HisSn               string // 电子申请单号
	PatientTypeCode     string // 患者类别code 1-门诊、2-急诊、3-住院、9-其他
	PatientTypeName     string // 患者类别名 门诊、急诊、住院、其他
	GenderCode          string // 性别code 0-未知的性别、1-男性、2-女性、9-未说明的性别
	GenderName          string // 性别名  0-未知、1-男性、2-女性、9-其他
	ApplyDepartmentName string // 检查申请科室
	ApplyDoctorName     string // 检查申请医生
	DeviceCode          string // 设备
	PatientSection      string // 病区名称
	SickbedIndex        string // 病床号
	ModalityName        string // 检查类别
	CheckItems          string // 检查项目
	CheckItemsCode      string // 检查项目代码
	BodyPart            string // 检查部位名称
	CheckDoctor         string // 检查医生
	StudyTime           string // 检查时间
	ReportDoctor        string // 报告医生
	ReportTime          string // 报告时间
	AuditDoctor         string // 审核医生
	AuditTime           string // 审核时间
	Finding             string // 检查报告结果-客观所见
	Conclusion          string // 检查报告结果-主观提示
	Status              string // 检查阳性标记
}

// 患者基本信息
type BasicPatientInfo struct {
	OrganizationCode         string //医疗机构组织机构代码 22位的医疗组织机构代码
	OrganizationName         string //机构名称
	Name                     string //患者姓名
	NamePY                   string //患者姓名拼音
	Gender                   string //性别
	Age                      int    //年龄
	BirthTime                string //出生日期 格式：yyyyMMdd格式为20120426
	Telecom                  string //联系电话
	StreetAddressLine        string //地址
	ResidentHealthCardNumber string //居民健康卡号
	InteractionId            string //病人所属医疗机构的域控代码
	IDNumber                 string //身份证号
	PatientID                string //病人号，在域控ID下唯一
	CreationTime             string //消息创建时间 格式为20190702152926
}

// CDA文档相关信息
type CDAInfo struct {
	BPInfo BasicPatientInfo
	RInfo  PatientReportInfo
}

// 文档信息
type DocInfo struct {
	Object    ObjectData
	CDA       CDAInfo
	CDAFileID string //患者CDA文档唯一id
	CDAData   string //CDA baseb4编码
	PDFFileID string //患者PDF文档唯一id
	PDFData   string //PDF baseb4编码
	KOSFileID string //患者KOS文档唯一id
	KOSData   string //KOS baseb4编码
}

var (
	ObjectDataChan chan ObjectData
	RunStatus      bool // 当前获取的数据是否运行完成
)
