package eboekhouden

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hooklift/gowsdl/soap"

	"github.com/sallandpioneers/go-eboekhouden/config"
	eboekhouden "github.com/sallandpioneers/go-eboekhouden/generated"
)

var (
	ErrRelationCodeAlreadyInUse = errors.New("relation code already in use")
	ErrSessionExpired           = errors.New("errors expired")
	ErrCreatingNewSession       = errors.New("error creating a new session")
	ErrEboekhoudenError         = errors.New("eboekhouding error")
	ErrResponseEmpty            = errors.New("response empty")
)

type Eboekhouden struct {
	client      eboekhouden.SoapAppSoap
	session     session
	cfg         *config.Config
	openSession eboekhouden.OpenSession
	logger      *slog.Logger
}

type session struct {
	createdAt     time.Time
	SessionID     string
	SecurityCode2 string
}

var (
	ErrLoggerNil = errors.New("logger cannot be nil")
	ErrConfigNil = errors.New("config cannot be nil")
)

func New(logger *slog.Logger, cfg *config.Config) (*Eboekhouden, error) {
	client := soap.NewClient(cfg.URL)

	if logger == nil {
		return nil, ErrLoggerNil
	}

	if cfg == nil {
		return nil, ErrConfigNil
	}

	return &Eboekhouden{
		client: eboekhouden.NewSoapAppSoap(client),
		session: session{
			SecurityCode2: cfg.SecurityCode2,
		},
		openSession: eboekhouden.OpenSession{
			Username:      cfg.Username,
			SecurityCode1: cfg.SecurityCode1,
			SecurityCode2: cfg.SecurityCode2,
			Source:        cfg.Source,
		},
		cfg:    cfg,
		logger: logger.WithGroup("go-eboekhouden"),
	}, nil
}

func (service *Eboekhouden) getSession(ctx context.Context) (*session, error) {
	if service.session.createdAt.Before(time.Now().Add(-5 * time.Minute)) {
		if err := service.newSession(ctx); err != nil {
			return nil, fmt.Errorf("new session: %w", err)
		}
	}

	return &service.session, nil
}

func (service *Eboekhouden) newSession(ctx context.Context) error {
	resp, err := service.client.OpenSessionContext(ctx, &service.openSession)
	if err != nil {
		return fmt.Errorf("open session: %w", err)
	}

	if resp.OpenSessionResult == nil {
		return fmt.Errorf("open session result: %w", ErrCreatingNewSession)
	}

	if resp.OpenSessionResult.ErrorMsg != nil {
		if err := service.handleError(*resp.OpenSessionResult.ErrorMsg); err != nil {
			return fmt.Errorf("open session error msg: %w", err)
		}
	}

	service.session.SessionID = resp.OpenSessionResult.SessionID
	service.session.createdAt = time.Now()

	service.logger.Debug(
		"new session created",
		"sessionID", service.session.SessionID,
		"createdAt", service.session.createdAt.Format(time.RFC3339),
	)

	return nil
}

func (service *Eboekhouden) do(ctx context.Context, fn func(*session) error) error {
	// Get session
	session, err := service.getSession(ctx)
	if err != nil {
		return err
	}

	// Execute fn, if there is no error, return
	err = fn(session)
	if err == nil {
		return nil
	}

	// If we have an error but it is not a session problem, return
	if !errors.Is(err, ErrSessionExpired) {
		return err
	}

	// We should only have errors that are about the session, retrieve a new one
	if err := service.newSession(ctx); err != nil {
		return err
	}

	// Get the session
	session, err = service.getSession(ctx)
	if err != nil {
		return err
	}

	// Try fn again if we still fail return
	if err := fn(session); err != nil {
		return err
	}

	return nil
}

func (service *Eboekhouden) handleError(err eboekhouden.CError) error {
	switch err.LastErrorCode {
	case "":
		return nil
	case "CREL001":
		return ErrRelationCodeAlreadyInUse
	case "E0023":
		service.logger.Debug("session expired, create new")

		return ErrSessionExpired
	default:
		service.logger.Debug("eboekhouden error occurred", "code", err.LastErrorCode, "description", err.LastErrorDescription)
		return fmt.Errorf("%w: code: %s description: %s", ErrEboekhoudenError, err.LastErrorCode, err.LastErrorDescription)
	}
}
