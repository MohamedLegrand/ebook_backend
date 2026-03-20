package handlers

import (
	"ebook-backend/models"
	"ebook-backend/repository"
	"ebook-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Register gère l'inscription d'un nouveau client
func Register(c *gin.Context) {

	// 1. On récupère et valide les données envoyées par le formulaire React
	//    Si un champ obligatoire manque, Gin renvoie automatiquement une erreur 400
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides : " + err.Error(),
		})
		return
	}

	// 2. On vérifie que les deux mots de passe correspondent
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Les mots de passe ne correspondent pas",
		})
		return
	}

	// 3. On crée le client dans la base de données
	//    Le repository va hasher le mot de passe avant de l'enregistrer
	client, err := repository.CreateClient(req)
	if err != nil {
		// Si l'email existe déjà, PostgreSQL renvoie une erreur de contrainte unique
		// On détecte cette erreur pour donner un message clair à l'utilisateur
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Cet email est déjà utilisé",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la création du compte",
		})
		return
	}

	// 4. Inscription réussie, on renvoie les infos du client créé
	//    Le mot de passe n'est jamais renvoyé grâce au json:"-" dans le modèle
	c.JSON(http.StatusCreated, gin.H{
		"message": "Compte créé avec succès",
		"client":  client,
	})
}

// Login gère la connexion d'un client existant
func Login(c *gin.Context) {

	// 1. On récupère email et mot de passe envoyés par le formulaire
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Données invalides : " + err.Error(),
		})
		return
	}

	// 2. On normalise l'email en minuscules
	//    Evite les problèmes si l'utilisateur tape "Test@Email.com" au lieu de "test@email.com"
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// 3. On cherche le client dans la base de données par son email
	client, err := repository.GetClientByEmail(req.Email)
	if err != nil {
		// On ne précise pas si c'est l'email ou le mot de passe qui est faux
		// C'est une bonne pratique de sécurité pour ne pas aider un attaquant
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ou mot de passe incorrect",
		})
		return
	}

	// 4. On compare le mot de passe envoyé avec le hash stocké en base
	if !repository.CheckPassword(client.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ou mot de passe incorrect",
		})
		return
	}

	// 5. Identifiants corrects, on génère le token JWT
	token, err := utils.GenerateToken(client.ID, client.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la génération du token",
		})
		return
	}

	// 6. On renvoie le token et les infos du client au frontend React
	//    React va stocker ce token et l'utiliser pour les prochaines requêtes
	c.JSON(http.StatusOK, models.LoginResponse{
		Token:  token,
		Client: *client,
	})
}