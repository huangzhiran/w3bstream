package account

import (
	"context"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/modules/account"
)

type UpdatePasswordByAccountID struct {
	httpx.MethodPut
	account.UpdatePasswordReq `in:"body"`
}

func (r *UpdatePasswordByAccountID) Path() string { return "/:accountID" }

func (r *UpdatePasswordByAccountID) Output(ctx context.Context) (interface{}, error) {
	ca := middleware.CurrentAccountFromContext(ctx)
	return nil, account.UpdateAccountPassword(ctx, ca.AccountID, &r.UpdatePasswordReq)
}
