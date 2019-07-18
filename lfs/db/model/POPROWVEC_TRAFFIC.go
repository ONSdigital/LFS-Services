package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type POPROWVECTRAFFIC struct {
	CGROUP sql.NullInt64   `gorm:"column:C_GROUP;primary_key"`
	T1     sql.NullFloat64 `gorm:"column:T_1"`
	T2     sql.NullFloat64 `gorm:"column:T_2"`
	T3     sql.NullFloat64 `gorm:"column:T_3"`
	T4     sql.NullFloat64 `gorm:"column:T_4"`
	T5     sql.NullFloat64 `gorm:"column:T_5"`
	T6     sql.NullFloat64 `gorm:"column:T_6"`
	T7     sql.NullFloat64 `gorm:"column:T_7"`
	T8     sql.NullFloat64 `gorm:"column:T_8"`
	T9     sql.NullFloat64 `gorm:"column:T_9"`
	T10    sql.NullFloat64 `gorm:"column:T_10"`
	T11    sql.NullFloat64 `gorm:"column:T_11"`
	T12    sql.NullFloat64 `gorm:"column:T_12"`
	T13    sql.NullFloat64 `gorm:"column:T_13"`
	T14    sql.NullFloat64 `gorm:"column:T_14"`
	T15    sql.NullFloat64 `gorm:"column:T_15"`
	T16    sql.NullFloat64 `gorm:"column:T_16"`
	T17    sql.NullFloat64 `gorm:"column:T_17"`
	T18    sql.NullFloat64 `gorm:"column:T_18"`
	T19    sql.NullFloat64 `gorm:"column:T_19"`
	T20    sql.NullFloat64 `gorm:"column:T_20"`
	T21    sql.NullFloat64 `gorm:"column:T_21"`
	T22    sql.NullFloat64 `gorm:"column:T_22"`
	T23    sql.NullFloat64 `gorm:"column:T_23"`
	T24    sql.NullFloat64 `gorm:"column:T_24"`
	T25    sql.NullFloat64 `gorm:"column:T_25"`
	T26    sql.NullFloat64 `gorm:"column:T_26"`
	T27    sql.NullFloat64 `gorm:"column:T_27"`
	T28    sql.NullFloat64 `gorm:"column:T_28"`
	T29    sql.NullFloat64 `gorm:"column:T_29"`
	T30    sql.NullFloat64 `gorm:"column:T_30"`
	T31    sql.NullFloat64 `gorm:"column:T_31"`
	T32    sql.NullFloat64 `gorm:"column:T_32"`
	T33    sql.NullFloat64 `gorm:"column:T_33"`
	T34    sql.NullFloat64 `gorm:"column:T_34"`
	T35    sql.NullFloat64 `gorm:"column:T_35"`
	T36    sql.NullFloat64 `gorm:"column:T_36"`
	T37    sql.NullFloat64 `gorm:"column:T_37"`
	T38    sql.NullFloat64 `gorm:"column:T_38"`
	T39    sql.NullFloat64 `gorm:"column:T_39"`
	T40    sql.NullFloat64 `gorm:"column:T_40"`
	T41    sql.NullFloat64 `gorm:"column:T_41"`
	T42    sql.NullFloat64 `gorm:"column:T_42"`
	T43    sql.NullFloat64 `gorm:"column:T_43"`
	T44    sql.NullFloat64 `gorm:"column:T_44"`
	T45    sql.NullFloat64 `gorm:"column:T_45"`
	T46    sql.NullFloat64 `gorm:"column:T_46"`
	T47    sql.NullFloat64 `gorm:"column:T_47"`
	T48    sql.NullFloat64 `gorm:"column:T_48"`
	T49    sql.NullFloat64 `gorm:"column:T_49"`
	T50    sql.NullFloat64 `gorm:"column:T_50"`
	T51    sql.NullFloat64 `gorm:"column:T_51"`
	T52    sql.NullFloat64 `gorm:"column:T_52"`
	T53    sql.NullFloat64 `gorm:"column:T_53"`
	T54    sql.NullFloat64 `gorm:"column:T_54"`
	T55    sql.NullFloat64 `gorm:"column:T_55"`
	T56    sql.NullFloat64 `gorm:"column:T_56"`
	T57    sql.NullFloat64 `gorm:"column:T_57"`
	T58    sql.NullFloat64 `gorm:"column:T_58"`
	T59    sql.NullFloat64 `gorm:"column:T_59"`
	T60    sql.NullFloat64 `gorm:"column:T_60"`
	T61    sql.NullFloat64 `gorm:"column:T_61"`
	T62    sql.NullFloat64 `gorm:"column:T_62"`
	T63    sql.NullFloat64 `gorm:"column:T_63"`
	T64    sql.NullFloat64 `gorm:"column:T_64"`
	T65    sql.NullFloat64 `gorm:"column:T_65"`
	T66    sql.NullFloat64 `gorm:"column:T_66"`
	T67    sql.NullFloat64 `gorm:"column:T_67"`
	T68    sql.NullFloat64 `gorm:"column:T_68"`
	T69    sql.NullFloat64 `gorm:"column:T_69"`
	T70    sql.NullFloat64 `gorm:"column:T_70"`
	T71    sql.NullFloat64 `gorm:"column:T_71"`
	T72    sql.NullFloat64 `gorm:"column:T_72"`
	T73    sql.NullFloat64 `gorm:"column:T_73"`
	T74    sql.NullFloat64 `gorm:"column:T_74"`
	T75    sql.NullFloat64 `gorm:"column:T_75"`
	T76    sql.NullFloat64 `gorm:"column:T_76"`
	T77    sql.NullFloat64 `gorm:"column:T_77"`
}

// TableName sets the insert table name for this struct type
func (p *POPROWVECTRAFFIC) TableName() string {
	return "POPROWVEC_TRAFFIC"
}
