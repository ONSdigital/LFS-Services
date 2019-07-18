package eurostat

import (
	"pds-go/lfs/util"
	"strings"
)

var lookupTable = util.LookupTable{
	{1, util.StringLookup{
		"  CH", "  CJ", "  CK", "  CL", "  CM",
	}},

	{2, util.StringLookup{
		"  EB", "  EC", "  EE", "  EF", "  EH", "16UB", "16UC", "16UD", "16UE",
		"16UF", "16UG", "20UB", "20UD", "20UE", "20UF", "20UG", "20UH", "20UJ", "35UB", "35UC", "35UD",
		"35UE", "35UF", "35UG",
	}},

	{3, util.StringLookup{
		"  CX", "  CY", "  CZ", "  DA", "  DB",
	}},

	{4, util.StringLookup{
		"  CX", "  CY", "  CZ", "  DA", "  DB",
	}},

	{5, util.StringLookup{
		"  FA", "  FB", "  FC", "  FD", "  FF", "36UB", "36UC", "36UD", "36UE",
		"36UF", "36UG", "36UH",
	}},

	{6, util.StringLookup{
		"  FK", "  FN", "  FP", "  FY", "17UB", "17UC", "17UD", "17UF", "17UG", "17UH", "17UJ", "17UK",
		"31UB", "31UC", "31UD", "31UE", "31UG", "31UH", "31UJ", "32UB", "32UC", "32UD", "32UE", "32UF",
		"32UG", "32UH", "34UB", "34UC", "34UD", "34UE", "34UF", "34UG", "34UH", "37UB", "37UC", "37UD",
		"37UE", "37UF", "37UG", "37UJ",
	}},

	{7, util.StringLookup{
		"  JA", "12UB", "12UC", "12UD", "12UE", "12UG", "33UB", "33UC", "33UD", "33UE", "33UF", "33UG", "33UH",
		"42UB", "42UC", "42UD", "42UE", "42UF", "42UG", "42UH",
	}},
}

