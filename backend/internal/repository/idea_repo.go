package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/nisha/gravehard/internal/domain"
)

type IdeaRepository struct {
	db DBTX
}

func NewIdeaRepository(db DBTX) *IdeaRepository {
	return &IdeaRepository{db: db}
}

const ideaColumns = `
	id, user_id, name, tagline, description, category::text,
	born_at, died_at, primary_cause::text, secondary_causes::text[],
	epitaph, autopsy, validation_level::text, money_spent_cents,
	hours_invested, status::text, would_resurrect, created_at, updated_at`

func scanIdea(row pgx.Row) (*domain.Idea, error) {
	var idea domain.Idea
	var category, primaryCause, validationLevel, status string
	var secondary []string

	err := row.Scan(
		&idea.ID, &idea.UserID, &idea.Name, &idea.Tagline, &idea.Description, &category,
		&idea.BornAt, &idea.DiedAt, &primaryCause, &secondary,
		&idea.Epitaph, &idea.Autopsy, &validationLevel, &idea.MoneySpentCents,
		&idea.HoursInvested, &status, &idea.WouldResurrect, &idea.CreatedAt, &idea.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	idea.Category = domain.IdeaCategory(category)
	idea.PrimaryCause = domain.DeathCause(primaryCause)
	idea.ValidationLevel = domain.ValidationLevel(validationLevel)
	idea.Status = domain.IdeaStatus(status)
	idea.SecondaryCauses = make([]domain.DeathCause, len(secondary))
	for i, s := range secondary {
		idea.SecondaryCauses[i] = domain.DeathCause(s)
	}
	return &idea, nil
}

func (r *IdeaRepository) Create(ctx context.Context, idea *domain.Idea) (*domain.Idea, error) {
	q := fmt.Sprintf(`
		INSERT INTO ideas (
			user_id, name, tagline, description, category,
			born_at, died_at, primary_cause, secondary_causes,
			epitaph, autopsy, validation_level, money_spent_cents,
			hours_invested, status, would_resurrect
		) VALUES (
			$1, $2, $3, $4, $5::idea_category,
			$6, $7, $8::death_cause, $9::death_cause[],
			$10, $11, $12::validation_level, $13,
			$14, $15::idea_status, $16
		)
		RETURNING %s`, ideaColumns)

	secondary := make([]string, len(idea.SecondaryCauses))
	for i, c := range idea.SecondaryCauses {
		secondary[i] = string(c)
	}

	created, err := scanIdea(r.db.QueryRow(ctx, q,
		idea.UserID, idea.Name, idea.Tagline, idea.Description, string(idea.Category),
		idea.BornAt, idea.DiedAt, string(idea.PrimaryCause), secondary,
		idea.Epitaph, idea.Autopsy, string(idea.ValidationLevel), idea.MoneySpentCents,
		idea.HoursInvested, string(idea.Status), idea.WouldResurrect,
	))
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *IdeaRepository) GetByID(ctx context.Context, userID, ideaID uuid.UUID) (*domain.Idea, error) {
	q := fmt.Sprintf(`
		SELECT %s
		FROM ideas
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		LIMIT 1`, ideaColumns)

	idea, err := scanIdea(r.db.QueryRow(ctx, q, ideaID, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return idea, nil
}

func (r *IdeaRepository) List(ctx context.Context, userID uuid.UUID, f domain.IdeaFilter) ([]domain.Idea, int64, error) {
	where := []string{"user_id = $1", "deleted_at IS NULL"}
	args := []any{userID}
	argIdx := 2

	if f.Cause != nil {
		where = append(where, fmt.Sprintf("primary_cause = $%d::death_cause", argIdx))
		args = append(args, string(*f.Cause))
		argIdx++
	}
	if f.Category != nil {
		where = append(where, fmt.Sprintf("category = $%d::idea_category", argIdx))
		args = append(args, string(*f.Category))
		argIdx++
	}
	if f.Status != nil {
		where = append(where, fmt.Sprintf("status = $%d::idea_status", argIdx))
		args = append(args, string(*f.Status))
		argIdx++
	}
	if f.Validation != nil {
		where = append(where, fmt.Sprintf("validation_level = $%d::validation_level", argIdx))
		args = append(args, string(*f.Validation))
		argIdx++
	}
	if f.From != nil {
		where = append(where, fmt.Sprintf("died_at >= $%d", argIdx))
		args = append(args, *f.From)
		argIdx++
	}
	if f.To != nil {
		where = append(where, fmt.Sprintf("died_at <= $%d", argIdx))
		args = append(args, *f.To)
		argIdx++
	}
	if qStr := strings.TrimSpace(f.Query); qStr != "" {
		where = append(where, fmt.Sprintf("search_vector @@ websearch_to_tsquery('english', $%d)", argIdx))
		args = append(args, qStr)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM ideas WHERE %s", whereClause)
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderDir := "DESC"
	if strings.EqualFold(f.Order, "asc") {
		orderDir = "ASC"
	}

	var orderBy string
	switch f.Sort {
	case "created_at":
		orderBy = fmt.Sprintf("created_at %s", orderDir)
	case "money_spent":
		orderBy = fmt.Sprintf("money_spent_cents %s, died_at DESC", orderDir)
	case "lifespan":
		orderBy = fmt.Sprintf("(died_at - born_at) %s, died_at DESC", orderDir)
	default:
		orderBy = fmt.Sprintf("died_at %s, created_at DESC", orderDir)
	}

	offset := (f.Page - 1) * f.Limit
	listSQL := fmt.Sprintf(`
		SELECT %s
		FROM ideas
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, ideaColumns, whereClause, orderBy, argIdx, argIdx+1)
	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]domain.Idea, 0)
	for rows.Next() {
		idea, err := scanIdea(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *idea)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *IdeaRepository) Update(ctx context.Context, idea *domain.Idea) (*domain.Idea, error) {
	q := fmt.Sprintf(`
		UPDATE ideas
		SET
			name = $3,
			tagline = $4,
			description = $5,
			category = $6::idea_category,
			born_at = $7,
			died_at = $8,
			primary_cause = $9::death_cause,
			secondary_causes = $10::death_cause[],
			epitaph = $11,
			autopsy = $12,
			validation_level = $13::validation_level,
			money_spent_cents = $14,
			hours_invested = $15,
			status = $16::idea_status,
			would_resurrect = $17,
			updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		RETURNING %s`, ideaColumns)

	secondary := make([]string, len(idea.SecondaryCauses))
	for i, c := range idea.SecondaryCauses {
		secondary[i] = string(c)
	}

	updated, err := scanIdea(r.db.QueryRow(ctx, q,
		idea.ID, idea.UserID, idea.Name, idea.Tagline, idea.Description, string(idea.Category),
		idea.BornAt, idea.DiedAt, string(idea.PrimaryCause), secondary,
		idea.Epitaph, idea.Autopsy, string(idea.ValidationLevel), idea.MoneySpentCents,
		idea.HoursInvested, string(idea.Status), idea.WouldResurrect,
	))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *IdeaRepository) SoftDelete(ctx context.Context, userID, ideaID uuid.UUID) error {
	const q = `
		UPDATE ideas
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	tag, err := r.db.Exec(ctx, q, ideaID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *IdeaRepository) UpdateStatus(ctx context.Context, userID, ideaID uuid.UUID, status domain.IdeaStatus) (*domain.Idea, error) {
	q := fmt.Sprintf(`
		UPDATE ideas
		SET status = $3::idea_status, updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
		RETURNING %s`, ideaColumns)

	updated, err := scanIdea(r.db.QueryRow(ctx, q, ideaID, userID, string(status)))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *IdeaRepository) ListAllByUser(ctx context.Context, userID uuid.UUID) ([]domain.Idea, error) {
	q := fmt.Sprintf(`
		SELECT %s
		FROM ideas
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY died_at ASC, created_at ASC`, ideaColumns)

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Idea, 0)
	for rows.Next() {
		idea, err := scanIdea(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *idea)
	}
	return items, rows.Err()
}
