package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// --- CONSTANTES ---
const FRONTEND_URL = "http://localhost:8081"

// --- STRUCTURES JSON ---
type Weapon struct {
	ID           int    `json:"id"`
	Category     string `json:"category"`
	Name         string `json:"name"`
	Manufacturer string `json:"manufacturer"`
	Rarity       string `json:"rarity"`
	FlavorText   string `json:"flavor_text"`
	Details      string `json:"details"`
	Source       string `json:"source"`
	ImageURL     string `json:"image_url"`
	IsFavorite   bool   `json:"is_favorite"`
}

// --- BDD ---
func dbConn() (db *sql.DB) {
	// Connexion ALWAYSDATA (Cloud)
	dbUser := "443067"
	dbPass := "giogio220706"
	dbHost := "mysql-borderlandsapi.alwaysdata.net"
	dbName := "borderlandsapi_database"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getUserID(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	return id, err
}

// ==========================================
// 1. API DATA (JSON)
// ==========================================

// Récupère la liste des armes en JSON (avec filtres)
func GetWeapons(w http.ResponseWriter, r *http.Request) {
	// Gestion Session (Pour savoir si c'est favori)
	var isLoggedIn bool
	var userID int
	cookie, err := r.Cookie("session_token")

	db := dbConn()
	defer db.Close()

	if err == nil {
		isLoggedIn = true
		userID, _ = getUserID(db, cookie.Value)
	}

	// Filtres URL
	cat := r.URL.Query().Get("category")
	rar := r.URL.Query().Get("rarity")
	nam := r.URL.Query().Get("name")
	man := r.URL.Query().Get("manufacturer") // Constructeur

	// Construction SQL
	query := "SELECT * FROM weapons WHERE 1=1"
	var args []interface{}

	if cat != "" {
		query += " AND category = ?"
		args = append(args, cat)
	}
	if rar != "" {
		query += " AND rarity = ?"
		args = append(args, rar)
	}
	if man != "" {
		query += " AND manufacturer = ?"
		args = append(args, man)
	}
	if nam != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+nam+"%")
	}

	rows, _ := db.Query(query, args...)
	defer rows.Close()

	var weapons []Weapon
	for rows.Next() {
		var wpn Weapon
		rows.Scan(&wpn.ID, &wpn.Category, &wpn.Name, &wpn.Manufacturer, &wpn.Rarity, &wpn.FlavorText, &wpn.Details, &wpn.Source, &wpn.ImageURL)

		if isLoggedIn {
			var count int
			db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = ? AND weapon_id = ?", userID, wpn.ID).Scan(&count)
			if count > 0 {
				wpn.IsFavorite = true
			}
		}
		weapons = append(weapons, wpn)
	}

	// Envoi JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weapons)
}

// Toggle Favoris (AJAX JSON)
func ToggleFavoriteAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error": "nologin"}`)
		return
	}

	db := dbConn()
	defer db.Close()

	weaponID := r.URL.Query().Get("id")
	userID, err := getUserID(db, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM favorites WHERE user_id = ? AND weapon_id = ?", userID, weaponID).Scan(&count)

	var weaponName string
	db.QueryRow("SELECT name FROM weapons WHERE id = ?", weaponID).Scan(&weaponName)

	if count > 0 {
		db.Exec("DELETE FROM favorites WHERE user_id = ? AND weapon_id = ?", userID, weaponID)
		fmt.Fprintf(w, `{"status": "removed", "name": "%s"}`, weaponName)
	} else {
		db.Exec("INSERT INTO favorites (user_id, weapon_id) VALUES (?, ?)", userID, weaponID)
		fmt.Fprintf(w, `{"status": "added", "name": "%s"}`, weaponName)
	}
}

// ==========================================
// 2. AUTHENTIFICATION (Redirections vers Frontend)
// ==========================================

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db := dbConn()
	defer db.Close()

	var dbPass string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPass)

	if err != nil || dbPass != password {
		// Erreur : On renvoie vers la page Login du Frontend
		http.Redirect(w, r, FRONTEND_URL+"/login.html?error=wrong", http.StatusSeeOther)
		return
	}
	createCookie(w, username)
	http.Redirect(w, r, FRONTEND_URL, http.StatusSeeOther) // Vers l'accueil
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db := dbConn()
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		http.Redirect(w, r, FRONTEND_URL+"/login.html?error=exists", http.StatusSeeOther)
		return
	}
	createCookie(w, username)
	http.Redirect(w, r, FRONTEND_URL, http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: "", Expires: time.Unix(0, 0), Path: "/",
	})
	http.Redirect(w, r, FRONTEND_URL, http.StatusSeeOther)
}

func createCookie(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name: "session_token", Value: username, Expires: time.Now().Add(24 * time.Hour), Path: "/",
	})
}
