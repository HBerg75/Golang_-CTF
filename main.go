package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
)

type User struct {
	Name string `json:"User"`
}

func testPort(serverIP string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", serverIP, port)

	// Tentative de connexion au serveur
	conn, err := net.Dial("tcp", address)
	if err == nil {
		conn.Close()

		// Faire une requête HTTP GET pour /ping
		pingURL := fmt.Sprintf("http://%s:%d/ping", serverIP, port)
		respPing, err := http.Get(pingURL)
		if err == nil {
			defer respPing.Body.Close()
			bodyPing, _ := io.ReadAll(respPing.Body)
			fmt.Printf("Port %d accessible - GET Response for /ping: %s\n", port, bodyPing)

			// Vérification de la réponse "pong"
			if string(bodyPing) == `{"message":"pong"}` {
				// Faire une requête HTTP POST pour /signup
				signupURL := fmt.Sprintf("http://%s:%d/signup", serverIP, port)
				user := User{Name: "Yanis"}
				userJSON, err := json.Marshal(user)
				if err != nil {
					fmt.Printf("Erreur lors de la conversion de l'utilisateur en JSON : %v\n", err)
					return
				}

				respPost, err := http.Post(signupURL, "application/json", bytes.NewBuffer(userJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", signupURL, err)
					return
				}
				defer respPost.Body.Close()
				bodyPost, _ := io.ReadAll(respPost.Body)
				fmt.Printf("Port %d accessible - POST Response for /signup: %s\n", port, bodyPost)
			}
		}
	}
}

func main() {
	serverIP := "10.49.122.144"
	minPort := 1
	maxPort := 8000

	var wg sync.WaitGroup

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		go testPort(serverIP, port, &wg)
	}

	// Attendre que toutes les goroutines se terminent
	wg.Wait()
}












// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// const (
// 	startPort   = 1025
// 	endPort     = 65535
// 	timeout     = 2 * time.Second
// 	serverAddr  = "10.49.122.144"
// 	httpTimeout = 5 * time.Second
// )

// type User struct {
// 	Name string `json:"User"`
// }

// type PingResponse struct {
// 	Message string `json:"message"`
// }

// func scanPort(ip string, port int, wg *sync.WaitGroup, openPorts chan int) {
// 	defer wg.Done()

// 	address := fmt.Sprintf("%s:%d", ip, port)
// 	conn, err := net.DialTimeout("tcp", address, timeout)

// 	if err == nil {
// 		conn.Close()
// 		openPorts <- port
// 	}
// }

// func makeHTTPRequest(port int) {
// 	url := fmt.Sprintf("http://%s:%d/ping", serverAddr, port)
// 	client := &http.Client{
// 		Timeout: httpTimeout,
// 	}
// 	resp, err := client.Get(url)
// 	if err != nil {
// 		fmt.Printf("Erreur lors de la requête GET vers %s: %v\n", url, err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusOK {
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("Erreur lors de la lecture du corps de la réponse : %v\n", err)
// 			return
// 		}

// 		fmt.Printf("Requête GET réussie vers %s sur le port %d avec le corps de la réponse : %s\n", url, port, body)

// 		var pingResponse PingResponse
// 		err = json.Unmarshal(body, &pingResponse)
// 		if err != nil {
// 			fmt.Printf("Erreur lors de la déserialization du corps de la réponse : %v\n", err)
// 			return
// 		}

// 		if pingResponse.Message == "pong" {
// 			signupURL := fmt.Sprintf("http://%s:%d/signup", serverAddr, port)
// 			user := User{Name: "Yanis"}
// 			userJSON, err := json.Marshal(user)
// 			if err != nil {
// 				fmt.Printf("Erreur lors de la conversion de l'utilisateur en JSON : %v\n", err)
// 				return
// 			}

// 			resp, err := client.Post(signupURL, "application/json", bytes.NewBuffer(userJSON))
// 			if err != nil {
// 				fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", signupURL, err)
// 				return
// 			}
// 			defer resp.Body.Close()

// 			if resp.StatusCode == http.StatusOK {
// 				postBody, err := io.ReadAll(resp.Body)
// 				if err != nil {
// 					fmt.Printf("Erreur lors de la lecture du corps de la réponse POST : %v\n", err)
// 					return
// 				}
// 				fmt.Printf("Requête POST réussie vers %s avec le corps de la réponse : %s\n", signupURL, postBody)
// 			} else {
// 				fmt.Printf("Réponse inattendue de la requête POST vers %s. Statut: %s\n", signupURL, resp.Status)
// 			}
// 		}
// 	} else {
// 		fmt.Printf("Réponse inattendue de la requête GET vers %s. Statut: %s\n", url, resp.Status)
// 	}
// }

// func main() {
// 	var wg sync.WaitGroup
// 	var httpWG sync.WaitGroup
// 	openPorts := make(chan int)

// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			for port := range openPorts {
// 				httpWG.Add(1)
// 				makeHTTPRequest(port)
// 				httpWG.Done()
// 			}
// 		}()
// 	}

// 	for port := startPort; port <= endPort; port++ {
// 		wg.Add(1)
// 		go scanPort(serverAddr, port, &wg, openPorts)
// 	}

// 	wg.Wait()

// 	httpWG.Wait()

// 	close(openPorts)
// }

