package model

import (
	"WowjoyProject/DataSharePlatForm/pkg/setting"
	"database/sql"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 任务数据
type DBData struct {
	uid_enc     sql.NullString // 患者唯一id
	update_time sql.NullString // 报告更新时间
}

// 患者基本信息
type PatientBaseInfo struct {
	name           sql.NullString // 患者姓名
	spell_name     sql.NullString // 患者姓名拼音
	sex_code       sql.NullString // 性别
	telephone      sql.NullString // 联系电话
	address        sql.NullString // 地址
	society_number sql.NullString // 医保号
	id_card        sql.NullString // 身份证
	patient_id     sql.NullString // 患者唯一id
}

type ReportInfo struct {
	accession_number      sql.NullString // 检查号
	study_instance_uid    sql.NullString // 影像检查唯一ID
	clinic_id             sql.NullString // 门( 急) 诊号 住院号
	his_sn                sql.NullString // 电子申请单号
	patient_type_code     sql.NullString // 患者类别
	apply_department_name sql.NullString // 检查申请科室
	apply_doctor_name     sql.NullString // 检查申请医生
	device_code           sql.NullString // 设备
	patient_section       sql.NullString // 病区名称
	sickbed_index         sql.NullString // 病床号
	modality_name         sql.NullString // 检查类别
	check_items           sql.NullString // 检查项目
	check_items_code      sql.NullString // 检查项目代码
	body_part             sql.NullString // 检查部位名称
	check_doctor          sql.NullString // 检查医生
	study_time            sql.NullString // 检查时间
	report_doctor         sql.NullString // 报告医生
	report_time           sql.NullString // 报告时间
	audit_doctor          sql.NullString // 审核医生
	audit_time            sql.NullString // 审核时间
	finding               sql.NullString // 检查报告结果-客观所见
	conclusion            sql.NullString // 检查报告结果-主观提示
	status                sql.NullString //检查阳性标记
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*sql.DB, error) {
	db, err := sql.Open(databaseSetting.DBType, databaseSetting.DBConn)
	if err != nil {
		return nil, err
	}
	// 数据库最大连接数
	db.SetConnMaxLifetime(time.Duration(databaseSetting.MaxLifetime) * time.Minute)
	db.SetMaxOpenConns(databaseSetting.MaxIdleConns)
	db.SetMaxIdleConns(databaseSetting.MaxIdleConns)

	return db, nil
}
