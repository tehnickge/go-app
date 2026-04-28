package httpdelivery

import (
	"net/http"

	"backend-app/internal/delivery/http/controllers"
	httpmiddleware "backend-app/internal/delivery/http/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, deps Dependencies) {
	// Healthcheck - no auth.
	r.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
	})

	authCtrl := controllers.NewAuthController(deps.Auth)
	usersCtrl := controllers.NewUsersController(deps.Users)
	rolesCtrl := controllers.NewRolesController(deps.Roles)
	groupsCtrl := controllers.NewGroupsController(deps.Groups)
	relCtrl := controllers.NewRelationsController(deps.Relations)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authCtrl.Register)
		r.Post("/login", authCtrl.Login)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(httpmiddleware.JWTAuthMiddleware(deps.JWTSecret))
		r.Get("/", usersCtrl.List)
		r.Post("/", usersCtrl.Create)
		r.Get("/me", usersCtrl.GetMe)
		r.Get("/{id}", usersCtrl.GetByID)
		r.Patch("/{id}", usersCtrl.Update)
		r.Delete("/{id}", usersCtrl.Delete)

		r.Route("/{id}/roles", func(r chi.Router) {
			r.Get("/", relCtrl.GetUserRoles)
			r.Post("/{roleId}", relCtrl.AssignRoleToUser)
			r.Delete("/{roleId}", relCtrl.RemoveRoleFromUser)
		})

		r.Route("/{id}/groups", func(r chi.Router) {
			r.Get("/", relCtrl.GetUserGroups)
		})
	})

	r.Route("/roles", func(r chi.Router) {
		r.Use(httpmiddleware.JWTAuthMiddleware(deps.JWTSecret))
		r.Get("/", rolesCtrl.List)
		r.Post("/", rolesCtrl.Create)
		r.Get("/{id}", rolesCtrl.GetByID)
		r.Patch("/{id}", rolesCtrl.Update)
		r.Delete("/{id}", rolesCtrl.SoftDelete)
	})

	r.Route("/groups", func(r chi.Router) {
		r.Use(httpmiddleware.JWTAuthMiddleware(deps.JWTSecret))
		r.Get("/", groupsCtrl.List)
		r.Post("/", groupsCtrl.Create)
		r.Get("/{id}", groupsCtrl.GetByID)
		r.Patch("/{id}", groupsCtrl.Update)
		r.Delete("/{id}", groupsCtrl.SoftDelete)

		r.Post("/{groupId}/users/{userId}", relCtrl.AddUserToGroup)
		r.Delete("/{groupId}/users/{userId}", relCtrl.RemoveUserFromGroup)
	})
}

