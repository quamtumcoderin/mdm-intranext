package http

import (
	"encoding/json"
	"net/http"

	"mdm-intranext/internal/core/ports"
)

type DeviceHandler struct {
	svc ports.DeviceService
}

func NewDeviceHandler(svc ports.DeviceService) *DeviceHandler {
	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/devices", h.handleRegisterDevice)
	mux.HandleFunc("GET /api/v1/devices/{id}", h.handleGetDevice)
}

func (h *DeviceHandler) handleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := r.Body.Close(); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to close request body")
		return
	}

	device := req.mapToDomain()

	if err := h.svc.RegisterDevice(r.Context(), device); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, mapToResponse(device))
}

func (h *DeviceHandler) handleGetDevice(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	device, err := h.svc.GetDevice(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Device not found")
		return
	}

	respondWithJSON(w, http.StatusOK, mapToResponse(device))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		panic(err)
	}
}
