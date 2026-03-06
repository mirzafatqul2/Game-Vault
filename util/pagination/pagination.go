package pagination

func GetPagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	return limit, offset
}
