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

func (o *Option) GetUserId() string {
	return o.userId
}

func (o *Option) SetUserId(userId string) {
	o.userId = userId
}

func (o *Option) SetLimit(limit int) {
	if limit < limitMin {
		limit = limitMin
	} else if limit > limitMax {
		limit = limitMax
	}
	o.limit = limit
}

func (o *Option) SetOffset(offset int) {
	if offset < offsetMin {
		offset = offsetMin
	}
	o.offset = offset
}

func (o *Option) Map() map[string]interface{} {
	return map[string]interface{}{
		"userId": o.userId,
		"limit":  o.limit,
		"offset": o.offset,
	}
}
