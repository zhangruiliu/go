package lsent

type Sampling struct {
	Id              int64  `json:"id,omitempty"`
	Server          int64  `json:"server,omitempty"` //0
	Uploaded        int64  `json:"uploaded,omitempty"`
	Region          int64  `json:"region,omitempty"`
	Category        int64  `json:"category,omitempty"`
	Perid           int64  `json:"perid,omitempty"`
	Commitid        int64  `json:"commitid,omitempty"`   //0
	Status          int64  `json:"status,omitempty"`     //0
	Code            string `json:"code,omitempty"`       //样本DNA编号
	Codebag         string `json:"codebag,omitempty"`    //索引编号
	Codeserial      string `json:"codeserial,omitempty"` //样本流水号
	Cabno           string `json:"cabno,omitempty"`
	Name            string `json:"name,omitempty"`
	Kind            string `json:"kind,omitempty"`            //词典表category=33的第一个
	Samplingaddr    string `json:"samplingaddr,omitempty"`    //采样地点
	Samplingunit    string `json:"samplingunit,omitempty"`    //采样单位
	Samplingoprator string `json:"samplingoprator,omitempty"` //采样人
	Samplingdate    string `json:"samplingdate,omitempty"`    //采样时间
	Samplingmemo    string `json:"samplingmemo,omitempty"`    //采样备注
	Checkingno      string `json:"checkingno,omitempty"`
	Checkingunit    string `json:"checkingunit,omitempty"`
	Checker         string `json:"checker,omitempty"`
	Checkingdate    string `json:"checkingdate,omitempty"`
	// JailLocation     string    `json:"jail_location,omitempty"`
	Jailno           string `json:"jailno,omitempty"` //物证编号
	Relation         string `json:"relation,omitempty"`
	Idno             string `json:"idno,omitempty"`
	Folk             string `json:"folk,omitempty"`
	Age              int64  `json:"age,omitempty"`
	Birthregion      string `json:"birthregion,omitempty"`
	Keeper           string `json:"keeper,omitempty"`
	Insertoperatorid int64  `json:"insertoperatorid,omitempty"` //插入人
	Insertdatetime   string `json:"insertdatetime,omitempty"`   //插入时间
	Updateoperatorid int64  `json:"updateoperatorid,omitempty"` //更新人
	Lastupdatetime   string `json:"lastupdatetime,omitempty"`   //更新时间
	Accepttime       string `json:"accepttime,omitempty"`
	Important        int64  `json:"important,omitempty"`
	Deleted          int64  `json:"deleted,omitempty"` //0
}

