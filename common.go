package zutils

import (
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"reflect"
	"strconv"
	"time"
)

// 判断string是否在slice内,使用map性能高
func IsInSliceByMap(s string, sl []string) bool {
	return isInMap(s, conVertSlice2Map(sl))
}

// 判断map里是否存在
func isInMap(s string, m map[string]struct{}) bool {
	_, ok := m[s]
	return ok
}

// 转换slice为map
func conVertSlice2Map(sl []string) map[string]struct{} {
	m := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		m[v] = struct{}{}
	}

	return m
}

// 判断string是否在slice内
func IsInSlice(t string, tr []string) bool {
	for _, v := range tr {
		if t == v {
			return true
		}
	}
	return false
}

// 判断是否uft8
func IsUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

// 判断是否中文编码
func IsGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// float64取指定位数
func DecimalTwo(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}

// struct2map
func StructToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	var out = make(map[string]interface{})
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// 字符串指定编码转换,如GBK转utf8
func ConvertStringByCode(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// DecideType 类型断言..
func DecideType(src interface{}) (string, error) {
	tmp := ""
	switch t := src.(type) {
	case nil:
		tmp = "null"
	case bool:
		if t {
			tmp = "True"
		} else {
			tmp = "False"
		}
	case []byte:
		tmp = string(t)
	case time.Time:
		tmp = t.Format("2006-01-02 15:04:05.999")
	case int:
		tmp = strconv.Itoa(t)
	case int32:
		tmp = strconv.Itoa(int(t))
	case int64:
		tmp = strconv.FormatInt(t, 10)
	case float32:
		tmp = strconv.FormatFloat(float64(t), 'f', -1, 32)
	case float64:
		tmp = strconv.FormatFloat(t, 'f', -1, 64)
	case string:
		tmp = t
	default:
		err := errors.New("no this type")
		return tmp, err
	}
	return tmp, nil
}
