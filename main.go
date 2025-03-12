package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type TransportMethod interface {
	DeliverPackage(destination string) (string, error)
	GetStatus() string
}

type Truck struct {
	ID       string
	Capacity int
}

// DeliverPackage simule la livraison par truck
func (t *Truck) DeliverPackage(destination string) (string, error) {
	log.Printf("Camion %s commence la livraison vers %s", t.ID, destination)
	time.Sleep(3 * time.Second) // Simulation d'un délai de livraison
	return fmt.Sprintf("Camion %s a livré vers %s", t.ID, destination), nil
}

// GetStatus renvoie l'état du truck
func (t *Truck) GetStatus() string {
	return fmt.Sprintf("Camion %s prêt (Capacité : %d)", t.ID, t.Capacity)
}

type Drone struct {
	ID      string
	Battery int
}

// DeliverPackage simule la livraison par drone
func (d *Drone) DeliverPackage(destination string) (string, error) {
	if d.Battery < 20 {
		return "", errors.New("batterie insuffisante")
	}
	d.Battery -= 20
	log.Printf("Drone %s commence la livraison vers %s", d.ID, destination)
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("Drone %s a livré vers %s", d.ID, destination), nil
}

// GetStatus renvoie l'état du drone
func (d *Drone) GetStatus() string {
	return fmt.Sprintf("Drone %s batterie : %d%%", d.ID, d.Battery)
}

// Bateau représente une méthode de livraison par boat
type Boat struct {
	ID      string
	Weather string
}

// DeliverPackage simule la livraison par boat
func (b *Boat) DeliverPackage(destination string) (string, error) {
	if b.Weather != "Clear" {
		return "", errors.New("conditions météo défavorables")
	}
	log.Printf("Bateau %s commence la livraison vers %s", b.ID, destination)
	time.Sleep(4 * time.Second)
	return fmt.Sprintf("Bateau %s a livré vers %s", b.ID, destination), nil
}

// GetStatus renvoie l'état du boat
func (b *Boat) GetStatus() string {
	return fmt.Sprintf("Bateau %s météo : %s", b.ID, b.Weather)
}

func main() {

}
