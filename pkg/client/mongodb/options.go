package mongodb

type Option func(*Mongo)

func WithHost(host string) Option {
	return func(postgres *Mongo) {
		postgres.host = host
	}
}

func WithPort(port string) Option {
	return func(postgres *Mongo) {
		postgres.port = port
	}
}

func WithUsername(username string) Option {
	return func(postgres *Mongo) {
		postgres.username = username
	}
}

func WithPassword(password string) Option {
	return func(postgres *Mongo) {
		postgres.password = password
	}
}

func WithDBName(dbName string) Option {
	return func(postgres *Mongo) {
		postgres.dbName = dbName
	}
}
