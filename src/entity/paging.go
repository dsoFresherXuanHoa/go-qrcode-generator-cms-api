package entity

type Paginate struct {
	Page  int `form:"page"`
	Size  int `form:"size"`
	Total int `form:"-" json:"total"`
}

func (p *Paginate) Standardized() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 5 {
		p.Size = 5
	}
	if p.Size > 15 {
		p.Size = 15
	}
}
