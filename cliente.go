package main

import (
	"fmt"
	"net"
	"time"
	"encoding/gob"
)

//
type Mensaje struct {
	ID int
	Cont int
}

func salida(salir chan bool) {
	var input string
	fmt.Scanln(&input)

	salir <- true
}

func main() {
	saliendo := make(chan bool)
	c, err := net.Dial("tcp", ":9999")
	var mensaje Mensaje

	if err != nil {
		fmt.Println(err)
		return
	}

	go salida(saliendo)

	err = gob.NewDecoder(c).Decode(&mensaje)

	if err != nil {
		fmt.Println(err)
		return
	}
	if mensaje.ID < 0 {
		fmt.Println("Vacio")
		return
	}

	for {
		select {
			case <-saliendo:
				err := gob.NewEncoder(c).Encode(mensaje)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Apagado")
				return

			default:
				fmt.Println(mensaje.ID, ":", mensaje.Cont)
				mensaje.Cont++
				time.Sleep(time.Millisecond * 500)
		}
	}
}