package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// User Routes
	mux.Get("/user/signup", dynamicMiddleware.Append(app.onlyGuestUsers).ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.Append(app.onlyGuestUsers).ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.Append(app.onlyGuestUsers).ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.Append(app.onlyGuestUsers).ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
	mux.Post("/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

	mux.Get("/ping", http.HandlerFunc(ping))

	apiMiddleware := alice.New(app.validateApiClient)
	mux.Post("/api/snippet/add", apiMiddleware.ThenFunc(app.createSnippetFromApi))
 	//mux.Post("/api/snippet/add", apiMiddleware.Then(http.HandlerFunc(app.createSnippetFromApi)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

