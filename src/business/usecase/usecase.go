package usecase

import (
	"loverly/lib/jwt"
	"loverly/lib/log"
	"loverly/src/business/domain"
	"loverly/src/business/usecase/dating"
	"loverly/src/business/usecase/match"
	"loverly/src/business/usecase/profile"
	"loverly/src/business/usecase/subscription"
	"loverly/src/business/usecase/user"
	"loverly/src/config"

	"loverly/lib/atomic"

	"go.opentelemetry.io/otel/trace"
)

type Usecases struct {
	User         user.Interface
	Dating       dating.Interface
	Subscription subscription.Interface
	Match        match.Interface
	Profile      profile.Interface
}

func Init(log log.Interface, cfg config.Configuration, jwt jwt.TokenProvider, dom domain.Domains, atomic atomic.AtomicSessionProvider, tr trace.Tracer) *Usecases {
	return &Usecases{
		User:         user.Init(log, &jwt, dom.User, dom.Profile, atomic),
		Dating:       dating.Init(log, dom.Subscription, dom.Profile, dom.Swipe, dom.Match),
		Subscription: subscription.Init(log, dom.Subscription),
		Match:        match.Init(log, dom.Match, dom.Profile),
		Profile:      profile.Init(log, dom.Profile),
	}
}
