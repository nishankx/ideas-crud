package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nisha/gravehard/internal/auth"
	"github.com/nisha/gravehard/internal/config"
	"github.com/nisha/gravehard/internal/domain"
	"github.com/nisha/gravehard/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		logger.Error("config_load_failed", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db_connect_failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	ideaRepo := repository.NewIdeaRepository(pool)

	demoUser, err := userRepo.GetDemoUser(ctx)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			logger.Error("get_demo_user_failed", slog.Any("error", err))
			os.Exit(1)
		}
		hash, err := auth.HashPassword("GraveyardDemo123!")
		if err != nil {
			logger.Error("hash_password_failed", slog.Any("error", err))
			os.Exit(1)
		}
		demoUser, err = userRepo.Create(ctx, "demo@gravehard.dev", hash, "Serial Pivoter", true)
		if err != nil {
			logger.Error("create_demo_user_failed", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("created_demo_user", slog.String("id", demoUser.ID.String()))
	}

	existing, err := ideaRepo.ListAllByUser(ctx, demoUser.ID)
	if err != nil {
		logger.Error("list_demo_ideas_failed", slog.Any("error", err))
		os.Exit(1)
	}
	if len(existing) > 0 {
		logger.Info("demo_ideas_already_seeded", slog.Int("count", len(existing)))
		return
	}

	seeds := buildDemoIdeas(demoUser.ID)
	for _, idea := range seeds {
		if _, err := ideaRepo.Create(ctx, &idea); err != nil {
			logger.Error("seed_idea_failed", slog.String("name", idea.Name), slog.Any("error", err))
			os.Exit(1)
		}
	}
	logger.Info("seeded_demo_ideas", slog.Int("count", len(seeds)))
}

