package extension

type Extension interface {
	CreateConnection() error
	CreateWorkspace() error
	DeleteConnection() error
	DeleteWorkspace() error
	GetConnection() error
	GetWorkspace() error
}
