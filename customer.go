package eboekhouden

import (
	"context"
	"fmt"
	"time"

	"github.com/hooklift/gowsdl/soap"

	eboekhouden "git.sallandpioneers.com/sallandpioneers/go-eboekhouden/generated"
	"git.sallandpioneers.com/sallandpioneers/go-eboekhouden/model"
)

func (service *Eboekhouden) CustomerCreate(ctx context.Context, customer *model.Customer) error {
	return service.do(ctx, func(session *session) error {
		addRelatie := &eboekhouden.AddRelatie{
			SessionID:     session.SessionID,
			SecurityCode2: session.SecurityCode2,
			ORel:          getRelation(customer),
		}

		resp, err := service.client.AddRelatieContext(ctx, addRelatie)
		if err != nil {
			return fmt.Errorf("add relatie: %w", err)
		}

		if resp.AddRelatieResult == nil {
			return fmt.Errorf("add relatie: %w", ErrResponseEmpty)
		}

		if resp.AddRelatieResult.ErrorMsg != nil {
			return fmt.Errorf("add relatie: %w", service.handleError(*resp.AddRelatieResult.ErrorMsg))
		}

		customer.ID = resp.AddRelatieResult.Rel_ID

		return nil
	})
}

func (service *Eboekhouden) CustomerUpdate(ctx context.Context, customer *model.Customer) error {
	return service.do(ctx, func(session *session) error {
		updateRelate := &eboekhouden.UpdateRelatie{
			SessionID:     session.SessionID,
			SecurityCode2: session.SecurityCode2,
			ORel:          getRelation(customer),
		}

		resp, err := service.client.UpdateRelatieContext(ctx, updateRelate)
		if err != nil {
			return fmt.Errorf("update relatie: %w", err)
		}

		if resp.UpdateRelatieResult == nil {
			return fmt.Errorf("update relatie: %w", ErrResponseEmpty)
		}

		if resp.UpdateRelatieResult != nil {
			return fmt.Errorf("update relatie: %w", service.handleError(*resp.UpdateRelatieResult))
		}

		return nil
	})
}

func getRelation(customer *model.Customer) *eboekhouden.CRelatie {
	return &eboekhouden.CRelatie{
		ID:        customer.ID,
		AddDatum:  soap.CreateXsdDateTime(time.Now(), true),
		Code:      customer.Code,
		Bedrijf:   customer.Business,
		Adres:     customer.Addresses.Business.Address,
		Postcode:  customer.Addresses.Business.ZipCode,
		Plaats:    customer.Addresses.Business.City,
		Land:      customer.Addresses.Business.Country,
		Adres2:    customer.Addresses.Mailing.Address,
		Postcode2: customer.Addresses.Mailing.ZipCode,
		Plaats2:   customer.Addresses.Mailing.City,
		Land2:     customer.Addresses.Mailing.Country,
		Telefoon:  customer.Phone,
		Email:     customer.Email,
		Site:      customer.Website,
		BTWNummer: customer.VAT,
		KvkNummer: customer.COC,
		IBAN:      customer.IBAN,
		BIC:       customer.BIC,
		BP:        string(customer.Type),
	}
}
