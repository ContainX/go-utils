package encoding

type personStruct struct {
	Name   string   `json:"name,omitempty"`
	Age    int      `json:"age,omitempty"`
	Score  float64  `json:"score,omitempty"`
	Colors []string `json:"colors,omitempty"`
	yaml   string
	json   string
}

var personTests = []personStruct{
	{
		Name: "Jack",
		yaml:  "name: Jack",
		json: `{"name":"Jack"}`,
	},
	{
		Name: "Jack",
		Age:  22,
		yaml:  "name: Jack\nage: 22",
		json: `{"name":"Jack","age":22}`,
	},
	{
		Colors: []string{"red", "blue"},
		yaml:    "colors:\n- red\n- blue",
		json: `{"colors":["red","blue"]}`,
	},
	{
		Score: 22.5,
		yaml:   "score: 22.5",
		json: `{"score":22.5}`,
	},
}
