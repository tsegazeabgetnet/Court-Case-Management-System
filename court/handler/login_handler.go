package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	entity "github.com/Surafeljava/Court-Case-Management-System/Entity"
	"github.com/Surafeljava/Court-Case-Management-System/caseUse"
	"github.com/Surafeljava/Court-Case-Management-System/form"
	"github.com/Surafeljava/Court-Case-Management-System/rtoken"
	"github.com/Surafeljava/Court-Case-Management-System/session"
)

type LoginHandler struct {
	tmpl           *template.Template
	loginSrv       caseUse.LoginService
	sessionService caseUse.SessionService
	userSess       *entity.Session
	csrfSignKey    []byte
}

func NewLoginHandler(t *template.Template, usrServ caseUse.LoginService,
	sessServ caseUse.SessionService, usrSess *entity.Session, csKey []byte) *LoginHandler {
	return &LoginHandler{tmpl: t, loginSrv: usrServ, sessionService: sessServ, userSess: usrSess, csrfSignKey: csKey}
}

type contextKey string

var ctxUserSessionKey = contextKey("signed_in_user_session")

// Authenticated checks if a user is authenticated to access a given route
func (lh *LoginHandler) AuthenticatedUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ok := lh.loggedIn(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxUserSessionKey, lh.userSess)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}

// Authorized checks if a user has proper authority to access a give route
// func (uh *LoginHandler) Authorized(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if uh.loggedInUser == nil {
// 			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 			return
// 		}
// 		roles, errs := uh.userService.UserRoles(uh.loggedInUser)
// 		if len(errs) > 0 {
// 			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 			return
// 		}

// 		for _, role := range roles {
// 			permitted := permission.HasPermission(r.URL.Path, role.Name, r.Method)
// 			if !permitted {
// 				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 				return
// 			}
// 		}
// 		if r.Method == http.MethodPost {
// 			ok, err := rtoken.ValidCSRF(r.FormValue("_csrf"), uh.csrfSignKey)
// 			if !ok || (err != nil) {
// 				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
// 				return
// 			}
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

func (lh *LoginHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	token, err := rtoken.CSRFToken(lh.csrfSignKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	if r.Method == http.MethodGet {
		loginForm := struct {
			Values  url.Values
			VErrors form.ValidationErrors
			CSRF    string
		}{
			Values:  nil,
			VErrors: nil,
			CSRF:    token,
		}
		lh.tmpl.ExecuteTemplate(w, "login.layout", loginForm)
		return
	}
	if r.Method == http.MethodPost {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user_id := r.FormValue("user_id")
		user_pwd := r.FormValue("user_pwd")

		//Checking the type of user trying to login
		who := CheckWho(user_id)

		error_message := entity.SuccessMessage{Status: "Error", Message: "Wrong ID or Password Try again"}
		//success_message := entity.SuccessMessage{Status: "Success", Message: "Login Success!"}

		if who == 0 {
			adm, err := lh.loginSrv.CheckAdmin(user_id, user_pwd)
			if adm != nil {

				//Creating admin session
				claims := rtoken.Claims(user_id, lh.userSess.Expires)
				session.Create(claims, lh.userSess.UUID, lh.userSess.SigningKey, w)
				newSess, errs := lh.sessionService.StoreSession(lh.userSess)
				if len(errs) > 0 || newSess == nil {
					lh.tmpl.ExecuteTemplate(w, "login.layout", nil)
					return
				}

				lh.tmpl.ExecuteTemplate(w, "admin.home.layout", adm)
			} else if len(err) > 0 {
				lh.tmpl.ExecuteTemplate(w, "login.layout", error_message)
			}
		} else if who == 1 {
			jud, err := lh.loginSrv.CheckJudge(user_id, user_pwd)
			if jud != nil {

				//Creating judge session
				claims := rtoken.Claims(user_id, lh.userSess.Expires)
				session.Create(claims, lh.userSess.UUID, lh.userSess.SigningKey, w)
				newSess, errs := lh.sessionService.StoreSession(lh.userSess)
				if len(errs) > 0 || newSess == nil {
					lh.tmpl.ExecuteTemplate(w, "login.layout", nil)
					return
				}
				lh.tmpl.ExecuteTemplate(w, "judge.home.layout", jud)
			} else if len(err) > 0 {
				lh.tmpl.ExecuteTemplate(w, "login.layout", error_message)
			}
		} else if who == 2 {
			opp, err := lh.loginSrv.CheckOpponent(user_id, user_pwd)
			if opp != nil {

				//Creating opponent session
				claims := rtoken.Claims(user_id, lh.userSess.Expires)
				session.Create(claims, lh.userSess.UUID, lh.userSess.SigningKey, w)
				newSess, errs := lh.sessionService.StoreSession(lh.userSess)
				if len(errs) > 0 || newSess == nil {
					lh.tmpl.ExecuteTemplate(w, "login.layout", nil)
					return
				}

				lh.tmpl.ExecuteTemplate(w, "opponent.home.layout", opp)
			} else if len(err) > 0 {
				lh.tmpl.ExecuteTemplate(w, "login.layout", error_message)
			}
		} else {
			lh.tmpl.ExecuteTemplate(w, "login.layout", error_message)
		}

	} else {
		lh.tmpl.ExecuteTemplate(w, "login.layout", nil)
	}

}

// Logout hanldes the POST /logout requests
func (uh *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userSess, er := r.Context().Value(ctxUserSessionKey).(entity.Session)
	if !er {
		fmt.Println("No user session registered with this Ctxusersessiokey!!")
	}
	fmt.Println("1 >>>>>>>>>>>>>>>>>>>")
	fmt.Println(userSess.SigningKey)
	// fmt.Println(userSess.UUID)
	fmt.Println("2 >>>>>>>>>>>>>>>>>>>")
	session.Remove(userSess.UUID, w)
	uh.sessionService.DeleteSession(userSess.UUID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (uh *LoginHandler) loggedIn(r *http.Request) bool {
	if uh.userSess == nil {
		return false
	}
	userSess := uh.userSess
	c, err := r.Cookie(userSess.UUID)
	if err != nil {
		return false
	}
	ok, err := session.Valid(c.Value, userSess.SigningKey)
	if !ok || (err != nil) {
		return false
	}
	return true
}

func CheckWho(id string) int {
	check := id[0:2]
	fmt.Println(check)
	if check == "AD" {
		return 0
	} else if check == "JU" {
		return 1
	} else if check == "OP" {
		return 2
	}
	return -1
}
