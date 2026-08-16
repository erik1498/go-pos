package utils

import (
	"go-pos/internal/domain"
	"strings"
)

var operatorMap = map[string]string{
	"eq":   "=",
	"neq":  "!=",
	"gt":   ">",
	"gte":  ">=",
	"lt":   "<",
	"lte":  "<=",
	"like": "ILIKE",
}

func SanitizeQuery(opts domain.QueryOptions, allowedFields, allowedSorts map[string]bool, defaultSort string) domain.QueryOptions {
	var validSorts []string

	if opts.Sort != "" {
		sorts := strings.Split(opts.Sort, ",")
		for _, s := range sorts {
			s = strings.TrimSpace(s)
			parts := strings.Split(s, " ")
			if len(parts) > 0 {
				field := parts[0]
				dir := "asc"
				if len(parts) == 2 && strings.ToLower(parts[1]) == "desc" {
					dir = "desc"
				}

				if allowedSorts[field] {
					validSorts = append(validSorts, field+" "+dir)
				}
			}
		}
	}

	opts.Sort = strings.Join(validSorts, ", ")
	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	var validFilter []domain.Filter
	for _, f := range opts.Filters {
		if !allowedFields[f.Field] {
			continue
		}

		sqlOperator, ok := operatorMap[f.Operator]
		if !ok {
			continue
		}

		if sqlOperator == "ILIKE" {
			f.Value = "%" + f.Value.(string) + "%"
		}

		validFilter = append(validFilter, domain.Filter{
			Field:    f.Field,
			Operator: sqlOperator,
			Value:    f.Value,
		})
	}

	opts.Filters = validFilter

	return opts
}
