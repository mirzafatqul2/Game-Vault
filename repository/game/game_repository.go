package game

import (
	"Mini-Project-Game-Vault-API/service/game"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type RAWGConfig struct {
	BaseURL string
	APIKey  string
}

type GameRepository struct {
	logger *slog.Logger
	cfg    RAWGConfig
	DB     *gorm.DB
}

func NewGameRepository(logger *slog.Logger, cfg RAWGConfig, db *gorm.DB) *GameRepository {
	return &GameRepository{
		logger: logger,
		cfg:    cfg,
		DB:     db,
	}
}

type rawgGame struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	BackgroundImage string      `json:"background_image"`
	Genres          []rawgGenre `json:"genres"`
}

type rawgResponse struct {
	Results []rawgGame `json:"results"`
}

type rawgGenre struct {
	Name string `json:"name"`
}

func (r *GameRepository) GetGames(ctx context.Context, page, limit int) ([]game.Game, error) {
	url := fmt.Sprintf(
		"%s/games?key=%s&page=%d&page_size=%d",
		r.cfg.BaseURL,
		r.cfg.APIKey,
		page,
		limit,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		r.logger.Error("create request failed", "err", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Error("rawg request failed", "err", err)
		return nil, err
	}

	defer resp.Body.Close()

	var rawgResp rawgResponse

	if err := json.NewDecoder(resp.Body).Decode(&rawgResp); err != nil {
		r.logger.Error("failed decode rawg response", "error", err)
		return nil, err
	}

	var games []game.Game

	for _, g := range rawgResp.Results {

		category := "unknown"

		if len(g.Genres) > 0 {
			category = g.Genres[0].Name
		}

		games = append(games, game.Game{
			ExternalID: fmt.Sprintf("%d", g.ID),
			Name:       g.Name,
			ImageUrl:   g.BackgroundImage,
			Category:   category,
		})
	}

	r.logger.Info("fetch games success", "total", len(games))

	return games, nil
}

func (r *GameRepository) GetGameDetail(ctx context.Context, externalID string) (game.Game, error) {

	url := fmt.Sprintf(
		"%s/games/%s?key=%s",
		r.cfg.BaseURL,
		externalID,
		r.cfg.APIKey,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		r.logger.Error("create request failed", "err", err)
		return game.Game{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Error("rawg request failed", "err", err)
		return game.Game{}, err
	}

	defer resp.Body.Close()

	var rawgResp rawgGame

	if err := json.NewDecoder(resp.Body).Decode(&rawgResp); err != nil {
		r.logger.Error("decode rawg response failed", "err", err)
		return game.Game{}, err
	}

	category := "unknown"
	if len(rawgResp.Genres) > 0 {
		category = rawgResp.Genres[0].Name
	}

	return game.Game{
		ExternalID: fmt.Sprintf("%d", rawgResp.ID),
		Name:       rawgResp.Name,
		ImageUrl:   rawgResp.BackgroundImage,
		Category:   category,
	}, nil
}

func (r *GameRepository) Create(ctx context.Context, game game.Game) error {
	return r.DB.WithContext(ctx).Create(&game).Error
}

func (r *GameRepository) GetGameByID(ctx context.Context, id string) (game.Game, error) {
	var game game.Game

	err := r.DB.WithContext(ctx).Where("id =?", id).First(&game).Error
	if err != nil {
		return game, err
	}

	return game, nil
}

func (r *GameRepository) GetGameByExternalID(ctx context.Context, externalID string) (game.Game, error) {
	var game game.Game

	err := r.DB.WithContext(ctx).First(&game, "external_id = ? ", externalID).Error
	if err != nil {
		return game, err
	}

	return game, nil
}

func (r *GameRepository) GetGameReady(ctx context.Context) ([]game.Game, error) {
	var games []game.Game

	err := r.DB.WithContext(ctx).Find(&games).Error
	if err != nil {
		return nil, err
	}

	return games, nil
}
