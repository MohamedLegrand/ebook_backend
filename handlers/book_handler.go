package handlers

import (
    "ebook-backend/config"
    "ebook-backend/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetAllBooks(c *gin.Context) {
    rows, err := config.DB.Query(c.Request.Context(), "SELECT id, titre, auteur, description, prix_fcfa, image, type, pages, stock, created_at, updated_at FROM livres ORDER BY id DESC")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var b models.Book
        if err := rows.Scan(&b.ID, &b.Titre, &b.Auteur, &b.Description, &b.PrixFCFA, &b.Image, &b.Type, &b.Pages, &b.Stock, &b.CreatedAt, &b.UpdatedAt); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        books = append(books, b)
    }
    c.JSON(http.StatusOK, books)
}