package internal

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lib/pq"

	corerand "github.com/authgear/authgear-server/pkg/util/rand"
	"github.com/authgear/authgear-server/pkg/util/uuid"
)

type CreateDefaultDomainOptions struct {
	DatabaseURL         string
	DatabaseSchema      string
	DefaultDomainSuffix string
}

func CreateDefaultDomain(opts CreateDefaultDomainOptions) (err error) {
	db := openDB(opts.DatabaseURL, opts.DatabaseSchema)
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
	}()

	allConfigSourceList, err := selectConfigSources(ctx, tx, nil)
	if err != nil {
		return
	}

	for _, configSource := range allConfigSourceList {
		appID := configSource.AppID
		domain := makeDefaultDomain(appID, opts.DefaultDomainSuffix)

		var exists bool
		exists, err = checkDefaultDomainExistance(ctx, tx, appID, domain)
		if err != nil {
			return
		}

		if !exists {
			err = createDefaultDomain(ctx, tx, appID, domain)
			if err != nil {
				return
			}
			fmt.Printf("created: %s\n", domain)
		} else {
			fmt.Printf("skipped: %s\n", domain)
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func makeVerificationNonce() string {
	nonce := make([]byte, 16)
	corerand.SecureRand.Read(nonce)
	verificationNonce := hex.EncodeToString(nonce)
	return verificationNonce
}

func makeDefaultDomain(appID string, suffix string) string {
	return appID + suffix
}

func checkDefaultDomainExistance(ctx context.Context, tx *sql.Tx, appID string, domain string) (exists bool, err error) {
	builder := newSQLBuilder().Select("id").
		From(pq.QuoteIdentifier("_portal_domain")).
		Where("app_id = ? AND domain = ?", appID, domain)

	q, args, err := builder.ToSql()
	if err != nil {
		return
	}

	rows, err := tx.QueryContext(ctx, q, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		exists = true
	}
	err = rows.Close()
	if err != nil {
		return
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func createDefaultDomain(ctx context.Context, tx *sql.Tx, appID string, domain string) error {
	isCustom := false
	// The apex domain of default domain is itself.
	apexDomain := domain
	verificationNonce := makeVerificationNonce()

	builder := newSQLBuilder().
		Insert(pq.QuoteIdentifier("_portal_domain")).
		Columns(
			"id", "app_id", "created_at", "domain", "apex_domain", "verification_nonce", "is_custom",
		).
		Values(
			uuid.New(),
			appID,
			time.Now().UTC(),
			domain,
			apexDomain,
			verificationNonce,
			isCustom,
		)

	q, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
