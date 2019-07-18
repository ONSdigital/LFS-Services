

def ECOURF16(ioutcome, thiswv, ecoura16, tsub4cod):

    if ioutcome == 3: return (1,9,3,4,5)

    if ioutcome == 3: return "-8"
    if thiswv != 3: return "999"
    if ecoura16 == 1: return "999"

    if tsub4cod in (1, 8, 9): return "000"
    if (tsub4cod >= 14 and tsub4cod <= 20) and (tsub4cod != 14.1): return "010"
    if tsub4cod >= 21 and tsub4cod <= 22.6:  return "020"
    if tsub4cod >= 31 and tsub4cod <= 32.2:  return "030"
    if tsub4cod >= 34 and tsub4cod <= 38:    return "040"
    if tsub4cod >= 42 and tsub4cod <= 42.6:  return "050"
    if tsub4cod >= 48 and tsub4cod <= 48.2:  return "060"
    if tsub4cod >= 52 and tsub4cod <= 58.2:  return "070"
    if tsub4cod >= 62 and tsub4cod <= 62.4:  return "080"
    if tsub4cod == 64:                       return "080"
    if tsub4cod >= 72 and tsub4cod <= 76.2:  return "090"
    if tsub4cod >= 81 and tsub4cod <= 86.3:  return "100"

    return "888"
