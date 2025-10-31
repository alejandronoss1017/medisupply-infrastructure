package domain

type Comparator string

const (
	Greater        Comparator = ">"
	Less           Comparator = "<"
	Equal          Comparator = "=="
	NotEqual       Comparator = "!="
	GreaterOrEqual Comparator = ">="
	LessOrEqual    Comparator = "<="
)

type SLA struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Target      int64      `json:"target"`
	Comparator  Comparator `json:"comparator"`
	Status      bool       `json:"status"`
}
