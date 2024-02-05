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

var parseTemplateError = "Failed to parse template"
var executeTemplateError = "Failed to execute template"
var expenseError = "Failed to get expenses"

type ExpenseHandler struct {
	Logger *slog.Logger
	database.Database
	webFS embed.FS
}

func NewExpenseHandler(logger *slog.Logger, db database.Database, webFS embed.FS) *ExpenseHandler {
	return &ExpenseHandler{
		Logger:   logger,
		Database: db,
		webFS:    webFS,
	}
}

func (e *ExpenseHandler) Routes(r chi.Router) {
	// api/v1/expense
	r.Get("/", e.HandleGetExpenses)
	r.Post("/", e.HandleAddExpense)
	r.Delete("/{id}", e.HandleDeleteExpense)
}

func (e *ExpenseHandler) HandleAddExpense(w http.ResponseWriter, r *http.Request) {
	amount := r.FormValue("amount")
	source := r.FormValue("expense")

	// add expense to database
	err := e.AddExpense(models.Expense{
		Amount: amount,
		Source: source,
	})
	if err != nil {
		e.Logger.Error("Failed to add expense", "error", err)
		http.Error(w, "Failed to add expense", http.StatusInternalServerError)
		return
	}

	err = executeGetExpenses(w, e)
}

func (e *ExpenseHandler) HandleGetExpenses(w http.ResponseWriter, r *http.Request) {
	err := executeGetExpenses(w, e)
	if err != nil {
		e.Logger.Error(expenseError, "error", err)
		http.Error(w, expenseError, http.StatusInternalServerError)
	}
}

func (e *ExpenseHandler) HandleDeleteExpense(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	expID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}
	// delete expense from database
	err = e.DeleteExpense(expID)
	if err != nil {
		e.Logger.Error("Failed to delete expense", "error", err)
		http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
		return
	}

	err = executeGetExpenses(w, e)
}

func executeGetExpenses(w http.ResponseWriter, e *ExpenseHandler) error {
	expenses, err := e.GetExpenses()
	if err != nil {
		e.Logger.Error(expenseError, "error", err)
		http.Error(w, expenseError, http.StatusInternalServerError)
		return err
	}

	// pass expenses to template
	tmpl, err := template.ParseFS(e.webFS, "web/components/expenses.html")
	if err != nil {
		e.Logger.Error(parseTemplateError, "error", err)
		http.Error(w, parseTemplateError, http.StatusInternalServerError)
		return err
	}
	err = tmpl.Execute(w, expenses)
	if err != nil {
		e.Logger.Error(executeTemplateError, "error", err)
		http.Error(w, executeTemplateError, http.StatusInternalServerError)
		return err
	}
	return nil
}
