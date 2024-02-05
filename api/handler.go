package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/Ewan-Greer09/finance-app/api/config"
	"github.com/Ewan-Greer09/finance-app/api/database"
)

var parseTemplateError = "Failed to parse template"
var executeTemplateError = "Failed to execute template"
var expenseError = "Failed to get expenses"
var incomeError = "Failed to get incomes"

type Handler struct {
	*slog.Logger
	database.Database
}

func NewHandler(logger *slog.Logger, cfg config.Config) *Handler {
	return &Handler{
		Logger:   logger,
		Database: database.NewDatabase(cfg),
	}
}

func (h *Handler) HandleGetExpensesAndIncomesGraph(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.GetExpenses()
	if err != nil {
		h.Logger.Error(expenseError, "error", err)
		http.Error(w, expenseError, http.StatusInternalServerError)
		return
	}

	incomes, err := h.GetIncomes()
	if err != nil {
		h.Logger.Error(incomeError, "error", err)
		http.Error(w, incomeError, http.StatusInternalServerError)
		return
	}

	expTotal, incTotal := 0.0, 0.0
	for _, expense := range expenses {
		val, err := strconv.ParseFloat(expense.Amount, 64)
		if err != nil {
			http.Error(w, "Failed to parse expense amount", http.StatusInternalServerError)
			return
		}
		expTotal += val
	}
	for _, income := range incomes {
		val, err := strconv.ParseFloat(income.Amount, 64)
		if err != nil {
			http.Error(w, "Failed to parse expense amount", http.StatusInternalServerError)
			return
		}
		incTotal += val
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Expenses and Incomes",
		Subtitle: "Your Expenses and Incomes",
	}))

	bar.SetXAxis([]string{"Expenses vs Incomes"}).
		AddSeries("Expenses", []opts.BarData{
			{Value: expTotal},
		}).
		AddSeries("Incomes", []opts.BarData{
			{Value: incTotal},
		})

	buff := bytes.NewBuffer([]byte{})
	err = bar.Render(buff)
	if err != nil {
		h.Logger.Error("Failed to render graph", "error", err)
		http.Error(w, "Failed to render graph", http.StatusInternalServerError)
	}

	// load graph and pass to template
	tmpl, err := template.ParseFS(webFS, "web/components/graph.html")
	if err != nil {
		h.Logger.Error(parseTemplateError, "error", err)
		http.Error(w, parseTemplateError, http.StatusInternalServerError)
		return
	}
	buff.Next(200) //remove the <head /> from the graph
	err = tmpl.Execute(w, buff.String())
}
