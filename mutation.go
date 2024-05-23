package eboekhouden

import (
	"context"
	"fmt"

	"github.com/hooklift/gowsdl/soap"

	eboekhouden "github.com/sallandpioneers/go-eboekhouden/generated"
	"github.com/sallandpioneers/go-eboekhouden/model"
)

func (service *Eboekhouden) MutationCreate(ctx context.Context, mutation *model.Mutation) error {
	return service.do(ctx, func(session *session) error {
		mutationType := eboekhouden.EnMutatieSoorten(mutation.Type)

		addMutation := &eboekhouden.AddMutatie{
			SessionID:     session.SessionID,
			SecurityCode2: session.SecurityCode2,
			OMut: &eboekhouden.CMutatie{
				MutatieNr:        0,
				Soort:            &mutationType,
				Datum:            soap.CreateXsdDateTime(mutation.Date, false),
				Rekening:         mutation.LedgerAccountCode,
				RelatieCode:      mutation.BoekhoudenCustomerID,
				Factuurnummer:    mutation.InvoiceNumber,
				Boekstuk:         mutation.LedgerAccountCode,
				Omschrijving:     mutation.Description,
				Betalingstermijn: mutation.PaymentTerm,
				Betalingskenmerk: mutation.PaymentFeature,
				InExBTW:          string(model.Exclusive),
				MutatieRegels: &eboekhouden.ArrayOfCMutatieRegel{
					CMutatieRegel: make([]*eboekhouden.CMutatieRegel, len(mutation.Items)),
				},
			},
		}

		for k, v := range mutation.Items {
			addMutation.OMut.MutatieRegels.CMutatieRegel[k] = &eboekhouden.CMutatieRegel{
				BedragInvoer:      v.Amount,
				BedragExclBTW:     v.AmountExVAT,
				BedragBTW:         v.AmountVAT,
				BedragInclBTW:     v.AmountInVAT,
				BTWCode:           string(v.VATCode),
				BTWPercentage:     v.VATPercentage,
				TegenrekeningCode: v.LedgerAccountCode,
				KostenplaatsID:    v.KostenplaatsID,
			}
		}

		resp, err := service.client.AddMutatieContext(ctx, addMutation)
		if err != nil {
			return err
		}

		if resp.AddMutatieResult == nil {
			return fmt.Errorf("add mutatie: %w", ErrResponseEmpty)
		}

		if resp.AddMutatieResult.ErrorMsg != nil {
			return fmt.Errorf("add mutatie: %w", service.handleError(*resp.AddMutatieResult.ErrorMsg))
		}

		return nil
	})
}
