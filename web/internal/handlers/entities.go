package handlers

import (
	"html/template"
	"net/http"

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
	rows, _ := dbpool.Query(
		r.Context(),
		`SELECT
			entities.*,
			addresses.postal_code,
			addresses.neighborhood,
			addresses.street_type,
			addresses.street_name,
			addresses.number
		 FROM entities
		 JOIN addresses ON addresses.id = entities.id`,
	)
	defer rows.Close()

	entities := []models.Entity{}

	for rows.Next() {
		entity := models.Entity{
			Address: &models.Address{},
		}
		err := rows.Scan(
			&entity.Id, &entity.Name, &entity.UserType, &entity.CpfCnpj, &entity.Ie, &entity.Email, &entity.Password,
			&entity.Address.Id, &entity.Address.PostalCode, &entity.Address.Neighborhood, &entity.Address.StreetType, &entity.Address.StreetName, &entity.Address.Number,
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
		"INSERT INTO addresses (postal_code, neighborhood, street_type, street_name, number) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
	)
	err = row.Scan(&entity.Address.Id)
	if err != nil {
		return
	}
	row = dbpool.QueryRow(
		r.Context(),
		"INSERT INTO entities (name, user_type, cpf_cnpj, ie, email, password, address_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, passwordHash, entity.Address.Id,
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
		`SELECT
			entities.*,
			addresses.postal_code,
			addresses.neighborhood,
			addresses.street_type,
			addresses.street_name,
			addresses.number
		 FROM entities
		 JOIN addresses ON addresses.id = entities.id
		 WHERE entities.id = $1`,
		chi.URLParam(r, "id"),
	)

	entity := models.Entity{
		Address: &models.Address{},
	}
	err := row.Scan(
		&entity.Id, &entity.Name, &entity.UserType, &entity.CpfCnpj, &entity.Ie, &entity.Email, &entity.Password,
		&entity.Address.Id, &entity.Address.PostalCode, &entity.Address.Neighborhood, &entity.Address.StreetType, &entity.Address.StreetName, &entity.Address.Number,
	)
	if err != nil {
		return
	}
	data := map[string]interface{}{
		"Entity": entity,
	}
	page.tmpl.ExecuteTemplate(w, "entity-form", data)
}

func (page *EntitiesPage) ParseTemplates() {
	page.tmpl = template.Must(template.ParseFiles(
		"internal/templates/layout.html", "internal/templates/entities.html", "internal/templates/entity-form.html",
	))
}
