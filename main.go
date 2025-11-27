package main

import (
	"borderlands_project/backend/router" // On importe le routeur du backend
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Le WaitGroup sert à dire au programme : "Attends que les 2 serveurs soient finis (jamais) avant de couper"
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("--------------------------------------------------")
	fmt.Println("🚀 GIOVANNI STREET ARMOURY - DÉMARRAGE GLOBAL")
	fmt.Println("--------------------------------------------------")

	// 1. LANCER LE BACKEND (Dans un "thread" séparé)
	go func() {
		defer wg.Done()
		// On lance le routeur API qu'on a déjà codé
		router.InitServer()
	}()

	// 2. LANCER LE FRONTEND (Dans un autre "thread")
	go func() {
		defer wg.Done()
		startFrontendServer()
	}()

	// On attend indéfiniment
	wg.Wait()
}

// Fonction pour lancer le site (HTML/CSS) sur le port 8081
func startFrontendServer() {
	// On sert le CSS (Chemin depuis la racine : frontend/stylecss)
	cssServer := http.FileServer(http.Dir("./frontend/stylecss"))
	http.Handle("/stylecss/", http.StripPrefix("/stylecss/", cssServer))

	// On sert les Templates HTML
	tmplServer := http.FileServer(http.Dir("./frontend/templates"))
	http.Handle("/", http.StripPrefix("/", tmplServer))

	// Petite pause pour laisser le temps au backend de s'afficher dans la console
	time.Sleep(1 * time.Second)

	fmt.Println("✅ FRONTEND (Site 8081) : Prêt")
	fmt.Println("👉 Accès Site : http://localhost:8081/index.html")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("❌ Erreur Frontend :", err)
	}
}
