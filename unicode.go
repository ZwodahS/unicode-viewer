package main

// func main() {
// 	for i := 0; i < 4000; i++ {
// 		if i%16 == 0 {
// 			fmt.Print("|" + fmt.Sprintf("%x", i) + "|")
// 		}
// 		if isControl(i) {
// 			fmt.Print(" ")
// 			continue
// 		}
// 		value := fmt.Sprintf("%U", i)[2:]
// 		quoted := "'\\u" + value + "'"
// 		c, err := strconv.Unquote(quoted)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			panic(err)
// 		}
// 		fmt.Print(c)
// 	}
//
// }
