package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"loverly/lib/appcontext"
	"loverly/lib/i18n"
	"loverly/lib/jwt"
	"loverly/lib/log"
	"net/http"
	"strings"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"loverly/lib/header"
	i18n_err "loverly/lib/i18n/errors"
)

type Response struct {
	Data     interface{} `json:"data"`
	Error    *Error      `json:"error"`
	Success  bool        `json:"success"`
	Metadata Meta        `json:"metadata"`
}

type Error struct {
	Code     string `json:"code"`
	Title    string `json:"message_title"`
	Message  string `json:"message"`
	Severity string `json:"message_severity"`
}

type Meta struct {
	RequestId string `json:"request_id"`
}

const (
	infoRequest  string = `httpclient Sent Request: uri=%v method=%v`
	infoResponse string = `httpclient Received Response: uri=%v method=%v resp_code=%v`
)

func JSONSuccess(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		Data:    data,
		Success: true,
		Metadata: Meta{
			RequestId: appcontext.GetRequestId(ctx),
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func JSONError(ctx context.Context, w http.ResponseWriter, code int, err i18n_err.I18nError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	lang := appcontext.GetAcceptLanguage(ctx)
	resp := Response{
		Error: &Error{
			Code:     err.Error(),
			Title:    i18n.Title(lang, err.Error()),
			Message:  i18n.Message(lang, err.Error()),
			Severity: "error",
		},
		Metadata: Meta{
			RequestId: appcontext.GetRequestId(ctx),
		},
	}

	json.NewEncoder(w).Encode(resp)
}

func addFieldsToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get(header.KeyRequestID)
		if reqId == "" {
			reqId = uuid.New().String()
		}

		c := r.Context()
		c = appcontext.SetRequestId(c, reqId)
		c = appcontext.SetUserAgent(c, r.Header.Get(header.KeyUserAgent))
		c = appcontext.SetAcceptLanguage(c, r.Header.Get(header.KeyAcceptLanguage))
		c = appcontext.SetDeviceType(c, r.Header.Get(header.KeyDeviceType))
		c = appcontext.SetCacheControl(c, r.Header.Get(header.KeyCacheControl))
		c = appcontext.SetServiceName(c, r.Header.Get(header.KeyServiceName))

		next.ServeHTTP(w, r.WithContext(c))
	})
}

func authentication(jwt *jwt.TokenProvider, log log.Interface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get(header.KeyAuthorization)
			if token == "" {
				JSONError(ctx, w, http.StatusUnauthorized, errors.New("missing authorization"))
				return
			}

			accessToken := strings.Split(token, " ")
			if len(accessToken) < 2 {
				JSONError(ctx, w, http.StatusUnauthorized, errors.New("invalid access token"))
				return
			}

			verify, err := jwt.DecodeAccessToken(ctx, accessToken[1])
			if err != nil {
				JSONError(ctx, w, http.StatusUnauthorized, errors.New("invalid or expired access token"))
				return
			}

			ctx = appcontext.SetUserId(ctx, int(verify.Data.UserId))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bodyLogger(log log.Interface) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info(r.Context(), fmt.Sprintf(infoRequest, r.RequestURI, r.Method))
			wrap := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			h.ServeHTTP(wrap, r)

			status := wrap.Status()
			if status < 300 {
				log.Info(r.Context(), fmt.Sprintf(infoResponse, r.RequestURI, r.Method, status))
			} else {
				log.Error(r.Context(), fmt.Sprintf(infoResponse, r.RequestURI, r.Method, status))
			}
		})
	}
}
