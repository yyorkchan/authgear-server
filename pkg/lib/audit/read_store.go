package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/authgear/authgear-server/pkg/api/model"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/auditdb"
	"github.com/authgear/authgear-server/pkg/util/graphqlutil"
)

type QueryPageOptions struct {
	RangeFrom     *time.Time
	RangeTo       *time.Time
	ActivityTypes []string
	UserIDs       []string
	SortDirection model.SortDirection
}

func (o QueryPageOptions) Apply(q db.SelectBuilder) db.SelectBuilder {
	if o.RangeFrom != nil {
		q = q.Where("created_at >= ?", o.RangeFrom)
	}

	if o.RangeTo != nil {
		q = q.Where("created_at < ?", o.RangeTo)
	}

	if len(o.ActivityTypes) > 0 {
		q = q.Where("activity_type = ANY (?)", pq.Array(o.ActivityTypes))
	}

	if len(o.UserIDs) > 0 {
		q = q.Where("user_id = ANY (?)", pq.Array(o.UserIDs))
	}

	return q
}

type ReadStore struct {
	SQLBuilder  *auditdb.SQLBuilderApp
	SQLExecutor *auditdb.ReadSQLExecutor
}

func (s *ReadStore) Count(opts QueryPageOptions) (uint64, error) {
	query := s.SQLBuilder.
		Select("count(*)").
		From(s.SQLBuilder.TableName("_audit_log"))

	query = opts.Apply(query)

	scanner, err := s.SQLExecutor.QueryRowWith(query)
	if err != nil {
		return 0, err
	}

	var count uint64
	err = scanner.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ReadStore) QueryPage(opts QueryPageOptions, pageArgs graphqlutil.PageArgs) ([]*Log, uint64, error) {
	query := s.selectQuery()

	query = opts.Apply(query)

	sortDirection := opts.SortDirection
	if sortDirection == model.SortDirectionDefault {
		sortDirection = model.SortDirectionDesc
	}

	query = query.OrderBy(fmt.Sprintf("created_at %s", sortDirection))

	query, offset, err := db.ApplyPageArgs(query, pageArgs)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.SQLExecutor.QueryWith(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*Log
	for rows.Next() {
		l, err := s.scan(rows)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}

	return logs, offset, nil
}

func (s *ReadStore) GetByIDs(ids []string) ([]*Log, error) {
	query := s.selectQuery().Where("id = ANY (?)", pq.Array(ids))

	rows, err := s.SQLExecutor.QueryWith(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*Log
	for rows.Next() {
		l, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}

func (s *ReadStore) selectQuery() db.SelectBuilder {
	return s.SQLBuilder.
		Select(
			"id",
			"created_at",
			"user_id",
			"activity_type",
			"ip_address",
			"user_agent",
			"client_id",
			"data",
		).
		From(s.SQLBuilder.TableName("_audit_log"))
}

func (s *ReadStore) scan(scn db.Scanner) (*Log, error) {
	l := &Log{}

	var data []byte

	err := scn.Scan(
		&l.ID,
		&l.CreatedAt,
		&l.UserID,
		&l.ActivityType,
		&l.IPAddress,
		&l.UserAgent,
		&l.ClientID,
		&data,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &l.Data)
	if err != nil {
		return nil, err
	}

	return l, nil
}
