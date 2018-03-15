package main

type RepositoryMock struct{}

// type RepositoryMock struct {
// 	FindByIDFunc func(id string) (User, error)
// 	SaveFunc     func(name string, surname string, email string) error
// 	UpdateFunc   func(id string, name string, surname string, email string) error
// 	DeleteFunc   func(id string) error
// 	GetAllFunc   func() ([]User, error)
// }

func (m *RepositoryMock) findByID(id string) (User, error) {
	return userSvc.repo.findByID(id)
	// return m.FindByIDFunc(id)
}

func (m *RepositoryMock) save(name, surname, email string) (err error) {
	return userSvc.repo.save(name, surname, email)
	// return m.SaveFunc(name, surname, email)
}

func (m *RepositoryMock) update(id, name, surname, email string) (err error) {
	return userSvc.repo.update(id, name, surname, email)
	// return m.UpdateFunc(id, name, surname, email)
}

func (m *RepositoryMock) delete(id string) (err error) {
	return userSvc.repo.delete(id)
	// return m.DeleteFunc(id)
}

func (m *RepositoryMock) getAll() ([]User, error) {
	return userSvc.repo.getAll()
	// return m.GetAllFunc()
}
