package repository

const (
	limitDefault = 10
	limitMin     = 0
	limitMax     = 20

	offsetDefault = 0
	offsetMin     = 0
)

type Option struct {
	userId string
	limit  int
	offset int
}

func DefaultOption() *Option {
	return &Option{
		userId: "",
		limit:  limitDefault,
		offset: offsetDefault,
	}
}

func (p *Option) SetUserId(userId string) {
	p.userId = userId
}

func (p *Option) SetLimit(limit int) {
	if limit < limitMin {
		limit = limitMin
	} else if limit > limitMax {
		limit = limitMax
	}
	p.limit = limit
}

func (p *Option) SetOffset(offset int) {
	if offset < offsetMin {
		offset = offsetMin
	}
	p.offset = offset
}

func (p *Option) Map() map[string]interface{} {
	return map[string]interface{}{
		"userId": p.userId,
		"limit":  p.limit,
		"offset": p.offset,
	}
}
