package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	time.Sleep(3 * time.Second)
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

func GetTransportMethod(method string) (TransportMethod, error) {
	switch method {
	case "truck":
		return &Truck{ID: "T123", Capacity: 10}, nil
	case "drone":
		return &Drone{ID: "D456", Battery: 100}, nil
	case "boat":
		return &Boat{ID: "B789", Weather: "Clear"}, nil
	default:
		return nil, errors.New("méthode de transport inconnue")
	}
}

func TrackDelivery(t TransportMethod, destination string, ch chan string) {
	status, err := t.DeliverPackage(destination)
	if err != nil {
		ch <- fmt.Sprintf("Échec de la livraison : %v", err)
		return
	}
	ch <- status
}

func main() {

	r := gin.Default()

	r.POST("/deliver", func(c *gin.Context) {
		var request struct {
			Destination string `json:"destination"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
			return
		}

		// Création d'un canal tamponné et d'un waitgroup pour les 3 méthodes de livraison.
		ch := make(chan string, 3)
		var wg sync.WaitGroup

		truck, _ := GetTransportMethod("truck")
		drone, _ := GetTransportMethod("drone")
		boat, _ := GetTransportMethod("boat")

		wg.Add(3)
		go func() { TrackDelivery(truck, request.Destination, ch); wg.Done() }()
		go func() { TrackDelivery(drone, request.Destination, ch); wg.Done() }()
		go func() { TrackDelivery(boat, request.Destination, ch); wg.Done() }()

		// Fermeture du canal après la fin des livraisons.
		go func() {
			wg.Wait()
			close(ch)
		}()

		results := []string{}
		for res := range ch {
			results = append(results, res)
		}

		c.JSON(http.StatusOK, gin.H{
			"destination": request.Destination,
			"results":     results,
		})
	})

	r.Run() // Lancement du serveur HTTP sur le port par défaut.
}
