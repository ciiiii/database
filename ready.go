package database

func (backend DBBackend) Ready() bool {
	if backend.DB.DB().Ping() != nil {
		return false
	}
	return true
}
