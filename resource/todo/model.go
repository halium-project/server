package todo

type Todo struct {
	Title string `json:"name"`
}

type GetAllCmd struct{}

type GetCmd struct {
	TodoID string
}

type DeleteCmd struct {
	TodoID string
}

type CreateCmd struct {
	Title string
}

var ValidTodoID = "8c21296d-fbe8-4ddd-aa09-a888a06d66b7"
var ValidTodo = Todo{
	Title: "Jane Doe",
}
