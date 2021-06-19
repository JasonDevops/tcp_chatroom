package main

// 存放公共方法

// 判断字符串是否存在切片里
func HasStrContain(str string, array []string) bool {
	for _, v := range array {
		if str == v {
			return true
		}
	}
	return false
}
