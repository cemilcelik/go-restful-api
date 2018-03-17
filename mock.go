package main

type RepositoryMock struct{}

var UserSvcMck = &UserService{&RepositoryMock{}}

func (m *RepositoryMock) findByID(id string) (User, error) {
	return UserSvcMck.repo.findByID(id)
}

func (m *RepositoryMock) save(name, surname, email string) (err error) {
	return UserSvcMck.repo.save(name, surname, email)
}

func (m *RepositoryMock) update(id, name, surname, email string) (err error) {
	return UserSvcMck.repo.update(id, name, surname, email)
}

func (m *RepositoryMock) delete(id string) (err error) {
	return UserSvcMck.repo.delete(id)
}

func (m *RepositoryMock) getAll() ([]User, error) {
	return UserSvcMck.repo.getAll()
}
