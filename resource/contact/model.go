package contact

type Contact struct {
	Name string `json:"name"`
}

type GetAllCmd struct{}

type GetCmd struct {
	ContactID string
}

type DeleteCmd struct {
	ContactID string
}

type CreateCmd struct {
	Name string
}

var ValidContactID = "8c21296d-fbe8-4ddd-aa09-a888a06d66b7"
var ValidContact = Contact{
	Name: "Jane Doe",
}
