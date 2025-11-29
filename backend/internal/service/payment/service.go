package payment

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"crypto/rand"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/example/ai-avatar-studio/internal/repository"
)

// Service orchestrates recharge flow through 易支付 (MD5) gateway.
type Service struct {
	payments      *repository.PaymentRepository
	assets        *repository.UserAssetRepository
	notifications *repository.NotificationRepository
	merchantID    string
	merchantKey   string
	gateway       string
	notifyURL     string
	returnURL     string
	coinsPerYuan  int64
}

func NewService(payments *repository.PaymentRepository, assets *repository.UserAssetRepository, notifications *repository.NotificationRepository, merchantID, merchantKey, gateway, notifyURL, returnURL string, coinsPerYuan int64) *Service {
	if coinsPerYuan <= 0 {
		coinsPerYuan = 1000
	}
	return &Service{
		payments:      payments,
		assets:        assets,
		notifications: notifications,
		merchantID:    merchantID,
		merchantKey:   merchantKey,
		gateway:       strings.TrimRight(gateway, "/"),
		notifyURL:     notifyURL,
		returnURL:     returnURL,
		coinsPerYuan:  coinsPerYuan,
	}
}

type CreateOrderResult struct {
	Order      *model.PaymentOrder
	PayURL     string
	OutTradeNo string
}

// CreateOrder builds a payment order and returns the pay URL to redirect the user.
func (s *Service) CreateOrder(ctx context.Context, userID string, amountYuan float64, payType string) (*CreateOrderResult, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("missing user")
	}
	if amountYuan <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if payType == "" {
		payType = "alipay"
	}
	outTradeNo := fmt.Sprintf("%d%s", time.Now().Unix(), strings.ReplaceAll(uuidV4(), "-", ""))
	moneyCents := int64(amountYuan*100 + 0.5)
	coins := int64(amountYuan*float64(s.coinsPerYuan) + 0.5)

	order := &model.PaymentOrder{
		UserID:     userID,
		OutTradeNo: outTradeNo,
		PayType:    payType,
		Status:     "pending",
		MoneyCents: moneyCents,
		Coins:      coins,
	}
	if err := s.payments.Create(ctx, order); err != nil {
		return nil, err
	}

	params := map[string]string{
		"name":         "账户充值",
		"money":        fmt.Sprintf("%.2f", amountYuan),
		"out_trade_no": outTradeNo,
		"notify_url":   s.notifyURL,
		"return_url":   s.returnURL,
		"pid":          s.merchantID,
		"type":         payType,
		"sign_type":    "MD5",
	}
	sign := buildMD5Sign(params, s.merchantKey)
	params["sign"] = sign
	payURL := s.gateway + "/submit.php?" + encodeParams(params)

	return &CreateOrderResult{
		Order:      order,
		PayURL:     payURL,
		OutTradeNo: outTradeNo,
	}, nil
}

// HandleNotify verifies the gateway callback and credits coins idempotently.
func (s *Service) HandleNotify(ctx context.Context, payload map[string]string) (*model.PaymentOrder, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("empty payload")
	}
	sign := payload["sign"]
	if sign == "" {
		return nil, fmt.Errorf("missing sign")
	}
	if payload["trade_status"] != "TRADE_SUCCESS" {
		return nil, fmt.Errorf("trade not successful")
	}
	if !s.verifySign(payload) {
		return nil, fmt.Errorf("invalid sign")
	}

	outTradeNo := payload["out_trade_no"]
	order, err := s.payments.FindByOutTradeNo(ctx, outTradeNo)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.Status == "paid" {
		return order, nil
	}
	// amount check
	paidMoney, _ := strconv.ParseFloat(payload["money"], 64)
	if int64(paidMoney*100+0.5) != order.MoneyCents {
		return nil, fmt.Errorf("amount mismatch")
	}
	updated, err := s.payments.MarkPaid(ctx, outTradeNo, payload["trade_no"], payload)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return order, nil
	}
	// credit coins
	asset, err := s.assets.GetByUser(ctx, updated.UserID)
	if err != nil {
		return nil, err
	}
	asset.Balance += updated.Coins
	if err := s.assets.Upsert(ctx, asset); err != nil {
		return nil, err
	}
	// Notify user
	if s.notifications != nil {
		_ = s.notifications.Create(ctx, &model.Notification{
			UserID: updated.UserID,
			Type:   "recharge",
			Title:  "充值成功",
			Content: fmt.Sprintf("订单 %s 支付成功，到账 %d 平台币。", updated.OutTradeNo, updated.Coins),
		})
	}
	return updated, nil
}

func (s *Service) Query(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error) {
	return s.payments.FindByOutTradeNo(ctx, outTradeNo)
}

func (s *Service) CoinsPerYuan() int64 {
	return s.coinsPerYuan
}

func (s *Service) ListUserOrders(ctx context.Context, userID string, limit, offset int, status string) ([]model.PaymentOrder, error) {
	return s.payments.ListByUser(ctx, userID, limit, offset, status)
}

func (s *Service) ListAll(ctx context.Context, limit, offset int) ([]model.PaymentOrder, error) {
	return s.payments.ListAll(ctx, limit, offset)
}

func (s *Service) verifySign(payload map[string]string) bool {
	given := payload["sign"]
	expected := buildMD5Sign(payload, s.merchantKey)
	return strings.EqualFold(given, expected)
}

// buildMD5Sign sorts params, skips sign/sign_type/empty, joins as a=b&c=d + key, md5 lower.
func buildMD5Sign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || k == "sign_type" || strings.TrimSpace(v) == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	raw := strings.Join(parts, "&") + key
	sum := md5.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func encodeParams(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

// uuidV4 returns a random UUID v4 string.
func uuidV4() string {
	u := [16]byte{}
	// cheap random source
	_, _ = rand.Read(u[:])
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
