package auth

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/google/uuid"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

// GoogleOAuth wraps the OAuth2 flow configuration.
type GoogleOAuth struct {
    config *oauth2.Config
}

// NewGoogleOAuth builds a GoogleOAuth helper.
func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
    return &GoogleOAuth{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes: []string{
                "https://www.googleapis.com/auth/userinfo.email",
                "https://www.googleapis.com/auth/userinfo.profile",
            },
            Endpoint: google.Endpoint,
        },
    }
}

// AuthURL creates a login URL with a state token.
func (g *GoogleOAuth) AuthURL(state string) string {
    return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges an authorization code for tokens.
func (g *GoogleOAuth) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
    return g.config.Exchange(ctx, code)
}

// FetchUser fetches Google profile info using the provided token.
func (g *GoogleOAuth) FetchUser(ctx context.Context, token *oauth2.Token) (GoogleUser, error) {
    client := g.config.Client(ctx, token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        return GoogleUser{}, err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return GoogleUser{}, fmt.Errorf("google userinfo status %d", resp.StatusCode)
    }
    var user GoogleUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return GoogleUser{}, err
    }
    if user.ID == "" {
        return GoogleUser{}, errors.New("google: empty user id")
    }
    return user, nil
}

// GoogleUser represents the minimal data we need from Google.
type GoogleUser struct {
    ID      string `json:"id"`
    Email   string `json:"email"`
    Name    string `json:"name"`
    Picture string `json:"picture"`
}

// NewStateToken creates a random string used for CSRF protection.
func NewStateToken() string {
    return uuid.NewString()
}

// CookieStateName is the cookie storing the OAuth state.
const CookieStateName = "google_oauth_state"

// SetStateCookie stores the OAuth state token.
func SetStateCookie(w http.ResponseWriter, state string) {
    http.SetCookie(w, &http.Cookie{
        Name:     CookieStateName,
        Value:    state,
        Path:     "/",
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
        Expires:  time.Now().Add(10 * time.Minute),
    })
}

// ValidateStateCookie ensures the incoming state matches the stored cookie.
func ValidateStateCookie(r *http.Request, state string) error {
    cookie, err := r.Cookie(CookieStateName)
    if err != nil {
        return err
    }
    if cookie.Value != state {
        return errors.New("invalid oauth state")
    }
    return nil
}
