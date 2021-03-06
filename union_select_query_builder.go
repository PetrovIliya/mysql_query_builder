package mysql_query_builder

import (
	"errors"
)

type unionSelectQueryBuilder struct {
	joinQueryBuilder
	whereQueryBuilder
	groupByQueryBuilder
	table     string
	alias     string
	selectStr string
	unions    []*unionSelectQueryBuilder
}

func (qb *unionSelectQueryBuilder) GetSql() (string, error) {
	if err := qb.validate(); err != nil {
		return "", err
	}

	unionPart, err := qb.getUnionPart()
	if err != nil {
		return "", err
	}

	sql := escape(qb.getSelectPart() + " " +
		qb.getFromPart() + " " +
		qb.getWherePart() + " " +
		qb.getGroupByPart() + " " +
		unionPart + " ",
	)
	return sql, nil
}

func (qb *unionSelectQueryBuilder) Union(union *unionSelectQueryBuilder) {
	qb.unions = append(qb.unions, union)
}

func (qb *unionSelectQueryBuilder) getSelectPart() string {
	if qb.selectStr == "" {
		qb.selectStr = "*"
	}
	return "SELECT " + qb.selectStr
}

func (qb *unionSelectQueryBuilder) getUnionPart() (string, error) {
	unionPart := ""
	for i := 0; i < len(qb.unions); i++ {
		sql, err := qb.unions[i].GetSql()
		if err != nil {
			return "", err
		}
		unionPart += " UNION " + sql
	}

	return unionPart, nil
}

func (qb *unionSelectQueryBuilder) getFromPart() string {
	return "FROM" + " `" + qb.table + "` " + qb.alias + " " + qb.getJoinsPart()
}

func (qb *unionSelectQueryBuilder) validate() error {
	err := qb.whereQueryBuilder.validate()
	if err != nil {
		return err
	}

	err = qb.joinQueryBuilder.validate()
	if err != nil {
		return err
	}

	for i := 0; i < len(qb.unions); i++ {
		if qb.unions[i].selectStr != qb.selectStr {
			return errors.New("different union select strings")
		}
	}

	if qb.table == "" {
		return errors.New("'table' param can not be empty")
	}

	return nil
}
