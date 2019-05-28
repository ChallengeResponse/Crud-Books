package dbConfig
//TODO read from a txt or conf file instead of requiring compilation / function call

type DbInfo struct{
	Driver, HostAndPort, DbName, User, Password string
}

func Main() DbInfo{
	return DbInfo{"mysql","127.0.0.1:3306","go","root","hellogo"}
}
