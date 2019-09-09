# bigint
自定义带单位数字与big.Int之间的转换


```
func exampleBigInt() {
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
	InitUnits(units, 10, 5)

	a := NewBigInt("10000000")
	b := NewBigInt("11300000000")
	b.Mul(b, a)
	fmt.Printf("%s is %s \n", b.String(), FormatUnit(b))

	strNum := FormatUnit(b)
	fmt.Println(NewBigInt(strNum))
	fmt.Println(NewBigInt("11300034500000"))

	fmt.Println(NewBigInt("1130s"))

}
```
