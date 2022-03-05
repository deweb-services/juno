package bridge

type WasmMsgExecuteContract struct {
	Sender   string
	Contract string
	Msg      interface{}
}