func d(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func buildDemoIdeas(uid [16]byte) []domain.Idea {
	return []domain.Idea{
		{
			UserID: uid, Name: "BarkCRM", Tagline: "Salesforce for neighborhood dog walkers",
			Description: "Full invoicing, GPS leash tracking, and treat inventory management.",
			Category: domain.CategorySaaS, BornAt: d("2023-02-01"), DiedAt: d("2023-03-10"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{domain.CauseLostMotivation},
			Epitaph: "Here lies a CRM for dog walkers. It never met a customer.",
			Autopsy: "Dog walkers already use Venmo and a group text.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 18000, HoursInvested: 65, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "SplitRentAI", Tagline: "Autonomous roommate chore & bill mediator",
			Description: "AI mediator that drafts passive-aggressive Slack messages about dirty dishes.",
			Category: domain.CategoryFintech, BornAt: d("2023-03-15"), DiedAt: d("2023-03-22"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{domain.CauseBadTiming},
			Epitaph: "Turns out roommates prefer passive-aggressive sticky notes.",
			Autopsy: "Nobody wanted to pay $9/mo to outsource awkwardness.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 9500, HoursInvested: 28, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "CoinJarMicro", Tagline: "Round-up crypto spare change for coffee",
			Description: "Connected to Plaid and automatically bought $0.35 of ETH per latte.",
			Category: domain.CategoryFintech, BornAt: d("2023-04-01"), DiedAt: d("2023-04-12"),
			PrimaryCause: domain.CauseLegalRegulatory, SecondaryCauses: []domain.DeathCause{domain.CauseRanOutOfMoney},
			Epitaph: "Gas fees cost 14x more than the spare change.",
			Autopsy: "Check unit economics before writing smart contracts.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 65000, HoursInvested: 90, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "TipTapTax", Tagline: "Instant tax withholding for freelance clowns",
			Description: "Hyper-niche bookkeeping for party entertainers.",
			Category: domain.CategoryFintech, BornAt: d("2023-05-01"), DiedAt: d("2023-05-08"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{},
			Epitaph: "The TAM was literally 12 clowns in Ohio.",
			Autopsy: "Do basic market sizing before registering the LLC.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 42000, HoursInvested: 19, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "TabGrave", Tagline: "Browser extension that shames you for 200 open tabs",
			Description: "Closed tabs automatically after 48 hours and emailed you a eulogy.",
			Category: domain.CategoryDevtool, BornAt: d("2023-05-20"), DiedAt: d("2023-07-18"),
			PrimaryCause: domain.CauseLostMotivation, SecondaryCauses: []domain.DeathCause{domain.CauseCompetitor},
			Epitaph: "People love hoarding tabs more than closing them.",
			Autopsy: "Great weekend hack, terrible SaaS business.", ValidationLevel: domain.ValidationTalkedToUsers,
			MoneySpentCents: 2400, HoursInvested: 165, Status: domain.StatusOnLifeSupport, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "FridgeWhisperer", Tagline: "IoT camera inside your fridge to track spoiled milk",
			Description: "Raspberry Pi + vision model that texted you when cilantro wilted.",
			Category: domain.CategoryHardware, BornAt: d("2023-06-01"), DiedAt: d("2023-10-15"),
			PrimaryCause: domain.CauseTechnicalDifficulty, SecondaryCauses: []domain.DeathCause{domain.CauseRanOutOfMoney},
			Epitaph: "Condensation killed three Raspberry Pis in a week.",
			Autopsy: "Hardware inside cold wet boxes is unforgiving.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 185000, HoursInvested: 240, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "GymBuddyBets", Tagline: "Lose $20 to your friends if you skip leg day",
			Description: "Geofenced gym check-ins with Stripe escrow.",
			Category: domain.CategoryConsumer, BornAt: d("2023-07-01"), DiedAt: d("2023-09-25"),
			PrimaryCause: domain.CauseLegalRegulatory, SecondaryCauses: []domain.DeathCause{domain.CauseLostMotivation},
			Epitaph: "Stripe flagged us as an illegal gambling ring.",
			Autopsy: "Holding user stakes in escrow requires real compliance.", ValidationLevel: domain.ValidationTalkedToUsers,
			MoneySpentCents: 54000, HoursInvested: 110, Status: domain.StatusBuried, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "PromptBazaar", Tagline: "Two-sided marketplace for artisanal ChatGPT prompts",
			Description: "Sell 12-line prompts for $4.99 with 20% take rate.",
			Category: domain.CategoryMarketplace, BornAt: d("2023-08-01"), DiedAt: d("2023-08-29"),
			PrimaryCause: domain.CauseCompetitor, SecondaryCauses: []domain.DeathCause{domain.CauseBadTiming},
			Epitaph: "Commoditized before the landing page finished deploying.",
			Autopsy: "Don't build a tollbooth on a free highway.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 12000, HoursInvested: 45, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "QuietDesk", Tagline: "White-noise generator for open-office panic",
			Description: "Adaptive brown noise that masked nearby product managers.",
			Category: domain.CategorySaaS, BornAt: d("2023-09-01"), DiedAt: d("2024-01-20"),
			PrimaryCause: domain.CauseCompetitor, SecondaryCauses: []domain.DeathCause{domain.CauseRanOutOfTime},
			Epitaph: "Spotify playlists and AirPods already won.",
			Autopsy: "Had 14 paying users, couldn't justify 10x CAC.", ValidationLevel: domain.ValidationPayingCustomers,
			MoneySpentCents: 31000, HoursInvested: 190, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "CofounderTinder", Tagline: "Swipe right on technical co-founders",
			Description: "Matching based on GitHub commit history and Myers-Briggs.",
			Category: domain.CategorySocial, BornAt: d("2023-10-01"), DiedAt: d("2023-11-12"),
			PrimaryCause: domain.CauseCofounderConflict, SecondaryCauses: []domain.DeathCause{domain.CauseNoMarketNeed},
			Epitaph: "Irony: our two co-founders broke up over equity.",
			Autopsy: "Always sign vesting schedules on day one.", ValidationLevel: domain.ValidationTalkedToUsers,
			MoneySpentCents: 15000, HoursInvested: 75, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "NapPodPass", Tagline: "Airbnb for 30-minute office naps",
			Description: "Bookable under-desk cots in downtown coworking spaces.",
			Category: domain.CategoryMarketplace, BornAt: d("2023-11-01"), DiedAt: d("2023-12-03"),
			PrimaryCause: domain.CauseBadTiming, SecondaryCauses: []domain.DeathCause{domain.CauseNoMarketNeed},
			Epitaph: "Remote work meant everyone already had a bed.",
			Autopsy: "Supply side was gross; demand side worked from home.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 22000, HoursInvested: 50, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "SQL2Poem", Tagline: "Turn slow Postgres queries into haikus",
			Description: "EXPLAIN ANALYZE parser that roasts your missing indexes in verse.",
			Category: domain.CategoryDevtool, BornAt: d("2023-12-01"), DiedAt: d("2024-02-14"),
			PrimaryCause: domain.CausePivotedAway, SecondaryCauses: []domain.DeathCause{domain.CauseNoMarketNeed},
			Epitaph: "Went viral on Hacker News, made $0.",
			Autopsy: "Developer memes are great top-of-funnel, not a product.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 4500, HoursInvested: 85, Status: domain.StatusResurrected, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "StudyStreak", Tagline: "Duolingo-style streaks for med school flashcards",
			Description: "Spaced repetition with aggressive owl-inspired guilt trips.",
			Category: domain.CategoryEducation, BornAt: d("2024-01-05"), DiedAt: d("2024-07-19"),
			PrimaryCause: domain.CauseCompetitor, SecondaryCauses: []domain.DeathCause{domain.CauseRanOutOfTime},
			Epitaph: "Anki is free and med students are loyal.",
			Autopsy: "Never underestimate entrenched open-source communities.", ValidationLevel: domain.ValidationPayingCustomers,
			MoneySpentCents: 78000, HoursInvested: 310, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "PosturePolice", Tagline: "Webcam AI that blurs your screen when you slouch",
			Description: "On-device pose estimation that held your IDE hostage.",
			Category: domain.CategoryHealth, BornAt: d("2024-02-01"), DiedAt: d("2024-03-02"),
			PrimaryCause: domain.CauseLostMotivation, SecondaryCauses: []domain.DeathCause{domain.CauseTechnicalDifficulty},
			Epitaph: "Users uninstalled it within 11 minutes.",
			Autopsy: "Punishing users is a terrible retention loop.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 6000, HoursInvested: 60, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "MeetupRoulette", Tagline: "Randomly pair strangers for 15-minute coffees",
			Description: "Calendar-integrated serendipity for lonely remote workers.",
			Category: domain.CategorySocial, BornAt: d("2024-03-01"), DiedAt: d("2024-05-20"),
			PrimaryCause: domain.CauseRanOutOfTime, SecondaryCauses: []domain.DeathCause{domain.CauseLostMotivation},
			Epitaph: "90% no-show rate killed the magic.",
			Autopsy: "Calendar invites without commitment are worthless.", ValidationLevel: domain.ValidationTalkedToUsers,
			MoneySpentCents: 19000, HoursInvested: 170, Status: domain.StatusBuried, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "InvoiceGhost", Tagline: "Automated polite follow-ups for unpaid invoices",
			Description: "Escalating email cadence from 'friendly bump' to 'collections'.",
			Category: domain.CategorySaaS, BornAt: d("2024-04-01"), DiedAt: d("2024-08-15"),
			PrimaryCause: domain.CausePivotedAway, SecondaryCauses: []domain.DeathCause{domain.CauseCompetitor},
			Epitaph: "Actually worked, but we got bored of accounting.",
			Autopsy: "Boring B2B pain points pay the bills if you stick around.", ValidationLevel: domain.ValidationPayingCustomers,
			MoneySpentCents: 29000, HoursInvested: 210, Status: domain.StatusOnLifeSupport, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "SnackBoxDAO", Tagline: "Community-governed office snack subscription",
			Description: "Token holders voted on Doritos vs. kale chips weekly.",
			Category: domain.CategoryConsumer, BornAt: d("2024-05-01"), DiedAt: d("2024-05-14"),
			PrimaryCause: domain.CauseBadTiming, SecondaryCauses: []domain.DeathCause{domain.CauseNoMarketNeed},
			Epitaph: "Nobody wanted governance overhead for trail mix.",
			Autopsy: "Don't put a blockchain on a vending machine.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 35000, HoursInvested: 35, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "DocuDiff", Tagline: "Git-style diff viewer for PDF contracts",
			Description: "Highlighted sneaky redlines in SaaS MSAs.",
			Category: domain.CategorySaaS, BornAt: d("2024-06-01"), DiedAt: d("2024-09-30"),
			PrimaryCause: domain.CauseRanOutOfTime, SecondaryCauses: []domain.DeathCause{domain.CauseCompetitor},
			Epitaph: "PDF layout parsing broke my spirit.",
			Autopsy: "Lawyers still live in Microsoft Word Track Changes.", ValidationLevel: domain.ValidationTalkedToUsers,
			MoneySpentCents: 41000, HoursInvested: 185, Status: domain.StatusBuried, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "TinyTutor", Tagline: "On-demand 5-minute math homework rescue",
			Description: "Uber for algebra help at 11pm before midnight deadlines.",
			Category: domain.CategoryEducation, BornAt: d("2024-07-01"), DiedAt: d("2024-08-05"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{domain.CauseCompetitor},
			Epitaph: "LLMs solved 5-minute homework help for free.",
			Autopsy: "Validate against AI substitutes before building.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 14000, HoursInvested: 42, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "APIWatchdog", Tagline: "Alerts when third-party APIs silently change schemas",
			Description: "Diffed JSON responses in staging and paged Slack.",
			Category: domain.CategoryDevtool, BornAt: d("2024-08-01"), DiedAt: d("2024-12-10"),
			PrimaryCause: domain.CauseRanOutOfMoney, SecondaryCauses: []domain.DeathCause{domain.CauseRanOutOfTime},
			Epitaph: "Great product, zero distribution strategy.",
			Autopsy: "Build distribution alongside the MVP.", ValidationLevel: domain.ValidationPayingCustomers,
			MoneySpentCents: 92000, HoursInvested: 260, Status: domain.StatusBuried, WouldResurrect: true,
		},
		{
			UserID: uid, Name: "PlantSitter", Tagline: "Neighborhood marketplace for vacation plant watering",
			Description: "Verified plant lovers water your monstera for $12/visit.",
			Category: domain.CategoryMarketplace, BornAt: d("2024-09-01"), DiedAt: d("2024-09-28"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{domain.CauseLostMotivation},
			Epitaph: "Neighbors do this for a six-pack of beer.",
			Autopsy: "Don't monetize favors people do for free.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 8000, HoursInvested: 33, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "SleepDebtBank", Tagline: "Gamified sleep tracking with social leaderboards",
			Description: "Compete with coworkers on who slept the least.",
			Category: domain.CategoryHealth, BornAt: d("2024-10-01"), DiedAt: d("2024-11-18"),
			PrimaryCause: domain.CauseLostMotivation, SecondaryCauses: []domain.DeathCause{domain.CauseNoMarketNeed},
			Epitaph: "Turns out bragging about 4 hours of sleep is toxic.",
			Autopsy: "Apple Watch already owns passive sleep tracking.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 11000, HoursInvested: 155, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "PaywallBypass", Tagline: "Micro-payments ($0.05) per news article",
			Description: "One wallet across 100 indie blogs and newsletters.",
			Category: domain.CategoryFintech, BornAt: d("2024-11-01"), DiedAt: d("2024-11-11"),
			PrimaryCause: domain.CauseNoMarketNeed, SecondaryCauses: []domain.DeathCause{domain.CauseBadTiming},
			Epitaph: "Everyone tried micropayments in 1999, 2011, and 2024.",
			Autopsy: "Readers hate mental friction per click.", ValidationLevel: domain.ValidationNone,
			MoneySpentCents: 16000, HoursInvested: 24, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "ZeroInboxClub", Tagline: "Human-in-the-loop email triage for $49/mo",
			Description: "Unsubscribed and archived newsletters every morning at 6am.",
			Category: domain.CategorySaaS, BornAt: d("2024-12-01"), DiedAt: d("2025-03-01"),
			PrimaryCause: domain.CauseTechnicalDifficulty, SecondaryCauses: []domain.DeathCause{domain.CauseLegalRegulatory},
			Epitaph: "Google OAuth CASA verification cost more than our revenue.",
			Autopsy: "Read Gmail API security review requirements first.", ValidationLevel: domain.ValidationPayingCustomers,
			MoneySpentCents: 145000, HoursInvested: 220, Status: domain.StatusBuried, WouldResurrect: false,
		},
		{
			UserID: uid, Name: "GraveyardStats", Tagline: "Benchmarking tool for failed side projects",
			Description: "Meta-analytics for serial builders who ship and abandon.",
			Category: domain.CategoryConsumer, BornAt: d("2025-01-10"), DiedAt: d("2025-02-20"),
			PrimaryCause: domain.CausePivotedAway, SecondaryCauses: []domain.DeathCause{},
			Epitaph: "Pivoted into Startup Idea Graveyard itself.",
			Autopsy: "Sometimes the graveyard IS the product.", ValidationLevel: domain.ValidationPrototype,
			MoneySpentCents: 5000, HoursInvested: 64, Status: domain.StatusResurrected, WouldResurrect: true,
		},
	}
}
