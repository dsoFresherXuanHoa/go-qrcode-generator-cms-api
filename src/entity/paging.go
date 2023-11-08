package entity

type Paginate struct {
	Page  int `form:"page"`
	Size  int `form:"size"`
	Total int `form:"-" json:"-"`
}

func (p *Paginate) Standardized() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 2 {
		p.Size = 2
	}
	if p.Size > 10 {
		p.Size = 10
	}
}
