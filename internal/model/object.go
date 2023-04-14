package model

import (
	"WowjoyProject/DataSharePlatForm/global"
	"strconv"
	"strings"
	"time"
)

// 获取需要处理的数据
func AutoGetObjectData() {
	if global.RunStatus {
		global.Logger.Info("上次获取的数据没有消耗完，等待消耗完，再获取数据....")
		return
	}
	global.RunStatus = true
	global.Logger.Info("***开始获取需要上传互认平台的数据***")
	// 互认数据只互认放射科的数据，只互认报告状态为已经审核的报告
	// 增加延时10分钟上传数据
	sql := `select r.uid_enc,r.update_time
	from report r
	left join platform_share_info psi ON r.uid_enc = psi.uid_enc
	where psi.uid_enc is null
	and r.report_status ='AUDITED' and r.uid_enc!='' and r.update_time > ?
	and TIMESTAMPDIFF(MINUTE,r.update_time,NOW()) > 10
	order by r.update_time DESC
	limit ?;`
	rows, err := global.ReadDBEngine.Query(sql, global.ObjectSetting.Start_Time, global.ObjectSetting.Object_MaxTasks)
	if err != nil {
		global.Logger.Error(err)
		global.RunStatus = false
		return
	}
	defer rows.Close()
	for rows.Next() {
		key := DBData{}
		_ = rows.Scan(&key.uid_enc, &key.update_time)
		if key.uid_enc.String != "" {
			data := global.ObjectData{
				Uid_Enc:            key.uid_enc.String,
				Report_update_time: key.update_time.String,
				Count:              1,
			}
			global.ObjectDataChan <- data
		}
	}
	global.RunStatus = false
}

func Test() {
	data := global.ObjectData{
		Uid_Enc:            "7faf67f1b43f065c85a0f3d7b2e37f73",
		Report_update_time: "2021-11-22 10:05:03",
		Count:              1,
	}
	global.ObjectDataChan <- data

}

// 更新注册基本信息状态
func UpdateStatus(enc string, code int) {
	global.Logger.Info("***开始跟新注册基本信息状态***")
	sql := `update platform_share_info psi set psi.status = ? where psi.uid_enc = ?;`
	global.WriteDBEngine.Exec(sql, code, enc)
}

// 更新文档id
func UpdateFileUid(enc, cdaid, pdfid, kosid string) {
	global.Logger.Info("***开始跟新文档唯一ID***")
	sql := `update platform_share_info psi set psi.cda_uuid = ?,psi.pdf_uuid = ?,psi.kos_uuid = ? where psi.uid_enc = ?;`
	global.WriteDBEngine.Exec(sql, cdaid, pdfid, kosid, enc)
}

// 插入需要处理的数据进入数据库
func InsertData(uidenc, updatetime string) {
	global.Logger.Info("***开始在数据库中插入处理的数据***", uidenc)
	sql := `insert into platform_share_info (uid_enc,report_update_time) values(?,?)`
	global.WriteDBEngine.Exec(sql, uidenc, updatetime)
}

// 获取患者基础信息
func GetBasicInfo(uidenc, updatetime string) global.BasicPatientInfo {
	obj := global.BasicPatientInfo{}
	global.Logger.Info("***开始获取患者基本信息***")
	sql := `select rp.name,rp.spell_name,rp.sex_code,rp.telephone,rp.address,rp.society_number,rp.id_card,rp.patient_number
	from register_patient rp
	inner join report r on r.patient_id = rp.patient_id
	where r.uid_enc = ? limit 1;`
	row := global.ReadDBEngine.QueryRow(sql, uidenc)
	key := PatientBaseInfo{}
	if err := row.Scan(&key.name, &key.spell_name, &key.sex_code, &key.telephone, &key.address, &key.society_number, &key.id_card, &key.patient_id); err != nil {
		global.Logger.Error(err)
		return obj
	}
	if key.name.String == "" || key.spell_name.String == "" || key.id_card.String == "" {
		global.Logger.Error("***患者基本信息有为空项***")
		return obj
	}
	updatetime = strings.Replace(updatetime, "-", "", -1)
	updatetime = strings.Replace(updatetime, ":", "", -1)
	updatetime = strings.Replace(updatetime, " ", "", -1)
	brith := key.id_card.String[6:14]
	gender := ""
	switch global.ObjectSetting.Object_SelectPlatForm {
	case global.PlatFormLaiDa:
		switch key.sex_code.String {
		case "MAN":
			gender = "M"
		case "FEMALE":
			gender = "F"
		case "FEMAIL":
			gender = "F"
		case "":
			gender = "O"
		}
	case global.PlatFormMingTian:
		switch key.sex_code.String {
		case "MAN":
			gender = "男"
		case "FEMALE":
			gender = "女"
		case "FEMAIL":
			gender = "未知"
		case "":
			gender = "未知"
		}
	}
	now := time.Now()
	now_year := now.Year()
	idCard := key.id_card.String                 // 年
	idcard_year, _ := strconv.Atoi(idCard[6:10]) // 年
	age := now_year - idcard_year
	obj = global.BasicPatientInfo{
		OrganizationCode:         global.ObjectSetting.Object_OrganizationCode,
		OrganizationName:         global.ObjectSetting.Object_OrganizationName,
		Name:                     key.name.String,
		NamePY:                   key.spell_name.String,
		Gender:                   gender,
		Age:                      age,
		BirthTime:                brith,
		Telecom:                  key.telephone.String,
		StreetAddressLine:        key.address.String,
		ResidentHealthCardNumber: key.society_number.String,
		InteractionId:            global.ObjectSetting.Object_InteractionId,
		IDNumber:                 key.id_card.String,
		PatientID:                key.patient_id.String,
		CreationTime:             updatetime,
	}
	return obj
}

