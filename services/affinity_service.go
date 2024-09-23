package services

import (
	"strings"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

// AffinityService est un service pour gérer l'affinité des utilisateurs
type AffinityService struct {
	discordSession *discordgo.Session
	userCtrl       *controllers.UserController
	wordCtrl       *controllers.WordController
}

// NewAffinityService crée une nouvelle instance de AffinityService
func NewAffinityService(discordSession *discordgo.Session) *AffinityService {
	return &AffinityService{
		discordSession: discordSession,
		userCtrl:       &controllers.UserController{DB: database.DB},
		wordCtrl:       &controllers.WordController{DB: database.DB},
	}
}

// AdjustAffinity ajuste l'affinité d'un utilisateur en fonction du contenu de son message
func (service *AffinityService) AdjustAffinity(m *discordgo.MessageCreate) {
	if m.Author.ID == service.discordSession.State.User.ID {
		return
	}

	// Récupérer l'utilisateur ou le créer s'il n'existe pas
	user, err := service.userCtrl.GetUserByDiscordID(m.Author.ID)
	if err != nil {
		user = &models.User{
			UserDiscordID: m.Author.ID,
			Username:      m.Author.Username,
			Affinity:      0, // Affinité de départ
		}
		_, err = service.userCtrl.CreateUser(user.UserDiscordID, user.Username, user.Affinity, 0, 0, "", 0, 0, 0, 0)
		if err != nil {
			// Gérer l'erreur de création de l'utilisateur
			return
		}
	}

	// Récupérer les listes de bons et mauvais mots
	goodwords, err := service.wordCtrl.GetGoodWords()
	if err != nil {
		// Gérer l'erreur de récupération des bons mots
		return
	}
	badwords, err := service.wordCtrl.GetBadWords()
	if err != nil {
		// Gérer l'erreur de récupération des mauvais mots
		return
	}

	// Vérifier les mots interdits
	for _, word := range badwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity-- // Diminuer l'affinité
			err := service.userCtrl.UpdateUser(user)
			if err != nil {
				// Gérer l'erreur de mise à jour de l'utilisateur
				return
			}
			return
		}
	}

	// Vérifier les bons mots
	for _, word := range goodwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity++ // Augmenter l'affinité
			err := service.userCtrl.UpdateUser(user)
			if err != nil {
				// Gérer l'erreur de mise à jour de l'utilisateur
				return
			}
			return
		}
	}
}

// GetUserAffinity récupère l'affinité d'un utilisateur
func (service *AffinityService) GetUserAffinity(userID string) (models.User, bool) {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		return models.User{}, false
	}
	return *user, true
}
