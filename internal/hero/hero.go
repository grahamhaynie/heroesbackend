package hero

type Hero struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Power    string `json:"power"`
	AlterEgo string `json:"alterEgo"`
	PhotoURL string `json:"photoURL"`
}
