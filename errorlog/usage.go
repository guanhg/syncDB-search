package errorlog

func CheckErr(err error, msg ...string)  {
	if err!=nil{
		errMsg := ""
		for _, m:= range msg {
			errMsg += m
		}
		panic(errMsg+"\n"+err.Error())
	}
}

