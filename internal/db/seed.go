package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/shanisharrma/gopher-social/internal/store"
)

var usernames = []string{
	"alpha_wolf", "beta_blast", "gamma_ray", "delta_force", "echo_shadow",
	"foxtrot_ninja", "giga_brain", "hacker_x", "ice_dragon", "joker_007",
	"king_cobra", "lucky_strike", "mystic_warrior", "neo_phantom", "omega_zeus",
	"pixel_pirate", "quantum_flux", "rogue_sniper", "shadow_hawk", "titan_vortex",
	"ultra_instinct", "venom_fang", "warp_speed", "xenon_blaze", "yankee_alpha",
	"zulu_ranger", "cosmic_knight", "digital_demon", "elemental_fury", "frostbite_x",
	"glitch_master", "hyper_nova", "iron_fist", "joker_card", "knightmare_king",
	"lunar_echo", "midnight_rebel", "nebula_storm", "obsidian_blade", "psycho_cyber",
	"quantum_ghost", "radioactive_x", "storm_bringer", "turbo_flash", "unstoppable_warrior",
	"void_walker", "wildfire_alpha", "xcalibur_shadow", "zenith_guardian", "zero_gravity",
}

var titles = []string{
	"Mastering Go Basics",
	"10 Tips for Remote Work",
	"The Future of AI",
	"How to Start a Blog",
	"Boost Your Productivity",
	"Why Minimalism Matters",
	"Learn JavaScript Fast",
	"Healthy Eating Habits",
	"The Power of Meditation",
	"Building a Startup",
	"Essential Linux Commands",
	"Understanding Crypto",
	"Writing Clean Code",
	"Time Management Hacks",
	"Travel on a Budget",
	"Investing for Beginners",
	"SEO Strategies 2025",
	"Best Books to Read",
	"How to Stay Motivated",
	"Fitness Tips for Beginners",
}

var contents = []string{
	"Go is a statically typed, compiled language designed for efficiency and ease of use.",
	"Remote work requires discipline, clear goals, and effective communication tools.",
	"AI is transforming industries, from healthcare to finance, with automation and deep learning.",
	"Starting a blog is easy—choose a niche, create content, and stay consistent.",
	"Boost productivity by managing your time, avoiding distractions, and setting priorities.",
	"Minimalism simplifies life, reduces stress, and enhances mental clarity.",
	"JavaScript is the backbone of web development; mastering it opens many career doors.",
	"A balanced diet and regular exercise are key to long-term health and well-being.",
	"Meditation helps reduce stress, increase focus, and improve emotional health.",
	"Building a startup requires a strong idea, market validation, and resilience.",
	"Mastering Linux commands enhances productivity and makes system administration easier.",
	"Cryptocurrency is reshaping digital finance, offering decentralized and secure transactions.",
	"Writing clean code improves maintainability, readability, and long-term project success.",
	"Effective time management allows you to accomplish more in less time with less stress.",
	"Traveling on a budget is possible with careful planning, deals, and off-season travel.",
	"Investing early and diversifying your portfolio can secure long-term financial growth.",
	"SEO strategies like keyword research and quality content improve website rankings.",
	"Reading books expands knowledge, enhances critical thinking, and improves vocabulary.",
	"Motivation comes from setting clear goals, maintaining discipline, and tracking progress.",
	"Fitness success starts with consistency, a balanced diet, and gradual progress.",
}

var tags = []string{
	"technology", "programming", "golang", "web development", "AI",
	"productivity", "startup", "finance", "minimalism", "health",
	"fitness", "remote work", "coding", "SEO", "self-improvement",
	"travel", "meditation", "linux", "cryptocurrency", "education",
}

var comments = []string{
	"Great article! Very informative.",
	"I totally agree with your point!",
	"Interesting perspective, but I have a different view.",
	"Can you elaborate on this topic?",
	"Thanks for sharing! This was helpful.",
	"Amazing content! Keep up the good work.",
	"I never thought about it this way before.",
	"This really helped me understand the concept better.",
	"Do you have any resources to dive deeper into this?",
	"Fantastic read! Looking forward to more posts.",
	"Not sure I completely agree, but I respect your opinion.",
	"This was exactly what I was looking for!",
	"I appreciate the effort put into writing this.",
	"Could you simplify this for beginners?",
	"Well explained! Thanks for breaking it down.",
	"Wow! This changed my perspective completely.",
	"I'll definitely try this out. Thanks!",
	"Your writing style is very engaging!",
	"This needs more discussion. Let's talk!",
	"Helpful insights! I’ll be sharing this with my friends.",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	tx, _ := db.BeginTx(ctx, nil)

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts:", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comments:", err)
			return
		}
	}

	fmt.Println("Seeding completed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
	}
	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := range num {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := range num {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		cms[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
