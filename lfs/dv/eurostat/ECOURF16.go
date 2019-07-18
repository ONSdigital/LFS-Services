package eurostat

func ECOURF16(ioutcome int, thiswv int, ecoura16 int, tsub4cod float64) int {

	if ioutcome == 3 {
		return -8
	}

	if thiswv != 3 {
		return 999
	}

	if ecoura16 == 1 {
		return 999
	}

	if tsub4cod == 1 || tsub4cod == 8 || tsub4cod == 9 {
		return 0
	}

	if (tsub4cod >= 14 && tsub4cod <= 20) && (tsub4cod != 14.1) {
		return 10
	}

	if tsub4cod >= 21 && tsub4cod <= 22.6 {
		return 20
	}
	if tsub4cod >= 31 && tsub4cod <= 32.2 {
		return 30
	}
	if tsub4cod >= 34 && tsub4cod <= 38 {
		return 40
	}
	if tsub4cod >= 42 && tsub4cod <= 42.6 {
		return 50
	}
	if tsub4cod >= 48 && tsub4cod <= 48.2 {
		return 60
	}
	if tsub4cod >= 52 && tsub4cod <= 58.2 {
		return 70
	}
	if tsub4cod >= 62 && tsub4cod <= 62.4 {
		return 80
	}
	if tsub4cod == 64 {
		return 80
	}
	if tsub4cod >= 72 && tsub4cod <= 76.2 {
		return 90
	}
	if tsub4cod >= 81 && tsub4cod <= 86.3 {
		return 100
	}

	return 888
}
