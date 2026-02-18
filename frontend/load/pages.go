package load

import (
	"html/template"
	"net/http"
)

type Pages struct {
	tmpl    *template.Template
	pathCss string
}

func (p *Pages) NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		r.Header.Del("If-Modified-Since")
		r.Header.Del("If-None-Match")
		next.ServeHTTP(w, r)
	})
}

func (p *Pages) Admin(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "admin.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) Login(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "home.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "unauthorized.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) ConfirmSignup(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "confirm_signup.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) ResetPass(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "resetpass.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) NewPass(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "newpass.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) Products(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "products.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (p *Pages) Ws(w http.ResponseWriter, r *http.Request) {
	if err := p.tmpl.ExecuteTemplate(w, "websocket.html", nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}
