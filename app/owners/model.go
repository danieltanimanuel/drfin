package owners

import (
	"fmt"
)

type Owner struct {
	ID             int
	Email          string
	Name           string
	BankHolderName string
}

// DisplayName is the sidebar label: name, else email, else "Owner N".
func (o Owner) DisplayName() string {
	switch {
	case o.Name != "":
		return o.Name
	case o.Email != "":
		return o.Email
	default:
		return fmt.Sprintf("Owner %d", o.ID)
	}
}