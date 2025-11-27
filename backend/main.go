package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Servir le CSS
	// Le dossier 'stylecss' doit être dans 'frontend/'
	cssServer := http.FileServer(http.Dir("./stylecss"))
	http.Handle("/stylecss/", http.StripPrefix("/stylecss/", cssServer))

	// 2. Servir les Pages HTML
	// Le dossier 'templates' doit être dans 'frontend/'
	// On dit que la racine "/" sert ce qu'il y a dans 'templates'
	fileServer := http.FileServer(http.Dir("./templates"))
	http.Handle("/", http.StripPrefix("/", fileServer))

	fmt.Println("🌍 FRONTEND (Site) démarré sur http://localhost:8081")
	fmt.Println("   -> Ouvre ton navigateur sur http://localhost:8081/index.html")

	http.ListenAndServe(":8081", nil)
}
