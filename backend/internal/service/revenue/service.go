package revenue

import (
	"context"
	"errors"
	"strings"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service encapsulates wallet accounting rules.
type Service struct {
	repo *repository.RevenueRepository
}

func NewService(repo *repository.RevenueRepository) *Service {
	return &Service{repo: repo}
}

// RecordEvent stores a revenue event and updates the creator wallet.
func (s *Service) RecordEvent(ctx context.Context, creatorID, userID, roleID, eventType string, amount int64) (*model.RevenueEvent, *model.CreatorWallet, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be > 0")
	}
	rules, _ := s.repo.ListRules(ctx)
	amountToCredit := applyRules(eventType, amount, rules)
	event := &model.RevenueEvent{CreatorID: creatorID, UserID: userID, RoleID: roleID, EventType: eventType, Amount: amountToCredit, Status: "confirmed"}
	if err := s.repo.CreateEvent(ctx, event); err != nil {
		return nil, nil, err
	}
	wallet, err := s.repo.GetWallet(ctx, creatorID)
	if err != nil {
		return nil, nil, err
	}
	wallet.AvailableBalance += amountToCredit
	wallet.TotalEarned += amountToCredit
	if err := s.repo.UpsertWallet(ctx, wallet); err != nil {
		return nil, nil, err
	}
	return event, wallet, nil
}

func (s *Service) Wallet(ctx context.Context, creatorID string) (*model.CreatorWallet, []model.RevenueEvent, []model.PayoutRecord, error) {
	wallet, err := s.repo.GetWallet(ctx, creatorID)
	if err != nil {
		return nil, nil, nil, err
	}
	events, err := s.repo.ListEvents(ctx, creatorID, 20)
	if err != nil {
		return nil, nil, nil, err
	}
	payouts, err := s.repo.ListPayouts(ctx, creatorID)
	if err != nil {
		return nil, nil, nil, err
	}
	return wallet, events, payouts, nil
}

func (s *Service) RequestPayout(ctx context.Context, creatorID string, amount int64, channel string) (*model.PayoutRecord, error) {
	wallet, err := s.repo.GetWallet(ctx, creatorID)
	if err != nil {
		return nil, err
	}
	if amount <= 0 || amount > wallet.AvailableBalance {
		return nil, errors.New("insufficient balance")
	}
	wallet.AvailableBalance -= amount
	wallet.FrozenBalance += amount
	if err := s.repo.UpsertWallet(ctx, wallet); err != nil {
		return nil, err
	}
	payout := &model.PayoutRecord{CreatorID: creatorID, Amount: amount, Channel: channel, Status: "requested"}
	if err := s.repo.CreatePayout(ctx, payout); err != nil {
		return nil, err
	}
	return payout, nil
}

func (s *Service) ListRules(ctx context.Context) ([]model.RevenueRule, error) {
	return s.repo.ListRules(ctx)
}

func (s *Service) SaveRule(ctx context.Context, rule *model.RevenueRule) error {
	return s.repo.SaveRule(ctx, rule)
}

// AdminListPayouts surfaces payout requests for backoffice review.
func (s *Service) AdminListPayouts(ctx context.Context, status string, limit, offset int) ([]model.PayoutRecord, error) {
	return s.repo.ListAllPayouts(ctx, status, limit, offset)
}

// AdminUpdatePayout updates status and syncs wallet balances.
// status should be "approved" or "rejected".
func (s *Service) AdminUpdatePayout(ctx context.Context, payoutID, status string) (*model.PayoutRecord, *model.CreatorWallet, error) {
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "approved" && status != "rejected" {
		return nil, nil, errors.New("invalid status")
	}
	payout, err := s.repo.FindPayout(ctx, payoutID)
	if err != nil {
		return nil, nil, err
	}
	if payout == nil {
		return nil, nil, errors.New("payout not found")
	}
	if payout.Status == "approved" || payout.Status == "rejected" {
		return payout, nil, nil
	}
	wallet, err := s.repo.GetWallet(ctx, payout.CreatorID)
	if err != nil {
		return nil, nil, err
	}
	// Adjust wallet: on approval -> move frozen to paid (reduce frozen only).
	// on rejection -> refund frozen back to available.
	if status == "approved" {
		if wallet.FrozenBalance >= payout.Amount {
			wallet.FrozenBalance -= payout.Amount
		} else {
			wallet.FrozenBalance = 0
		}
	} else if status == "rejected" {
		wallet.FrozenBalance -= payout.Amount
		if wallet.FrozenBalance < 0 {
			wallet.FrozenBalance = 0
		}
		wallet.AvailableBalance += payout.Amount
	}
	if err := s.repo.UpsertWallet(ctx, wallet); err != nil {
		return nil, nil, err
	}
	updated, err := s.repo.UpdatePayoutStatus(ctx, payoutID, status)
	if err != nil {
		return nil, nil, err
	}
	return updated, wallet, nil
}

func applyRules(eventType string, base int64, rules []model.RevenueRule) int64 {
	var matched *model.RevenueRule
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if strings.EqualFold(rule.EventType, eventType) {
			matched = &rule
			break
		}
	}
	if matched == nil {
		return base
	}
	if matched.Rate > 0 {
		return int64(float64(base) * matched.Rate)
	}
	if matched.Amount > 0 {
		return matched.Amount
	}
	return base
}
