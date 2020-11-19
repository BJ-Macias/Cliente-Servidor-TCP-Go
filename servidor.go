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

func procesado(dato chan Mensaje, cerrar chan bool, aux chan Mensaje) {
	activos := [5]bool{true, true, true, true, true}
	procesos := [5]int{0, 0, 0, 0, 0}

	for {
		select {
			case <-cerrar:
				terminado := false
				for i := 0; i <= 4; i++ {
					if activos[i] {
						aux <- Mensaje{
							ID:    i,
							Cont: procesos[i],
						}
						activos[i] = false
						terminado = true
						break
					}
				}
				if !terminado {
					aux <- Mensaje{
						ID: -1,
						Cont: 0,
					}
				}
				break

			case contenido := <-dato:
				if contenido.ID < 0 || contenido.ID > 4 {
					break
				}
				activos[contenido.ID] = true
				procesos[contenido.ID] = contenido.Cont
				break

			default:
		}

		fmt.Println("\n********\n")
		for i := 0; i <= 4; i++ {
			if activos[i] {
				fmt.Println(i, ":", procesos[i])
				procesos[i]++
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func peticiones(c net.Conn, dato chan Mensaje, aux chan Mensaje) {
	cliente := <-aux
	err := gob.NewEncoder(c).Encode(cliente)

	if err != nil {
		fmt.Println(err)
		dato <- cliente
		c.Close()
		return
	}

	if cliente.ID < 0 {
		return
	}

	var msg Mensaje
	err = gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		dato <- cliente
	} else {
		dato <- msg
	}
}

func main() {
	desconectado := make(chan Mensaje)
	proceso := make(chan Mensaje)
	detener := make(chan bool)
	
	c, err := net.Listen("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	go procesado(proceso, detener, desconectado)
	defer c.Close()

	for {
		c, err := c.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		detener <- true
		go peticiones(c, proceso, desconectado)
	}
}