// type Sampling struct {
// 	ResidenceRegionalism   string    `json:"residence_regionalism,omitempty" sql:"field" xml:"residenceRegionalism"`      //现住址行政区划
// 	ResidenceAddr          string    `json:"residence_addr,omitempty" sql:"field" xml:"residenceAddr"`                    //现住址名称
// 	FingerprintNo          string    `json:"fingerprint_no,omitempty" sql:"field" xml:"fingerprintNo"`                    //指纹编号
// 	BloodType              string    `json:"blood_type,omitempty" sql:"field" xml:"bloodType"`                            //血型
// 	Height                 string    `json:"height,omitempty" sql:"field" xml:"height"`                                   //身高
// 	BodilyForm             string    `json:"bodily_form,omitempty" sql:"field" xml:"bodilyForm"`                          //体型
// 	BirthDateFrom          string `json:"birth_date_from,omitempty" sql:"field" xml:"birthDateFrom"`                   //出生日期（起）
// 	BirthDateTo            string `json:"birth_date_to,omitempty" sql:"field" xml:"birthDateTo"`                       //出生日期（止）
// 	RoughAge               string    `json:"rough_age,omitempty" sql:"field" xml:"roughAge"`                              //大致年龄
// 	RoughlyBirthdate       string `json:"roughly_birthdate,omitempty" sql:"field" xml:"roughlyBirthdate"`              //大致出生日期
// 	MissingType            string    `json:"missing_type,omitempty" sql:"field" xml:"missingType"`                        //失踪类型
// 	MissingTime            string `json:"missing_time,omitempty" sql:"field" xml:"missingTime"`                        //失踪时间
// 	MissingPlace           string    `json:"missing_place,omitempty" sql:"field" xml:"missingPlace"`                      //失踪地点
// 	ExtrinsicSign          string    `json:"extrinsic_sign,omitempty" sql:"field" xml:"extrinsicSign"`                    //体表标记
// 	SpecialSign            string    `json:"special_sign,omitempty" sql:"field" xml:"specialSign"`                        //特殊特征
// 	InvolvedCaseName       string    `json:"involved_case_name,omitempty" sql:"field" xml:"involvedCaseName"`             //涉案名称
// 	InvolvedcaseNo         string    `json:"involved_case_no,omitempty" sql:"field" xml:"involvedCaseNo"`                 //涉案编号
// 	CaseProperty           string    `json:"case_property,omitempty" sql:"field" xml:"caseProperty"`                      //涉案性质
// 	PrisonType             string    `json:"prison_type,omitempty" sql:"field" xml:"prisonType"`                          //关押地类型
// 	PrisonNo               string    `json:"prison_no,omitempty" sql:"field" xml:"prisonNo"`                              //科所队编号
// 	DeathFlag              int64     `json:"death_flag" sql:"field" xml:"deathFlag"`                                      //死亡标识
// 	DeleteFlag             int64     `json:"delete_flag,omitempty" sql:"field" xml:"deleteFlag"`                          //删除标识
// 	IndexFlag              string    `json:"index_flag,omitempty" sql:"field" xml:"indexFlag"`                            //索引标识
// 	TransferFlag           int64     `json:"transfer_flag,omitempty" sql:"field" xml:"transferFlag"`                      //上报标识，0不上报，1待上报，2成功，3等待重试，4错误，5正在上报
// 	TransferUser           string    `json:"transfer_user,omitempty" sql:"field" xml:"transferUser"`                      //上报人
// 	TransferDatetime       string `json:"transfer_datetime,omitempty" sql:"field" xml:"transferDatetime"`              //上报时间
// 	DataSource             string    `json:"data_source,omitempty" sql:"field" xml:"dataSource"`                          //数据来源
// 	DataLevel              int64     `json:"data_level,omitempty" sql:"field" xml:"dataLevel"`                            //数据级别
// 	Reserve1               string    `json:"reserve1,omitempty" sql:"field" xml:"reserve1"`                               //备用字段1
// 	Reserve2               string    `json:"reserve2,omitempty" sql:"field" xml:"reserve2"`                               //备用字段2
// 	Reserve3               string    `json:"reserve3,omitempty" sql:"field" xml:"reserve3"`                               //备用字段3
// 	Reserve4               string    `json:"reserve4,omitempty" sql:"field" xml:"reserve4"`                               //备用字段4
// 	Reserve5               string    `json:"reserve5,omitempty" sql:"field" xml:"reserve5"`                               //备用字段5
// 	Reserve6               string    `json:"reserve6,omitempty" sql:"field" xml:"reserve6"`                               //备用字段6
// 	Remark                 string    `json:"remark,omitempty" sql:"field" xml:"remark"`                                   //备注
// 	CreateUser             string    `json:"create_user,omitempty" sql:"field" xml:"createUser"`                          //创建人
// 	CreateDatetime         string `json:"create_datetime,omitempty" sql:"field" xml:"createDatetime"`                  //创建时间
// 	UpdateUser             string    `json:"update_user,omitempty" sql:"field" xml:"updateUser"`                          //更新人
// 	UpdateDatetime         string `json:"update_datetime,omitempty" sql:"field" xml:"updateDatetime"`                  //更新时间
// 	PersonLabel            string    `json:"person_label,omitempty" sql:"field" xml:"personLabel"`                        //人员标签
// 	AbductType             string    `json:"abduct_type,omitempty" sql:"field" xml:"abductType"`                          //被拐儿童类型
// 	IfSampling             int64     `json:"if_sampling,omitempty" sql:"field" xml:"ifSampling"`                          //是否在公安机关采过学样
// 	SamplingDatetime       string    `json:"sampling_datetime,omitempty" sql:"field" xml:"samplingDatetime"`              //采样日期
// 	SamplingRegionalism    string    `json:"sampling_regionalism,omitempty" sql:"field" xml:"samplingRegionalism"`        //采样单位
// 	IfTestdna              int64     `json:"if_testdna,omitempty" sql:"field" xml:"ifTestdna"`                            //是否在公安机关检验DNA
// 	TestDatetime           string    `json:"test_datetime,omitempty" sql:"field" xml:"testDatetime"`                      //检验日期
// 	TestRegionalism        string    `json:"test_regionalism,omitempty" sql:"field" xml:"testRegionalism"`                //检验单位
// 	Unit                   string    `json:"unit,omitempty" sql:"field" xml:"unit"`                                       //单位
// 	IsArrested             string    `json:"is_arrested,omitempty" sql:"field" xml:"isArrested"`                          //是否抓获(DG)
// 	ArrestDaettime         string `json:"arrest_daettime,omitempty" sql:"field" xml:"arrestDaettime"`                  //抓获时间(DG)
// 	IsOffender             string    `json:"is_offender,omitempty" sql:"field" xml:"isOffender"`                          //是否为违法前科(KM)
// 	Va                     string    `json:"va,omitempty" sql:"field" xml:"va"`                                           //通用字段
// 	ExtId                  string    `json:"ext_id,omitempty" sql:"field" xml:"extId"`                                    //外部系统主键
// 	FamilyNo               string    `json:"family_no,omitempty" sql:"field" xml:"familyNo"`                              //家系编号
// 	FamilyName             string    `json:"family_name,omitempty" sql:"field" xml:"familyName"`                          //家系名称
// 	LocalStoreDatetime     string `json:"local_store_datetime,omitempty" sql:"field" xml:"localStoreDatetime"`         //本地存储时间
// 	Id                     string    `json:"id,omitempty" sql:"field" xml:"id"`                                           //主键ID
// 	InitServerNo           string    `json:"init_server_no,omitempty" sql:"field" xml:"initServerNo"`                     //原始服务器编号
// 	LabId                  string    `json:"lab_id,omitempty" sql:"field" xml:"labId"`                                    //实验室ID
// 	ConsignmentId          string    `json:"consignment_id,omitempty" sql:"field" xml:"consignmentId"`                    //委托ID
// 	ConsignOrgCode         string    `json:"consign_org_code,omitempty" sql:"field" xml:"consignOrgCode"`                 //委托单位编号
// 	InputCategory          string    `json:"input_category,omitempty" sql:"field" xml:"inputCategory"`                    //鉴定类别
// 	DbCategory             string    `json:"db_category,omitempty" sql:"field" xml:"dbCategory"`                          //建库类别
// 	SubCategory            string    `json:"sub_category,omitempty" sql:"field" xml:"subCategory"`                        //
// 	PersonNo               string    `json:"person_no,omitempty" sql:"field" xml:"personNo"`                              //人员编号
// 	PersonName             string    `json:"person_name,omitempty" sql:"field" xml:"personName"`                          //人员名称
// 	GenerateMode           string    `json:"generate_mode,omitempty" sql:"field" xml:"generateMode"`                      //对象生成模式（手工、自建）
// 	Alias                  string    `json:"alias,omitempty" sql:"field" xml:"alias"`                                     //别名、绰号
// 	Gender                 string    `json:"gender,omitempty" sql:"field" xml:"gender"`                                   //性别
// 	BirthDatetime          string `json:"birth_datetime,omitempty" sql:"field" xml:"birthDatetime"`                    //出生日期
// 	Age                    int64     `json:"age,omitempty" sql:"field" xml:"age"`                                         //年龄
// 	IdCardNo               string    `json:"id_card_no,omitempty" sql:"field" xml:"idCardNo"`                             //身份证号
// 	CertificateType        string    `json:"certificate_type,omitempty" sql:"field" xml:"certificateType"`                //其他证件类型
// 	CertificateNo          string    `json:"certificate_no,omitempty" sql:"field" xml:"certificateNo"`                    //其他证件号码
// 	Race                   string    `json:"race,omitempty" sql:"field" xml:"race"`                                       //民族
// 	Nationality            string    `json:"nationality,omitempty" sql:"field" xml:"nationality"`                         //国籍
// 	MobilePhone            string    `json:"mobile_phone,omitempty" sql:"field" xml:"mobilePhone"`                        //移动电话
// 	HomePhone              string    `json:"home_phone,omitempty" sql:"field" xml:"homePhone"`                            //家庭电话
// 	Email                  string    `json:"email,omitempty" sql:"field" xml:"email"`                                     //电子邮件
// 	EducationLevel         string    `json:"education_level,omitempty" sql:"field" xml:"educationLevel"`                  //教育程度
// 	Identity               string    `json:"identity,omitempty" sql:"field" xml:"identity"`                               //身份
// 	Occupation             string    `json:"occupation,omitempty" sql:"field" xml:"occupation"`                           //职业
// 	NativePlaceRegionalism string    `json:"native_place_regionalism,omitempty" sql:"field" xml:"nativePlaceRegionalism"` //户籍地行政区划
// 	NativePlaceAddr        string    `json:"native_place_addr,omitempty" sql:"field" xml:"nativePlaceAddr"`               //户籍地名称

// 	TestOr    string `json:"test_or"`    //是否检验
// 	GenoType  string `json:"genotype"`   //是否有基因分型
// 	YstrFlag  string `json:"ystr_flag"`  //是否有Ystr
// 	StrFlag   string `json:"str_flag"`   //是否有str
// 	YstrCount int64  `json:"ystr_count"` //Ystr数量
// 	StrCount  int64  `json:"str_count"`  //str数量
// }
