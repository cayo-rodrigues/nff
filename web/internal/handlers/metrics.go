package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type MetricsPage struct {
	service       interfaces.MetricsService
	entityService interfaces.EntityService
	siareBGWorker interfaces.SiareBGWorker
}

func NewMetricsPage(service interfaces.MetricsService, entityService interfaces.EntityService, siareBGWorker interfaces.SiareBGWorker) *MetricsPage {
	return &MetricsPage{
		entityService: entityService,
		service:       service,
		siareBGWorker: siareBGWorker,
	}
}

type MetricsPageData struct {
	IsAuthenticated  bool
	GeneralError     string
	FormMsg          string
	FormSelectFields *models.MetricsFormSelectFields
	MetricsQuery     *models.MetricsQuery
	QueriesHistory   []*models.MetricsQuery
	ResourceName     string
}

func (p *MetricsPage) NewEmptyData() *MetricsPageData {
	return &MetricsPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.MetricsFormSelectFields{
			Entities: []*models.Entity{},
		},
		ResourceName: "metrics",
	}
}

func (p *MetricsPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	pageData.MetricsQuery = models.NewEmptyMetricsQuery()
	userID := c.Locals("UserID").(int)

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("metrics", pageData, "layouts/base")
	}

	pageData.FormSelectFields.Entities = entities

	// get the latest 10 metricsHistory
	metricsHistory, err := p.service.ListMetrics(c.Context(), userID)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("metrics", pageData, "layouts/base")
	}

	pageData.QueriesHistory = metricsHistory

	return c.Render("metrics", pageData, "layouts/base")
}

func (p *MetricsPage) GenerateMetrics(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	// TODO async data aggregation with go routines

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	entityId, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		log.Println("Error converting entity id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	entity, err := p.entityService.RetrieveEntity(c.Context(), entityId, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	query := models.NewMetricsQueryFromForm(c)
	query.Entity = entity
	query.CreatedBy = userID

	if !query.IsValid() {
		pageData.MetricsQuery = query
		pageData.FormMsg = "Corrija os campos abaixo."
		return utils.RetargetToForm(c, "metrics", pageData)
	}

	err = p.service.CreateMetrics(c.Context(), query)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	go p.siareBGWorker.GetMetrics(query)

	c.Set("HX-Trigger-After-Settle", "metrics-query-started")
	return c.Render("partials/request-card", query)
}

func (p *MetricsPage) GetRequestCardDetails(c *fiber.Ctx) error {
	queryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.MetricsNotFoundErr)
	}
	userID := c.Locals("UserID").(int)
	query, err := p.service.RetrieveMetrics(c.Context(), queryID, userID)

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", query)
}

func (p *MetricsPage) GetMetricsForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	queryID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.MetricsNotFoundErr)
	}
	query, err := p.service.RetrieveMetrics(c.Context(), queryID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	pageData.MetricsQuery = query

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/metrics-form", pageData)
}

func (p *MetricsPage) GetRequestStatus(c *fiber.Ctx) error {
	queryId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.MetricsNotFoundErr)
	}

	key := fmt.Sprintf("reqstatus:metrics:%v", queryId)
	err = db.Redis.GetDel(c.Context(), key).Err()
	if err == redis.Nil {
		return c.Render("partials/request-card-status", "pending")
	}
	if err != nil {
		log.Printf("Error reading redis key %v: %v\n", key, err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	userID := c.Locals("UserID").(int)

	query, err := p.service.RetrieveMetrics(c.Context(), queryId, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	targetId := fmt.Sprintf("#request-card-%v", c.Params("id"))
	c.Set("HX-Retarget", targetId)
	c.Set("HX-Reswap", "outerHTML")
	return c.Render("partials/request-card", query)
}
