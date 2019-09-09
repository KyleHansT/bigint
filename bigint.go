package bigint

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

// ExampleBigInt 例子
func ExampleBigInt() {
	// 字符单位换算表
	units := map[string]string{
		"K": "1000",
		"M": "1000K",
		"B": "1000M",
		"T": "1000B",
		"q": "1000T",
		"Q": "1000q",
		"s": "1000Q",
		"S": "1000s",
		"O": "1000S",
		"N": "1000O",
		"d": "1000N",
		"U": "1000d",
		"D": "1000U",
	}
	InitBigIntUnit(units, 10, 5)

	// InitDefaultUnits()

	a := NewBigInt("10000000")
	b := NewBigInt("11300000000")
	b.Mul(b, a)
	fmt.Printf("%s is %s \n", b.String(), FormatUnit(b))

	strNum := FormatUnit(b)
	fmt.Println(NewBigInt(strNum))
	fmt.Println(NewBigInt("11300034500000"))

	fmt.Println(NewBigInt("1130s"))
	// fmt.Println(NewBigInt("1130;2"))

}

/***********************************************
**	大数字单位换算
***********************************************/
// 数值单位配置
type bigIntConfig struct {
	base      int               // 进制数 10
	signDigit int               // 有效数字significant digit. eg. 5位有效数字
	flag      bool              // 量级标识 eg. true: k=1000, false: k=1024
	step      int               // 一个单位梯度的数量级数 eg. k=1000为3个数量级
	unitType  int               // 单位类型 0: 数字单位(默认) 如 100,2 ; 1: 字符单位,如 100K
	units     map[string]string // 单位换算表 eg. k:1000, m:1000000
	len2Unit  map[int]string    // 数值位数与单位对照表
}

var uc *bigIntConfig

// NewBigInt 可带可不带单位的数值字符串
func NewBigInt(strNum string) *big.Int {
	num, unit := splitUnit(strNum) // 数字与单位拆分
	return NewBigIntWithUnit(num, unit)
}

// NewBigIntWithUnit 不带单位的数字和单位作为入参
func NewBigIntWithUnit(num, unit string) *big.Int {
	b := &big.Int{}

	if len(unit) == 0 { // 无单位
		b.SetString(num, uc.base)
	} else {
		magnitude := ParseUnit(unit)
		if uc.flag {
			// 量级单位换算 eg. k=1000
			b.SetString(num+magnitude[1:], uc.base)
		} else {
			// 非量级单位换算 eg. k=1024
			m := &big.Int{}
			m.SetString(magnitude, uc.base)
			b.SetString(num, uc.base)
			b.Mul(b, m)
		}
	}
	return b
}

// InitDefaultUnits 初始化默认单位量级配置
func InitDefaultUnits() {
	// 数字单位换算表
	units := map[string]string{
		"1":  "10",
		"2":  "10;1",
		"3":  "10;2",
		"4":  "10;3",
		"5":  "10;4",
		"6":  "10;5",
		"7":  "10;6",
		"8":  "10;7",
		"9":  "10;8",
		"10": "10;9",
		"11": "10;10",
		"12": "10;11",
		"13": "10;12",
		"14": "10;13",
		"15": "10;14",
		"16": "10;15",
		"17": "10;16",
		"18": "10;17",
		"19": "10;18",
		"20": "10;19",
		"21": "10;20",
		"22": "10;21",
		"23": "10;22",
		"24": "10;23",
		"25": "10;24",
		"26": "10;25",
		"27": "10;26",
		"28": "10;27",
		"29": "10;28",
		"30": "10;29",
		"31": "10;30",
		"32": "10;31",
		"33": "10;32",
		"34": "10;33",
		"35": "10;34",
		"36": "10;35",
		"37": "10;36",
		"38": "10;37",
		"39": "10;38",
		"40": "10;39",
		"41": "10;40",
		"42": "10;41",
		"43": "10;42",
		"44": "10;43",
		"45": "10;44",
		"46": "10;45",
		"47": "10;46",
		"48": "10;47",
		"49": "10;48",
		"50": "10;49",
		"51": "10;50",
		"52": "10;51",
		"53": "10;52",
		"54": "10;53",
		"55": "10;54",
		"56": "10;55",
		"57": "10;56",
		"58": "10;57",
		"59": "10;58",
		"60": "10;59",
	}

	InitBigIntUnit(units, 10, 5)
}

// InitBigIntUnit 加载单位量级配置
//	加载单位量级配置
//	units:		单位换算表
//	base:		进制数 eg. base=10 -> 十进制
//	signDigit: 	有效数字位数 eg. signDigit=5 -> 21200k
func InitBigIntUnit(unitsTable map[string]string, base, signDigit int) {
	uc = &bigIntConfig{
		base:      base,
		signDigit: signDigit,
		units:     unitsTable,
		len2Unit:  map[int]string{},
	}

	// 解析单位类型
	for unit := range unitsTable {
		lu := len(unit)
		for i := 0; i < lu; i++ {
			if unit[i] != ';' && (unit[i] > '9' || unit[i] < '0') {
				uc.unitType = 1
			}
		}
	}

	lenOfDigital := 999
	units := map[string]string{}
	// 遍历单位换算表，将值转为无单位的数值
	for key, value := range unitsTable {
		_, last := splitUnit(value) // 取数值单位
		if len(last) != 0 {         // 带单位 eg. 1000k
			value = parseR(value)
		}

		// 清理单位常数
		u := trimDigit(key)
		units[u] = value

		// 量级单位标识
		if len(value) < lenOfDigital {
			lenOfDigital = len(value)
			num, _ := splitUnit(value)
			pow := math.Pow(float64(base), float64(lenOfDigital-1))
			strPow := fmt.Sprintf("%d", int64(pow))
			if num == strPow {
				uc.flag = true
			}
		}
	}

	//替换为无单位的转换表
	uc.units = units
	uc.step = lenOfDigital - 1

	maxLenOfNum := len(uc.units) * uc.step
	for i := 1; i <= maxLenOfNum; i++ {
		unit := ""
		icursor := math.Ceil(float64(i) / float64(uc.step))
		for k, v := range uc.units {
			lcursor := len(v) / uc.step
			if uc.step == 1 { //特殊处理以一个量级为一个单位的情况
				lcursor--
			}
			if lcursor == int(icursor) {
				unit = k
			}
		}

		index := i + uc.signDigit
		uc.len2Unit[index] = unit
	}
}

