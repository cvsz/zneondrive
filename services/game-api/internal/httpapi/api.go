package httpapi

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cvsz/zneondrive/services/game-api/internal/core"
	"github.com/cvsz/zneondrive/services/game-api/internal/store"
)

const gameplayTicketTTL = 60 * time.Second

type API struct {
	store         store.Store
	gameServerKey string
}

func New(s store.Store, gameServerKey string) http.Handler {
	api := &API{store: s, gameServerKey: gameServerKey}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.health)
	mux.HandleFunc("POST /v1/sessions/bootstrap", api.bootstrap)
	mux.HandleFunc("GET /v1/state", api.state)
	mux.HandleFunc("POST /v1/game-tickets", api.issueGameTicket)
	mux.HandleFunc("POST /v1/internal/game-tickets/redeem", api.redeemGameTicket)
	mux.HandleFunc("POST /v1/quests/{questID}/complete", api.completeQuest)
	mux.HandleFunc("POST /v1/vehicles/{vehicleID}/builds", api.reviseBuild)
	return mux
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Ping(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "dependency_unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type bootstrapRequest struct {
	ResumeKey string `json:"resume_key"`
}

type bootstrapResponse struct {
	SessionToken string        `json:"session_token"`
	ResumeKey    string        `json:"resume_key,omitempty"`
	Snapshot     core.Snapshot `json:"snapshot"`
}

func (a *API) bootstrap(w http.ResponseWriter, r *http.Request) {
	var req bootstrapRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	resumeKey := strings.TrimSpace(req.ResumeKey)
	generatedResumeKey := false
	if resumeKey == "" {
		var err error
		resumeKey, err = core.NewSecret()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "secret_generation_failed")
			return
		}
		generatedResumeKey = true
	}
	if len(resumeKey) < 32 {
		writeError(w, http.StatusBadRequest, "invalid_resume_key")
		return
	}

	snapshot, _, err := a.store.Bootstrap(r.Context(), core.HashSecret(resumeKey))
	if err != nil {
		writeStoreError(w, err)
		return
	}

	sessionToken, err := core.NewSecret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "secret_generation_failed")
		return
	}
	if err = a.store.CreateSession(
		r.Context(),
		snapshot.AccountID,
		core.HashSecret(sessionToken),
		time.Now().UTC().Add(24*time.Hour),
	); err != nil {
		writeStoreError(w, err)
		return
	}

	resp := bootstrapResponse{SessionToken: sessionToken, Snapshot: snapshot}
	if generatedResumeKey {
		resp.ResumeKey = resumeKey
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (a *API) state(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_session")
		return
	}
	snapshot, err := a.store.SnapshotBySession(r.Context(), core.HashSecret(token))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

type gameplayTicketResponse struct {
	Ticket           string `json:"ticket"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (a *API) issueGameTicket(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_session")
		return
	}

	ticket, err := core.NewSecret()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "secret_generation_failed")
		return
	}
	if err = a.store.IssueGameTicket(
		r.Context(),
		core.HashSecret(token),
		core.HashSecret(ticket),
		time.Now().UTC().Add(gameplayTicketTTL),
	); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, gameplayTicketResponse{
		Ticket:           ticket,
		ExpiresInSeconds: int(gameplayTicketTTL / time.Second),
	})
}

type redeemTicketRequest struct {
	Ticket string `json:"ticket"`
}

func (a *API) redeemGameTicket(w http.ResponseWriter, r *http.Request) {
	if !constantTimeSecretEqual(r.Header.Get("X-Game-Server-Key"), a.gameServerKey) {
		writeError(w, http.StatusUnauthorized, "unauthorized_game_server")
		return
	}

	var req redeemTicketRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	ticket := strings.TrimSpace(req.Ticket)
	if len(ticket) != 64 {
		writeError(w, http.StatusBadRequest, "invalid_game_ticket")
		return
	}

	snapshot, err := a.store.RedeemGameTicket(r.Context(), core.HashSecret(ticket))
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

type mutationRequest struct {
	OperationID string `json:"operation_id"`
}

type questResponse struct {
	Snapshot core.Snapshot      `json:"snapshot"`
	Receipt  core.RewardReceipt `json:"receipt"`
}

func (a *API) completeQuest(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_session")
		return
	}
	questID := r.PathValue("questID")
	if _, err := core.ParseQuestID(questID); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_quest_id")
		return
	}
	var req mutationRequest
	if err := decodeJSON(r, &req); err != nil || strings.TrimSpace(req.OperationID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_operation")
		return
	}
	snapshot, receipt, err := a.store.CompleteQuest(
		r.Context(),
		core.HashSecret(token),
		questID,
		strings.TrimSpace(req.OperationID),
	)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, questResponse{Snapshot: snapshot, Receipt: receipt})
}

type buildRequest struct {
	ExpectedRevision int      `json:"expected_revision"`
	PartIDs          []string `json:"part_ids"`
	OperationID      string   `json:"operation_id"`
}

func (a *API) reviseBuild(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing_session")
		return
	}
	var req buildRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if req.ExpectedRevision < 1 || strings.TrimSpace(req.OperationID) == "" {
		writeError(w, http.StatusBadRequest, "invalid_build_mutation")
		return
	}
	snapshot, err := a.store.ReviseBuild(
		r.Context(),
		core.HashSecret(token),
		r.PathValue("vehicleID"),
		req.ExpectedRevision,
		req.PartIDs,
		strings.TrimSpace(req.OperationID),
	)
	if err != nil {
		writeStoreError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

func bearerToken(r *http.Request) (string, bool) {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "
	if !strings.HasPrefix(value, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(value, prefix))
	return token, token != ""
}

func constantTimeSecretEqual(provided, expected string) bool {
	providedHash := sha256.Sum256([]byte(strings.TrimSpace(provided)))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(providedHash[:], expectedHash[:]) == 1
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, "unauthorized_session")
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found")
	case errors.Is(err, store.ErrConflict):
		writeError(w, http.StatusConflict, "state_conflict")
	case errors.Is(err, store.ErrOutOfOrder):
		writeError(w, http.StatusConflict, "quest_prerequisite_incomplete")
	case errors.Is(err, store.ErrOperationKey):
		writeError(w, http.StatusConflict, "operation_id_conflict")
	case errors.Is(err, core.ErrInvalidQuest), errors.Is(err, core.ErrInvalidParts):
		writeError(w, http.StatusBadRequest, "invalid_domain_input")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
