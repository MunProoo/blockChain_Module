package types

import "github.com/hacpy/go-ethereum/core/types"

type NeedToCatchEventFunc func(e *types.Log, tx *types.Transaction)

type NeedToCatchEvent struct {
	NeedToCatchEventFunc NeedToCatchEventFunc
}
