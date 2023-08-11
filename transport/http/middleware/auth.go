package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/oauth"
	"github.com/evermos/boilerplate-go/transport/http/response"
)

type ResponseValidate struct {
    Data shared.Claims `json:"data"`
}

type ClaimsKey string
type Authentication struct {
	db *infras.MySQLConn
	config *configs.Config
}

const (
	HeaderAuthorization = "Authorization"
)

func ProvideAuthentication(db *infras.MySQLConn, conf *configs.Config) *Authentication {
	return &Authentication{
		db: db,
		config: conf,
	}
}

func (a *Authentication) ClientCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) ClientCredentialWithQueryParameter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		token := params.Get("token")
		tokenType := params.Get("token_type")
		accessToken := tokenType + " " + token

		auth := oauth.New(a.db.Read, oauth.Config{})
		parseToken, err := auth.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *Authentication) Password(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		token := oauth.New(a.db.Read, oauth.Config{})

		parseToken, err := token.ParseWithAccessToken(accessToken)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyExpireIn() {
			response.WithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !parseToken.VerifyUserLoggedIn() {
			response.WithMessage(w, http.StatusUnauthorized, oauth.ErrorInvalidPassword)
			return
		}

		next.ServeHTTP(w, r)
	})
}



func (a *Authentication) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(HeaderAuthorization)
		if accessToken == "" {
			response.WithMessage(w, http.StatusUnauthorized, "No Header")
			return
		}

		accessToken = accessToken[(len("Bearer ")):]
		client := http.Client{}
		req , err := http.NewRequest("GET",configs.Get().App.AuthServiceBaseURL + configs.Get().App.AuthServiceValidatePath, nil)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		req.Header.Set("Authorization", "Bearer " +  accessToken)

		res , err := client.Do(req)
		if err != nil {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if res.StatusCode != 200 {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		decoder := json.NewDecoder(res.Body)
		var requestFormat ResponseValidate
		err = decoder.Decode(&requestFormat)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), requestFormat.Data)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}


func (a *Authentication) RoleAdminCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey("claims")).(shared.Claims)
		if !ok {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return

		}

		if claims.Role != "admin" {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


