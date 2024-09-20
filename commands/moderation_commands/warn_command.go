package moderation_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// WarnCommand avertit un utilisateur
func WarnCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorer les messages du bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if !isMod {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%swarn", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupère la mention et la raison
	parts := strings.Fields(m.Content)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?warn <@mention_user> <raison>")
		return
	}

	// Récupérer l'utilisateur mentionné
	mention := parts[1]
	if !strings.HasPrefix(mention, "<@") || !strings.HasSuffix(mention, ">") {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide.")
		return
	}
	mentionID := mention[2 : len(mention)-1] // Extrait l'ID de la mention

	// Combiner les autres parties comme la raison
	reason := strings.Join(parts[2:], " ")

	// Créer une instance de WarnController
	warnController := controllers.NewWarnController()

	// Créer une instance de WarnService avec la session Discord et l'ID du serveur
	warnService := services.NewWarnService(warnController, s, m.GuildID)

	// Ajouter l'avertissement à la base de données
	err := warnService.AddWarn(mentionID, reason, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'enregistrement de l'avertissement.")
		return
	}

	// Envoyer un message confirmant l'avertissement
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissement enregistré pour <@%s> : %s", mentionID, reason))
}
