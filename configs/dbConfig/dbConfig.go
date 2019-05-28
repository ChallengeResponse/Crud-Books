package dbConfig
//TODO read from a txt or conf file instead of requiring compilation / function call

type DbInfo struct{
	Driver, HostAndPort, DbName, User, Password string
}

func Main(host string) DbInfo{
	return DbInfo{"mysql",host + ":3306","go","root","hellogo"}
}
