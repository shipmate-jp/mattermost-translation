package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// ServeHTTP handles HTTP requests for the plugin.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Middleware to require that the user is logged in
	router.Use(p.MattermostAuthorizationRequired)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/translate", p.handleTranslate).Methods(http.MethodPost)

	router.ServeHTTP(w, r)
}

func (p *Plugin) MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// translateRequest is the payload received from the webapp.
type translateRequest struct {
	PostID     string `json:"post_id"`
	TargetLang string `json:"target"`
}

// googleTranslateResponse is the response received from Google Translate API v2.
type googleTranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

// translateResponse is the payload returned to the webapp.
type translateResponse struct {
	Translated string `json:"translated"`
}

func (p *Plugin) handleTranslate(w http.ResponseWriter, r *http.Request) {
	cfg := p.getConfiguration()
	if cfg == nil || cfg.GoogleAPIKey == "" {
		http.Error(w, "Google API key not configured", http.StatusBadRequest)
		return
	}

	var req translateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.PostID == "" {
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}
	if req.TargetLang == "" {
		req.TargetLang = cfg.DefaultTargetLang
		if req.TargetLang == "" {
			req.TargetLang = "ja"
		}
	}

	// Fetch the post contents using pluginapi client
	post, err := p.client.Post.GetPost(req.PostID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get post: %v", err), http.StatusInternalServerError)
		return
	}

	// Call Google Translate API v2 using application/x-www-form-urlencoded as per docs
	params := url.Values{}
	params.Set("q", post.Message)
	params.Set("target", req.TargetLang)
	params.Set("format", "text")

	urlStr := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", cfg.GoogleAPIKey)
	httpReq, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(params.Encode()))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		http.Error(w, "failed to call translate API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("translate API error: %s", string(b)), http.StatusBadGateway)
		return
	}

	var gtResp googleTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gtResp); err != nil {
		http.Error(w, "failed to decode translate response", http.StatusBadGateway)
		return
	}
	if len(gtResp.Data.Translations) == 0 {
		http.Error(w, "no translation returned", http.StatusBadGateway)
		return
	}

	translated := html.UnescapeString(gtResp.Data.Translations[0].TranslatedText)

	// Optionally send an ephemeral post to the requesting user with the translation
	if p.client != nil {
		if userID := r.Header.Get("Mattermost-User-ID"); userID != "" {
			rootID := ""
			if post.RootId != "" {
				rootID = post.RootId
			}
			p.client.Post.SendEphemeralPost(userID, &model.Post{
				ChannelId: post.ChannelId,
				RootId:    rootID,
				Message:   fmt.Sprintf("Translation (%s):\n%s", req.TargetLang, translated),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(translateResponse{Translated: translated})
}
