package constants

const (
	NoteTableName         = "note"
	UserTableName         = "user"
	SecretKey             = "secret key"
	IdentityKey           = "id"
	Total                 = "total"
	Notes                 = "notes"
	NoteID                = "note_id"
	ApiServiceName        = "demoapi"
	NoteServiceName       = "demonote"
	UserServiceName       = "demouser"
	MySQLDefaultDSN       = "root:Ling@tcp(localhost:3306)/mynote?charset=utf8&parseTime=True&loc=Local"
	EtcdAddress           = "127.0.0.1:2379"
	DefaultLimit          = 10
	Ttltime         int64 = 1000
	UserAddress           = "127.0.0.1:8889"
	NoteAddress           = "127.0.0.1:8888"
	Schema                = "Etcd"
)
