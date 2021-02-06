package web

/*
// RunClient - NNGSクライアントを走らせます。
func RunClient(host string, port uint16) error {
	connectionString := fmt.Sprintf("%s:%d", host, port)
	return telnet.DialToAndCall(connectionString, clientListener{})
	// return telnet.DialToAndCall("localhost:9696", clientListener{})
}

type clientListener struct{}

// CallTELNET - 決まった形のメソッド。
func (c clientListener) CallTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var carry [1]byte  // このバッファーが埋まるまで待ってしまう？ だから 1文字ずつ。
	pCarry := carry[:] // スライス？

	var chars [1024]byte // 1byteでは使いづらいんで 溜めます。
	i := 0

	for {
		// いつ通信が切れるのか分からん☆（＾～＾）
		n, err := r.Read(pCarry)
		print(fmt.Sprintf("(n=%d)", n))

		// 相手が切断したときなど。
		if n <= 0 || err != nil {
			break
		}

		if n == 1 {
			bytes := pCarry[:n] // 結局 1バイトだが。

			chars[i] = bytes[0]
			i += n

			// print(string(bytes))
			print("[")
			print(string(chars[:i]))
			print("]\n")
		} else {
			panic(fmt.Sprintf("I want an 1 byte, but %d bytes.", n))
		}
	}

	print("Done [")
	print(string(chars[:i]))
	print("]\n")

}
*/
