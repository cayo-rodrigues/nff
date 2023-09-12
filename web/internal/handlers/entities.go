package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/go-chi/chi/v5"
)

type EntitiesPage struct {
	tmpl *template.Template
}

func (page *EntitiesPage) Render(w http.ResponseWriter, r *http.Request) {
	dbpool := sql.GetDatabasePool()
	rows, _ := dbpool.Query(r.Context(), "SELECT * FROM entities")
	defer rows.Close()

	entities := []models.Entity{}

	for rows.Next() {
		entity := models.Entity{
			Address: &models.Address{},
		}
		err := rows.Scan(
			&entity.Id, &entity.Name, &entity.UserType, &entity.CpfCnpj, &entity.Ie, &entity.Email, &entity.Password,
			&entity.Address.PostalCode, &entity.Address.Neighborhood, &entity.Address.StreetType, &entity.Address.StreetName, &entity.Address.Number,
		)
		if err != nil {
			return
		}
		entities = append(entities, entity)
	}

	data := map[string]interface{}{
		"IsAuthenticated": true,
		"Entities":        entities,
	}

	page.tmpl.ExecuteTemplate(w, "layout", data)
}

func (page *EntitiesPage) CreateEntity(w http.ResponseWriter, r *http.Request) {
	// TODO validate data
	entity := models.NewEntityFromForm(r)

	passwordHash, err := utils.HashPassword(entity.Password)
	if err != nil {
		return
	}

	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
		r.Context(),
		"INSERT INTO entities (name, user_type, cpf_cnpj, ie, email, password, postal_code, neighborhood, street_type, street_name, number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, passwordHash,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
	)
	err = row.Scan(&entity.Id)
	if err != nil {
		return
	}

	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) GetEntityForm(w http.ResponseWriter, r *http.Request) {
	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
		r.Context(),
		"SELECT * FROM entities WHERE entities.id = $1",
		chi.URLParam(r, "id"),
	)

	entity := models.Entity{
		Address: &models.Address{},
	}
	err := row.Scan(
		&entity.Id, &entity.Name, &entity.UserType, &entity.CpfCnpj, &entity.Ie, &entity.Email, &entity.Password,
		&entity.Address.PostalCode, &entity.Address.Neighborhood, &entity.Address.StreetType, &entity.Address.StreetName, &entity.Address.Number,
	)
	if err != nil {
		return
	}
	data := map[string]interface{}{
		"Entity": entity,
	}
	page.tmpl.ExecuteTemplate(w, "entity-form", data)
}

func (page *EntitiesPage) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	entity := models.NewEntityFromForm(r)
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}
	entity.Id = entityId

	dbpool := sql.GetDatabasePool()
	_, err = dbpool.Exec(
		r.Context(),
		"UPDATE entities SET name = $1, user_type = $2, cpf_cnpj = $3, ie = $4, email = $5, password = $6, postal_code = $7, neighborhood = $8, street_type = $9, street_name = $10, number = $11 WHERE id = $12",
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
		entity.Id,
	)
	if err != nil {
		return
	}

	eventMsg := fmt.Sprintf("{\"entityUpdated\": \"%v\"}", entityId)
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)
	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	dbpool := sql.GetDatabasePool()
	entityId := chi.URLParam(r, "id")
	_, err := dbpool.Exec(
		r.Context(),
		"DELETE FROM entities WHERE id = $1",
		entityId,
	)
	if err != nil {
		return
	}

	eventMsg := fmt.Sprintf("{\"entityDeleted\": \"%v\"}", entityId)
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)
	page.tmpl.ExecuteTemplate(w, "entity-form", nil)
}

func (page *EntitiesPage) ParseTemplates() {
	page.tmpl = template.Must(template.ParseFiles(
		"internal/templates/layout.html", "internal/templates/entities.html", "internal/templates/entity-form.html",
	))
}
