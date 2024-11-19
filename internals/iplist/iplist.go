package iplist

import (
	"context"
	"fmt"
	"log/slog"
)

type IPList struct {
	logger  *slog.Logger
	storage Storage
}

type Storage interface {
	WhitelistAdd(ctx context.Context, subnet string) error
	WhitelistDelete(ctx context.Context, subnet string) error
	BlacklistAdd(ctx context.Context, subnet string) error
	BlacklistDelete(ctx context.Context, subnet string) error
	WhitelistCheckSubnet(ctx context.Context, ip string) (bool, error)
	BlacklistCheckSubnet(ctx context.Context, ip string) (bool, error)
}

func NewIPList(logger *slog.Logger, storage Storage) *IPList {
	return &IPList{
		logger:  logger,
		storage: storage,
	}
}

func (i *IPList) WhitelistAdd(ctx context.Context, subnet string) error {
	logg := i.logger.With("op", "WhitelistAdd")
	err := i.storage.WhitelistAdd(ctx, subnet)
	if err != nil {
		logg.Error("failed to add to whitelist", "err", err)
		return fmt.Errorf("failed to add to whitelist: %w", err)
	}

	return nil
}

func (i *IPList) WhitelistDelete(ctx context.Context, subnet string) error {
	logg := i.logger.With("op", "WhitelistDelete")
	err := i.storage.WhitelistDelete(ctx, subnet)
	if err != nil {
		logg.Error("failed to delete from whitelist", "err", err)
		return fmt.Errorf("failed to delete from whitelist: %w", err)
	}

	return nil
}

func (i *IPList) BlacklistAdd(ctx context.Context, subnet string) error {
	logg := i.logger.With("op", "BlacklistAdd")
	err := i.storage.BlacklistAdd(ctx, subnet)
	if err != nil {
		logg.Error("failed to add to blacklist", "err", err)
		return fmt.Errorf("failed to add to blacklist: %w", err)
	}

	return nil
}

func (i *IPList) BlacklistDelete(ctx context.Context, subnet string) error {
	logg := i.logger.With("op", "BlacklistDelete")
	err := i.storage.BlacklistDelete(ctx, subnet)
	if err != nil {
		logg.Error("failed to delete from blacklist", "err", err)
		return fmt.Errorf("failed to delete from blacklist: %w", err)
	}

	return nil
}

func (i *IPList) WhitelistCheckSubnet(ctx context.Context, subnet string) (bool, error) {
	logg := i.logger.With("op", "WhitelistCheckIP")
	exist, err := i.storage.WhitelistCheckSubnet(ctx, subnet)
	if err != nil {
		logg.Error("failed to check whitelist", "err", err)
		return false, fmt.Errorf("failed to check whitelist: %w", err)
	}

	return exist, nil
}

func (i *IPList) BlacklistCheckSubnet(ctx context.Context, subnet string) (bool, error) {
	logg := i.logger.With("op", "BlacklistCheckIP")
	exist, err := i.storage.BlacklistCheckSubnet(ctx, subnet)
	if err != nil {
		logg.Error("failed to check blacklist", "err", err)
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	return exist, nil
}
