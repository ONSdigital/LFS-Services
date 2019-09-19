package dataset

import (
	log "github.com/sirupsen/logrus"
	"os"
	conf "pds-go/lfs/config"
	"pds-go/lfs/io/spss"
	"testing"
)

func setupTable(logger *log.Logger) (d *Dataset, err error) {

	_ = os.Remove("LFS.db")

	d, err = NewDataset("address", logger)

	if err != nil {
		logger.Panic("Cannot create database")
	}

	_ = d.AddColumn("Name", spss.STRING)
	_ = d.AddColumn("Address", spss.STRING)
	_ = d.AddColumn("PostCode", spss.INT)
	_ = d.AddColumn("HowMany", spss.FLOAT)

	row1 := map[string]interface{}{
		"Name":     "Boss Lady",
		"Address":  "123 the Valleys Newport Wales",
		"PostCode": 1908,
		"HowMany":  10.24,
	}

	row2 := map[string]interface{}{
		"Name":     "Thorny El",
		"Address":  "Down the pub, as usual",
		"PostCode": 666,
		"HowMany":  11.24,
	}
	row3 := map[string]interface{}{
		"Name":     "George the Dragon",
		"Address":  "With El down the pub",
		"PostCode": 667,
		"HowMany":  12.24,
	}
	_ = d.Insert(row1)
	_ = d.Insert(row2)
	_ = d.Insert(row3)

	return
}

func TestDeleteWhere(t *testing.T) {
	logger := log.New()
	dataset, err := setupTable(logger)
	if err != nil {
		panic(err)
	}

	defer dataset.Close()

	err = dataset.DeleteWhere("PostCode = ? and HowMany = ?", 667, 0)
	rows := dataset.NumRows()
	if rows != 3 {
		t.Errorf("DeleteWhere failed as NumRows is incorrect, got: %d, want: %d.", rows, 3)
	}

	err = dataset.DeleteWhere("PostCode", 667, "HowMany", 12.24)
	rows = dataset.NumRows()
	if rows != 2 {
		t.Errorf("DeleteWhere failed as NumRows is incorrect, got: %d, want: %d.", rows, 2)
	}

}

func TestNumberRowsColumns(t *testing.T) {
	logger := log.New()
	dataset, err := setupTable(logger)
	if err != nil {
		panic(err)
	}
	defer dataset.Close()

	rows := dataset.NumRows()
	cols := dataset.NumColumns()
	if rows != 3 {
		t.Errorf("NumRows was incorrect, got: %d, want: %d.", rows, 3)
	}
	if cols != 5 {
		t.Errorf("NumColumns was incorrect, got: %d, want: %d.", cols, 5)
	}
}

func TestDropByColumn(t *testing.T) {
	logger := log.New()
	dataset, err := setupTable(logger)
	if err != nil {
		panic(err)
	}
	defer dataset.Close()

	err = dataset.DropColumn("Address")
	cols := dataset.NumColumns()
	if cols != 4 {
		t.Errorf("DropByColumn failed as NumColumns is incorrect, got: %d, want: %d.", cols, 4)
	}
}

func TestMean(t *testing.T) {

	logger := log.New()

	dataset, err := setupTable(logger)
	if err != nil {
		panic(err)
	}

	mean, err := dataset.Mean("HowMany")
	if err != nil {
		panic(err)
	}

	if mean != 11.24 {
		t.Errorf("TestMean failed as mean value is incorrect, got: %f, want: %f.", mean, 11.24)
	}
}

func TestFromCSV(t *testing.T) {

	type TestDataset struct {
		Shiftno      float64 `csv:"Shiftno"`
		Serial       float64 `csv:"Serial"`
		Version      string  `csv:"Version"`
		PortRoute2   float64 `csv:"PortRoute2"`
		Baseport     string  `csv:"Baseport"`
		PRouteLatDeg float64 `csv:"PRouteLatDeg"`
		PRouteLonEW  string  `csv:"PRouteLonEW"`
		DVLineName   string  `csv:"DVLineName"`
		DVPortName   string  `csv:"DVPortName"`
	}

	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	dataset, err := d.FromCSV(testDirectory()+"out.csv", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}
	defer dataset.Close()

	logger.Printf("dataset contains %d row(s)\n", dataset.NumRows())
	_ = dataset.Head(5)
}

func TestFromSav(t *testing.T) {

	type TestDataset struct {
		Shiftno      float64 `spss:"Shiftno"`
		Serial       float64 `spss:"Serial"`
		Version      string  `spss:"Version"`
		PortRoute2   float64 `spss:"PortRoute2"`
		Baseport     string  `spss:"Baseport"`
		PRouteLatDeg float64 `spss:"PRouteLatDeg"`
		PRouteLonEW  string  `spss:"PRouteLonEW"`
		DVLineName   string  `spss:"DVLineName"`
		DVPortName   string  `spss:"DVPortName"`
	}

	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	dataset, err := d.FromSav(testDirectory()+"LFSwk18PERS_non_confidential.sav.sav", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}
	defer dataset.Close()

	logger.Printf("dataset contains %d row(s)\n", dataset.NumRows())
	_ = dataset.Head(5)
}

func TestToSav(t *testing.T) {
	type TestDataset struct {
		Shiftno      float64 `spss:"Shiftno"`
		Serial       float64 `spss:"Serial"`
		Version      string  `spss:"Version"`
		PortRoute2   float64 `spss:"PortRoute2"`
		Baseport     string  `spss:"Baseport"`
		PRouteLatDeg float64 `spss:"PRouteLatDeg"`
		PRouteLonEW  string  `spss:"PRouteLonEW"`
		DVLineName   string  `spss:"DVLineName"`
		DVPortName   string  `spss:"DVPortName"`
	}

	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	dataset, err := d.FromSav(testDirectory()+"ips1710bv2.sav", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}
	defer dataset.Close()

	err = dataset.ToSpss(testDirectory() + "dataset-export.sav")
	if err != nil {
		logger.Panic(err)
	}
}

func TestToCSV(t *testing.T) {

	type TestDataset struct {
		Shiftno      float64 `csv:"Shiftno"`
		Serial       float64 `csv:"Serial"`
		Version      string  `csv:"Version"`
		PortRoute2   float64 `csv:"PortRoute2"`
		Baseport     string  `csv:"Baseport"`
		PRouteLatDeg float64 `csv:"PRouteLatDeg"`
		PRouteLonEW  string  `csv:"PRouteLonEW"`
		DVLineName   string  `csv:"DVLineName"`
		DVPortName   string  `csv:"DVPortName"`
	}
	logger := log.New()

	d, err := NewDataset("test", logger)
	if err != nil {
		logger.Panic(err)
	}

	dataset, err := d.FromSav(testDirectory()+"ips1710bv2.sav", TestDataset{})
	if err != nil {
		logger.Panic(err)
	}
	defer dataset.Close()

	err = dataset.ToCSV("out.csv")
	if err != nil {
		logger.Panic(err)
	}

	t.Logf("Dataset Size: %d\n", dataset.NumRows())
	_ = dataset.Head(5)
}

