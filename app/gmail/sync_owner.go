package gmail

import (
	"context"
	"drfin/app/owners"
)

type OwnerSyncResult struct {
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
}

// SyncOwners syncs one owner, or every owner when ownerID is nil.
// Owners are independent: one failure does not stop the rest.
func (s *Service) SyncOwners(ctx context.Context, ownerID *int) ([]OwnerSyncResult, error) {
	var list []owners.Owner
	if ownerID != nil {
		o, err := s.ownerRepository.GetByID(*ownerID)
		if err != nil {
			return nil, err
		}
		list = []owners.Owner{*o}
	} else {
		var err error
		if list, err = s.ownerRepository.List(); err != nil {
			return nil, err
		}
	}

	results := make([]OwnerSyncResult, 0, len(list))
	for _, o := range list {
		res := OwnerSyncResult{OwnerID: o.ID, Name: o.Name, OK: true}
		if err := s.syncMail(ctx, o); err != nil {
			res.OK = false
			res.Error = err.Error()
		}
		results = append(results, res)
	}
	return results, nil
}