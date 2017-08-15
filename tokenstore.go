package yljwt

type TokenStore interface {
	Create(token string) error //创建一个token
	Discard(tid string) error  //丢弃一个tid
	Check(tid string) bool     //检查tid是否有效
}

type DefaultTokenStore struct {
	invalid map[string]bool
}

func (d DefaultTokenStore) Create(tid string) error {
	return nil
}

func (d DefaultTokenStore) Discard(tid string) error {
	d.invalid[tid] = true
	return nil

}

func (d DefaultTokenStore) Check(tid string) bool {
	if d.invalid[tid] == true {
		return false
	}
	return true
}
