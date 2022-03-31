package errorlog

import (
	"errors"
)

var (
	ParameterError = errors.New(" [ParameterError] ")
	CanalSqlTypeError = errors.New(" [CanalSqlTypeError] ")
)

func CheckErr(err error, msg ...string)  {
	if err!=nil{
		errMsg := ""
		for _, m:= range msg {
			errMsg += m
		}
		panic(err.Error()+errMsg+"\n")
	}
}

