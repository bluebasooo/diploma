package config

type MongoConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	DatabaseName     string
	CollectionsNames []string
}

func (mongo *MongoConfig) Uri() string {
	return mongo.Host + ":" + mongo.Port
}
