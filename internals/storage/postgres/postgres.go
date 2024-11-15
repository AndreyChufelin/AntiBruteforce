package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db       *sqlx.DB
	user     string
	password string
	name     string
	host     string
	port     string
}

func New(user, password, name, host, port string) *Storage {
	return &Storage{
		user:     user,
		password: password,
		name:     name,
		host:     host,
		port:     port,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
			s.user,
			s.name,
			s.password,
			s.host,
			s.port,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("no connection to close")
	}

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.db.Close()
		s.db = nil
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("falied to close postgres: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("postgres close operation canceled: %w", ctx.Err())
	}
}

func (s *Storage) WhitelistAdd(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO whitelist (subnet) VALUES ($1)", subnet)
	if err != nil {
		return fmt.Errorf("failed execute insert whitelist query: %w", err)
	}

	return nil
}

func (s *Storage) WhitelistDelete(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM whitelist WHERE subnet=$1", subnet)
	if err != nil {
		return fmt.Errorf("failed execute delete whitelist query: %w", err)
	}

	return nil
}

func (s *Storage) BlacklistAdd(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO blacklist (subnet) VALUES ($1)", subnet)
	if err != nil {
		return fmt.Errorf("failed execute insert blacklist query: %w", err)
	}

	return nil
}

func (s *Storage) BlacklistDelete(ctx context.Context, subnet string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM blacklist WHERE subnet=$1", subnet)
	if err != nil {
		return fmt.Errorf("failed execute delete blacklist query: %w", err)
	}

	return nil
}

func (s *Storage) WhitelistCheckSubnet(ctx context.Context, ip string) (bool, error) {
	var exists bool
	err := s.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM whitelist WHERE subnet >> $1::inet)", ip)
	if err != nil {
		return false, fmt.Errorf("failed to query check exist subnet: %w", err)
	}

	return exists, err
}

func (s *Storage) BlacklistCheckSubnet(ctx context.Context, ip string) (bool, error) {
	var exists bool
	err := s.db.GetContext(ctx, &exists, "SELECT EXISTS (SELECT 1 FROM blacklist WHERE subnet >> $1::inet)", ip)
	if err != nil {
		return false, fmt.Errorf("failed to query check exist subnet: %w", err)
	}

	return exists, err
}