// ParseUnit 解析单位、转为纯数值
func ParseUnit(u string) string {
	if v, ok := uc.units[u]; ok {
		return v
	} else {
		return "1"
	}
}

// FormatUnit 将大数字格式化为带单位的数值
// 注意: 不进行四舍五入
func FormatUnit(bigInt *big.Int) string {
	if uc.unitType == 1 { // 字符单位
		return format2CharUnit(bigInt)
	}

	// 默认数字单位
	return format2DigitalUnit(bigInt)
}

func format2CharUnit(bigInt *big.Int) string {
	num := bigInt.String()
	lenOfNum := len(num)
	unit, ok := uc.len2Unit[lenOfNum]
	if !ok {
		return num
	}

	digit, ok := uc.units[unit]
	if !ok {
		return num
	}

	if uc.flag {
		lenOfUnit := len(digit)
		numWithUnit := num[0 : lenOfNum-lenOfUnit+1]
		numWithUnit += unit
		return numWithUnit
	} else {
		d := NewBigInt(digit)
		bigInt.Div(bigInt, d)
		strNum := bigInt.String()
		return strNum + unit
	}
}

func format2DigitalUnit(bigInt *big.Int) string {
	num := bigInt.String()
	lenOfNum := len(num)
	unit, ok := uc.len2Unit[lenOfNum]
	if !ok {
		return num
	}

	digit, ok := uc.units[unit]
	if !ok {
		return num
	}

	if uc.flag {
		lenOfUnit := len(digit)
		numWithUnit := num[0 : lenOfNum-lenOfUnit+1]
		numWithUnit += ";" + unit
		return numWithUnit
	} else {
		d := NewBigInt(digit)
		bigInt.Div(bigInt, d)
		strNum := bigInt.String()
		return strNum + ";" + unit
	}
}

// 递归解析带单位的数值、转为无单位数字
func parseR(strNum string) string {
	num, last := splitUnit(strNum) // 数字与单位拆分
	if len(last) == 0 {
		return strNum
	}
	if len(num) == 0 {
		num = "1"
	}

	v := ParseUnit(last)
	magnitude := parseR(v)

	b := &big.Int{}
	m := &big.Int{}
	b.SetString(num, uc.base)
	m.SetString(magnitude, uc.base)
	b.Mul(b, m)
	return b.String()

}

func parse(strNum string) *big.Int {
	b := &big.Int{}

	num, unit := splitUnit(strNum) // 数字与单位拆分
	if len(unit) == 0 {            // 无单位
		b.SetString(strNum, uc.base)
	} else {
		magnitude := ParseUnit(unit) // 带单位 eg. 1k
		if uc.flag {
			// 量级单位换算 eg. k=1000
			b.SetString(num+magnitude[1:], uc.base)
		} else {
			// 非量级单位换算 eg. k=1024
			m := &big.Int{}
			m.SetString(magnitude, uc.base)
			b.SetString(num, uc.base)
			b.Mul(b, m)
		}
	}

	return b
}

// 拆分数值和单位
func splitUnit(strNum string) (num, unit string) {
	if uc.unitType == 1 { // 字符单位
		return splitWithCharUnit(strNum)
	}

	// 默认数字单位
	return splitWithDigitalUnit(strNum)
}

// 拆分数值和单位 eg. 1000K => 1000 K
func splitWithCharUnit(strNum string) (num, unit string) {
	l := len(strNum)
	for i := l; i > 0; i-- {
		if strNum[i-1] <= '9' && strNum[i-1] >= '0' {
			num = strNum[:i]
			unit = strNum[i:]
			return
		}
	}
	return "", strNum
}

// 拆分数值和单位 eg. 1000,1 => 1000 1
func splitWithDigitalUnit(strNum string) (num, unit string) {
	s := strings.Split(strNum, ";")
	if len(s) == 2 {
		return s[0], s[1]
	}
	return strNum, ""
}

// 清理单位常数
func trimDigit(unit string) string {
	if uc.unitType == 1 { // 字符单位
		return trimDigitWithCharUnit(unit)
	}

	// 默认数字单位
	return trimDigitWithDigitalUnit(unit)
}

// 清理单位常数 eg. 1000K => K or K => K
func trimDigitWithCharUnit(unit string) string {
	_, u := splitWithCharUnit(unit)
	return u
}

// 清理单位常数 eg. 1000,1 => 1 or 1 => 1
func trimDigitWithDigitalUnit(unit string) string {
	s := strings.Split(unit, ";")
	if len(s) == 2 {
		return s[1]
	}
	return unit
}

// 当单位是10的n次方时，传输时使用标准的e格式
func outputSep2E(strNum string) string {
	newStr := strings.Replace(strNum, ";", "e", -1)
	return newStr
}

func inputE2Sep(strNum string) string {
	newStr := strings.Replace(strNum, "e", ";", -1)
	return newStr
}
