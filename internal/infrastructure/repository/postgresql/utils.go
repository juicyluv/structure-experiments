package postgresql

import "fmt"

func paginationQuery(page, objectsPerPage int) string {
	if page == 0 || objectsPerPage == 0 {
		return ""
	}

	return fmt.Sprintf("LIMIT %d OFFSET %d", objectsPerPage, (page-1)*objectsPerPage)
}
