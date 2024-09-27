package controllers

import (
	"time"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/models"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

// LogController est un contrôleur pour gérer les opérations sur les entrées de journal
type LogController struct {
	DB *gorm.DB
}

// NewLogController crée une nouvelle instance de LogController avec une connexion à la base de données
func NewLogController() *LogController {
	return &LogController{
		DB: database.DB,
	}
}

// CreateLog crée une nouvelle entrée de journal
func (ctrl *LogController) CreateLog(timestamp time.Time, serverID, serverName, channelID, channelName, userID, username, message string) (*models.Log, error) {
	logEntry := models.Log{
		Timestamp:     timestamp,
		ServerID:      serverID,
		ServerName:    serverName,
		ChannelID:     channelID,
		ChannelName:   channelName,
		UserDiscordID: userID,
		Username:      username,
		Message:       message,
	}

	if err := ctrl.DB.Create(&logEntry).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// GetLog récupère une entrée de journal par ID
func (ctrl *LogController) GetLog(id uint) (*models.Log, error) {
	var logEntry models.Log
	if err := ctrl.DB.First(&logEntry, id).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// UpdateLog met à jour une entrée de journal existante
func (ctrl *LogController) UpdateLog(id uint, timestamp time.Time, serverID, serverName, channelID, channelName, userID, username, message string) (*models.Log, error) {
	var logEntry models.Log
	if err := ctrl.DB.First(&logEntry, id).Error; err != nil {
		return nil, err
	}

	// Mise à jour des champs de l'entrée de journal
	logEntry.Timestamp = timestamp
	logEntry.ServerID = serverID
	logEntry.ServerName = serverName
	logEntry.ChannelID = channelID
	logEntry.ChannelName = channelName
	logEntry.UserDiscordID = userID
	logEntry.Username = username
	logEntry.Message = message

	if err := ctrl.DB.Save(&logEntry).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// GetLastLogs récupère les X derniers logs
func (ctrl *LogController) GetLastLogs(limit int) ([]models.Log, error) {
	var logs []models.Log
	if err := ctrl.DB.Order("timestamp desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLogsByUser récupère les logs associés à un utilisateur Discord spécifique
func (ctrl *LogController) GetLogsByUser(userID string, limit int) ([]models.Log, error) {
	var logs []models.Log
	if err := ctrl.DB.Where("user_discord_id = ?", userID).Order("timestamp desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLogsByUserAndChannel récupère les logs associés à un utilisateur dans un canal spécifique
func (ctrl *LogController) GetLogsByUserAndChannel(userID, channelID string, limit int) ([]models.Log, error) {
	var logs []models.Log
	if err := ctrl.DB.Where("user_discord_id = ? AND channel_id = ?", userID, channelID).
		Order("timestamp desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// GetLogsByChannel récupère les logs d'un canal spécifique
func (ctrl *LogController) GetLogsByChannel(channelID string, limit int) ([]models.Log, error) {
	var logs []models.Log
	if err := ctrl.DB.Where("channel_id = ?", channelID).
		Order("timestamp desc").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// DeleteLog supprime une entrée de journal par ID
func (ctrl *LogController) DeleteLog(id uint) error {
	return ctrl.DB.Delete(&models.Log{}, id).Error
}

// InsertLog insère un log dans la base de données
func (c *LogController) InsertLog(s *discordgo.Session, msg *discordgo.Message) error {
	serverName, err := discord.GetServerName(s, msg.GuildID)
	if err != nil {
		return err // Retourne l'erreur si la récupération échoue
	}

	channelName, err := discord.GetChannelName(s, msg.ChannelID)
	if err != nil {
		return err // Retourne l'erreur si la récupération échoue
	}

	logEntry := models.Log{
		Timestamp:     time.Now(),
		ServerID:      msg.GuildID,
		ServerName:    serverName,
		ChannelID:     msg.ChannelID,
		ChannelName:   channelName,
		UserDiscordID: msg.Author.ID,
		Username:      msg.Author.Username,
		Message:       msg.Content,
	}

	return c.DB.Create(&logEntry).Error
}

// SaveLog enregistre ou met à jour une entrée de journal dans la base de données
func (ctrl *LogController) SaveLog(entry *models.Log) error {
	if entry.ID == 0 {
		// Créer une nouvelle entrée si ID est zéro
		return ctrl.DB.Create(entry).Error
	}
	// Sinon, mettre à jour l'entrée existante
	return ctrl.DB.Save(entry).Error
}
