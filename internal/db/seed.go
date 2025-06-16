package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/theluminousartemis/socialnews/internal/store"
)

var usernames = []string{
	"coolCat123",
	"happyPenguin",
	"starryNight",
	"mightyLion",
	"swiftEagle",
	"charmingFox",
	"braveBear",
	"cleverRabbit",
	"gentleGiraffe",
	"fierceTiger",
	"playfulDolphin",
	"wiseOwl",
	"boldShark",
	"friendlyPanda",
	"quickSquirrel",
	"loyalWolf",
	"jollyKoala",
	"curiousOtter",
	"eagerBeaver",
	"gracefulSwan",
	"cheerfulParrot",
	"nobleStag",
	"wittyRaccoon",
	"zanyZebra",
	"radiantPeacock",
	"mellowMoose",
	"fuzzyBunny",
	"sereneSeal",
	"playfulPuppy",
	"happyHedgehog",
	"charmingChinchilla",
	"mischievousMonkey",
	"jumpyKangaroo",
	"snazzyLizard",
	"spunkyFerret",
	"quirkyQuokka",
	"bubblyGoldfish",
	"shinyButterfly",
	"vibrantChameleon",
	"elegantFlamingo",
	"zestyZucchini",
	"jovialJellyfish",
	"bouncyBumblebee",
	"fancyFalcon",
	"snappyTurtle",
	"chillyChinchilla",
	"giddyGopher",
	"spicySalsa",
	"whimsicalWombat",
	"frostyFawn",
}

var title = []string{
	"The Benefits of Mindfulness Meditation",
	"10 Tips for a Healthier Lifestyle",
	"Exploring the Wonders of Nature",
	"How to Boost Your Productivity",
	"Understanding the Basics of Cryptocurrency",
	"Traveling on a Budget: Tips and Tricks",
	"The Importance of Mental Health Awareness",
	"Delicious and Healthy Recipes for Busy People",
	"How to Create a Sustainable Garden",
	"Mastering the Art of Public Speaking",
	"Top 5 Books to Read This Year",
	"How to Build Stronger Relationships",
	"The Future of Remote Work",
	"Essential Skills for the Modern Workplace",
	"Exploring Different Cultures Through Food",
	"How to Stay Motivated During Tough Times",
	"Understanding Climate Change and Its Impact",
	"Creative Hobbies to Try in Your Free Time",
	"How to Manage Stress Effectively",
	"Tips for a Successful Job Interview",
	"The Power of Positive Thinking",
	"How to Start a Successful Blog",
	"Exploring the World of Digital Marketing",
	"How to Cultivate a Growth Mindset",
	"Understanding the Basics of Personal Finance",
}

