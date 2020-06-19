package main

type Payment interface {
	Pay(money float64, data map[string]interface{}) error
}

type Operations struct {
	payment Payment
}

func (op *Operations) Pay(money float64, data map[string]interface{}) error  {
	return op.payment.Pay(money, data)
}


func CreateOperation(payment Payment) Operations {
	return Operations{payment: payment}
}

type Alipay struct {}
func (ali *Alipay) Pay(money float64, data map[string]interface{}) error {
	return nil
}

type Wechat struct {}
func (wechat *Wechat) Pay(money float64, data map[string]interface{}) error {
	return nil
}

type Jidong struct {}
func (jidong *Jidong) Pay(money float64, data map[string]interface{}) error {
	return nil
}


type Yibao struct {}
func (yibao *Yibao) Pay(money float64, data map[string]interface{}) error {
	return nil
}



func main() {
	alipay := new(Alipay)
	operator := CreateOperation(alipay)
	operator.Pay(1.0, map[string]interface{}{
		"hello": "world",
	})
}