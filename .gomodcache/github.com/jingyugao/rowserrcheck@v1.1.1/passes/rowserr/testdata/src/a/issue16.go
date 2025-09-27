package a

func issue16() error {
	rows, err := db.Query("select 1")
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
	}
	return rows.Err()
}
