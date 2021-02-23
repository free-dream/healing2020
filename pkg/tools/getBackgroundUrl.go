package tools

//返回对应序号背景的url
func GetBackgroundUrl(n int) string {
	switch n {
		case 1: return "url"	//B1
		case 2:	return "url"	//B2
		case 3: return "url"	//B3
		case 4: return "url"	//B4
		case 5: return "url"	//B5
		default: return "error"
	}
}