var content = []string{
	"Mindfulness meditation can help reduce stress, improve focus, and enhance overall well-being. In this post, we explore various techniques and tips to get started with mindfulness meditation.",
	"Living a healthier lifestyle doesn't have to be complicated. Here are 10 simple tips that can help you improve your diet, exercise routine, and overall health.",
	"Nature has a way of rejuvenating our spirits. Join us as we explore some of the most breathtaking natural wonders around the world and the benefits of spending time outdoors.",
	"Productivity is key to achieving your goals. In this post, we share effective strategies to help you manage your time better and increase your productivity.",
	"Cryptocurrency is changing the financial landscape. This article breaks down the basics of cryptocurrency, how it works, and what you need to know to get started.",
	"Traveling doesn't have to break the bank. Discover practical tips and tricks for exploring new destinations without overspending.",
	"Mental health is just as important as physical health. This post discusses the significance of mental health awareness and how we can support ourselves and others.",
	"Busy schedules can make healthy eating challenging. Here are some quick and nutritious recipes that are perfect for those on the go.",
	"Sustainable gardening is not only good for the environment but also rewarding. Learn how to create a garden that thrives while being eco-friendly.",
	"Public speaking can be daunting, but with practice, anyone can master it. This post offers tips on how to become a confident and effective speaker.",
	"Reading can expand your horizons. Here are five must-read books that will inspire and motivate you this year.",
	"Building strong relationships takes effort. In this article, we discuss key strategies for nurturing and maintaining meaningful connections with others.",
	"Remote work is becoming the norm. Explore the future of remote work and how it is reshaping the way we approach our careers.",
	"Essential skills for the modern workplace include adaptability, communication, and problem-solving. This post highlights the skills you need to thrive in today's job market.",
	"Food is a gateway to understanding different cultures. Join us as we explore various cuisines and the stories behind them.",
	"Staying motivated can be challenging, especially during tough times. Here are some strategies to help you maintain your motivation and resilience.",
	"Climate change is a pressing issue that affects us all. This article explains the basics of climate change and its impact on our planet.",
	"Hobbies can be a great way to express creativity and relieve stress. Discover some creative hobbies you can try in your free time.",
	"Stress management is crucial for maintaining mental health. This post offers practical tips for managing stress effectively.",
	"Job interviews can be nerve-wracking. Here are some essential tips to help you prepare and succeed in your next interview.",
	"Positive thinking can transform your life. Learn about the power of positive thinking and how to cultivate a more optimistic mindset.",
	"Starting a blog can be a fulfilling endeavor. This post provides a step-by-step guide on how to launch a successful blog.",
	"Digital marketing is an ever-evolving field. Explore the latest trends and strategies to enhance your digital marketing skills.",
	"A growth mindset can lead to greater success. This article discusses how to cultivate a growth mindset and embrace challenges.",
	"Creating a morning routine can set a positive tone for your day. In this post, we discuss the elements of a successful morning routine and how to implement it effectively.",
	"Minimalism is about simplifying your life and focusing on what truly matters. This article explores the benefits of living with less and how to start your minimalist journey.",
	"Yoga offers numerous physical and mental health benefits. Join us as we explore different styles of yoga and how they can enhance your well-being.",
	"Improving your writing skills can open up new opportunities. This post provides practical tips and exercises to help you become a better writer.",
	"Artificial intelligence is transforming various industries. In this article, we break down the basics of AI and its potential impact on our future.",
	"Effective time management is crucial for productivity. Discover strategies and tools that can help you manage your time more efficiently.",
	"Working from home presents unique challenges for maintaining health. This post offers tips on how to stay active and healthy while working remotely.",
	"Freelancing can provide flexibility and independence. Explore the pros and cons of freelancing and how to get started in this growing field.",
	"A positive work environment fosters collaboration and creativity. Learn how to cultivate a supportive atmosphere in your workplace.",
	"Emotional intelligence is key to personal and professional success. This article discusses the importance of emotional intelligence and how to develop it.",
}

var tags = []string{
	"Mindfulness",
	"Health",
	"Productivity",
	"Travel",
	"Personal Finance",
	"Self-Improvement",
	"Yoga",
	"Minimalism",
	"Freelancing",
	"Digital Marketing",
	"Emotional Intelligence",
	"Time Management",
	"Sustainability",
	"Public Speaking",
	"Creativity",
	"Motivation",
	"Remote Work",
	"Cooking",
	"Gardening",
	"Technology",
}

var comment = []string{
	"Great post! I learned a lot from it.",
	"I enjoyed reading this article. It was informative and engaging.",
	"This is a great resource for anyone looking to improve their writing skills.",
	"The author did an excellent job of explaining the concepts.",
	"I found the information in this article very helpful.",
	"Great post! Very informative.",
	"I love this topic, thanks for sharing!",
	"Really helpful tips, I appreciate it.",
	"Interesting perspective, I hadn't thought of that.",
	"Thanks for the insights!",
	"Can't wait to try these suggestions!",
	"Very well written, thank you!",
	"I completely agree with this!",
	"Such a valuable read!",
	"Thanks for the motivation!",
	"Looking forward to more posts like this.",
	"Great advice, I'll implement it!",
	"Love the examples you provided.",
	"Such a timely article, thanks!",
	"Really enjoyed this, keep it up!",
	"Fantastic information, very useful!",
	"Thanks for breaking it down so clearly.",
	"Your writing style is engaging!",
	"Appreciate the effort you put into this.",
	"Very inspiring, thank you for sharing!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post: ", err)
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment: ", err)
		}
	}

	log.Println("Data seeded successfully!")
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := 0; i < n; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Role: store.Role{
				Name: "user",
			},
		}
	}
	return users

}

func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		user := users[rand.IntN(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   title[rand.IntN(len(title))],
			Content: content[rand.IntN(len(content))],
			Tags: []string{
				tags[rand.IntN(len(tags))],
				tags[rand.IntN(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, n)
	for i := 0; i < n; i++ {
		user := users[rand.IntN(len(users))]
		post := posts[rand.IntN(len(posts))]
		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: comment[rand.IntN(len(comment))],
		}

	}
	return comments

}
