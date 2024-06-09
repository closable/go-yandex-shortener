package handlers

import "fmt"

func ExampleGetUserID() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTI3MzgzOTAsIlVzZXJJRCI6NH0.V9WdWdJWeU1qqVCGDfTGu0asPZhiFUPmtnsfpN0GPro"
	//userID := 5
	out1 := GetUserID(tokenString)
	fmt.Println(out1)

	//Output
	//5

}
