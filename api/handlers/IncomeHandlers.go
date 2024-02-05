package handlers

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Ewan-Greer09/finance-app/api/database"
	"github.com/Ewan-Greer09/finance-app/api/models"
)

var incomeError = "Failed to get incomes"

type IncomeHandler struct {
	Logger *slog.Logger
	database.Database
	webFS embed.FS
}

func NewIncomeHandler(logger *slog.Logger, db database.Database, fs embed.FS) *IncomeHandler {
	return &IncomeHandler{
		Logger:   logger,
		Database: db,
		webFS:    fs,
	}
}

func (h *IncomeHandler) Routes(r chi.Router) {
	// api/v1/income
	r.Post("/", h.HandleAddIncome)
	r.Get("/", h.HandleGetIncomes)
	r.Delete("/{id}", h.HandleDeleteIncome)
}

func (h *IncomeHandler) HandleAddIncome(w http.ResponseWriter, r *http.Request) {
	amount := r.FormValue("amount")
	source := r.FormValue("income")

	// add income to database
	err := h.AddIncome(models.Income{
		Amount: amount,
		Source: source,
	})
	if err != nil {
		h.Logger.Error("Failed to add income", "error", err)
		http.Error(w, "Failed to add income", http.StatusInternalServerError)
		return
	}

	err = executeGetIncomes(w, h)
	if err != nil {
		h.Logger.Error(incomeError, "error", err)
		http.Error(w, incomeError, http.StatusInternalServerError)
	}
}

func (h *IncomeHandler) HandleGetIncomes(w http.ResponseWriter, r *http.Request) {
	err := executeGetIncomes(w, h)
	if err != nil {
		h.Logger.Error(incomeError, "error", err)
		http.Error(w, incomeError, http.StatusInternalServerError)
	}
}

func (h *IncomeHandler) HandleDeleteIncome(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	incID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid income ID", http.StatusBadRequest)
		return
	}

	err = h.DeleteIncome(incID)
	if err != nil {
		h.Logger.Error("Failed to delete income", "error", err)
		http.Error(w, "Failed to delete income", http.StatusInternalServerError)
		return
	}

	err = executeGetIncomes(w, h)
	if err != nil {
		h.Logger.Error(incomeError, "error", err)
		http.Error(w, incomeError, http.StatusInternalServerError)
	}
}

// reads incomes from database and passes them to the template
func executeGetIncomes(w http.ResponseWriter, h *IncomeHandler) error {
	incomes, err := h.GetIncomes()
	if err != nil {
		h.Logger.Error(incomeError, "error", err)
		http.Error(w, incomeError, http.StatusInternalServerError)
		return err
	}

	tmpl, err := template.ParseFS(h.webFS, "web/components/incomes.html")
	if err != nil {
		h.Logger.Error(parseTemplateError, "error", err)
		http.Error(w, parseTemplateError, http.StatusInternalServerError)
		return err
	}

	err = tmpl.Execute(w, incomes)
	if err != nil {
		h.Logger.Error(executeTemplateError, "error", err)
		http.Error(w, executeTemplateError, http.StatusInternalServerError)
		return err
	}
	return nil
}
