package utils

import "github.com/spf13/viper"

var Origins []string

//为何需要中间变量origin?

//func InitOrigin() {
//	var origin []string
//	viper.UnmarshalKey("server.origin", &origin)
//	for _, v := range origin {
//		Origins = append(Origins, v)
//	}
//}

func InitOrigin() {
	viper.UnmarshalKey("server.origin", &Origins)
}
