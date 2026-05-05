package user

type Entity struct {
	ID        int64          `json:"id" pg:"id,pk"`
	FirstName string         `json:"firstName" pg:"first_name"`
	LastName  string         `json:"lastName" pg:"last_name"`
	Nickname  string         `json:"nickname" pg:"nickname"`
	Phone     string         `json:"phone" pg:"phone"`
	Email     string         `json:"email" pg:"email"`
	Status    string         `json:"status" pg:"status"`
	Role      string         `json:"role" pg:"role"`
	Version   int64          `json:"version" pg:"version"`
	Metadata  map[string]any `json:"metadata" pg:"metadata,scanonly"`
}