// 获取患者报告相关信息
func GetReportInfo(uidenc string) global.PatientReportInfo {
	obj := global.PatientReportInfo{}
	global.Logger.Info("***开始获取患者报告信息***")
	sql := `select roi.accession_number,s.study_instance_uid,ri.clinic_id,roi.his_sn,ri.patient_type_code,roi.apply_department_name,roi.apply_doctor_name,
	roi.device_code,ri.patient_section,ri.sickbed_index,roi.modality_name,roi.check_items,roi.check_items_code,roi.bodypart,rir.inspect_doctor_name,r.study_time,
	r.report_doctor,r.report_time,r.audit_doctor,r.audit_time,r.finding,r.conclusion,r.status
	from report r 
	left join register_info ri on r.uid_enc = ri.uid_enc
	left join register_order_info roi on ri.register_uid_enc = roi.uid_enc
	left join study s on r.uid_enc = s.uid_enc
	left join register_info_relation rir on rir.register_uid_enc = r.uid_enc
	where r.uid_enc = ?;`
	row := global.ReadDBEngine.QueryRow(sql, uidenc)
	key := ReportInfo{}
	if err := row.Scan(&key.accession_number, &key.study_instance_uid, &key.clinic_id, &key.his_sn, &key.patient_type_code, &key.apply_department_name, &key.apply_doctor_name,
		&key.device_code, &key.patient_section, &key.sickbed_index, &key.modality_name, &key.check_items, &key.check_items_code, &key.body_part, &key.check_doctor, &key.study_time, &key.report_doctor,
		&key.report_time, &key.audit_doctor, &key.audit_time, &key.finding, &key.conclusion, &key.status); err != nil {
		global.Logger.Error(err)
		return obj
	}
	patienttypecode := key.patient_type_code.String
	patienttypename := ""
	switch patienttypecode {
	case "OP": // 门诊
		patienttypecode = "1"
		patienttypename = "门诊"
	case "EM": // 急诊
		patienttypecode = "2"
		patienttypename = "急诊"

	case "IH": // 住院
		patienttypecode = "3"
		patienttypename = "住院"
	default:
		patienttypecode = "9"
		patienttypename = "其他"
	}
	modalityname := key.modality_name.String

	switch global.ObjectSetting.Object_SelectPlatForm {
	case global.PlatFormLaiDa:
		switch modalityname {
		case "X-Ray":
			modalityname = "1"
		case "DR":
			modalityname = "2"
		case "CT":
			modalityname = "3"
		case "MR":
			modalityname = "4"
		case "DSA":
			modalityname = "5"
		case "US":
			modalityname = "6"
		case "ES":
			modalityname = "7"
		case "PA":
			modalityname = "8"
		case "NM":
			modalityname = "9"
		case "PET":
			modalityname = "10"
		default:
			modalityname = "99"
		}
	case global.PlatFormMingTian:
	}

	studytime := key.study_time.String
	if studytime != "" {
		studytime = strings.Replace(studytime, "-", "", -1)
		studytime = strings.Replace(studytime, ":", "", -1)
		studytime = strings.Replace(studytime, " ", "", -1)
	}
	reporttime := key.study_time.String
	if studytime != "" {
		reporttime = strings.Replace(reporttime, "-", "", -1)
		reporttime = strings.Replace(reporttime, ":", "", -1)
		reporttime = strings.Replace(reporttime, " ", "", -1)
	}
	audittime := key.study_time.String
	if studytime != "" {
		audittime = strings.Replace(audittime, "-", "", -1)
		audittime = strings.Replace(audittime, ":", "", -1)
		audittime = strings.Replace(audittime, " ", "", -1)
	}

	// 2022/08/16修改 增加互认项目，多个以";"间隔
	check_item := key.check_items.String
	if check_item != "" {
		check_item = strings.Replace(check_item, ",", ";", -1)
		check_item = strings.Replace(check_item, "|", ";", -1)
	}

	// 增加互认项目字典表
	check_items_str := strings.Split(check_item, ";")
	check_items_len := len(check_items_str)
	var check_items_code, check_items string
	// 通过his 检查项目获取互认标准检查项目
	for i := 0; i < check_items_len; i++ {
		if i > 0 {
			check_items_code += ";"
			check_items += ";"
		}
		check_items_code, check_items = GetCheckItem(check_items_str[i])
	}
	global.Logger.Info("互认标准项目编码code: ", check_items_code, " name: ", check_items)

	// check_items_code := key.check_items_code.String
	// if check_items_code != "" {
	// 	check_items_code = strings.Replace(check_items_code, ",", ";", -1)
	// 	check_items_code = strings.Replace(check_items_code, "|", ";", -1)
	// }

	if !key.check_doctor.Valid {
		key.check_doctor.String = key.audit_doctor.String
	}
	obj = global.PatientReportInfo{
		AccessionNumber:     key.accession_number.String,
		StudyInstanceUid:    key.study_instance_uid.String,
		Clinic_Id:           key.clinic_id.String,
		HisSn:               key.his_sn.String,
		PatientTypeCode:     patienttypecode,
		PatientTypeName:     patienttypename,
		ApplyDepartmentName: key.apply_department_name.String,
		ApplyDoctorName:     key.apply_doctor_name.String,
		DeviceCode:          key.device_code.String,
		PatientSection:      key.patient_section.String,
		SickbedIndex:        key.sickbed_index.String,
		ModalityName:        modalityname,
		CheckItems:          check_items,
		CheckItemsCode:      check_items_code,
		BodyPart:            key.body_part.String,
		CheckDoctor:         key.check_doctor.String,
		StudyTime:           studytime,
		ReportDoctor:        key.report_doctor.String,
		ReportTime:          reporttime,
		AuditDoctor:         key.audit_doctor.String,
		AuditTime:           audittime,
		Finding:             key.finding.String,
		Conclusion:          key.conclusion.String,
		Status:              key.status.String,
	}
	return obj
}

