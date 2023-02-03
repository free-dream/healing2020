package tools

//返回对应序号头像的url
func GetAvatarUrl(n int) string {
	switch n {
	case 0:
		return "http://cdn.healing2020.100steps.top/static/personal/avatarFemale.png"
	case 1:
		return "http://cdn.healing2020.100steps.top/static/personal/avatarMale.png"
	case 2:
		return "http://cdn.healing2020.100steps.top/static/personal/avatarSecret.png"
	default:
		return "error"
	}
}
