package tenant

import (
	"encoding/json"
	"net/http"

	"github.com/ShalArl/trip-manager/backend/shared/authclient"
	"github.com/ShalArl/trip-manager/backend/shared/tenantdb"
	"github.com/ShalArl/trip-manager/backend/users/repository"
)

type MemberResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func ListMembersHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		if tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}
		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		users, err := repo.ListByTenant(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		var members []MemberResponse
		for _, u := range users {
			members = append(members, MemberResponse{
				ID: u.ID, Email: u.Email, Name: u.Name, Role: u.Role,
			})
		}
		if members == nil {
			members = []MemberResponse{}
		}
		respondJSON(w, http.StatusOK, members)
	}
}

func RemoveMemberHandler(repo repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		role := authclient.GetUserRole(r)
		if tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}
		if role != "tenant_owner" && role != "platform_admin" {
			respondError(w, http.StatusForbidden, "permission denied")
			return
		}
		userID := r.PathValue("userId")
		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		if err := repo.RemoveFromTenant(ctx, userID); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type CreateInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func CreateInvitationHandler(invRepo InvitationRepository, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		role := authclient.GetUserRole(r)
		if tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}
		if role != "tenant_owner" && role != "tenant_admin" && role != "platform_admin" {
			respondError(w, http.StatusForbidden, "permission denied")
			return
		}
		firebaseUID, _ := authclient.GetUserID(r)

		var req CreateInvitationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if req.Email == "" {
			respondError(w, http.StatusBadRequest, "email is required")
			return
		}
		if req.Role == "" {
			req.Role = "tenant_member"
		}

		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		inv, err := invRepo.Create(ctx, tenantID, req.Email, req.Role, firebaseUID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondJSON(w, http.StatusCreated, map[string]interface{}{
			"id":         inv.ID,
			"email":      inv.Email,
			"role":       inv.Role,
			"inviteLink": baseURL + "/join?token=" + inv.Token,
			"expiresAt":  inv.ExpiresAt,
		})
	}
}

func ListInvitationsHandler(invRepo InvitationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		if tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}
		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		invitations, err := invRepo.ListByTenant(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if invitations == nil {
			invitations = []*Invitation{}
		}
		respondJSON(w, http.StatusOK, invitations)
	}
}

func DeleteInvitationHandler(invRepo InvitationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := authclient.GetTenantID(r)
		role := authclient.GetUserRole(r)
		if tenantID == "default" {
			respondError(w, http.StatusNotFound, "no tenant found")
			return
		}
		if role != "tenant_owner" && role != "tenant_admin" && role != "platform_admin" {
			respondError(w, http.StatusForbidden, "permission denied")
			return
		}
		invID := r.PathValue("invitationId")
		ctx := tenantdb.WithTenantID(r.Context(), tenantID)
		if err := invRepo.Delete(ctx, invID); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
