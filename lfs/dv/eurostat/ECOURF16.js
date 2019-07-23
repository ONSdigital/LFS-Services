/**
 * @return {string}
 */
function ECOURF16(ioutcome, thiswv, ecoura16, tsub4cod) {

    if (ioutcome === 3) return "-8";
    if (thiswv !== 3) return "999";
    if (ecoura16 === 1) return "999";

    var tsub4codList = [1, 8, 9];
    if (tsub4codList.indexOf(tsub4cod) !== -1) return "000";

    if ((tsub4cod >= 14 && tsub4cod <= 20) && (tsub4cod !== 14.1)) return "010";

    if (tsub4cod >= 21 && tsub4cod <= 22.6) return "020";
    if (tsub4cod >= 31 && tsub4cod <= 32.2) return "030";
    if (tsub4cod >= 34 && tsub4cod <= 38) return "040";
    if (tsub4cod >= 42 && tsub4cod <= 42.6) return "050";
    if (tsub4cod >= 48 && tsub4cod <= 48.2) return "060";
    if (tsub4cod >= 52 && tsub4cod <= 58.2) return "070";
    if (tsub4cod >= 62 && tsub4cod <= 62.4) return "080";
    if (tsub4cod === 64) return "080";
    if (tsub4cod >= 72 && tsub4cod <= 76.2) return "090";
    if (tsub4cod >= 81 && tsub4cod <= 86.3) return "100";

    return "888"
}