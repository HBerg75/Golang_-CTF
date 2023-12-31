package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

type User struct {
	Name   string `json:"User"`
	Secret string `json:"Secret,omitempty"`
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
			if string(bodyPing) == "pong" {
				user := User{Name: "Yanis"}
				userJSON, err := json.Marshal(user)
				if err != nil {
					fmt.Printf("Erreur lors de la conversion de l'utilisateur en JSON : %v\n", err)
					return
				}

				// Faire une requête HTTP POST pour /signup
				signupURL := fmt.Sprintf("http://%s:%d/signup", serverIP, port)
				respPost, err := http.Post(signupURL, "application/json", bytes.NewBuffer(userJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", signupURL, err)
					return
				}
				defer respPost.Body.Close()

				// Vérifiez le statut de la réponse
				if respPost.StatusCode == http.StatusOK {
					bodyPost, _ := io.ReadAll(respPost.Body)
					fmt.Printf("Port %d accessible - POST Response for /signup: %s\n", port, bodyPost)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respPost.Status)
				}

				// Faire une requête HTTP POST pour /check
				checkURL := fmt.Sprintf("http://%s:%d/check", serverIP, port)
				respCheck, err := http.Post(checkURL, "application/json", bytes.NewBuffer(userJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", checkURL, err)
					return
				}
				defer respCheck.Body.Close()

				// Vérifiez le statut de la réponse
				if respCheck.StatusCode == http.StatusOK {
					bodyCheck, _ := io.ReadAll(respCheck.Body)
					fmt.Printf("Port %d accessible - POST Response for /check: %s\n", port, bodyCheck)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respCheck.Status)
				}

				// Boucle pour la requête POST /getUserSecret
				for {
					getUserSecretURL := fmt.Sprintf("http://%s:%d/getUserSecret", serverIP, port)
					respGetUserSecret, err := http.Post(getUserSecretURL, "application/json", bytes.NewBuffer(userJSON))
					if err != nil {
						fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", getUserSecretURL, err)
						return
					}
					defer respGetUserSecret.Body.Close()
			
					// Vérifiez le statut de la réponse
					if respGetUserSecret.StatusCode == http.StatusOK {
						bodyGetUserSecret, _ := io.ReadAll(respGetUserSecret.Body)
						if string(bodyGetUserSecret) != "Really don't feel like working today huh..." {
							fmt.Printf("Port %d accessible - POST Response for /getUserSecret: %s\n", port, bodyGetUserSecret)
							userSecret := strings.TrimSpace(strings.Split(string(bodyGetUserSecret), "User secret: ")[1])  // Cette ligne a été modifiée
							user.Secret = userSecret
							fmt.Printf("%s\n", userSecret)
			
							break
						}
					} else {
						fmt.Printf("Erreur, statut de la réponse : %s\n", respGetUserSecret.Status)
						break // Sortir de la boucle si le statut n'est pas OK
					}
				}
				// Recréer userJSON avec le Secret inclus
				userLevelJSON, err := json.Marshal(user)  // Cette ligne a été modifiée
				if err != nil {
					fmt.Printf("Erreur lors de la conversion de l'utilisateur en JSON : %v\n", err)
					return
				}

				// Faire une requête HTTP POST pour /getUserLevel
				getUserLevelURL := fmt.Sprintf("http://%s:%d/getUserLevel", serverIP, port)
				respGetUserLevel, err := http.Post(getUserLevelURL, "application/json", bytes.NewBuffer(userLevelJSON))
				fmt.Println(bytes.NewBuffer(userLevelJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", getUserLevelURL, err)
					return
				}
				defer respGetUserLevel.Body.Close()

				// Vérifiez le statut de la réponse
				if respGetUserLevel.StatusCode == http.StatusOK {
					bodyGetUserLevel, _ := io.ReadAll(respGetUserLevel.Body)
					fmt.Printf("Port %d accessible - POST Response for /getUserLevel: %s\n", port, bodyGetUserLevel)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respGetUserLevel.Status)
				}
				// Faire une requête HTTP POST pour /getUserPoints
				getUserPointsURL := fmt.Sprintf("http://%s:%d/getUserPoints", serverIP, port)
				respGetUserPoints, err := http.Post(getUserPointsURL, "application/json", bytes.NewBuffer(userLevelJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", getUserPointsURL, err)
					return
				}
				defer respGetUserPoints.Body.Close()

				// Vérifiez le statut de la réponse
				if respGetUserPoints.StatusCode == http.StatusOK {
					bodyGetUserPoints, _ := io.ReadAll(respGetUserPoints.Body)
					fmt.Printf("Port %d accessible - POST Response for /getUserPoints: %s\n", port, bodyGetUserPoints)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respGetUserPoints.Status)
				}

				// Faire une requête HTTP POST pour /iNeedAHint
				iNeedAHintURL := fmt.Sprintf("http://%s:%d/iNeedAHint", serverIP, port)
				respINeedAHint, err := http.Post(iNeedAHintURL, "application/json", bytes.NewBuffer(userLevelJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", iNeedAHintURL, err)
					return
				}
				defer respINeedAHint.Body.Close()

				// Vérifiez le statut de la réponse
				if respINeedAHint.StatusCode == http.StatusOK {
					bodyINeedAHint, _ := io.ReadAll(respINeedAHint.Body)
					fmt.Printf("Port %d accessible - POST Response for /iNeedAHint: %s\n", port, bodyINeedAHint)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respINeedAHint.Status)
				}
				
				// Faire une requête HTTP POST pour /enterChallenge
				enterChallengeURL := fmt.Sprintf("http://%s:%d/enterChallenge", serverIP, port)
				respEnterChallenge, err := http.Post(enterChallengeURL, "application/json", bytes.NewBuffer(userLevelJSON))
				if err != nil {
					fmt.Printf("Erreur lors de la requête POST vers %s: %v\n", enterChallengeURL, err)
					return
				}
				defer respEnterChallenge.Body.Close()

				// Vérifiez le statut de la réponse
				if respEnterChallenge.StatusCode == http.StatusOK {
					bodyEnterChallenge, _ := io.ReadAll(respEnterChallenge.Body)
					fmt.Printf("Port %d accessible - POST Response for /enterChallenge: %s\n", port, bodyEnterChallenge)
				} else {
					fmt.Printf("Erreur, statut de la réponse : %s\n", respEnterChallenge.Status)
				}

			}
		}
	}
}

func main() {
	serverIP := "10.49.122.144"
	minPort := 1024
	maxPort := 8192

	var wg sync.WaitGroup

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		go testPort(serverIP, port, &wg)
	}

	// Attendre que toutes les goroutines se terminent
	wg.Wait()
}
