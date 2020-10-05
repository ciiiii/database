package database

func (backend DBBackend) Ready() bool {
	return backend.DB.DB().Ping() == nil
}
