package db

func (db *DB) CreateSubject(s *Subject) error {
	var count int64
	db.Engine.Model(&Subject{}).Where(s).Count(&count)
	if count != 0 {
		return AlreadyExistsError
	}

	return db.Engine.Create(s).Error
}
