package page

import "errors"

func calcTotalPages(pageSize, totalRecords int64) (totalPages int64, err error) {
	if pageSize <= 0 {
		return 0, errors.New("page size should be positive integer")
	}

	if totalRecords < 0 {
		return 0, errors.New("total records should not be negative integer")
	}

	if totalRecords == 0 {
		return 0, nil
	}

	if totalRecords%pageSize == 0 {
		totalPages = totalRecords / pageSize
	} else {
		totalPages = totalRecords/pageSize + 1
	}
	return totalPages, nil
}
