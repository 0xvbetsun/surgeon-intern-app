package organisations

type (
	APIOrganisation struct {
		ID          string `json:"id,omitempty"`
		DisplayName string `json:"display_name,omitempty"`
	}
)

func (o *APIOrganisation) Valid() bool {
	return o.DisplayName != ""
}
