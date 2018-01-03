package transactor

// Transactor
// Storage test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "fmt"
import "testing"

func TestStorage(t *testing.T) {
	st := newStorage()
	//fmt.Print(st.data2[st.id(7)])
	st.addUnit(7)
	//fmt.Print(st.id(7))
	//fmt.Print(st.data2[st.id(7)])
	if st.getUnit(7) == nil {
		t.Error("Error storage 1")
	}

}
