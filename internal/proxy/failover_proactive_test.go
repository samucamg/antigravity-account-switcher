package proxy

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samucamg/antigravity-account-switcher/internal/domain"
)

func TestFailoverEngine_RotateProactively_Success(t *testing.T) {
	repo := newMockAccountRepo()
	now := time.Now().UTC()
	accA := &domain.Account{
		ID:          "acc-A",
		Email:       "a@example.com",
		AccessToken: "token-A",
		IsActive:    true,
		Status:      domain.AccountStatusActive,
		UpdatedAt:   now.Add(-10 * time.Minute),
	}
	accB := &domain.Account{
		ID:          "acc-B",
		Email:       "b@example.com",
		AccessToken: "token-B",
		IsActive:    false,
		Status:      domain.AccountStatusActive,
		UpdatedAt:   now.Add(-5 * time.Minute),
	}
	repo.addAccount(accA)
	repo.addAccount(accB)

	broadcaster := NewBroadcaster(10)
	eventsCh, unsubscribe := broadcaster.Subscribe()
	defer unsubscribe()

	eventRepo := &mockEventRepo{}
	engine := NewFailoverEngine(repo, broadcaster, eventRepo)

	nextAcc, err := engine.RotateProactively(context.Background(), accA, 0.87)
	if err != nil {
		t.Fatalf("expected successful proactive rotation, got error: %v", err)
	}
	if nextAcc.ID != "acc-B" {
		t.Errorf("expected rotated account acc-B, got %s", nextAcc.ID)
	}

	// Proactive rotation must NOT mark the source account exhausted
	updatedA, _ := repo.GetByID(context.Background(), "acc-A")
	if updatedA.Status != domain.AccountStatusActive {
		t.Errorf("proactive rotation should keep acc-A active, got %s", updatedA.Status)
	}

	// Account B must now be the active one
	activeAcc, _ := repo.GetActive(context.Background())
	if activeAcc.ID != "acc-B" {
		t.Errorf("expected active account acc-B, got %s", activeAcc.ID)
	}

	// Verify EventTypeProactiveSwitch was broadcast exactly once
	select {
	case ev := <-eventsCh:
		if ev.Type != domain.EventTypeProactiveSwitch {
			t.Errorf("expected event type %s, got %s", domain.EventTypeProactiveSwitch, ev.Type)
		}
		if ev.AccountID != "acc-B" {
			t.Errorf("expected event AccountID acc-B, got %s", ev.AccountID)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for proactive switch event")
	}

	// No extra events should follow
	select {
	case extra := <-eventsCh:
		t.Errorf("unexpected extra event: %+v", extra)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestFailoverEngine_RotateProactively_NoAlternate_KeepsCurrent(t *testing.T) {
	repo := newMockAccountRepo()
	accA := &domain.Account{
		ID:          "acc-sole",
		Email:       "sole@example.com",
		AccessToken: "token-sole",
		IsActive:    true,
		Status:      domain.AccountStatusActive,
		UpdatedAt:   time.Now().UTC(),
	}
	repo.addAccount(accA)

	broadcaster := NewBroadcaster(10)
	eventsCh, unsubscribe := broadcaster.Subscribe()
	defer unsubscribe()

	engine := NewFailoverEngine(repo, broadcaster, nil)

	// No healthier alternate: graceful no-op returning current account, nil error
	nextAcc, err := engine.RotateProactively(context.Background(), accA, 0.92)
	if err != nil {
		t.Fatalf("expected nil error when no alternate available, got %v", err)
	}
	if nextAcc == nil || nextAcc.ID != "acc-sole" {
		t.Fatalf("expected current account to be returned, got %+v", nextAcc)
	}

	// Account must remain active and not exhausted
	updatedA, _ := repo.GetByID(context.Background(), "acc-sole")
	if updatedA.Status != domain.AccountStatusActive || !updatedA.IsActive {
		t.Errorf("expected acc-sole to stay active, got status=%s isActive=%v", updatedA.Status, updatedA.IsActive)
	}

	// No proactive event should have been emitted
	select {
	case ev := <-eventsCh:
		t.Errorf("unexpected event emitted: %+v", ev)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestFailoverEngine_RotateProactively_NilAccount(t *testing.T) {
	repo := newMockAccountRepo()
	engine := NewFailoverEngine(repo, nil, nil)

	_, err := engine.RotateProactively(context.Background(), nil, 0.85)
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestFailoverEngine_RotateProactively_CanceledContext(t *testing.T) {
	repo := newMockAccountRepo()
	now := time.Now().UTC()
	accA := &domain.Account{
		ID:          "acc-A",
		Email:       "a@example.com",
		AccessToken: "token-A",
		IsActive:    true,
		Status:      domain.AccountStatusActive,
		UpdatedAt:   now.Add(-10 * time.Minute),
	}
	accB := &domain.Account{
		ID:          "acc-B",
		Email:       "b@example.com",
		AccessToken: "token-B",
		IsActive:    false,
		Status:      domain.AccountStatusActive,
		UpdatedAt:   now.Add(-5 * time.Minute),
	}
	repo.addAccount(accA)
	repo.addAccount(accB)

	engine := NewFailoverEngine(repo, nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := engine.RotateProactively(ctx, accA, 0.85)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}

	// Rotation must not have happened despite accounts being available
	activeAcc, _ := repo.GetActive(context.Background())
	if activeAcc.ID != "acc-A" {
		t.Errorf("canceled rotation should not change active account, got %s", activeAcc.ID)
	}
}