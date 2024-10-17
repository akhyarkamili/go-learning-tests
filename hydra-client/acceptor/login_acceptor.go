package acceptor

import (
	"context"
	api "learning-tests/hydra-client/hydra-client"
	"net/http"

	client "github.com/ory/hydra-client-go"
)

type HydraAuthFlowAcceptorServer struct {
	authCode chan<- string
}

func NewHydraAuthFlowAcceptorServer(authCode chan string) *HydraAuthFlowAcceptorServer {
	return &HydraAuthFlowAcceptorServer{
		authCode: authCode,
	}
}

// Serve serves dummy hydra login & consent acceptor HTTP server for handling login & consent requests
func Serve(ctx context.Context, port string, s *HydraAuthFlowAcceptorServer) {
	// create a http server on given port
	// handle login & consent requests

	server := &http.Server{
		Addr:    ":" + port,
		Handler: s,
	}
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()
	_ = server.ListenAndServe()
}

// ServeHTTP handles login & consent requests
func (h *HydraAuthFlowAcceptorServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hydraClient := api.CreateHydraApiClient("4445", nil)
	switch r.URL.Path {
	case "/login":
		// parse query params
		lc := r.URL.Query().Get("login_challenge")
		req := hydraClient.AdminApi.AcceptLoginRequest(context.Background())
		cr, _, err := req.LoginChallenge(lc).AcceptLoginRequest(client.AcceptLoginRequest{
			Subject: "user123",
		}).Execute()

		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		redirectUrl := cr.GetRedirectTo()
		w.Header().Set("Location", redirectUrl)
		w.WriteHeader(302)
	case "/consent":
		// handle consent request
		cc := r.URL.Query().Get("consent_challenge")
		req := hydraClient.AdminApi.AcceptConsentRequest(context.Background())
		cr, _, err := req.ConsentChallenge(cc).AcceptConsentRequest(client.AcceptConsentRequest{
			GrantScope: []string{"openid"},
		}).Execute()
		if err != nil {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		redirectUrl := cr.GetRedirectTo()
		// just return ok
		w.Header().Set("Location", redirectUrl)
		w.WriteHeader(302)
		return
	case "/callback":
		code := r.URL.Query().Get("code")
		go func() {
			h.authCode <- code
		}()
		w.WriteHeader(200)
	default:
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Not implemented"))
		return
	}

}
