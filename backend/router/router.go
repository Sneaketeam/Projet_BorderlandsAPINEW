package router

import (
	"backend/controller"
	"fmt"
	"net/http"
)

// Middleware CORS : Indispensable pour que le Frontend (8081) puisse parler au Backend (8000)
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// On autorise le Frontend (http://localhost:8081)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Pour les cookies de session

		// Si le navigateur demande "Est-ce que je peux ?", on dit OUI tout de suite
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

func InitServer() {
	// 1. API (Données JSON)
	http.HandleFunc("/api/weapons", enableCORS(controller.GetWeapons))
	http.HandleFunc("/api/fav/toggle", enableCORS(controller.ToggleFavoriteAPI))

	// 2. AUTH (Login/Signup) via API
	http.HandleFunc("/login", enableCORS(controller.LoginHandler))
	http.HandleFunc("/signup", enableCORS(controller.SignupHandler))
	http.HandleFunc("/logout", enableCORS(controller.LogoutHandler))

	// 3. IMAGES (Servies par le backend)
	// Les images doivent être dans le dossier 'backend/images'
	imgServer := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", imgServer))

	fmt.Println("📡 BACKEND API démarré sur le port 8000")
	http.ListenAndServe(":8000", nil)
}
