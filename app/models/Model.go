package models

type Model interface {
	Insert() (Id uint64, err error)
	Update(Id uint) (Rows int64, err error)
	Delete(Ids string) (Rows int64, err error)
	GetAll(RequestData map[string]interface{}, Offset uint64, Limit uint8) (Data interface{}, err error)
	GetOne(Id uint) (Data interface{}, err error)
}
