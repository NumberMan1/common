package summer

// GetSDBMHash SDBM 这个算法在开源的SDBM中使用，似乎对很多不同类型的数据都能得到不错的分布
func GetSDBMHash(text string) int32 {
	hash := int32(0)
	for i, l := 0, len(text); i < l; i += 1 {
		hash = int32(text[i]) + (hash << 6) + (hash << 16) - hash
	}
	//这里需要注意的是，如果hash值为负数，那么hash & 0x7FFFFFFF将会得到一个正数
	return hash & 0x7FFFFFFF
}

// GetBKDRHash 这个算法来自Brian Kernighan 和 Dennis Ritchie的 The C Programming Language。
// 这是一个很简单的哈希算法,使用了一系列奇怪的数字,形式如31,3131,31...31,
// 看上去和DJB算法很相似
func GetBKDRHash(text string) int32 {
	seed := int32(131) // 31 131 1313 13131 131313 etc..
	hash := int32(0)
	for i, l := 0, len(text); i < l; i += 1 {
		hash = (hash * seed) + int32(text[i])
	}
	return hash & 0x7FFFFFFF
}

// GetDJBHash DJB 这个算法是Daniel J.Bernstein 教授发明的，是目前公布的最有效的哈希函数。
func GetDJBHash(text string) int32 {
	hash := int32(0)
	for i, l := 0, len(text); i < l; i += 1 {
		hash += (hash << 5) + int32(text[i])
	}
	return hash & 0x7FFFFFFF
}
