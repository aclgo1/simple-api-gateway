package models

import "errors"

func(p *ParamPixWebHookInput)Validate() error{
	if p.GatewayTransactionId == "" {
		return errors.New("gateway transactio id is empty")
	}

	return nil
}