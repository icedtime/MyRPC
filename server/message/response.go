package message

//先这样封装，因为我不知道一个response到底还要什么额外的字段
type Response struct {
	Msg *ReciveMessage
}

func NewResponse() *Response {
	return &Response{
		Msg: NewReciveMessage(""),
	}
}

type DisResponse struct {
}
