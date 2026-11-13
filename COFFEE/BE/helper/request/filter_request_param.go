package request_helper

import "strings"

// BuildOrderClause returns a safe ORDER BY clause.
// allowedFields is a whitelist map of frontend key -> SQL column expression.
// Falls back to defaultOrder if sortBy is empty or not in whitelist.
func BuildOrderClause(sortBy, sortOrder string, allowedFields map[string]string, defaultOrder string) string {
	col, ok := allowedFields[sortBy]
	if !ok || col == "" {
		return defaultOrder
	}

	order := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		order = "DESC"
	}

	return " ORDER BY " + col + " " + order
}
