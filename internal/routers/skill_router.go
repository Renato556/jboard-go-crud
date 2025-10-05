package routers

import (
	"jboard-go-crud/internal/controllers"
	"net/http"
)

func NewSkillsController(skillHandler *controllers.SkillHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/skills", skillHandler.GetAllSkills)
	mux.HandleFunc("POST /v1/skills", skillHandler.AddSkill)
	mux.HandleFunc("PUT /v1/skills", skillHandler.RemoveSkill)
	mux.HandleFunc("DELETE /v1/skills", skillHandler.DeleteUserSkills)

	return mux
}
