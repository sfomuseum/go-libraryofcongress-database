package docstore

import (
	"fmt"
)

type Document struct {
	Id     string `json:"id"`
	Label  string `json:"label"`
	Source string `json:"source"`
}

func (d *Document) String() string {
	return fmt.Sprintf("%s:%s %s", d.Source, d.Id, d.Label)
}