func regwkr(ioutcome int, a int, b int, c string, d int, lad96, wad96 string) int {
	if ioutcome == 3 {
		return -9
	}

	regwp := -9
	LADWAD := strings.TrimSpace(lad96) + strings.TrimSpace(wad96)

	if (a == 4) || (a == -9 || a == -8) {

		if d == 999997 {
			return 23
		}

		found, ret := lookupTable.Contains(c)
		if !found {
			return -9
		}
		return ret



		//	ELSE IF (&c.in ("  AA", "  AG", "  AM", "  AN", "  AP", "  AU", "  AW", "  AY", "  AZ", "  BB", "  BE", "  BG", "  BJ", "  BK")) THEN DO;
		//	 IF (&d.in (100201, 101024, 101691,101903, 103531, 113676, 107109, 110328,112213, 117919, 123040, 104323, 123043,
		//	  101397, 102024, 106920, 109969, 117987,118370, 118440, 120121, 107277, 111036,123038, 115512, 118408,
		//	  104426, 117876,111380, 101556, 102181, 112268, 114253,118218, 118294, 100948, 101458, 102135,103906,
		//	  105045, 123039, 111400, 111690,111891, 112938, 113098, 113153, 115207,115229, 115257, 115588, 115648,
		//	  116379,117972, 118294, 118385, 119029, 120883,121687, 121763, 121996, 123046)) then REGWP = (8);
		//	ELSE REGWP = (9);
		//	END;
		//	ELSE IF (&c.in ("  AB", "  AC", "  AD", "  AE", "  AF", "  AH", "  AJ", "  AK", "  AL", "  AQ", "  AR", "  AS", "  AT", "  AX",
		//	"  BA", "  BC", "  BD", "  BF", "  BH")) THEN REGWP = (11);
		//	ELSE IF (&c.in ("  KA", "  KF", "  KG", "  LC", "  MA", "  MB", "  MC", "  MD", "  ME", "  MF", "  MG", "  ML", "  MR", "  MS",
		//	"  MW", "09UC", "09UD", "09UE", "11UB", "11UC", "11UE", "11UF", "21UC", "21UD", "21UF", "21UG",
		//	"21UH", "22UB", "22UC", "22UD", "22UE", "22UF", "22UG", "22UH", "22UJ",
		//	"22UK", "22UL", "22UN", "22UQ", "24UB", "24UC", "24UD", "24UE", "24UF", "24UG",
		//	"24UH", "24UJ", "24UL", "24UN", "24UP", "26UB", "26UC", "26UD", "26UE", "26UF",
		//	"26UG", "26UH", "26UJ", "26UK", "26UL", "29UB", "29UC", "29UD",
		//	"29UE", "29UG", "29UH", "29UK", "29UL", "29UM", "29UN", "29UP", "29UQ", "38UB",
		//	"38UC", "38UD", "38UE", "38UF", "43UB", "43UC", "43UD", "43UE", "43UF", "43UG",
		//	"43UH", "43UJ", "43UK", "43UL", "43UM", "45UB", "45UC", "45UD", "45UE", "45UF",
		//	"45UG", "45UH")) then REGWP= (12);
		//	ELSE IF (&c.in ("  HA", "  HB", "  HC", "  HD", "  HG", "  HH", "  HN", "  HP", "  HX",
		//	"15UB", "15UC", "15UD", "15UE", "15UF", "15UG", "15UH",
		//	"18UB", "18UC", "18UD", "18UE", "18UG", "18UH", "18UK", "18UL",
		//	"19UC", "19UD", "19UE", "19UG", "19UH", "19UJ",
		//	"23UB", "23UC", "23UD", "23UE", "23UF", "23UG",
		//	"40UB", "40UC", "40UD", "40UE", "40UF",
		//	"46UB", "46UC", "46UD", "46UF")) then REGWP = (13) ;
		//	ELSE IF (&c.in ("  CN", "  CQ", "  CR", "  CS", "  CT", "  CU", "  CW")) THEN REGWP = (14) ;
		//	ELSE IF (&c.in ("  GA", "  GF", "  GL", "39UB", "39UC", "39UD", "39UE", "39UF",
		//	"41UB", "41UC", "41UD", "41UE", "41UF", "41UG", "41UH", "41UK",
		//	"44UB", "44UC", "44UD", "44UE", "44UF",
		//	"47UB", "47UC", "47UD", "47UE", "47UF", "47UG")) THEN REGWP = (15) ;
		//	ELSE IF (&c.in ("  BL", "  BM", "  BN", "  BP", "  BQ", "  BR", "  BS", "  BT", "  BU", "  BW")) THEN REGWP = (16);
		//	ELSE IF (&c.in ("  BX", "  BY", "  BZ", "  CA", "  CB")) THEN REGWP = (17);
		//	ELSE IF (&c.in ("  ET", "  EU", "  EX", "  EY", "13UB", "13UC", "13UD", "13UE", "13UG", "13UH",
		//	"30UD", "30UE", "30UF", "30UG", "30UH", "30UJ", "30UK", "30UL", "30UM", "30UN", "30UP", "30UQ")) THEN REGWP = (18);
		//	ELSE IF (&c.in ("  NA", "  NC", "  NE", "  NG", "  NJ", "  NL", "  NN", "  NQ", "  NU", "  NS", "  NX", "  NZ",
		//	"  PB", "  PD", "  PF", "  PH", "  PK", "  PL", "  PM", "  PP", "  PR", "  PT")) THEN REGWP = (19);
		//	ELSE IF (&c.in ("  QD", "  QG", "  QK", "  QL", "  QN", "  QS", "  QU", "  QY", "  QZ", "  RC", "  RE", "  RF")) THEN REGWP = (20);
		//	ELSE IF (&c.in ("  QA", "  QB", "  QC", "  QE", "  QF", "  QH", "  QJ", "  QM", "  QP", "  QQ", "  QR", "  QT", "  QW", "  QX",
		//	"  RA", "  RB", "  RD", "  RG", "  RH", "  RJ")) THEN REGWP = (21);
		//	ELSE IF (&c.in (" 010", " 020", " 030", " 040", " 050", " 060", " 070", " 080", " 090",
		//	" 100", " 110", " 120", " 130", " 140", " 150", " 160", " 170", " 180", " 190",
		//	" 200", " 210", " 220", " 230", " 240", " 250", " 260", " 460")) THEN REGWP = (22);
		//	ELSE REGWP = (-8);
		//	}
		//  ELSE DO ;
		//     IF (URESMC in (8,9)) THEN do ;
		//         IF (URESMC EQ 9) THEN REGWP =(11) ;
		//         ELSE IF (UALAD99 EQ "  AA") THEN REGWP =(8) ;
		//         ELSE IF (UALAD99 in ("  AG","  AU","  AW","  BE","  BK")) THEN DO ;
		//             IF (LADWAD in ("AGFT","AGFC","AGFR","AGFD","AGFZ","AUFE","AUFB","AWFL","BEFJ","BEFK","BEFU","BKFA","BKFC",
		//                                       "BKFD","BKFE","BKFF","BKFL","BKFK","BKFR","BKFU","BKFW","BKFX","BKFZ")) THEN regwp=(8) ;
		//             else regwp =(9) ;
		//         end;
		//         else regwp =(9) ;
		//     end;
		//     else if (uresmc ge 1 and uresmc le 7) then regwp = uresmc ;
		//   else regwp = uresmc + 2 ;
		//end;
	}

	if (statr in(1, 2, 4)) then
	do;
	%RegProc(Home, indm92m, ualdwk, wkpl99);
	RegWkr = RegWp;
	end;
	else RegWkr = (-9);
	end;
}
