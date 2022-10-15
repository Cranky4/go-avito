package sqlstorage

type DatabaseConf struct {
	Dsn, ConnectionTryDelay string
	MaxConnectionTries      int
}