func TestToSQL(t *testing.T) {

	type TestDataset struct {
		Quota         float64 `spss:"Quota"`
		Week          float64 `spss:"Week"`
		W1Yr          float64 `spss:"W1Yr"`
		Qrtr          float64 `spss:"Qrtr"`
		Addr          float64 `spss:"Addr"`
		WavFnd        float64 `spss:"WavFnd"`
		Hhld          float64 `spss:"Hhld"`
		Person        float64 `spss:"Person"`
		ExpStartInt   float64 `spss:"ExpStartInt"`
		ExpEndInt     float64 `spss:"ExpEndInt"`
		SurveyYear    string  `spss:"SurveyYear"`
		Year          float64 `spss:"Year"`
		Month         float64 `spss:"Month"`
		StageAttempt  string  `spss:"StageAttempt"`
		SurvID        string  `spss:"SurvID"`
		SubSample     float64 `spss:"SubSample"`
		Quota2        string  `spss:"Quota2"`
		AddressKey    string  `spss:"AddressKey"`
		MO            float64 `spss:"MO"`
		SampArea      string  `spss:"SampArea"`
		OSGridRef     float64 `spss:"OSGridRef"`
		LEA           string  `spss:"LEA"`
		MajorStrat    float64 `spss:"MajorStrat"`
		GOR           float64 `spss:"GOR"`
		Acorn         string  `spss:"Acorn"`
		GFFMU         float64 `spss:"GFFMU"`
		PrevIssSerNo  string  `spss:"PrevIssSerNo"`
		OldSerial     string  `spss:"OldSerial"`
		PIDNo         string  `spss:"PIDNo"`
		Rand          string  `spss:"Rand"`
		OSCty         string  `spss:"OSCty"`
		OSLAUA        string  `spss:"OSLAUA"`
		OSWard        string  `spss:"OSWard"`
		OSHlthAu      string  `spss:"OSHlthAu"`
		Ctry          string  `spss:"Ctry"`
		PCon          string  `spss:"PCon"`
		TECLEC        string  `spss:"TECLEC"`
		TTWA          string  `spss:"TTWA"`
		PCT           string  `spss:"PCT"`
		NUTS          string  `spss:"NUTS"`
		PSED          string  `spss:"PSED"`
		WardC91       string  `spss:"WardC91"`
		WardO91       string  `spss:"WardO91"`
		Ward98        string  `spss:"Ward98"`
		StatsWard     string  `spss:"StatsWard"`
		OACode        string  `spss:"OACode"`
		CASWard       string  `spss:"CASWard"`
		SOA1          string  `spss:"SOA1"`
		DZone1        string  `spss:"DZone1"`
		SOA2          string  `spss:"SOA2"`
		DZone2        string  `spss:"DZone2"`
		OAC           string  `spss:"OAC"`
		MOCOUNT       string  `spss:"MOCOUNT"`
		Attempt       float64 `spss:"Attempt"`
		LA            string  `spss:"LA"`
		SampQtr       float64 `spss:"SampQtr"`
		HInd          float64 `spss:"HInd"`
		Nurse         float64 `spss:"Nurse"`
		CargoH        string  `spss:"CargoH"`
		Country       float64 `spss:"Country"`
		GB            float64 `spss:"GB"`
		LFSSamp       float64 `spss:"LFSSamp"`
		PCode         string  `spss:"PCode"`
		CID96         string  `spss:"CID96"`
		LAD96         string  `spss:"LAD96"`
		WAD96         string  `spss:"WAD96"`
		TLEC99        string  `spss:"TLEC99"`
		PCA           string  `spss:"PCA"`
		Teleph        float64 `spss:"Teleph"`
		Train         float64 `spss:"Train"`
		ThisWv        float64 `spss:"ThisWv"`
		ThisQtr       float64 `spss:"ThisQtr"`
		LstHO         float64 `spss:"LstHO"`
		LstHO4        float64 `spss:"LstHO4"`
		LstNC         float64 `spss:"LstNC"`
		LstRfDt       float64 `spss:"LstRfDt"`
		LstRfDy       float64 `spss:"LstRfDy"`
		LstRfM        string  `spss:"LstRfM"`
		LstHRPId1     float64 `spss:"LstHRPId1"`
		OldSNo        string  `spss:"OldSNo"`
		OnAddr        float64 `spss:"OnAddr"`
		Workload      float64 `spss:"Workload"`
		Cont          float64 `spss:"Cont"`
		NewHld        float64 `spss:"NewHld"`
		RefDte        float64 `spss:"RefDte"`
		RefDay        float64 `spss:"RefDay"`
		RefMnth       string  `spss:"RefMnth"`
		CalWeek       float64 `spss:"CalWeek"`
		Calqrtr       float64 `spss:"Calqrtr"`
		CalW1YR       float64 `spss:"CalW1YR"`
		IntvNo        float64 `spss:"IntvNo"`
		Lintnum       float64 `spss:"Lintnum"`
		RespHH        float64 `spss:"RespHH"`
		OldNper       float64 `spss:"OldNper"`
		Rels1         string  `spss:"Rels1"`
		HHComp        float64 `spss:"HHComp"`
		HHNew         float64 `spss:"HHNew"`
		Wv1Num        float64 `spss:"Wv1Num"`
		ChkSt         float64 `spss:"ChkSt"`
		NumPer        float64 `spss:"NumPer"`
		PerNo         float64 `spss:"PerNo"`
		RelTxt        string  `spss:"RelTxt"`
		LEstimte      string  `spss:"LEstimte"`
		HBNow         float64 `spss:"HBNow"`
		Sex           float64 `spss:"Sex"`
		Age           float64 `spss:"Age"`
		HallRes       float64 `spss:"HallRes"`
		MarStt        float64 `spss:"MarStt"`
		XMarSta       float64 `spss:"xMarSta"`
		MarSta        float64 `spss:"MarSta"`
		MarChk        float64 `spss:"MarChk"`
		LivTog        float64 `spss:"LivTog"`
		Liv12W        float64 `spss:"Liv12W"`
		HRPId         float64 `spss:"HRPId"`
		Estimate      float64 `spss:"Estimate"`
		MarComp12     float64 `spss:"MarComp12"`
		DVMrDF12      float64 `spss:"DVMrDF12"`
		HldCount      float64 `spss:"HldCount"`
		HH1           float64 `spss:"HH1"`
		RelTxt2       string  `spss:"RelTxt2"`
		R01           float64 `spss:"R01"`
		R02           float64 `spss:"R02"`
		R03           float64 `spss:"R03"`
		R04           float64 `spss:"R04"`
		R05           float64 `spss:"R05"`
		R06           float64 `spss:"R06"`
		R07           float64 `spss:"R07"`
		R08           float64 `spss:"R08"`
		R09           float64 `spss:"R09"`
		R10           float64 `spss:"R10"`
		R11           float64 `spss:"R11"`
		R12           float64 `spss:"R12"`
		R13           float64 `spss:"R13"`
		R14           float64 `spss:"R14"`
		R15           float64 `spss:"R15"`
		R16           float64 `spss:"R16"`
		Parent1       float64 `spss:"Parent1"`
		NoUnits       float64 `spss:"NoUnits"`
		AFAM1         float64 `spss:"AFAM1"`
		ShowFAM       float64 `spss:"ShowFAM"`
		U14kids1      float64 `spss:"U14kids1"`
		U18kids1      float64 `spss:"U18kids1"`
		ParentA1      float64 `spss:"ParentA1"`
		ParentB1      float64 `spss:"ParentB1"`
		Parent2       float64 `spss:"Parent2"`
		NoUnitsSS     float64 `spss:"NoUnitsSS"`
		FAMUNIT1      float64 `spss:"FAMUNIT1"`
		ShowFAM2      float64 `spss:"ShowFAM2"`
		U14kids2      float64 `spss:"U14kids2"`
		U18kids2      float64 `spss:"U18kids2"`
		ParentA2      float64 `spss:"ParentA2"`
		ParentB2      float64 `spss:"ParentB2"`
		SSC1          float64 `spss:"SSC1"`
		SSCType1      float64 `spss:"SSCType1"`
		RDay1         float64 `spss:"RDay1"`
		RInd1         float64 `spss:"RInd1"`
		NDOB1         float64 `spss:"NDOB1"`
		NumAds        float64 `spss:"NumAds"`
		PenFlag       float64 `spss:"PenFlag"`
		RanDob        float64 `spss:"RanDob"`
		DumDob        float64 `spss:"DumDob"`
		LonePar1      float64 `spss:"LonePar1"`
		DVMrDF1       float64 `spss:"DVMrDF1"`
		MarComp1      float64 `spss:"MarComp1"`
		SSC2          float64 `spss:"SSC2"`
		SSCType2      float64 `spss:"SSCType2"`
		Ten1          float64 `spss:"Ten1"`
		Tied          float64 `spss:"Tied"`
		LLord         float64 `spss:"LLord"`
		Furn          float64 `spss:"Furn"`
		EndHHInf      float64 `spss:"EndHHInf"`
		Pind          float64 `spss:"Pind"`
		CargoP        string  `spss:"CargoP"`
		PersNo        float64 `spss:"PersNo"`
		LINow         float64 `spss:"LINow"`
		LIOut         float64 `spss:"LIOut"`
		IntNow        float64 `spss:"IntNow"`
		IOutDate      string  `spss:"IOutDate"`
		FamUnit       float64 `spss:"FamUnit"`
		RespNo        float64 `spss:"RespNo"`
		NTNLTY12      float64 `spss:"NTNLTY12"`
		NatSpec       string  `spss:"NatSpec"`
		Nltyspc       string  `spss:"Nltyspc"`
		NatO7         float64 `spss:"NatO7"`
		Cry12         float64 `spss:"Cry12"`
		CrySpec       string  `spss:"CrySpec"`
		CryO7         float64 `spss:"CryO7"`
		CameYr        float64 `spss:"CameYr"`
		ContUK        float64 `spss:"ContUK"`
		CameYr2       float64 `spss:"CameYr2"`
		CameMT        float64 `spss:"CameMT"`
		WhyUK10       float64 `spss:"WhyUK10"`
		NatldE111     float64 `spss:"NatldE111"`
		NatldE112     float64 `spss:"NatldE112"`
		NatldE113     float64 `spss:"NatldE113"`
		NatldE114     float64 `spss:"NatldE114"`
		NatldE115     float64 `spss:"NatldE115"`
		NatldE116     float64 `spss:"NatldE116"`
		NatldS111     float64 `spss:"NatldS111"`
		NatldS112     float64 `spss:"NatldS112"`
		NatldS113     float64 `spss:"NatldS113"`
		NatldS114     float64 `spss:"NatldS114"`
		NatldS115     float64 `spss:"NatldS115"`
		NatldS116     float64 `spss:"NatldS116"`
		NatldW111     float64 `spss:"NatldW111"`
		NatldW112     float64 `spss:"NatldW112"`
		NatldW113     float64 `spss:"NatldW113"`
		NatldW114     float64 `spss:"NatldW114"`
		NatldW115     float64 `spss:"NatldW115"`
		NatldW116     float64 `spss:"NatldW116"`
		NatldN111     float64 `spss:"NatldN111"`
		NatldN112     float64 `spss:"NatldN112"`
		NatldN113     float64 `spss:"NatldN113"`
		NatldN114     float64 `spss:"NatldN114"`
		NatldN115     float64 `spss:"NatldN115"`
		NatldN116     float64 `spss:"NatldN116"`
		NatldN117     float64 `spss:"NatldN117"`
		NatldO        string  `spss:"NatldO"`
		NatIdCod      float64 `spss:"NatIdCod"`
		CymU          float64 `spss:"CymU"`
		CymS          float64 `spss:"CymS"`
		CymsF         float64 `spss:"CymsF"`
		CymR          float64 `spss:"CymR"`
		CymW          float64 `spss:"CymW"`
		Eth11EW       float64 `spss:"Eth11EW"`
		Eth11S        float64 `spss:"Eth11S"`
		Eth11NI       float64 `spss:"Eth11NI"`
		EthWhe        float64 `spss:"EthWhe"`
		EthWhW        float64 `spss:"EthWhW"`
		EthWSC        float64 `spss:"EthWSC"`
		EthMx11       float64 `spss:"EthMx11"`
		EthAs11       float64 `spss:"EthAs11"`
		EthAs11S      float64 `spss:"EthAs11S"`
		EthBl11       float64 `spss:"EthBl11"`
		EthAFS        float64 `spss:"EthAFS"`
		EthCBS        float64 `spss:"EthCBS"`
		EthOth11      string  `spss:"EthOth11"`
		EthDes        string  `spss:"EthDes"`
		EthOCod       string  `spss:"EthOCod"`
		ResTme        float64 `spss:"ResTme"`
		ResMth        float64 `spss:"ResMth"`
		ResBby        float64 `spss:"ResBby"`
		M3Cry         float64 `spss:"M3Cry"`
		M3CrySpec     string  `spss:"M3CrySpec"`
		M3CryO        float64 `spss:"M3CryO"`
		M3Area        string  `spss:"M3Area"`
		M3Cty         string  `spss:"M3Cty"`
		M3ResC        float64 `spss:"M3ResC"`
		OYEqM3        float64 `spss:"OYEqM3"`
		OYCry         float64 `spss:"OYCry"`
		OYCrySpec     string  `spss:"OYCrySpec"`
		OYCryO        float64 `spss:"OYCryO"`
		OYArea        string  `spss:"OYArea"`
		OYCty         string  `spss:"OYCty"`
		OYResC        float64 `spss:"OYResC"`
		Eth02         string  `spss:"Eth02"`
		Ethc          float64 `spss:"Ethc"`
		SIDFtFQn      float64 `spss:"SIDFtFQn"`
		SIDTUQn       float64 `spss:"SIDTUQn"`
		SIDV          float64 `spss:"SIDV"`
		CardNo        float64 `spss:"CardNo"`
		ReligE        float64 `spss:"ReligE"`
		ReligW        float64 `spss:"ReligW"`
		ReligS        float64 `spss:"ReligS"`
		RelOth        string  `spss:"RelOth"`
		RelOcod       float64 `spss:"RelOcod"`
		Lang          float64 `spss:"Lang"`
		LangD1        float64 `spss:"LangD1"`
		LangD2        float64 `spss:"LangD2"`
		Intuse        float64 `spss:"Intuse"`
		SUBWLLFS      float64 `spss:"SUBWLLFS"`
		Satis         float64 `spss:"Satis"`
		Worth         float64 `spss:"Worth"`
		Happy         float64 `spss:"Happy"`
		Anxious       float64 `spss:"Anxious"`
		SUBPOST       float64 `spss:"SUBPOST"`
		VetServ       float64 `spss:"VetServ"`
		VetCurr       float64 `spss:"VetCurr"`
		VetYrLft1     float64 `spss:"VetYrLft1"`
		VtYrLft2      float64 `spss:"VtYrLft2"`
		VtYrLft3      float64 `spss:"VtYrLft3"`
		VETYEARLFT    float64 `spss:"VETYEARLFT"`
		VLFT2CHK      float64 `spss:"VLFT2CHK"`
		Schm12        float64 `spss:"Schm12"`
		NewDea10      float64 `spss:"NewDea10"`
		FUND12        float64 `spss:"FUND12"`
		YTEtMp        float64 `spss:"YTEtMp"`
		TYPSCH12      float64 `spss:"TYPSCH12"`
		HELPSE12      float64 `spss:"HELPSE12"`
		YTEtJb        float64 `spss:"YTEtJb"`
		LastJb        float64 `spss:"LastJb"`
		Wrking        float64 `spss:"Wrking"`
		JbAway        float64 `spss:"JbAway"`
		OwnBus        float64 `spss:"OwnBus"`
		RelBus        float64 `spss:"RelBus"`
		EverWk        float64 `spss:"EverWk"`
		Caswrk        float64 `spss:"Caswrk"`
		LeftYr        float64 `spss:"LeftYr"`
		LeftM         float64 `spss:"LeftM"`
		LeftW         float64 `spss:"LeftW"`
		LIndD         string  `spss:"LIndD"`
		LIndT         string  `spss:"LIndT"`
		LOccD         string  `spss:"LOccD"`
		LOccT         string  `spss:"LOccT"`
		LStat         float64 `spss:"LStat"`
		LManage       float64 `spss:"LManage"`
		LSupvis       float64 `spss:"LSupvis"`
		LMpnE01       float64 `spss:"LMpnE01"`
		LSolo         float64 `spss:"LSolo"`
		LMpnS01       float64 `spss:"LMpnS01"`
		IState        float64 `spss:"IState"`
		IndD          string  `spss:"IndD"`
		IndT          string  `spss:"IndT"`
		Sector        float64 `spss:"Sector"`
		Sectro03      float64 `spss:"Sectro03"`
		SocEnt        float64 `spss:"SocEnt"`
		SecSoc        float64 `spss:"SecSoc"`
		SecOth        float64 `spss:"SecOth"`
		OccT          string  `spss:"OccT"`
		OccD          string  `spss:"OccD"`
		DVOccT        float64 `spss:"DVOccT"`
		RecJob        float64 `spss:"RecJob"`
		Stat          float64 `spss:"Stat"`
		AGWRK         float64 `spss:"AGWRK"`
		PdWg10        float64 `spss:"PdWg10"`
		Self1         float64 `spss:"Self1"`
		Self2         float64 `spss:"Self2"`
		Self3         float64 `spss:"Self3"`
		Self4         float64 `spss:"Self4"`
		NITax         float64 `spss:"NITax"`
		HwLng         float64 `spss:"HwLng"`
		FifSal        float64 `spss:"FifSal"`
		Supvis        float64 `spss:"Supvis"`
		Manage        float64 `spss:"Manage"`
		MpnE01        float64 `spss:"MpnE01"`
		MpnE02        float64 `spss:"MpnE02"`
		Solo          float64 `spss:"Solo"`
		MpnS01        float64 `spss:"MpnS01"`
		MpnS02        float64 `spss:"MpnS02"`
		OneTen        float64 `spss:"OneTen"`
		OEmpStat      float64 `spss:"OEmpStat"`
		OMCont        float64 `spss:"OMCont"`
		NoCust        float64 `spss:"NoCust"`
		FtPtWk        float64 `spss:"FtPtWk"`
		YPTCIA        float64 `spss:"YPTCIA"`
		YPtJob        float64 `spss:"YPtJob"`
		PTNCre71      float64 `spss:"PTNCre71"`
		PTNCre72      float64 `spss:"PTNCre72"`
		YNotFt        float64 `spss:"YNotFt"`
		JobTyp        float64 `spss:"JobTyp"`
		JbTp101       float64 `spss:"JbTp101"`
		JbTp102       float64 `spss:"JbTp102"`
		JbTp103       float64 `spss:"JbTp103"`
		JbTp104       float64 `spss:"JbTp104"`
		JbTp105       float64 `spss:"JbTp105"`
		WhyTmp6       float64 `spss:"WhyTmp6"`
		TemLen        float64 `spss:"TemLen"`
		ConMpY        float64 `spss:"ConMpY"`
		ConSEY        float64 `spss:"ConSEY"`
		ConMon        float64 `spss:"ConMon"`
		HowGet        float64 `spss:"HowGet"`
		CONPRE        float64 `spss:"CONPRE"`
		CONPRY        float64 `spss:"CONPRY"`
		CONPRM        float64 `spss:"CONPRM"`
		CONPRR        float64 `spss:"CONPRR"`
		TmpCon        float64 `spss:"TmpCon"`
		WRKLNG1       float64 `spss:"WRKLNG1"`
		WRKLNG2       float64 `spss:"WRKLNG2"`
		WRKLNG3       float64 `spss:"WRKLNG3"`
		WRKLNG4       float64 `spss:"WRKLNG4"`
		WRKLNG5       float64 `spss:"WRKLNG5"`
		WRKLNG6       float64 `spss:"WRKLNG6"`
		WRKLNG7       float64 `spss:"WRKLNG7"`
		MAINRET       float64 `spss:"MAINRET"`
		TmpPay        float64 `spss:"TmpPay"`
		RedPaid       float64 `spss:"RedPaid"`
		RedYL13       float64 `spss:"RedYL13"`
		HthDis        float64 `spss:"HthDis"`
		HthRet        float64 `spss:"HthRet"`
		HthRes        float64 `spss:"HthRes"`
		RedAny        float64 `spss:"RedAny"`
		RedYRs        float64 `spss:"RedYRs"`
		HthOth        float64 `spss:"HthOth"`
		RedStat       float64 `spss:"RedStat"`
		RedMpNo       float64 `spss:"RedMpNo"`
		RdMpNo2       float64 `spss:"RdMpNo2"`
		RedMpN        float64 `spss:"RedMpN"`
		RedClos       float64 `spss:"RedClos"`
		RedP1         float64 `spss:"RedP1"`
		RedP2         float64 `spss:"RedP2"`
		RedP3         float64 `spss:"RedP3"`
		RedInd        float64 `spss:"RedInd"`
		RedOcc        float64 `spss:"RedOcc"`
		RdIndD        string  `spss:"RdIndD"`
		RdIndT        string  `spss:"RdIndT"`
		RdOccT        string  `spss:"RdOccT"`
		RdOccD        string  `spss:"RdOccD"`
		Home          float64 `spss:"Home"`
		EvHm98        float64 `spss:"EvHm98"`
		HomeD1        float64 `spss:"HomeD1"`
		HomeD2        float64 `spss:"HomeD2"`
		HomeD3        float64 `spss:"HomeD3"`
		TeleQA        float64 `spss:"TeleQA"`
		TeleQB        float64 `spss:"TeleQB"`
		AtFrom        float64 `spss:"AtFrom"`
		SmeSit        float64 `spss:"SmeSit"`
		WkTown        string  `spss:"WkTown"`
		WkCty         string  `spss:"WkCty"`
		WkPl99        float64 `spss:"WkPl99"`
		WkAbrC        float64 `spss:"WkAbrC"`
		TrvTme        float64 `spss:"TrvTme"`
		TrvMth        float64 `spss:"TrvMth"`
		TrvDrv        float64 `spss:"TrvDrv"`
		HWW4WK        float64 `spss:"HWW4WK"`
		HWWRET        float64 `spss:"HWWRET"`
		ActWkDy1      float64 `spss:"ActWkDy1"`
		ActWkDy2      float64 `spss:"ActWkDy2"`
		ActWkDy3      float64 `spss:"ActWkDy3"`
		ActWkDy4      float64 `spss:"ActWkDy4"`
		ActWkDy5      float64 `spss:"ActWkDy5"`
		ActWkDy6      float64 `spss:"ActWkDy6"`
		ActWkDy7      float64 `spss:"ActWkDy7"`
		ActLow        float64 `spss:"ActLow"`
		ActHgh        float64 `spss:"ActHgh"`
		IllWk         float64 `spss:"IllWk"`
		IllDays1      float64 `spss:"IllDays1"`
		IllDays2      float64 `spss:"IllDays2"`
		IllDays3      float64 `spss:"IllDays3"`
		IllDays4      float64 `spss:"IllDays4"`
		IllDays5      float64 `spss:"IllDays5"`
		IllDays6      float64 `spss:"IllDays6"`
		IllDays7      float64 `spss:"IllDays7"`
		IllLow        float64 `spss:"IllLow"`
		IllHgh        float64 `spss:"IllHgh"`
		ILL1Pd        float64 `spss:"ILL1Pd"`
		IL1Bef        float64 `spss:"IL1Bef"`
		ILLNE11       float64 `spss:"ILLNE11"`
		IL2Bef        float64 `spss:"IL2Bef"`
		ILLFst11      float64 `spss:"ILLFst11"`
		ILLSt         float64 `spss:"ILLSt"`
		ILNxSm        float64 `spss:"ILNxSm"`
		ILLNxt11      float64 `spss:"ILLNxt11"`
		EverOT        float64 `spss:"EverOT"`
		TotUs1        float64 `spss:"TotUs1"`
		UsuHr         float64 `spss:"UsuHr"`
		POtHr         float64 `spss:"POtHr"`
		UOtHr         float64 `spss:"UOtHr"`
		TotUs2        float64 `spss:"TotUs2"`
		TotAc1        float64 `spss:"TotAc1"`
		ActHr         float64 `spss:"ActHr"`
		ActPOt        float64 `spss:"ActPOt"`
		ActUOt        float64 `spss:"ActUOt"`
		TotAc2        float64 `spss:"TotAc2"`
		YLess6        float64 `spss:"YLess6"`
		MatLve        float64 `spss:"MatLve"`
		YMore         float64 `spss:"YMore"`
		VaryHr        float64 `spss:"VaryHr"`
		ShftWk99      float64 `spss:"ShftWk99"`
		ShfTyp        float64 `spss:"ShfTyp"`
		Flex101       float64 `spss:"Flex101"`
		Flex102       float64 `spss:"Flex102"`
		Flex103       float64 `spss:"Flex103"`
		LssOth        float64 `spss:"LssOth"`
		DaysPZ        float64 `spss:"DaysPZ"`
		UsuWrkM1      float64 `spss:"UsuWrkM1"`
		UsuWrkM2      float64 `spss:"UsuWrkM2"`
		UsuWrkM3      float64 `spss:"UsuWrkM3"`
		UsuWrkM4      float64 `spss:"UsuWrkM4"`
		EvEve         float64 `spss:"EvEve"`
		EVENG         float64 `spss:"EVENG"`
		EvNght        float64 `spss:"EvNght"`
		NIGHT         float64 `spss:"NIGHT"`
		UsuWrk1       float64 `spss:"UsuWrk1"`
		UsuWrk2       float64 `spss:"UsuWrk2"`
		UsuWrk3       float64 `spss:"UsuWrk3"`
		EvDay         float64 `spss:"EvDay"`
		WchDay1       float64 `spss:"WchDay1"`
		WchDay2       float64 `spss:"WchDay2"`
		WchDay3       float64 `spss:"WchDay3"`
		WchDay4       float64 `spss:"WchDay4"`
		WchDay5       float64 `spss:"WchDay5"`
		WchDay6       float64 `spss:"WchDay6"`
		WchDay7       float64 `spss:"WchDay7"`
		EvSat         float64 `spss:"EvSat"`
		SATDY         float64 `spss:"SATDY"`
		EvSun         float64 `spss:"EvSun"`
		SUNDY         float64 `spss:"SUNDY"`
		Hols          float64 `spss:"Hols"`
		HolsB         float64 `spss:"HolsB"`
		BHolChk       float64 `spss:"BHolChk"`
		BHolCor       float64 `spss:"BHolCor"`
		BHolPlc       float64 `spss:"BHolPlc"`
		Bank          float64 `spss:"Bank"`
		BnkH1101      float64 `spss:"BnkH1101"`
		BnkH1102      float64 `spss:"BnkH1102"`
		BnkH1103      float64 `spss:"BnkH1103"`
		BnkH1104      float64 `spss:"BnkH1104"`
		BnkH1105      float64 `spss:"BnkH1105"`
		BnkH1106      float64 `spss:"BnkH1106"`
		BnkH1107      float64 `spss:"BnkH1107"`
		BnkH1108      float64 `spss:"BnkH1108"`
		BnkH1109      float64 `spss:"BnkH1109"`
		BnkH1110      float64 `spss:"BnkH1110"`
		BnkH1111      float64 `spss:"BnkH1111"`
		BnkHolF       float64 `spss:"BnkHolF"`
		BHPaid        float64 `spss:"BHPaid"`
		BHNotA        float64 `spss:"BHNotA"`
		BHNotB        float64 `spss:"BHNotB"`
		BHNotC        float64 `spss:"BHNotC"`
		Union         float64 `spss:"Union"`
		TUPres        float64 `spss:"TUPres"`
		TUCov         float64 `spss:"TUCov"`
		LIndD2        string  `spss:"LIndD2"`
		LIndT2        string  `spss:"LIndT2"`
		LOccD2        string  `spss:"LOccD2"`
		LOccT2        string  `spss:"LOccT2"`
		LStat2        float64 `spss:"LStat2"`
		LManag2       float64 `spss:"LManag2"`
		LSupvis2      float64 `spss:"LSupvis2"`
		LMpnES01      float64 `spss:"LMpnES01"`
		LSolo2        float64 `spss:"LSolo2"`
		LMpnSS01      float64 `spss:"LMpnSS01"`
		SecJob        float64 `spss:"SecJob"`
		Y2Job         float64 `spss:"Y2Job"`
		JobTyp2       float64 `spss:"JobTyp2"`
		Jb2T101       float64 `spss:"Jb2T101"`
		Jb2T102       float64 `spss:"Jb2T102"`
		Jb2T103       float64 `spss:"Jb2T103"`
		Jb2T104       float64 `spss:"Jb2T104"`
		Jb2T105       float64 `spss:"Jb2T105"`
		IndD2         string  `spss:"IndD2"`
		IndT2         string  `spss:"IndT2"`
		OccT2         string  `spss:"OccT2"`
		OccD2         string  `spss:"OccD2"`
		RecJb2        float64 `spss:"RecJb2"`
		Stat2         float64 `spss:"Stat2"`
		PdWg102       float64 `spss:"PdWg102"`
		Self21        float64 `spss:"Self21"`
		Self22        float64 `spss:"Self22"`
		Self23        float64 `spss:"Self23"`
		Self24        float64 `spss:"Self24"`
		NITax2        float64 `spss:"NITax2"`
		Supvis2       float64 `spss:"Supvis2"`
		Manag2        float64 `spss:"Manag2"`
		MpnES01       float64 `spss:"MpnES01"`
		MpnES02       float64 `spss:"MpnES02"`
		Solo2         float64 `spss:"Solo2"`
		MpnSS01       float64 `spss:"MpnSS01"`
		MpnSS02       float64 `spss:"MpnSS02"`
		OEmpStat2     float64 `spss:"OEmpStat2"`
		Acthr2        float64 `spss:"Acthr2"`
		Home2         float64 `spss:"Home2"`
		WkTow2        string  `spss:"WkTow2"`
		WkCty2        string  `spss:"WkCty2"`
		WkPl299       float64 `spss:"WkPl299"`
		DifJob        float64 `spss:"DifJob"`
		AddJob        float64 `spss:"AddJob"`
		LookM111      float64 `spss:"LookM111"`
		LookM112      float64 `spss:"LookM112"`
		LookM113      float64 `spss:"LookM113"`
		PrefHr        float64 `spss:"PrefHr"`
		LesPay        float64 `spss:"LesPay"`
		UndEmp        float64 `spss:"UndEmp"`
		UndHrs        float64 `spss:"UndHrs"`
		UndY981       float64 `spss:"UndY981"`
		UndY982       float64 `spss:"UndY982"`
		UndY983       float64 `spss:"UndY983"`
		UndY984       float64 `spss:"UndY984"`
		UndY985       float64 `spss:"UndY985"`
		UndY986       float64 `spss:"UndY986"`
		UndY987       float64 `spss:"UndY987"`
		UndY988       float64 `spss:"UndY988"`
		UndY989       float64 `spss:"UndY989"`
		UndSt         float64 `spss:"UndSt"`
		ExtOth        string  `spss:"ExtOth"`
		LesPay2       float64 `spss:"LesPay2"`
		LesPay3       float64 `spss:"LesPay3"`
		OvHrs         float64 `spss:"OvHrs"`
		Look4         float64 `spss:"Look4"`
		LkYt4         float64 `spss:"LkYt4"`
		Wait          float64 `spss:"Wait"`
		JobBeg        float64 `spss:"JobBeg"`
		LikeWk        float64 `spss:"LikeWk"`
		NoLoWa01      float64 `spss:"NoLoWa01"`
		NoLoWa02      float64 `spss:"NoLoWa02"`
		NoLoWa03      float64 `spss:"NoLoWa03"`
		NoLoWa04      float64 `spss:"NoLoWa04"`
		NoLoWa05      float64 `spss:"NoLoWa05"`
		NoLoWa06      float64 `spss:"NoLoWa06"`
		NoLoWa07      float64 `spss:"NoLoWa07"`
		NoLoWa08      float64 `spss:"NoLoWa08"`
		NoLoWa09      float64 `spss:"NoLoWa09"`
		NoLoWa10      float64 `spss:"NoLoWa10"`
		NoLWM         float64 `spss:"NoLWM"`
		NoLWF         float64 `spss:"NoLWF"`
		NWNCre1       float64 `spss:"NWNCre1"`
		NWNCre2       float64 `spss:"NWNCre2"`
		FutWk         float64 `spss:"FutWk"`
		FWkWen        float64 `spss:"FWkWen"`
		LkSelA        float64 `spss:"LkSelA"`
		LkSelC        float64 `spss:"LkSelC"`
		LkFtPA        float64 `spss:"LkFtPA"`
		LkFtPC        float64 `spss:"LkFtPC"`
		AxPA          float64 `spss:"AxPA"`
		AxPB          float64 `spss:"AxPB"`
		AxFA          float64 `spss:"AxFA"`
		AxFB          float64 `spss:"AxFB"`
		MethMp01      float64 `spss:"MethMp01"`
		MethMp02      float64 `spss:"MethMp02"`
		MethMp03      float64 `spss:"MethMp03"`
		MethMp04      float64 `spss:"MethMp04"`
		MethMp05      float64 `spss:"MethMp05"`
		MethMp06      float64 `spss:"MethMp06"`
		MethMp07      float64 `spss:"MethMp07"`
		MethMp08      float64 `spss:"MethMp08"`
		MethMp09      float64 `spss:"MethMp09"`
		MethMp10      float64 `spss:"MethMp10"`
		MethMp11      float64 `spss:"MethMp11"`
		MethSE1       float64 `spss:"MethSE1"`
		MethSE2       float64 `spss:"MethSE2"`
		MethSE3       float64 `spss:"MethSE3"`
		MethSE4       float64 `spss:"MethSE4"`
		MethSE5       float64 `spss:"MethSE5"`
		MethSE6       float64 `spss:"MethSE6"`
		MethAl01      float64 `spss:"MethAl01"`
		MethAl02      float64 `spss:"MethAl02"`
		MethAl03      float64 `spss:"MethAl03"`
		MethAl04      float64 `spss:"MethAl04"`
		MethAl05      float64 `spss:"MethAl05"`
		MethAl06      float64 `spss:"MethAl06"`
		MethAl07      float64 `spss:"MethAl07"`
		MethAl08      float64 `spss:"MethAl08"`
		MethAl09      float64 `spss:"MethAl09"`
		MethAl10      float64 `spss:"MethAl10"`
		MethAl11      float64 `spss:"MethAl11"`
		MethAl12      float64 `spss:"MethAl12"`
		MethAl13      float64 `spss:"MethAl13"`
		MethAl14      float64 `spss:"MethAl14"`
		MainMe        float64 `spss:"MainMe"`
		MainMA        float64 `spss:"MainMA"`
		MainMs        float64 `spss:"MainMs"`
		MethM         float64 `spss:"MethM"`
		Start         float64 `spss:"Start"`
		YStart        float64 `spss:"YStart"`
		LkTimA        float64 `spss:"LkTimA"`
		LLkTimA       float64 `spss:"LLkTimA"`
		LkTimB        float64 `spss:"LkTimB"`
		Befor         float64 `spss:"Befor"`
		SttBen        float64 `spss:"SttBen"`
		Benfts        float64 `spss:"Benfts"`
		TpBn1301      float64 `spss:"TpBn1301"`
		TpBn1302      float64 `spss:"TpBn1302"`
		TpBn1303      float64 `spss:"TpBn1303"`
		TpBn1304      float64 `spss:"TpBn1304"`
		TpBn1305      float64 `spss:"TpBn1305"`
		TpBn1306      float64 `spss:"TpBn1306"`
		TpBn1307      float64 `spss:"TpBn1307"`
		TpBn1308      float64 `spss:"TpBn1308"`
		TpBn1309      float64 `spss:"TpBn1309"`
		TpBn1310      float64 `spss:"TpBn1310"`
		BenOth        string  `spss:"BenOth"`
		UnEmBn1       float64 `spss:"UnEmBn1"`
		UnEmBn2       float64 `spss:"UnEmBn2"`
		UCredit       float64 `spss:"UCredit"`
		JSATyp        float64 `spss:"JSATyp"`
		JSADur        float64 `spss:"JSADur"`
		IncSup        float64 `spss:"IncSup"`
		DisBen1       float64 `spss:"DisBen1"`
		DisBen2       float64 `spss:"DisBen2"`
		DisBen3       float64 `spss:"DisBen3"`
		DisBen4       float64 `spss:"DisBen4"`
		DisBen5       float64 `spss:"DisBen5"`
		DisBen6       float64 `spss:"DisBen6"`
		DisBen7       float64 `spss:"DisBen7"`
		DisBen8       float64 `spss:"DisBen8"`
		PenBn131      float64 `spss:"PenBn131"`
		PenBn132      float64 `spss:"PenBn132"`
		PenBn133      float64 `spss:"PenBn133"`
		HsngGB1       float64 `spss:"HsngGB1"`
		HsngGB2       float64 `spss:"HsngGB2"`
		MarrCohab     float64 `spss:"MarrCohab"`
		MCNum         float64 `spss:"MCNum"`
		OYCirc        float64 `spss:"OYCirc"`
		OYSInd        float64 `spss:"OYSInd"`
		OYIndD        string  `spss:"OYIndD"`
		OYIndT        string  `spss:"OYIndT"`
		OYSOcc        float64 `spss:"OYSOcc"`
		OYOccT        string  `spss:"OYOccT"`
		OYOccD        string  `spss:"OYOccD"`
		OYStat        float64 `spss:"OYStat"`
		OYSupvi       float64 `spss:"OYSupvi"`
		Oymnge        float64 `spss:"Oymnge"`
		OYSolo        float64 `spss:"OYSolo"`
		OYMPE01       float64 `spss:"OYMPE01"`
		OYMPE02       float64 `spss:"OYMPE02"`
		OYMPS01       float64 `spss:"OYMPS01"`
		OYMPS02       float64 `spss:"OYMPS02"`
		OYOEmpStat    float64 `spss:"OYOEmpStat"`
		OYFtPt        float64 `spss:"OYFtPt"`
		SubjQ         string  `spss:"SubjQ"`
		SubCode       string  `spss:"SubCode"`
		Subqal        string  `spss:"Subqal"`
		Subno         float64 `spss:"Subno"`
		SubjQ2        string  `spss:"SubjQ2"`
		SubCode2      string  `spss:"SubCode2"`
		Subqal2       string  `spss:"Subqal2"`
		Subno2        float64 `spss:"Subno2"`
		SubjQ3        string  `spss:"SubjQ3"`
		SubCode3      string  `spss:"SubCode3"`
		Subqal3       string  `spss:"Subqal3"`
		Subno3        float64 `spss:"Subno3"`
		SubjQ4        string  `spss:"SubjQ4"`
		SubCode4      string  `spss:"SubCode4"`
		Subqal4       string  `spss:"Subqal4"`
		Subno4        float64 `spss:"Subno4"`
		SubjQ5        string  `spss:"SubjQ5"`
		SubCode5      string  `spss:"SubCode5"`
		Subqal5       string  `spss:"Subqal5"`
		Subno5        float64 `spss:"Subno5"`
		SubjQ6        string  `spss:"SubjQ6"`
		SubCode6      string  `spss:"SubCode6"`
		Subqal6       string  `spss:"Subqal6"`
		Subno6        float64 `spss:"Subno6"`
		SubjQ7        string  `spss:"SubjQ7"`
		SubCode7      string  `spss:"SubCode7"`
		Subqal7       string  `spss:"Subqal7"`
		Subno7        float64 `spss:"Subno7"`
		SubjQ8        string  `spss:"SubjQ8"`
		SubCode8      string  `spss:"SubCode8"`
		Subqal8       string  `spss:"Subqal8"`
		Subno8        float64 `spss:"Subno8"`
		FORQUAL       float64 `spss:"FORQUAL"`
		QLFOR111      float64 `spss:"QLFOR111"`
		QLFOR112      float64 `spss:"QLFOR112"`
		QLFOR113      float64 `spss:"QLFOR113"`
		QLFOR114      float64 `spss:"QLFOR114"`
		QLFOR115      float64 `spss:"QLFOR115"`
		QLFOR116      float64 `spss:"QLFOR116"`
		QualUK        float64 `spss:"QualUK"`
		QualRe        float64 `spss:"QualRe"`
		ForTyp11      float64 `spss:"ForTyp11"`
		FORSUB        string  `spss:"FORSUB"`
		QULCH111      float64 `spss:"QULCH111"`
		QULCH112      float64 `spss:"QULCH112"`
		QULCH113      float64 `spss:"QULCH113"`
		QULCH114      float64 `spss:"QULCH114"`
		QULCH115      float64 `spss:"QULCH115"`
		QULCH116      float64 `spss:"QULCH116"`
		SCQUAL01      float64 `spss:"SCQUAL01"`
		SCQUAL02      float64 `spss:"SCQUAL02"`
		SCQUAL03      float64 `spss:"SCQUAL03"`
		SCQUAL04      float64 `spss:"SCQUAL04"`
		SCQUAL05      float64 `spss:"SCQUAL05"`
		SCQUAL06      float64 `spss:"SCQUAL06"`
		SCQUAL07      float64 `spss:"SCQUAL07"`
		SCQUAL08      float64 `spss:"SCQUAL08"`
		SCQUAL09      float64 `spss:"SCQUAL09"`
		SCQUAL10      float64 `spss:"SCQUAL10"`
		SCQUAL11      float64 `spss:"SCQUAL11"`
		SCQUAL12      float64 `spss:"SCQUAL12"`
		SCQUAL13      float64 `spss:"SCQUAL13"`
		SCQUAL14      float64 `spss:"SCQUAL14"`
		SCQUAL15      float64 `spss:"SCQUAL15"`
		SCQUAL16      float64 `spss:"SCQUAL16"`
		SCQUAL17      float64 `spss:"SCQUAL17"`
		SCQUAL18      float64 `spss:"SCQUAL18"`
		SCQUAL19      float64 `spss:"SCQUAL19"`
		SCQUAL20      float64 `spss:"SCQUAL20"`
		SCQUAL21      float64 `spss:"SCQUAL21"`
		SCQUAL22      float64 `spss:"SCQUAL22"`
		SCQUAL23      float64 `spss:"SCQUAL23"`
		SCQUAL24      float64 `spss:"SCQUAL24"`
		SCQUAL25      float64 `spss:"SCQUAL25"`
		UNIQUAL01     float64 `spss:"UNIQUAL01"`
		UNIQUAL02     float64 `spss:"UNIQUAL02"`
		UNIQUAL03     float64 `spss:"UNIQUAL03"`
		UNIQUAL04     float64 `spss:"UNIQUAL04"`
		UNIQUAL05     float64 `spss:"UNIQUAL05"`
		UNIQUAL06     float64 `spss:"UNIQUAL06"`
		UNIQUAL07     float64 `spss:"UNIQUAL07"`
		UNIQUAL08     float64 `spss:"UNIQUAL08"`
		UNIQUAL09     float64 `spss:"UNIQUAL09"`
		UNIQUAL10     float64 `spss:"UNIQUAL10"`
		UNIQUAL11     float64 `spss:"UNIQUAL11"`
		UNIQUAL12     float64 `spss:"UNIQUAL12"`
		UNIQUAL13     float64 `spss:"UNIQUAL13"`
		UNIQUAL14     float64 `spss:"UNIQUAL14"`
		UNIQUAL15     float64 `spss:"UNIQUAL15"`
		UNIQUAL16     float64 `spss:"UNIQUAL16"`
		UNIQUAL17     float64 `spss:"UNIQUAL17"`
		UNIQUAL18     float64 `spss:"UNIQUAL18"`
		UNIQUAL19     float64 `spss:"UNIQUAL19"`
		UNIQUAL20     float64 `spss:"UNIQUAL20"`
		UNIQUAL21     float64 `spss:"UNIQUAL21"`
		UNIQUAL22     float64 `spss:"UNIQUAL22"`
		UNIQUAL23     float64 `spss:"UNIQUAL23"`
		UNIQUAL24     float64 `spss:"UNIQUAL24"`
		UNIQUAL25     float64 `spss:"UNIQUAL25"`
		UNIQUAL26     float64 `spss:"UNIQUAL26"`
		UNIQUAL27     float64 `spss:"UNIQUAL27"`
		UNIQUAL28     float64 `spss:"UNIQUAL28"`
		UNIQUAL29     float64 `spss:"UNIQUAL29"`
		UNIQUAL30     float64 `spss:"UNIQUAL30"`
		UNIQUAL31     float64 `spss:"UNIQUAL31"`
		WOQUAL01      float64 `spss:"WOQUAL01"`
		WOQUAL02      float64 `spss:"WOQUAL02"`
		WOQUAL03      float64 `spss:"WOQUAL03"`
		WOQUAL04      float64 `spss:"WOQUAL04"`
		WOQUAL05      float64 `spss:"WOQUAL05"`
		WOQUAL06      float64 `spss:"WOQUAL06"`
		WOQUAL07      float64 `spss:"WOQUAL07"`
		WOQUAL08      float64 `spss:"WOQUAL08"`
		WOQUAL09      float64 `spss:"WOQUAL09"`
		WOQUAL10      float64 `spss:"WOQUAL10"`
		WOQUAL11      float64 `spss:"WOQUAL11"`
		WOQUAL12      float64 `spss:"WOQUAL12"`
		WOQUAL13      float64 `spss:"WOQUAL13"`
		WOQUAL14      float64 `spss:"WOQUAL14"`
		WOQUAL15      float64 `spss:"WOQUAL15"`
		WOQUAL16      float64 `spss:"WOQUAL16"`
		WOQUAL17      float64 `spss:"WOQUAL17"`
		WOQUAL18      float64 `spss:"WOQUAL18"`
		WOQUAL19      float64 `spss:"WOQUAL19"`
		WOQUAL20      float64 `spss:"WOQUAL20"`
		WOQUAL21      float64 `spss:"WOQUAL21"`
		WOQUAL22      float64 `spss:"WOQUAL22"`
		WOQUAL23      float64 `spss:"WOQUAL23"`
		WOQUAL24      float64 `spss:"WOQUAL24"`
		WOQUAL25      float64 `spss:"WOQUAL25"`
		WOQUAL26      float64 `spss:"WOQUAL26"`
		WOQUAL27      float64 `spss:"WOQUAL27"`
		WOQUAL28      float64 `spss:"WOQUAL28"`
		WOQUAL29      float64 `spss:"WOQUAL29"`
		WOQUAL30      float64 `spss:"WOQUAL30"`
		WOQUAL31      float64 `spss:"WOQUAL31"`
		GSQUAL01      float64 `spss:"GSQUAL01"`
		GSQUAL02      float64 `spss:"GSQUAL02"`
		GSQUAL03      float64 `spss:"GSQUAL03"`
		GSQUAL04      float64 `spss:"GSQUAL04"`
		GSQUAL05      float64 `spss:"GSQUAL05"`
		GSQUAL06      float64 `spss:"GSQUAL06"`
		GSQUAL07      float64 `spss:"GSQUAL07"`
		GSQUAL08      float64 `spss:"GSQUAL08"`
		GSQUAL09      float64 `spss:"GSQUAL09"`
		GSQUAL10      float64 `spss:"GSQUAL10"`
		GSQUAL11      float64 `spss:"GSQUAL11"`
		GSQUAL12      float64 `spss:"GSQUAL12"`
		GSQUAL13      float64 `spss:"GSQUAL13"`
		GSQUAL14      float64 `spss:"GSQUAL14"`
		GSQUAL15      float64 `spss:"GSQUAL15"`
		GSQUAL16      float64 `spss:"GSQUAL16"`
		GSQUAL17      float64 `spss:"GSQUAL17"`
		GSQUAL18      float64 `spss:"GSQUAL18"`
		GSQUAL19      float64 `spss:"GSQUAL19"`
		GSQUAL20      float64 `spss:"GSQUAL20"`
		GSQUAL21      float64 `spss:"GSQUAL21"`
		GSQUAL22      float64 `spss:"GSQUAL22"`
		GSQUAL23      float64 `spss:"GSQUAL23"`
		GSQUAL24      float64 `spss:"GSQUAL24"`
		GSQUAL25      float64 `spss:"GSQUAL25"`
		GSQUAL26      float64 `spss:"GSQUAL26"`
		GSQUAL27      float64 `spss:"GSQUAL27"`
		GSQUAL28      float64 `spss:"GSQUAL28"`
		GSQUAL29      float64 `spss:"GSQUAL29"`
		GSQUAL30      float64 `spss:"GSQUAL30"`
		GSQUAL31      float64 `spss:"GSQUAL31"`
		OTQUAL01      float64 `spss:"OTQUAL01"`
		OTQUAL02      float64 `spss:"OTQUAL02"`
		OTQUAL03      float64 `spss:"OTQUAL03"`
		OTQUAL04      float64 `spss:"OTQUAL04"`
		OTQUAL05      float64 `spss:"OTQUAL05"`
		OTQUAL06      float64 `spss:"OTQUAL06"`
		OTQUAL07      float64 `spss:"OTQUAL07"`
		OTQUAL08      float64 `spss:"OTQUAL08"`
		OTQUAL09      float64 `spss:"OTQUAL09"`
		OTQUAL10      float64 `spss:"OTQUAL10"`
		OTQUAL11      float64 `spss:"OTQUAL11"`
		OTQUAL12      float64 `spss:"OTQUAL12"`
		OTQUAL13      float64 `spss:"OTQUAL13"`
		OTQUAL14      float64 `spss:"OTQUAL14"`
		OTQUAL15      float64 `spss:"OTQUAL15"`
		OTQUAL16      float64 `spss:"OTQUAL16"`
		OTQUAL17      float64 `spss:"OTQUAL17"`
		OTQUAL18      float64 `spss:"OTQUAL18"`
		OTQUAL19      float64 `spss:"OTQUAL19"`
		OTQUAL20      float64 `spss:"OTQUAL20"`
		OTQUAL21      float64 `spss:"OTQUAL21"`
		OTQUAL22      float64 `spss:"OTQUAL22"`
		OTQUAL23      float64 `spss:"OTQUAL23"`
		OTQUAL24      float64 `spss:"OTQUAL24"`
		OTQUAL25      float64 `spss:"OTQUAL25"`
		OTQUAL26      float64 `spss:"OTQUAL26"`
		OTQUAL27      float64 `spss:"OTQUAL27"`
		OTQUAL28      float64 `spss:"OTQUAL28"`
		OTQUAL29      float64 `spss:"OTQUAL29"`
		OTQUAL30      float64 `spss:"OTQUAL30"`
		OTQUAL31      float64 `spss:"OTQUAL31"`
		QLYr1101      float64 `spss:"QLYr1101"`
		QLYr1102      float64 `spss:"QLYr1102"`
		QLYr1103      float64 `spss:"QLYr1103"`
		QLYr1104      float64 `spss:"QLYr1104"`
		QLYr1105      float64 `spss:"QLYr1105"`
		QLYr1106      float64 `spss:"QLYr1106"`
		QLYr1107      float64 `spss:"QLYr1107"`
		QLYr1108      float64 `spss:"QLYr1108"`
		QLYr1109      float64 `spss:"QLYr1109"`
		QLYr1110      float64 `spss:"QLYr1110"`
		QLYr1111      float64 `spss:"QLYr1111"`
		IntroLev      float64 `spss:"IntroLev"`
		Degree71      float64 `spss:"Degree71"`
		Degree72      float64 `spss:"Degree72"`
		Degree73      float64 `spss:"Degree73"`
		Degree74      float64 `spss:"Degree74"`
		Degree75      float64 `spss:"Degree75"`
		HighO         float64 `spss:"HighO"`
		SubjctN       string  `spss:"SubjctN"`
		SinComN       float64 `spss:"SinComN"`
		SngDegN       string  `spss:"SngDegN"`
		CmbDegN01     float64 `spss:"CmbDegN01"`
		CmbDegN02     float64 `spss:"CmbDegN02"`
		CmbDegN03     float64 `spss:"CmbDegN03"`
		CmbDegN04     float64 `spss:"CmbDegN04"`
		CmbDegN05     float64 `spss:"CmbDegN05"`
		CmbDegN06     float64 `spss:"CmbDegN06"`
		CmbDegN07     float64 `spss:"CmbDegN07"`
		CmbDegN08     float64 `spss:"CmbDegN08"`
		CmbDegN09     float64 `spss:"CmbDegN09"`
		CmbDegN10     float64 `spss:"CmbDegN10"`
		CmbDegN11     float64 `spss:"CmbDegN11"`
		CmbDegN12     float64 `spss:"CmbDegN12"`
		CmbMainN      float64 `spss:"CmbMainN"`
		FDSUBJ        string  `spss:"FDSUBJ"`
		FDSinCom      float64 `spss:"FDSinCom"`
		FDSngDeg      string  `spss:"FDSngDeg"`
		FDCMBD01      float64 `spss:"FDCMBD01"`
		FDCMBD02      float64 `spss:"FDCMBD02"`
		FDCMBD03      float64 `spss:"FDCMBD03"`
		FDCMBD04      float64 `spss:"FDCMBD04"`
		FDCMBD05      float64 `spss:"FDCMBD05"`
		FDCMBD06      float64 `spss:"FDCMBD06"`
		FDCMBD07      float64 `spss:"FDCMBD07"`
		FDCMBD08      float64 `spss:"FDCMBD08"`
		FDCMBD09      float64 `spss:"FDCMBD09"`
		FDCMBD10      float64 `spss:"FDCMBD10"`
		FDCMBD11      float64 `spss:"FDCMBD11"`
		FDCMBD12      float64 `spss:"FDCMBD12"`
		FDCmbMa       float64 `spss:"FDCmbMa"`
		FDINST        string  `spss:"FDINST"`
		UGINST        string  `spss:"UGINST"`
		DegCls7       float64 `spss:"DegCls7"`
		HDSUBJCT      string  `spss:"HDSUBJCT"`
		HDSINCOM      float64 `spss:"HDSINCOM"`
		SNGHD         string  `spss:"SNGHD"`
		CMBHD01       float64 `spss:"CMBHD01"`
		CMBHD02       float64 `spss:"CMBHD02"`
		CMBHD03       float64 `spss:"CMBHD03"`
		CMBHD04       float64 `spss:"CMBHD04"`
		CMBHD05       float64 `spss:"CMBHD05"`
		CMBHD06       float64 `spss:"CMBHD06"`
		CMBHD07       float64 `spss:"CMBHD07"`
		CMBHD08       float64 `spss:"CMBHD08"`
		CMBHD09       float64 `spss:"CMBHD09"`
		CMBHD10       float64 `spss:"CMBHD10"`
		CMBHD11       float64 `spss:"CMBHD11"`
		CMBHD12       float64 `spss:"CMBHD12"`
		CMBHDMA       float64 `spss:"CMBHDMA"`
		HDINST        string  `spss:"HDINST"`
		PGINST        string  `spss:"PGINST"`
		CRYDEG        float64 `spss:"CRYDEG"`
		Teach41       float64 `spss:"Teach41"`
		Teach42       float64 `spss:"Teach42"`
		Teach43       float64 `spss:"Teach43"`
		Teach44       float64 `spss:"Teach44"`
		Teach45       float64 `spss:"Teach45"`
		Teach46       float64 `spss:"Teach46"`
		NumAL         float64 `spss:"NumAL"`
		NumAS         float64 `spss:"NumAS"`
		TypHST1       float64 `spss:"TypHST1"`
		TypHST2       float64 `spss:"TypHST2"`
		TypHST3       float64 `spss:"TypHST3"`
		TypHST4       float64 `spss:"TypHST4"`
		TypHST5       float64 `spss:"TypHST5"`
		AdvHST        float64 `spss:"AdvHST"`
		HST           float64 `spss:"HST"`
		WlshBc8       float64 `spss:"WlshBc8"`
		QGCSE41       float64 `spss:"QGCSE41"`
		QGCSE42       float64 `spss:"QGCSE42"`
		QGCSE43       float64 `spss:"QGCSE43"`
		QGCSE44       float64 `spss:"QGCSE44"`
		QGCSE45       float64 `spss:"QGCSE45"`
		GCSE41        float64 `spss:"GCSE41"`
		GCSE42        float64 `spss:"GCSE42"`
		GCSE43        float64 `spss:"GCSE43"`
		GCSE44        float64 `spss:"GCSE44"`
		GCSE45        float64 `spss:"GCSE45"`
		NumOl5        float64 `spss:"NumOl5"`
		MeGCSE        float64 `spss:"MeGCSE"`
		NumOl5O       float64 `spss:"NumOl5O"`
		NumOl5F       float64 `spss:"NumOl5F"`
		QDipTyp       float64 `spss:"QDipTyp"`
		VOCYRA        float64 `spss:"VOCYRA"`
		BTE11         float64 `spss:"BTE11"`
		BTACD         float64 `spss:"BTACD"`
		BTLEV         float64 `spss:"BTLEV"`
		BTSUBJ        string  `spss:"BTSUBJ"`
		BTCTH111      float64 `spss:"BTCTH111"`
		BTCTH112      float64 `spss:"BTCTH112"`
		BTCTH113      float64 `spss:"BTCTH113"`
		BTCTH114      float64 `spss:"BTCTH114"`
		BTCOTLA1      float64 `spss:"BTCOTLA1"`
		BTCOTLA2      float64 `spss:"BTCOTLA2"`
		BTCOTLA3      float64 `spss:"BTCOTLA3"`
		BTCOTLB1      float64 `spss:"BTCOTLB1"`
		BTCOTLB2      float64 `spss:"BTCOTLB2"`
		BTCOTLB3      float64 `spss:"BTCOTLB3"`
		BTCOTLB4      float64 `spss:"BTCOTLB4"`
		BTCOTLB5      float64 `spss:"BTCOTLB5"`
		BTCOTLB6      float64 `spss:"BTCOTLB6"`
		BTCOTLB7      float64 `spss:"BTCOTLB7"`
		BTCOTLB8      float64 `spss:"BTCOTLB8"`
		BTCOTLB9      float64 `spss:"BTCOTLB9"`
		VOCYRB        float64 `spss:"VOCYRB"`
		SCTVC11       float64 `spss:"SCTVC11"`
		SCACD         float64 `spss:"SCACD"`
		SCLEV         float64 `spss:"SCLEV"`
		SCSUBJ        string  `spss:"SCSUBJ"`
		STCOT111      float64 `spss:"STCOT111"`
		STCOT112      float64 `spss:"STCOT112"`
		STCOT113      float64 `spss:"STCOT113"`
		STCOT114      float64 `spss:"STCOT114"`
		STCOT115      float64 `spss:"STCOT115"`
		STCOTLA1      float64 `spss:"STCOTLA1"`
		STCOTLA2      float64 `spss:"STCOTLA2"`
		STCOTLA3      float64 `spss:"STCOTLA3"`
		STCOTLB1      float64 `spss:"STCOTLB1"`
		STCOTLB2      float64 `spss:"STCOTLB2"`
		STCOTLB3      float64 `spss:"STCOTLB3"`
		STCOTLB4      float64 `spss:"STCOTLB4"`
		STCOTLB5      float64 `spss:"STCOTLB5"`
		STCOTLB6      float64 `spss:"STCOTLB6"`
		STCOTLB7      float64 `spss:"STCOTLB7"`
		STCOTLB8      float64 `spss:"STCOTLB8"`
		STCOTLB9      float64 `spss:"STCOTLB9"`
		VOCYRC        float64 `spss:"VOCYRC"`
		RSA11         float64 `spss:"RSA11"`
		RSACD         float64 `spss:"RSACD"`
		RSLEV         float64 `spss:"RSLEV"`
		RSASUBJ       string  `spss:"RSASUBJ"`
		RSAOT111      float64 `spss:"RSAOT111"`
		RSAOT112      float64 `spss:"RSAOT112"`
		RSAOT113      float64 `spss:"RSAOT113"`
		RSAOT114      float64 `spss:"RSAOT114"`
		RSAOTLA1      float64 `spss:"RSAOTLA1"`
		RSAOTLA2      float64 `spss:"RSAOTLA2"`
		RSAOTLA3      float64 `spss:"RSAOTLA3"`
		RSAOTLB1      float64 `spss:"RSAOTLB1"`
		RSAOTLB2      float64 `spss:"RSAOTLB2"`
		RSAOTLB3      float64 `spss:"RSAOTLB3"`
		RSAOTLB4      float64 `spss:"RSAOTLB4"`
		RSAOTLB5      float64 `spss:"RSAOTLB5"`
		RSAOTLB6      float64 `spss:"RSAOTLB6"`
		RSAOTLB7      float64 `spss:"RSAOTLB7"`
		RSAOTLB8      float64 `spss:"RSAOTLB8"`
		RSAOTLB9      float64 `spss:"RSAOTLB9"`
		VOCYRD        float64 `spss:"VOCYRD"`
		CaG11         float64 `spss:"CaG11"`
		CAGACD        float64 `spss:"CAGACD"`
		CAGLEV        float64 `spss:"CAGLEV"`
		CGSUBJ        string  `spss:"CGSUBJ"`
		CaGOT111      float64 `spss:"CaGOT111"`
		CaGOT112      float64 `spss:"CaGOT112"`
		CaGOT113      float64 `spss:"CaGOT113"`
		CaGOTLA1      float64 `spss:"CaGOTLA1"`
		CaGOTLA2      float64 `spss:"CaGOTLA2"`
		CaGOTLA3      float64 `spss:"CaGOTLA3"`
		CAGOTLB1      float64 `spss:"CAGOTLB1"`
		CAGOTLB2      float64 `spss:"CAGOTLB2"`
		CAGOTLB3      float64 `spss:"CAGOTLB3"`
		CAGOTLB4      float64 `spss:"CAGOTLB4"`
		CAGOTLB5      float64 `spss:"CAGOTLB5"`
		CAGOTLB6      float64 `spss:"CAGOTLB6"`
		CAGOTLB7      float64 `spss:"CAGOTLB7"`
		CAGOTLB8      float64 `spss:"CAGOTLB8"`
		CAGOTLB9      float64 `spss:"CAGOTLB9"`
		QGNVQ         float64 `spss:"QGNVQ"`
		VOCYRE        float64 `spss:"VOCYRE"`
		GNVQ11        float64 `spss:"GNVQ11"`
		GNACD         float64 `spss:"GNACD"`
		GNLEV         float64 `spss:"GNLEV"`
		GNVQSUBJ      string  `spss:"GNVQSUBJ"`
		GNVQO111      float64 `spss:"GNVQO111"`
		GNVQO112      float64 `spss:"GNVQO112"`
		GNVQO113      float64 `spss:"GNVQO113"`
		GNVQO114      float64 `spss:"GNVQO114"`
		GNVQO115      float64 `spss:"GNVQO115"`
		GNVQOLA1      float64 `spss:"GNVQOLA1"`
		GNVQOLA2      float64 `spss:"GNVQOLA2"`
		GNVQOLA3      float64 `spss:"GNVQOLA3"`
		GNVQOLB1      float64 `spss:"GNVQOLB1"`
		GNVQOLB2      float64 `spss:"GNVQOLB2"`
		GNVQOLB3      float64 `spss:"GNVQOLB3"`
		GNVQOLB4      float64 `spss:"GNVQOLB4"`
		GNVQOLB5      float64 `spss:"GNVQOLB5"`
		GNVQOLB6      float64 `spss:"GNVQOLB6"`
		GNVQOLB7      float64 `spss:"GNVQOLB7"`
		GNVQOLB8      float64 `spss:"GNVQOLB8"`
		GNVQOLB9      float64 `spss:"GNVQOLB9"`
		NVQSVQ        float64 `spss:"NVQSVQ"`
		VOCYRF        float64 `spss:"VOCYRF"`
		NVQ11         float64 `spss:"NVQ11"`
		NVACD         float64 `spss:"NVACD"`
		NVLEV         float64 `spss:"NVLEV"`
		NVQSUBJ       string  `spss:"NVQSUBJ"`
		NVQO111       float64 `spss:"NVQO111"`
		NVQO112       float64 `spss:"NVQO112"`
		NVQO113       float64 `spss:"NVQO113"`
		NVQO114       float64 `spss:"NVQO114"`
		NVQO115       float64 `spss:"NVQO115"`
		NVQO116       float64 `spss:"NVQO116"`
		NVOTLEA1      float64 `spss:"NVOTLEA1"`
		NVOTLEA2      float64 `spss:"NVOTLEA2"`
		NVOTLEA3      float64 `spss:"NVOTLEA3"`
		NVOTLEB1      float64 `spss:"NVOTLEB1"`
		NVOTLEB2      float64 `spss:"NVOTLEB2"`
		NVOTLEB3      float64 `spss:"NVOTLEB3"`
		NVOTLEB4      float64 `spss:"NVOTLEB4"`
		NVOTLEB5      float64 `spss:"NVOTLEB5"`
		NVOTLEB6      float64 `spss:"NVOTLEB6"`
		NVOTLEB7      float64 `spss:"NVOTLEB7"`
		NVOTLEB8      float64 `spss:"NVOTLEB8"`
		NVOTLEB9      float64 `spss:"NVOTLEB9"`
		NVQun         float64 `spss:"NVQun"`
		VOCYRG        float64 `spss:"VOCYRG"`
		QCFACD        float64 `spss:"QCFACD"`
		QCFLEV        float64 `spss:"QCFLEV"`
		QCFSUBJ       string  `spss:"QCFSUBJ"`
		QCFOTHA1      float64 `spss:"QCFOTHA1"`
		QCFOTHA2      float64 `spss:"QCFOTHA2"`
		QCFOTHA3      float64 `spss:"QCFOTHA3"`
		QCFOTHB1      float64 `spss:"QCFOTHB1"`
		QCFOTHB2      float64 `spss:"QCFOTHB2"`
		QCFOTHB3      float64 `spss:"QCFOTHB3"`
		QCFOTHB4      float64 `spss:"QCFOTHB4"`
		QCFOTHB5      float64 `spss:"QCFOTHB5"`
		QCFOTHB6      float64 `spss:"QCFOTHB6"`
		QCFOTHB7      float64 `spss:"QCFOTHB7"`
		QCFOTHB8      float64 `spss:"QCFOTHB8"`
		QCFOTHB9      float64 `spss:"QCFOTHB9"`
		TpQl111       float64 `spss:"TpQl111"`
		TpQl112       float64 `spss:"TpQl112"`
		TpQl113       float64 `spss:"TpQl113"`
		OthQu91       float64 `spss:"OthQu91"`
		OthQu92       float64 `spss:"OthQu92"`
		OthQu93       float64 `spss:"OthQu93"`
		OthQu94       float64 `spss:"OthQu94"`
		OthQu95       float64 `spss:"OthQu95"`
		VOCYRH        float64 `spss:"VOCYRH"`
		OTHQAL11      string  `spss:"OTHQAL11"`
		OTHQLEV       float64 `spss:"OTHQLEV"`
		QalPl11       float64 `spss:"QalPl11"`
		VocQPl11      float64 `spss:"VocQPl11"`
		YerQal1       float64 `spss:"YerQal1"`
		YerQal2       float64 `spss:"YerQal2"`
		YerQal3       float64 `spss:"YerQal3"`
		WchQGot1      float64 `spss:"WchQGot1"`
		WchQGot2      float64 `spss:"WchQGot2"`
		WchQGot3      float64 `spss:"WchQGot3"`
		WchQGot4      float64 `spss:"WchQGot4"`
		WchQGot5      float64 `spss:"WchQGot5"`
		WchQGot6      float64 `spss:"WchQGot6"`
		WchQGot7      float64 `spss:"WchQGot7"`
		WchQGot8      float64 `spss:"WchQGot8"`
		DVQual01      float64 `spss:"DVQual01"`
		DVQual02      float64 `spss:"DVQual02"`
		DVQual03      float64 `spss:"DVQual03"`
		DVQual04      float64 `spss:"DVQual04"`
		DVQual05      float64 `spss:"DVQual05"`
		DVQual06      float64 `spss:"DVQual06"`
		DVQual07      float64 `spss:"DVQual07"`
		DVQual08      float64 `spss:"DVQual08"`
		DVQual09      float64 `spss:"DVQual09"`
		DVQual10      float64 `spss:"DVQual10"`
		DVQual11      float64 `spss:"DVQual11"`
		DVQual12      float64 `spss:"DVQual12"`
		QHighY        string  `spss:"QHighY"`
		LAppD         string  `spss:"LAppD"`
		LAppT         string  `spss:"LAppT"`
		EdAge         float64 `spss:"EdAge"`
		QulNow        float64 `spss:"QulNow"`
		QulHi11       float64 `spss:"QulHi11"`
		DegNow        float64 `spss:"DegNow"`
		HghNow        float64 `spss:"HghNow"`
		TcNw11        float64 `spss:"TcNw11"`
		TCNWACD       float64 `spss:"TCNWACD"`
		TCNWLEV       float64 `spss:"TCNWLEV"`
		SCNow11       float64 `spss:"SCNow11"`
		SCNWACD       float64 `spss:"SCNWACD"`
		SCNWLEV       float64 `spss:"SCNWLEV"`
		DipTyp        float64 `spss:"DipTyp"`
		OCRN11        float64 `spss:"OCRN11"`
		OCRNACD       float64 `spss:"OCRNACD"`
		OCRNLEV       float64 `spss:"OCRNLEV"`
		CGNw11        float64 `spss:"CGNw11"`
		CGNWACD       float64 `spss:"CGNWACD"`
		CGNWLEV       float64 `spss:"CGNWLEV"`
		HSTNow        float64 `spss:"HSTNow"`
		WBac          float64 `spss:"WBac"`
		NVQKn2        float64 `spss:"NVQKn2"`
		NVQLe11       float64 `spss:"NVQLe11"`
		NVNWACD       float64 `spss:"NVNWACD"`
		NVNWLEV       float64 `spss:"NVNWLEV"`
		QCFNOW        float64 `spss:"QCFNOW"`
		QCFLVNW       float64 `spss:"QCFLVNW"`
		CurSub        string  `spss:"CurSub"`
		CurCode       string  `spss:"CurCode"`
		CurQal        string  `spss:"CurQal"`
		Enroll        float64 `spss:"Enroll"`
		Attend        float64 `spss:"Attend"`
		Course        float64 `spss:"Course"`
		EdIns11       float64 `spss:"EdIns11"`
		APPR12        float64 `spss:"APPR12"`
		APPRCURR      float64 `spss:"APPRCURR"`
		AppSam        float64 `spss:"AppSam"`
		AppD          string  `spss:"AppD"`
		AppT          string  `spss:"AppT"`
		AppInD        string  `spss:"AppInD"`
		AppInT        string  `spss:"AppInT"`
		APPST12       float64 `spss:"APPST12"`
		APPRLEV       float64 `spss:"APPRLEV"`
		Ref3Mths      float64 `spss:"Ref3Mths"`
		Ref3Day       float64 `spss:"Ref3Day"`
		Ref3Mnth      string  `spss:"Ref3Mnth"`
		Ed13Wk        float64 `spss:"Ed13Wk"`
		Ed4Wk         float64 `spss:"Ed4Wk"`
		Futur13       float64 `spss:"Futur13"`
		Futur4        float64 `spss:"Futur4"`
		JobEd         float64 `spss:"JobEd"`
		TrnOpp11      float64 `spss:"TrnOpp11"`
		JobTrn        float64 `spss:"JobTrn"`
		TSte10        float64 `spss:"TSte10"`
		TrNI10        float64 `spss:"TrNI10"`
		TFee101       float64 `spss:"TFee101"`
		TFee102       float64 `spss:"TFee102"`
		TFee103       float64 `spss:"TFee103"`
		TFee104       float64 `spss:"TFee104"`
		TFee105       float64 `spss:"TFee105"`
		FeeIr1        float64 `spss:"FeeIr1"`
		FeeIr2        float64 `spss:"FeeIr2"`
		FeeIr3        float64 `spss:"FeeIr3"`
		FeeIr4        float64 `spss:"FeeIr4"`
		FeeIr5        float64 `spss:"FeeIr5"`
		TrnLen        float64 `spss:"TrnLen"`
		TrHr11        float64 `spss:"TrHr11"`
		TrOnJB        float64 `spss:"TrOnJB"`
		NFE13WK       float64 `spss:"NFE13WK"`
		NFE4WK        float64 `spss:"NFE4WK"`
		TAUT4WK       float64 `spss:"TAUT4WK"`
		T4Purp        float64 `spss:"T4Purp"`
		T4Work        float64 `spss:"T4Work"`
		TAUT3M        float64 `spss:"TAUT3M"`
		TSUBJ4WK      string  `spss:"TSUBJ4WK"`
		TSUB4Cod      string  `spss:"TSUB4Cod"`
		TSUBJ3M       string  `spss:"TSUBJ3M"`
		TSUB3Cod      string  `spss:"TSUB3Cod"`
		TautHrs       float64 `spss:"TautHrs"`
		TLRN4WK       float64 `spss:"TLRN4WK"`
		TLRN3M        float64 `spss:"TLRN3M"`
		Neets         float64 `spss:"Neets"`
		HPrMb         float64 `spss:"HPrMb"`
		HPrMb2        float64 `spss:"HPrMb2"`
		AccCalc       float64 `spss:"AccCalc"`
		QHealth1      float64 `spss:"QHealth1"`
		LngLst        float64 `spss:"LngLst"`
		LimitK        float64 `spss:"LimitK"`
		LimitA        float64 `spss:"LimitA"`
		Dintro        float64 `spss:"Dintro"`
		DisLmK        float64 `spss:"DisLmK"`
		DisLmA        float64 `spss:"DisLmA"`
		DOnset        float64 `spss:"DOnset"`
		DCause        float64 `spss:"DCause"`
		WkSSEmp       float64 `spss:"WkSSEmp"`
		DisMobl       float64 `spss:"DisMobl"`
		AsistPv       float64 `spss:"AsistPv"`
		AsistNd       float64 `spss:"AsistNd"`
		AsisFm1       float64 `spss:"AsisFm1"`
		AsisFm2       float64 `spss:"AsisFm2"`
		AsisFm3       float64 `spss:"AsisFm3"`
		AsisFm4       float64 `spss:"AsisFm4"`
		AsisFm5       float64 `spss:"AsisFm5"`
		AsisFm6       float64 `spss:"AsisFm6"`
		AsisFm7       float64 `spss:"AsisFm7"`
		AsisFm8       float64 `spss:"AsisFm8"`
		Heal01        float64 `spss:"Heal01"`
		Heal02        float64 `spss:"Heal02"`
		Heal03        float64 `spss:"Heal03"`
		Heal04        float64 `spss:"Heal04"`
		Heal05        float64 `spss:"Heal05"`
		Heal06        float64 `spss:"Heal06"`
		Heal07        float64 `spss:"Heal07"`
		Heal08        float64 `spss:"Heal08"`
		Heal09        float64 `spss:"Heal09"`
		Heal10        float64 `spss:"Heal10"`
		Heal11        float64 `spss:"Heal11"`
		Heal12        float64 `spss:"Heal12"`
		Heal13        float64 `spss:"Heal13"`
		Heal14        float64 `spss:"Heal14"`
		Heal15        float64 `spss:"Heal15"`
		Heal16        float64 `spss:"Heal16"`
		Heal17        float64 `spss:"Heal17"`
		LernD         float64 `spss:"LernD"`
		Health        float64 `spss:"Health"`
		LimAct        float64 `spss:"LimAct"`
		HealYr        float64 `spss:"HealYr"`
		HealPB01      float64 `spss:"HealPB01"`
		HealPB02      float64 `spss:"HealPB02"`
		HealPB03      float64 `spss:"HealPB03"`
		HealPB04      float64 `spss:"HealPB04"`
		HealPB05      float64 `spss:"HealPB05"`
		HealPB06      float64 `spss:"HealPB06"`
		HealPB07      float64 `spss:"HealPB07"`
		HealPB08      float64 `spss:"HealPB08"`
		HealPB09      float64 `spss:"HealPB09"`
		HealPB10      float64 `spss:"HealPB10"`
		LernDB        float64 `spss:"LernDB"`
		HealYL        float64 `spss:"HealYL"`
		Accdnt        float64 `spss:"Accdnt"`
		Road          float64 `spss:"Road"`
		WchJb         float64 `spss:"WchJb"`
		GoBack        float64 `spss:"GoBack"`
		GoBck9        float64 `spss:"GoBck9"`
		TimeDays      string  `spss:"TimeDays"`
		TimeCode      string  `spss:"TimeCode"`
		AccDay4       float64 `spss:"AccDay4"`
		TypInj        float64 `spss:"TypInj"`
		SiteFr1       float64 `spss:"SiteFr1"`
		SiteFr2       float64 `spss:"SiteFr2"`
		SiteFr3       float64 `spss:"SiteFr3"`
		SiteFr4       float64 `spss:"SiteFr4"`
		SiteFr5       float64 `spss:"SiteFr5"`
		SiteFr6       float64 `spss:"SiteFr6"`
		SiteFr7       float64 `spss:"SiteFr7"`
		SiteDi1       float64 `spss:"SiteDi1"`
		SiteDi2       float64 `spss:"SiteDi2"`
		SiteDi3       float64 `spss:"SiteDi3"`
		SiteDi4       float64 `spss:"SiteDi4"`
		SiteDi5       float64 `spss:"SiteDi5"`
		SiteDi6       float64 `spss:"SiteDi6"`
		AccurH1       float64 `spss:"AccurH1"`
		AccurH2       float64 `spss:"AccurH2"`
		AccurH3       float64 `spss:"AccurH3"`
		AccurH4       float64 `spss:"AccurH4"`
		AccKind       float64 `spss:"AccKind"`
		IllWrk        float64 `spss:"IllWrk"`
		NumIll        float64 `spss:"NumIll"`
		TypIll        float64 `spss:"TypIll"`
		TmeOff        float64 `spss:"TmeOff"`
		ILCurr        float64 `spss:"ILCurr"`
		WchJb3        float64 `spss:"WchJb3"`
		Aware         float64 `spss:"Aware"`
		ReasOff9      float64 `spss:"ReasOff9"`
		NoBack9       float64 `spss:"NoBack9"`
		SmokEver      float64 `spss:"SmokEver"`
		CigNow        float64 `spss:"CigNow"`
		RedAct        float64 `spss:"RedAct"`
		NRINT         float64 `spss:"NRINT"`
		CONCERN       float64 `spss:"CONCERN"`
		UNDERSTA      float64 `spss:"UNDERSTA"`
		AFFECT        float64 `spss:"AFFECT"`
		INTQ3         float64 `spss:"INTQ3"`
		WALK          float64 `spss:"WALK"`
		CONFBED       float64 `spss:"CONFBED"`
		WASH          float64 `spss:"WASH"`
		NOWASH        float64 `spss:"NOWASH"`
		ACT           float64 `spss:"ACT"`
		NOACT         float64 `spss:"NOACT"`
		PAIN          float64 `spss:"PAIN"`
		ANXIETY       float64 `spss:"ANXIETY"`
		SCORE1        float64 `spss:"SCORE1"`
		SCORE2        float64 `spss:"SCORE2"`
		IncNow        float64 `spss:"IncNow"`
		PayDoc        float64 `spss:"PayDoc"`
		DocNot        string  `spss:"DocNot"`
		PayIntro      float64 `spss:"PayIntro"`
		Gross99       float64 `spss:"Gross99"`
		GrsExp        float64 `spss:"GrsExp"`
		GrsPrd        float64 `spss:"GrsPrd"`
		BandG         string  `spss:"BandG"`
		UsGrs99       float64 `spss:"UsGrs99"`
		UsuGPay       float64 `spss:"UsuGPay"`
		UsBandG       string  `spss:"UsBandG"`
		Net99         float64 `spss:"Net99"`
		NetPrd        float64 `spss:"NetPrd"`
		BandN         string  `spss:"BandN"`
		IncChk        float64 `spss:"IncChk"`
		UsNet99       float64 `spss:"UsNet99"`
		UsuNPay       float64 `spss:"UsuNPay"`
		UsBandN       string  `spss:"UsBandN"`
		YVary99       float64 `spss:"YVary99"`
		YPayL         float64 `spss:"YPayL"`
		PaySSp        float64 `spss:"PaySSp"`
		YPayM         float64 `spss:"YPayM"`
		ErnFilt       float64 `spss:"ErnFilt"`
		Erncm0101     float64 `spss:"Erncm0101"`
		Erncm0102     float64 `spss:"Erncm0102"`
		Erncm0103     float64 `spss:"Erncm0103"`
		Erncm0104     float64 `spss:"Erncm0104"`
		Erncm0105     float64 `spss:"Erncm0105"`
		Erncm0106     float64 `spss:"Erncm0106"`
		Erncm0107     float64 `spss:"Erncm0107"`
		Erncm0108     float64 `spss:"Erncm0108"`
		Erncm0109     float64 `spss:"Erncm0109"`
		Erncm0110     float64 `spss:"Erncm0110"`
		Erncm0111     float64 `spss:"Erncm0111"`
		BonCmp1       float64 `spss:"BonCmp1"`
		BonCmp2       float64 `spss:"BonCmp2"`
		BonCmp3       float64 `spss:"BonCmp3"`
		BonCmp4       float64 `spss:"BonCmp4"`
		Hourly        float64 `spss:"Hourly"`
		HrRate        float64 `spss:"HrRate"`
		OvrTme        float64 `spss:"OvrTme"`
		UseSlp        float64 `spss:"UseSlp"`
		SecSta        float64 `spss:"SecSta"`
		Hourly2       float64 `spss:"Hourly2"`
		HrRate2       float64 `spss:"HrRate2"`
		SecGro        float64 `spss:"SecGro"`
		SecNet        float64 `spss:"SecNet"`
		ScNtGa        float64 `spss:"ScNtGa"`
		SecChk        float64 `spss:"SecChk"`
		SecEx         float64 `spss:"SecEx"`
		SecGA         float64 `spss:"SecGA"`
		SecGB         float64 `spss:"SecGB"`
		BandG2        string  `spss:"BandG2"`
		BandN2        string  `spss:"BandN2"`
		IrEnd2        float64 `spss:"IrEnd2"`
		RelBUp        float64 `spss:"RelBUp"`
		CryFth        float64 `spss:"CryFth"`
		CryFSpc       string  `spss:"CryFSpc"`
		CryFFrm       float64 `spss:"CryFFrm"`
		CryMth        float64 `spss:"CryMth"`
		CryMSpc       string  `spss:"CryMSpc"`
		CryMFrm       float64 `spss:"CryMFrm"`
		FathEdu       float64 `spss:"FathEdu"`
		MothEdu       float64 `spss:"MothEdu"`
		WkOthCry      float64 `spss:"WkOthCry"`
		CrySix        float64 `spss:"CrySix"`
		WchCry        string  `spss:"WchCry"`
		WchCryFr      float64 `spss:"WchCryFr"`
		ComeUK1       float64 `spss:"ComeUK1"`
		ComeUK2       float64 `spss:"ComeUK2"`
		ComeUK3       float64 `spss:"ComeUK3"`
		ComeUK4       float64 `spss:"ComeUK4"`
		ComeUK5       float64 `spss:"ComeUK5"`
		ComeUKMn      float64 `spss:"ComeUKMn"`
		JobInUK       float64 `spss:"JobInUK"`
		OvrQual       float64 `spss:"OvrQual"`
		ObSkilM       float64 `spss:"ObSkilM"`
		ObSkilS       float64 `spss:"ObSkilS"`
		ObJobM        float64 `spss:"ObJobM"`
		ObJobS        float64 `spss:"ObJobS"`
		LangSkil      float64 `spss:"LangSkil"`
		LangCour      float64 `spss:"LangCour"`
		FindJob       float64 `spss:"FindJob"`
		Soc10MU       string  `spss:"Soc10MU"`
		S2kStMU_index float64 `spss:"S2kStMU_index"`
		Soc2kMU_index string  `spss:"Soc2kMU_index"`
		SerNoMU_index float64 `spss:"SerNoMU_index"`
		Soc10MC       string  `spss:"Soc10MC"`
		S2kStMC_index float64 `spss:"S2kStMC_index"`
		Soc2kMC_index string  `spss:"Soc2kMC_index"`
		SerNoMC_index float64 `spss:"SerNoMC_index"`
		Soc102U       string  `spss:"Soc102U"`
		S2kSt2U_index float64 `spss:"S2kSt2U_index"`
		Soc2k2U_index string  `spss:"Soc2k2U_index"`
		SerNo2U_index float64 `spss:"SerNo2U_index"`
		Soc102C       string  `spss:"Soc102C"`
		S2kSt2C_index float64 `spss:"S2kSt2C_index"`
		Soc2k2C_index string  `spss:"Soc2k2C_index"`
		SerNo2C_index float64 `spss:"SerNo2C_index"`
		Soc10OY       string  `spss:"Soc10OY"`
		S2kStOY_index float64 `spss:"S2kStOY_index"`
		Soc2kOY_index string  `spss:"Soc2kOY_index"`
		SerNoOY_index float64 `spss:"SerNoOY_index"`
		Soc10Rd       string  `spss:"Soc10Rd"`
		S2kStRd_index float64 `spss:"S2kStRd_index"`
		Soc2kRd_index string  `spss:"Soc2kRd_index"`
		SerNoRd_index float64 `spss:"SerNoRd_index"`
		Soc10AU       string  `spss:"Soc10AU"`
		S2kStAU_index float64 `spss:"S2kStAU_index"`
		Soc2kAU_index string  `spss:"Soc2kAU_index"`
		SerNoAU_index float64 `spss:"SerNoAU_index"`
		Soc10AC       string  `spss:"Soc10AC"`
		S2kStAC_index float64 `spss:"S2kStAC_index"`
		Soc2kAC_index string  `spss:"Soc2kAC_index"`
		SerNoAC_index float64 `spss:"SerNoAC_index"`
		SearchTxtM    string  `spss:"SearchTxtM"`
		SearchTxt2    string  `spss:"SearchTxt2"`
		SearchTxtO    string  `spss:"SearchTxtO"`
		SearchTxtR    string  `spss:"SearchTxtR"`
		SearchTxtA    string  `spss:"SearchTxtA"`
		SearchTxtF    string  `spss:"SearchTxtF"`
		SearchTxtW    string  `spss:"SearchTxtW"`
		Phase         float64 `spss:"Phase"`
		ChangeOccM    float64 `spss:"ChangeOccM"`
		ChangeOcc2    float64 `spss:"ChangeOcc2"`
		EquivOccO     float64 `spss:"EquivOccO"`
		EquivOccR     float64 `spss:"EquivOccR"`
		EquivOccOY    float64 `spss:"EquivOccOY"`
		ChangeOccA    float64 `spss:"ChangeOccA"`
		CodeIntroM    float64 `spss:"CodeIntroM"`
		CodeIntro2    float64 `spss:"CodeIntro2"`
		CodeIntroO    float64 `spss:"CodeIntroO"`
		CodeIntroR    float64 `spss:"CodeIntroR"`
		CodeIntroA    float64 `spss:"CodeIntroA"`
		CodeIntroW    float64 `spss:"CodeIntroW"`
		Ocod10M       float64 `spss:"ocod10M"`
		S2kStM_index  float64 `spss:"S2kStM_index"`
		Soc2kM_index  string  `spss:"Soc2kM_index"`
		SerNoM_index  float64 `spss:"SerNoM_index"`
		Ocod102       float64 `spss:"ocod102"`
		S2kSt2_index  float64 `spss:"S2kSt2_index"`
		Soc2k2_index  string  `spss:"Soc2k2_index"`
		SerNo2_index  float64 `spss:"SerNo2_index"`
		Ocod10O       float64 `spss:"ocod10O"`
		S2kStO_index  float64 `spss:"S2kStO_index"`
		Soc2kO_index  string  `spss:"Soc2kO_index"`
		SerNoO_index  float64 `spss:"SerNoO_index"`
		Ocod10R       float64 `spss:"ocod10R"`
		S2kStR_index  float64 `spss:"S2kStR_index"`
		Soc2kR_index  string  `spss:"Soc2kR_index"`
		SerNoR_index  float64 `spss:"SerNoR_index"`
		Ocod10A       float64 `spss:"ocod10A"`
		S2kStA_index  float64 `spss:"S2kStA_index"`
		Soc2kA_index  string  `spss:"Soc2kA_index"`
		SerNoA_index  float64 `spss:"SerNoA_index"`
		S2kStW_index  float64 `spss:"S2kStW_index"`
		Soc2kW_index  string  `spss:"Soc2kW_index"`
		SerNoW_index  float64 `spss:"SerNoW_index"`
		Ocod2KM       float64 `spss:"ocod2KM"`
		Ocod2K2       float64 `spss:"ocod2K2"`
		Ocod2KO       float64 `spss:"ocod2KO"`
		Ocod2KR       float64 `spss:"ocod2KR"`
		Ocod2KA       float64 `spss:"ocod2KA"`
		SECM          float64 `spss:"SECM"`
		SECFlagM      float64 `spss:"SECFlagM"`
		ES2000M       float64 `spss:"ES2000M"`
		SECMIES       float64 `spss:"SECMIES"`
		SEC2          float64 `spss:"SEC2"`
		SECFlag2      float64 `spss:"SECFlag2"`
		ES20002       float64 `spss:"ES20002"`
		SEC2IES       float64 `spss:"SEC2IES"`
		SECO          float64 `spss:"SECO"`
		SECFlagO      float64 `spss:"SECFlagO"`
		ES2000O       float64 `spss:"ES2000O"`
		SECOIES       float64 `spss:"SECOIES"`
		SECR          float64 `spss:"SECR"`
		SECFlagR      float64 `spss:"SECFlagR"`
		ES2000R       float64 `spss:"ES2000R"`
		SECRIES       float64 `spss:"SECRIES"`
		ICdM          float64 `spss:"ICdM"`
		ICod2007      string  `spss:"ICod2007"`
		ICod07        float64 `spss:"ICod07"`
		ICdC2007      string  `spss:"ICdC2007"`
		ICdC07        float64 `spss:"ICdC07"`
		IEmpStat      float64 `spss:"IEmpStat"`
		IEmpStatC     float64 `spss:"IEmpStatC"`
		IMPSTM        float64 `spss:"IMPSTM"`
		ICd2          float64 `spss:"ICd2"`
		ICd22007      string  `spss:"ICd22007"`
		ICd207        float64 `spss:"ICd207"`
		ICC22007      string  `spss:"ICC22007"`
		ICC207        float64 `spss:"ICC207"`
		IEmpStat2     float64 `spss:"IEmpStat2"`
		IEmpStat2C    float64 `spss:"IEmpStat2C"`
		IMPST2        float64 `spss:"IMPST2"`
		OYIC2007      string  `spss:"OYIC2007"`
		OYICd07       float64 `spss:"OYICd07"`
		OYIEmpStat    float64 `spss:"OYIEmpStat"`
		RdIC2007      string  `spss:"RdIC2007"`
		RdICd07       float64 `spss:"RdICd07"`
		CodChk        float64 `spss:"CodChk"`
		IndOut        float64 `spss:"IndOut"`
		PERSFLAG      float64 `spss:"PERSFLAG"`
		RPERSFLAG     float64 `spss:"RPERSFLAG"`
		RECALLP       float64 `spss:"RECALLP"`
		CodeNow       float64 `spss:"CodeNow"`
		HiHNum        string  `spss:"HiHNum"`
		IsHIC1        float64 `spss:"isHIC1"`
		JntEldA       float64 `spss:"JntEldA"`
		JntEldB       float64 `spss:"JntEldB"`
		DVHRPNUM      float64 `spss:"DVHRPNUM"`
		HRPCheck      float64 `spss:"HRPCheck"`
		ShowBen       float64 `spss:"ShowBen"`
		DVBenU1       float64 `spss:"DVBenU1"`
		DVnumBU       float64 `spss:"DVnumBU"`
		Thanks        float64 `spss:"Thanks"`
		ThankE        float64 `spss:"ThankE"`
		ThankWvF      float64 `spss:"ThankWvF"`
		ThankEth      float64 `spss:"ThankEth"`
		Flag75        float64 `spss:"Flag75"`
		Thank75a      float64 `spss:"Thank75a"`
		Thank75b      float64 `spss:"Thank75b"`
		RecallH       float64 `spss:"RecallH"`
		GotPhone      float64 `spss:"GotPhone"`
		RecPhone      float64 `spss:"RecPhone"`
		Chk_Num1      float64 `spss:"Chk_Num1"`
		PrefNo        float64 `spss:"PrefNo"`
		Chk_Num2      float64 `spss:"Chk_Num2"`
		AltNo         float64 `spss:"AltNo"`
		Display       float64 `spss:"Display"`
		AppointType   float64 `spss:"AppointType"`
		WeekDays1     float64 `spss:"WeekDays1"`
		WeekDays2     float64 `spss:"WeekDays2"`
		WeekDays3     float64 `spss:"WeekDays3"`
		WeekDays4     float64 `spss:"WeekDays4"`
		WeekDays5     float64 `spss:"WeekDays5"`
		WeekDays6     float64 `spss:"WeekDays6"`
		AppointTime   float64 `spss:"AppointTime"`
		TimeStart     float64 `spss:"TimeStart"`
		TimeEnd       float64 `spss:"TimeEnd"`
		CalSun        float64 `spss:"CalSun"`
		MultHh11      float64 `spss:"MultHh11"`
		MHHType       float64 `spss:"MHHType"`
		AddRes        float64 `spss:"AddRes"`
		AccSh         float64 `spss:"AccSh"`
		ShaLiv        float64 `spss:"ShaLiv"`
		MealSh        float64 `spss:"MealSh"`
		Shacook       float64 `spss:"Shacook"`
		MltiOld1      float64 `spss:"MltiOld1"`
		MltiOld2      float64 `spss:"MltiOld2"`
		MultiNew      float64 `spss:"MultiNew"`
		NumHhld       float64 `spss:"NumHhld"`
		HhldDesc      string  `spss:"HhldDesc"`
		Multhhld      float64 `spss:"Multhhld"`
		IntIntrp      float64 `spss:"IntIntrp"`
		NonEng        float64 `spss:"NonEng"`
		WhLang01      float64 `spss:"WhLang01"`
		WhLang02      float64 `spss:"WhLang02"`
		WhLang03      float64 `spss:"WhLang03"`
		WhLang04      float64 `spss:"WhLang04"`
		WhLang05      float64 `spss:"WhLang05"`
		WhLang06      float64 `spss:"WhLang06"`
		WhLang07      float64 `spss:"WhLang07"`
		WhLang08      float64 `spss:"WhLang08"`
		WhLang09      float64 `spss:"WhLang09"`
		WhLang10      float64 `spss:"WhLang10"`
		WhLang11      float64 `spss:"WhLang11"`
		WhlangO       string  `spss:"WhlangO"`
		WhoTrans1     float64 `spss:"WhoTrans1"`
		WhoTrans2     float64 `spss:"WhoTrans2"`
		WhoTrans3     float64 `spss:"WhoTrans3"`
		WhoTrans4     float64 `spss:"WhoTrans4"`
		WhoTrans5     float64 `spss:"WhoTrans5"`
		WhoTrans6     float64 `spss:"WhoTrans6"`
		NmTrans       float64 `spss:"NmTrans"`
		Iout1         float64 `spss:"Iout1"`
		Iout2         float64 `spss:"Iout2"`
		Iout4         float64 `spss:"Iout4"`
		ProxPers      float64 `spss:"ProxPers"`
		AgeCheck      float64 `spss:"AgeCheck"`
		NumFuPer      float64 `spss:"NumFuPer"`
		NumPaPry      float64 `spss:"NumPaPry"`
		NumFuPry      float64 `spss:"NumFuPry"`
		NumPaPer      float64 `spss:"NumPaPer"`
		NumNoElg      float64 `spss:"NumNoElg"`
		NumRefus      float64 `spss:"NumRefus"`
		NumNonco      float64 `spss:"NumNonco"`
		NumParls      float64 `spss:"NumParls"`
		NumFull       float64 `spss:"NumFull"`
		NumInt        float64 `spss:"NumInt"`
		NumNotEl      float64 `spss:"NumNotEl"`
		NumSevIn      float64 `spss:"NumSevIn"`
		IndCheck      float64 `spss:"IndCheck"`
		HarmIntr      float64 `spss:"HarmIntr"`
		Intsome       float64 `spss:"Intsome"`
		IntFin        float64 `spss:"IntFin"`
		Outsum        float64 `spss:"Outsum"`
		IndOut4       float64 `spss:"IndOut4"`
		Inelig1       float64 `spss:"Inelig1"`
		Uncer1        float64 `spss:"Uncer1"`
		NonSum        float64 `spss:"NonSum"`
		Ref1          float64 `spss:"Ref1"`
		Ref2          float64 `spss:"Ref2"`
		Ref3          float64 `spss:"Ref3"`
		Refreas1      float64 `spss:"Refreas1"`
		Refreas2      float64 `spss:"Refreas2"`
		Refreas3      float64 `spss:"Refreas3"`
		Nonc1         float64 `spss:"Nonc1"`
		Nonreas1      float64 `spss:"Nonreas1"`
		Nonreas2      float64 `spss:"Nonreas2"`
		Nonreas3      float64 `spss:"Nonreas3"`
		Othr1         float64 `spss:"Othr1"`
		Othr2         float64 `spss:"Othr2"`
		Othr3         float64 `spss:"Othr3"`
		Hout04        float64 `spss:"Hout04"`
		HoutLFS       float64 `spss:"HoutLFS"`
		AnyVisit      float64 `spss:"AnyVisit"`
		RtypHH        float64 `spss:"RtypHH"`
		RTypOth       string  `spss:"RTypOth"`
		DwellTyp      float64 `spss:"DwellTyp"`
		FloorN        float64 `spss:"FloorN"`
		EntryN1       float64 `spss:"EntryN1"`
		EntryN2       float64 `spss:"EntryN2"`
		EntryN3       float64 `spss:"EntryN3"`
		EntryN4       float64 `spss:"EntryN4"`
		EntryN5       float64 `spss:"EntryN5"`
		EntryN6       float64 `spss:"EntryN6"`
		TotTime       float64 `spss:"TotTime"`
		Main          float64 `spss:"Main"`
		IntvLang      float64 `spss:"IntvLang"`
		FTFphone      float64 `spss:"FTFphone"`
		Direction     string  `spss:"Direction"`
		BriefSDC1     float64 `spss:"BriefSDC1"`
		BriefSDC2     float64 `spss:"BriefSDC2"`
		BriefSDC3     float64 `spss:"BriefSDC3"`
		Brief1        string  `spss:"Brief1"`
		LBrief1       string  `spss:"LBrief1"`
		Brief2        string  `spss:"Brief2"`
		AnyLeft       float64 `spss:"AnyLeft"`
		DoneCode      float64 `spss:"DoneCode"`
		RefNon        float64 `spss:"RefNon"`
		Refuse1       float64 `spss:"Refuse1"`
		Refuse2       float64 `spss:"Refuse2"`
		Refuse3       float64 `spss:"Refuse3"`
		RefOth        string  `spss:"RefOth"`
		NnCont        float64 `spss:"NnCont"`
		NnCOth        string  `spss:"NnCOth"`
		ReIssue       float64 `spss:"ReIssue"`
		ReOther       string  `spss:"ReOther"`
		Iss1Int       float64 `spss:"Iss1Int"`
		Iss1HOut      float64 `spss:"Iss1HOut"`
		Iss1NC        float64 `spss:"Iss1NC"`
		Iss1CRef      float64 `spss:"Iss1CRef"`
		ReIss         float64 `spss:"ReIss"`
		HOut          float64 `spss:"HOut"`
		HOutC         float64 `spss:"HOutC"`
		HOutDate      float64 `spss:"HOutDate"`
		CRef          float64 `spss:"CRef"`
		TOIntOut      float64 `spss:"TOIntOut"`
		ChkLet        string  `spss:"ChkLet"`
		DivAddInd     string  `spss:"DivAddInd"`
		CaseOjectInd  string  `spss:"CaseOjectInd"`
		GORA          string  `spss:"GORA"`
		CountryCod    string  `spss:"CountryCod"`
		Wave          string  `spss:"Wave"`
		ChkLet2       string  `spss:"ChkLet2"`
		OAInd         string  `spss:"OAInd"`
		URIndEW       string  `spss:"URIndEW"`
		URIndSc       string  `spss:"URIndSc"`
		NHSAcc        string  `spss:"NHSAcc"`
		GOR99         string  `spss:"GOR99"`
		TelFTF        string  `spss:"TelFTF"`
		SupVsMU_index string  `spss:"SupVsMU_index"`
		SupVsMC_index string  `spss:"SupVsMC_index"`
		SupVs2U_index string  `spss:"SupVs2U_index"`
		SupVs2C_index string  `spss:"SupVs2C_index"`
		SupVsOY_index string  `spss:"SupVsOY_index"`
		SupVsRd_index string  `spss:"SupVsRd_index"`
		SupVsAU_index string  `spss:"SupVsAU_index"`
		SupVsAC_index string  `spss:"SupVsAC_index"`
		SupVsM_index  string  `spss:"SupVsM_index"`
		SupVs2_index  string  `spss:"SupVs2_index"`
		SupVsO_index  string  `spss:"SupVsO_index"`
		SupVsR_index  string  `spss:"SupVsR_index"`
		SupVsA_index  string  `spss:"SupVsA_index"`
		SupVsW_index  string  `spss:"SupVsW_index"`
		BULater       string  `spss:"BULater"`
		FldReg        string  `spss:"FldReg"`
		Dteofbth      string  `spss:"dteofbth"`
		CID96_new     string  `spss:"CID96_new"`
		LAD96_new     string  `spss:"LAD96_new"`
		PrimaryFirst  float64 `spss:"PrimaryFirst"`
		UALAD99       string  `spss:"UALAD99"`
		Urindew_new   string  `spss:"urindew_new"`
		Urindsc_new   string  `spss:"urindsc_new"`
		Urind         string  `spss:"urind"`
	}

	logger := log.New()

	log.Info("created dataset")
	d, err := NewDataset("test", logger)
	if err != nil {
		t.Fatalf("NewDataset failed: %s\n", err)
	}

	log.Info("loading dataset")
	dataset, err := d.FromSav(testDirectory()+"LFSwk18PERS_non_confidential.sav", TestDataset{})
	log.Info("Got dataset from sav")
	if err != nil {
		t.Fatalf("fromSav failed: %s\n", err)
	}
	defer dataset.Close()

	err = dataset.ToSQL()
	if err != nil {
		t.Fatalf("toSQL failed: %s\n", err)
	}

	t.Logf("Dataset Size: %d\n", dataset.NumRows())
	_ = dataset.Head(5)
}

func testDirectory() (testDirectory string) {
	testDirectory = conf.Config.TestDirectory

	if testDirectory == "" {
		panic("Add the TEST_DIRECTORY in configuration")
	}
	return
}
