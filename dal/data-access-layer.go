package dal

import "github.com/clockworksoul/cog2/data/rest"

// DataAccess represents a common DataAccessObject, backed either by a
// database or an in-memory datastore.
type DataAccess interface {
	Initialize() error
	GroupAddUser(string, string) error
	GroupCreate(rest.Group) error
	GroupDelete(string) error
	GroupExists(string) (bool, error)
	GroupGet(string) (rest.Group, error)
	GroupGrantRole() error
	GroupList() ([]rest.Group, error)
	GroupRemoveUser(string, string) error
	GroupRevokeRole() error
	GroupUpdate(rest.Group) error
	UserCreate(rest.User) error
	UserDelete(string) error
	UserExists(string) (bool, error)
	UserGet(string) (rest.User, error)
	UserList() ([]rest.User, error)
	UserUpdate(rest.User) error
}
