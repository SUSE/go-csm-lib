package provisioner

type RabbitmqProvisionerInterface interface {
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	CreateUser(string, string, string) (map[string]string, error)
	DeleteUser(string, string) error
	UserExists(string, string) (bool, error)
}
