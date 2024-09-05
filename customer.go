package eboekhouden

import (
	"context"
	"fmt"
	"time"

	"github.com/hooklift/gowsdl/soap"

	eboekhouden "github.com/sallandpioneers/go-eboekhouden/generated"
	"github.com/sallandpioneers/go-eboekhouden/model"
)

func (service *Eboekhouden) CustomerCreate(ctx context.Context, customer model.Customer) error {
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

		if err := service.handleError(*resp.AddRelatieResult.ErrorMsg); err != nil {
			return fmt.Errorf("add relatie: %w", err)
		}

		customer.ID = resp.AddRelatieResult.Rel_ID

		return nil
	})
}

func (service *Eboekhouden) CustomerUpdate(ctx context.Context, customer model.Customer) error {
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

		if err := service.handleError(*resp.UpdateRelatieResult); err != nil {
			return fmt.Errorf("update relatie: %w", err)
		}

		return nil
	})
}

func (service *Eboekhouden) CustomerGet(ctx context.Context) ([]model.Customer, error) {
	var customers []model.Customer
	err := service.do(ctx, func(session *session) error {
		getRelatie := &eboekhouden.GetRelaties{
			SessionID:     session.SessionID,
			SecurityCode2: session.SecurityCode2,
			CFilter:       &eboekhouden.CRelatieFilter{},
		}

		resp, err := service.client.GetRelatiesContext(ctx, getRelatie)
		if err != nil {
			return fmt.Errorf("get relation: %w", err)
		}

		if resp.GetRelatiesResult == nil {
			return fmt.Errorf("get relation: %s", "relation result is empty")
		}

		if resp.GetRelatiesResult.Relaties == nil {
			if resp.GetRelatiesResult.ErrorMsg != nil {
				return fmt.Errorf("get relation: %s", resp.GetRelatiesResult.ErrorMsg.LastErrorCode)
			}
			return fmt.Errorf("get relation: %s", "relations empty")
		}

		for _, rel := range resp.GetRelatiesResult.Relaties.CRelatie {
			customers = append(customers, model.Customer{
				ID:         rel.ID,
				Code:       rel.Code,
				Business:   rel.Bedrijf,
				Addresses:  model.CustomerAddresses{},
				Phone:      "",
				Email:      rel.Email,
				Website:    "",
				Notition:   "",
				VAT:        "",
				COC:        "",
				Salutation: "",
				IBAN:       "",
				BIC:        "",
				Type:       "",
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func getRelation(customer model.Customer) *eboekhouden.CRelatie {
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
