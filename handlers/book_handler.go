package handlers

import (
    "ebook-backend/config"
    "ebook-backend/models"
    "net/http"

    "github.com/gin-gonic/gin"
)

// GetAllBooks godoc
// @Summary      Liste tous les livres
// @Description  Récupère la liste complète des livres (titre, auteur, prix, stock, etc.)
// @Tags         Livres
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Book
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/books [get]
func GetAllBooks(c *gin.Context) {
    rows, err := config.DB.Query(c.Request.Context(), 
        "SELECT id, titre, auteur, description, prix_fcfa, image, type, pages, stock, created_at, updated_at FROM livres ORDER BY id DESC")
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

// CreateBook godoc
// @Summary      Ajouter un livre
// @Description  Crée un nouveau livre (admin uniquement)
// @Tags         Livres
// @Accept       json
// @Produce      json
// @Param        book body models.Book true "Informations du livre"
// @Success      201  {object}  models.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/books [post]
func CreateBook(c *gin.Context) {
    var book models.Book
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `
        INSERT INTO livres (titre, auteur, description, prix_fcfa, image, type, pages, stock)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
    err := config.DB.QueryRow(
        c.Request.Context(),
        query,
        book.Titre, book.Auteur, book.Description, book.PrixFCFA,
        book.Image, book.Type, book.Pages, book.Stock,
    ).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, book)
}

// DeleteBook godoc
// @Summary      Supprimer un livre
// @Description  Supprime un livre par son ID (admin uniquement)
// @Tags         Livres
// @Accept       json
// @Produce      json
// @Param        id path int true "ID du livre"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/books/{id} [delete]
func DeleteBook(c *gin.Context) {
    id := c.Param("id")
    query := `DELETE FROM livres WHERE id = $1`
    result, err := config.DB.Exec(c.Request.Context(), query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    rowsAffected := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Livre non trouvé"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Livre supprimé avec succès"})
}