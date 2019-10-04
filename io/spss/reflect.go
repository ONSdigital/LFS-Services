package spss

import (
	"reflect"
	"strings"
	"sync"
)

type StructInfo struct {
	Fields []FieldInfo
}

type FieldInfo struct {
	Keys       []string
	FieldType  reflect.Kind
	OmitEmpty  bool
	IndexChain []int
}

func (f FieldInfo) getFirstKey() string {
	return f.Keys[0]
}

func (f FieldInfo) MatchesKey(key string) bool {
	for _, k := range f.Keys {
		if key == k || strings.TrimSpace(key) == k {
			return true
		}
	}
	return false
}

var structMap = make(map[reflect.Type]*StructInfo)
var structMapMutex sync.RWMutex

func GetStructInfo(rType reflect.Type) *StructInfo {
	structMapMutex.RLock()
	stInfo, ok := structMap[rType]
	structMapMutex.RUnlock()
	if ok {
		return stInfo
	}
	fieldsList := GetFieldInfos(rType, []int{})
	stInfo = &StructInfo{fieldsList}
	return stInfo
}

func GetFieldInfos(rType reflect.Type, parentIndexChain []int) []FieldInfo {
	fieldsCount := rType.NumField()
	fieldsList := make([]FieldInfo, 0, fieldsCount)
	for i := 0; i < fieldsCount; i++ {
		field := rType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		var cpy = make([]int, len(parentIndexChain))
		copy(cpy, parentIndexChain)
		indexChain := append(cpy, i)

		// if the field is a pointer to a struct, follow the pointer then create fieldinfo for each field
		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if !(canMarshal(field.Type.Elem())) {
				fieldsList = append(fieldsList, GetFieldInfos(field.Type.Elem(), indexChain)...)
			}
		}
		// if the field is a struct, create a fieldInfo for each of its fields
		if field.Type.Kind() == reflect.Struct {
			if !(canMarshal(field.Type)) {
				fieldsList = append(fieldsList, GetFieldInfos(field.Type, indexChain)...)
			}
		}

		// if the field is an embedded struct, ignore the spss tag
		if field.Anonymous {
			continue
		}

		fieldInfo := FieldInfo{IndexChain: indexChain}
		fieldTag := field.Tag.Get("spss")
		fieldTags := strings.Split(fieldTag, TagSeparator)
		filteredTags := []string{}
		for _, fieldTagEntry := range fieldTags {
			if fieldTagEntry != "omitempty" {
				filteredTags = append(filteredTags, fieldTagEntry)
			} else {
				fieldInfo.OmitEmpty = true
			}
		}

		if len(filteredTags) == 1 && filteredTags[0] == "-" {
			continue
		} else if len(filteredTags) > 0 && filteredTags[0] != "" {
			fieldInfo.Keys = filteredTags
		} else {
			fieldInfo.Keys = []string{field.Name}
		}
		fieldInfo.FieldType = field.Type.Kind() // for writing
		fieldsList = append(fieldsList, fieldInfo)
	}
	return fieldsList
}

func GetConcreteContainerInnerType(in reflect.Type) (inInnerWasPointer bool, inInnerType reflect.Type) {
	inInnerType = in.Elem()
	inInnerWasPointer = false
	if inInnerType.Kind() == reflect.Ptr {
		inInnerWasPointer = true
		inInnerType = inInnerType.Elem()
	}
	return inInnerWasPointer, inInnerType
}

func GetConcreteReflectValueAndType(in interface{}) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(in)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value, value.Type()
}
