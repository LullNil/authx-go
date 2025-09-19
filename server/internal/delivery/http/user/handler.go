package user

import (
	"log/slog"
	"net/http"

	"authx-go/domain/user"

	"github.com/LullNil/go-http-utils/httputils"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	userService user.Service
	log         *slog.Logger
	validator   *validator.Validate
}

// New returns a new user handler.
func New(userService user.Service, log *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		log:         log,
		validator:   validator.New(),
	}
}

// RegisterUser saves a new user to the database.
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.http.user.RegisterUser"

	// Decode request
	req, ok := httputils.DecodeRequest[user.RegisterUserRequest](w, r, h.log, op)
	if !ok {
		return
	}

	// Validate request
	if !httputils.ValidateRequest(w, r, h.log, op, req) {
		return
	}

	// Call service
	id, err := h.userService.RegisterUser(r.Context(), req)
	if err != nil {
		httputils.WriteHTTPError(w, h.log, op, err)
	}

	// Send successful response
	httputils.SendDataOK(w, r, h.log, op, id)
}

// LoginUser authenticates user and returns JWT stub
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.http.user.LoginUser"

	// Decode request
	req, ok := httputils.DecodeRequest[user.LoginRequest](w, r, h.log, op)
	if !ok {
		return
	}

	// Validate request
	if !httputils.ValidateRequest(w, r, h.log, op, req) {
		return
	}

	// Call service
	token, err := h.userService.LoginUser(r.Context(), req)
	if err != nil {
		httputils.WriteHTTPError(w, h.log, op, err)
	}

	// TODO: send JWT in cookies
	// Send successful response
	httputils.SendDataOK(w, r, h.log, op, token)
}

// GetUserInfo retrieves user info from the database.
func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	const op = "delivery.http.user.GetUserInfo"

	// TODO: get user id from JWT
	var uId int64 = 1

	// Call service
	user, err := h.userService.GetUserByID(r.Context(), uId)
	if err != nil {
		httputils.WriteHTTPError(w, h.log, op, err)
		return
	}

	// Send successful response
	httputils.SendDataOK(w, r, h.log, op, user)
}

// // GetUserByID retrieves an user by ID from the database.
// func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
// 	const op = "delivery.http.user.GetUserByID"

// 	// Get id from path parameter
// 	idParam := chi.URLParam(r, "id")
// 	if idParam == "" {
// 		http.Error(w, "missing id", http.StatusBadRequest)
// 		return
// 	}

// 	// Convert id to int64
// 	id, err := strconv.ParseInt(idParam, 10, 64)
// 	if err != nil {
// 		h.log.Error(op, "invalid id", err)
// 		http.Error(w, "id must be an integer", http.StatusBadRequest)
// 		return
// 	}

// 	// Call service
// 	user, err := h.userService.GetUserByID(r.Context(), id)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrNotFound) {
// 			http.Error(w, "user not found", http.StatusNotFound)
// 			return
// 		}

// 		h.log.Error(op, "service error", err)
// 		httputils.WriteServiceError(w, err)
// 		return
// 	}

// 	// Send successful response
// 	httputils.WriteJSON(w, http.StatusOK, user)
// }
