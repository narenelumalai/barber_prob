package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	Barbers     = 2
	Chairs      = 5
	ClosingTime = 30 // in seconds
	Delay       = 2  // in seconds, time between customer arrivals
	CuttingTime = 5  // in seconds to finish cutting
)

var (
	waitingRoom    = make(chan struct{}, Chairs)
	barberChair    = make(chan struct{}, Barbers)
	wg             sync.WaitGroup
	closingChannel = make(chan struct{})
)

func main() {
	fmt.Println("Barbershop simulation started...")

	for i := 1; i <= Barbers; i++ {
		go barber(i)
	}

	go openShop()

	time.Sleep(ClosingTime * time.Second)

	closeShop()

	wg.Wait()
	fmt.Println("Barbershop closed.")
}

func openShop() {
	for {
		select {
		case <-time.After(time.Duration(Delay) * time.Second):
			wg.Add(1)
			go customerArrives()
		case <-closingChannel:
			return
		}
	}
}

func closeShop() {
	close(closingChannel)
	close(waitingRoom)
}

func customerArrives() {
	defer wg.Done()

	select {
	case waitingRoom <- struct{}{}:
		fmt.Println("Customer arrived and took a seat.")
		select {
		case barberChair <- struct{}{}:
			fmt.Println("Barber starts cutting hair.")
			time.Sleep(time.Duration(CuttingTime) * time.Second)
			fmt.Println("Barber finished cutting hair.")
			<-barberChair
		default:
			fmt.Println("All barber chairs are occupied, customer waiting.")
			<-waitingRoom
		}
	default:
		fmt.Println("No available seats in the waiting room, customer leaves.")
	}
}

func barber(id int) {
	for {
		select {
		case <-closingChannel:
			return
		default:
			fmt.Printf("Barber %d is waiting for a customer.\n", id)
			select {
			case <-waitingRoom:
				fmt.Printf("Barber %d woke up and started cutting hair.\n", id)
				time.Sleep(time.Duration(CuttingTime) * time.Second)
				fmt.Printf("Barber %d finished cutting hair.\n", id)
			default:
				fmt.Printf("Barber %d is sleeping.\n", id)
				time.Sleep(1 * time.Second)
			}
		}
	}
}
