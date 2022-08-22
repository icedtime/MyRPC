package impl

import (
	"fmt"
	"server/message"
)

type Test struct {
}

func NewTest() *Test {
	return &Test{}
}
func (t *Test) SayHello(req *message.Request, res *message.Response) error {
	data := req.Msg.Data
	msg := fmt.Sprintf("%s,风止意难平！", data)

	res.Msg.Data = msg

	return nil
}
