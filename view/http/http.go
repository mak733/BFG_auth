package http

import (
	"BFG_auth/session"
	"fmt"
	"github.com/swaggo/http-swagger"
	"net/http"
	"os"
)

type ViewHttp struct {
	Server  http.Server
	Manager *session.Manager
}

func (h *ViewHttp) StartServer(address string, sessionManager *session.Manager) error {
	h.Server = http.Server{Addr: address}
	h.Manager = sessionManager
	http.HandleFunc("/api", h.handleApiRequests)
	http.HandleFunc("/login", h.handleLogin)
	http.HandleFunc("/", h.handleLogin)

	http.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	err := h.Server.ListenAndServe()
	return err
}

// @Summary API Endpoint
// @Description Get API details or execute a command
// @ID get-api
// @Produce  html
// @Param command query string false "Command to execute"
// @Success 200 {string} string "Success"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api [get]
func (h *ViewHttp) handleApiRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Проверка токена из кук
	cookie, err := r.Cookie("token")
	if err != nil {
		// http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	token := cookie.Value
	// Ваша функция проверки токена (например, проверка подписи, проверка срока действия и т.д.)
	isValid, err := h.Manager.ValidateToken(token)
	if err != nil || !isValid {
		fmt.Printf("vaslidate is %t, err : %s\n", isValid, err.Error())
		// http.Error(w, "Invalid request method", http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := h.Manager.GetUser(token)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	command := r.URL.Query().Get("command")
	if command != "" {
		answer, err := h.Manager.ExecuteCommand(user, command)
		fmt.Printf(answer)
		if err != nil {
			fmt.Printf("execute %s err : %s\n", command, err)

			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				// handle error
				fmt.Printf("Error writing response: %v\n", err)
				return
			}
			http.Redirect(w, r, "/login", http.StatusForbidden)
		}

		_, err = w.Write([]byte(answer))
		if err != nil {
			// handle error
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
		return
	} else {
		html, err := os.ReadFile("view/http/html/api.html")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		_, err = w.Write(html)
		if err != nil {
			// handle error
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
	}
}

// @Summary Login Endpoint
// @Description Get login page or authenticate a user
// @ID login
// @Accept  json
// @Produce  json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {string} string "Success"
// @Failure 401 {string} string "Unauthorized"
// @Failure 400 {string} string "Bad Request"
// @Router /login [get]
// @Router /login [post]
func (h *ViewHttp) handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("handle login method is %s\n", r.Method)

	switch r.Method {
	case http.MethodGet:
		//чекни куку, проверь что она жива и если да то редирект на пи страницу
		html, err := os.ReadFile("view/http/html/login.html")
		if err != nil {
			http.Error(w, "File reading error", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(html)
		if err != nil {
			// handle error
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		token, err := h.Manager.Authenticate(username, password)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}

		err = h.Manager.Authorize(username, token)
		if err != nil {
			fmt.Printf("%s", err)

			http.Error(w, "Authorization failed", http.StatusUnauthorized)
			return
		}

		// Установка куки с токеном
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
			// Используйте защищённые параметры для куки
			HttpOnly: true,
			Secure:   true, // если вы используете HTTPS
			Path:     "/",
		})

		http.Redirect(w, r, "/api", http.StatusSeeOther)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}