// 获取影像的路径
func GetImagePath(uidenc string) []string {
	filepath := []string{}
	global.Logger.Info("***开始获取影像的路径***")
	sql := `select ins.file_name,sl.ip,sl.s_virtual_dir
	from instance ins 
	left join study s on s.study_key = ins.study_key
	left join study_location sl on s.location_code = sl.n_station_code
	where s.uid_enc = ?;`
	rows, err := global.ReadDBEngine.Query(sql, uidenc)
	if err != nil {
		global.Logger.Error(err)
		return filepath
	}
	defer rows.Close()
	for rows.Next() {
		filename := ""
		ip := ""
		s_virtual_dir := ""
		_ = rows.Scan(&filename, &ip, &s_virtual_dir)

		if filename != "" && ip != "" && s_virtual_dir != "" {
			filename = "\\\\" + ip + "\\" + s_virtual_dir + "\\" + filename
			index := strings.Index(filename, uidenc)
			substr := filename[:index+len(uidenc)]
			filepath = append(filepath, substr)
		}
	}
	return filepath
}

// 获取PDF文件路径
func GetPDFFilePath(uidenc string) string {
	global.Logger.Info("***开始获取PDF文件的路径***")
	filename := ""
	ip := ""
	s_virtual_dir := ""
	sql := `select pf.file_name,sl.ip,sl.s_virtual_dir from pdf_file pf
	left join study_location sl on pf.localtion_code = sl.n_station_code
	where pf.report_id = ?;`
	row := global.ReadDBEngine.QueryRow(sql, uidenc)
	if err := row.Scan(&filename, &ip, &s_virtual_dir); err != nil {
		global.Logger.Error(err)
		return filename
	}
	if filename != "" && ip != "" && s_virtual_dir != "" {
		filename = "\\\\" + ip + "\\" + s_virtual_dir + filename
	}
	return filename
}

// 获取检查项目
func GetCheckItem(checkItem string) (code, name string) {
	global.Logger.Info("***获取互认检查项目***")
	sql := `SELECT dic.standard_check_item,dic.standard_check_item_code FROM dict_checkitem dic
	WHERE dic.pacs_check_item = ?;`
	row := global.ReadDBEngine.QueryRow(sql, checkItem)
	if err := row.Scan(&name, &code); err != nil {
		global.Logger.Error(err)
		return
	}
	return
}
