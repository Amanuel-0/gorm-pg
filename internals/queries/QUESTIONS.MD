To help you **progressively build ORM mastery**, here’s a structured **skill roadmap**:

---

## 🧩 **LEVEL 1: Basic ORM CRUD & Relations**

*Focus: get comfortable with models, simple queries, inserts, updates, and foreign keys.*

### ✅ Simple tasks

 - [x] Create a new user with a profile.
 - [x] Fetch a user by email.
 - [x] Update a user’s display name and avatar.
 - [x] Soft delete a user (set `deleted_at`).
 - [x] Create an author and a few genres.
 - [x] Create a book with an author, owner and assign multiple genres.
 - [x] Retrieve all books owned by a specific user.
 - [x] Fetch a book including its author and genres.
 - [x] List all users who prefer a certain genre.
 - [x] Add a new subscription plan.

### ⚙️ Intermediate (1–N and N–M)

 - [x] Create a user and link them to multiple preferred genres.
 - [ ] Create a book and attach multiple `book_images`. (kind of same with other question)
 - [x] Find all books with “like_new” condition by a specific author.
 - [x] List all books in a specific city or country.
 - [x] Get all genres associated with a given book (many-to-many).

---

## 🧭 **LEVEL 2: Querying & Filtering**

*Focus: filtering, joins, eager loading, query builders.*

 - [x] Fetch all books available between two dates.
 - [ ] Find all users who haven’t verified their email. (I don't think there is enough info to perform this, so passed.)
 - [x] Retrieve all active subscriptions and their plans.
 - [x] Find all users whose subscription is expired.
 - [ ] List all users who have ever made a successful payment. (I don't think there is enough ifno to perform this, so passed).
 - [x] Count how many books each user owns.
 - [x] Retrieve all books with their review averages.
 - [x] List all exchanges in the “requested” state with book and user info.
 - [x] Fetch messages in a chat thread, sorted newest → oldest.
 - [x] Find all users who haven’t logged in for 30+ days.

---

## 💳 **LEVEL 3: Transactions & Lifecycle**

*Focus: multi-table operations, transactions, and complex updates.*

 - [x] When a new user subscribes:
   - [x] Create a new `subscription` record.
   - [x] Create an initial `payment` record. --  (I don't think we have this)
   - [x] Update user’s role to `premium` (or store in metadata). -- (we don't have such kind of relation / but, an alternative is a subscription plan )
 - [x] On book deletion:
   - [x] Soft delete the book (`archived_at`).
   - [x] Cascade delete its images and related `book_genres`. (i.e. cascade works if it is not gorm.DeletedAt(this by default gives a soft delete))
 - [x] Create a transaction that handles an exchange:
   - [x] Mark exchange as `completed`.
   - [x] Update both books as unavailable.
   - [x] Insert two user ratings.
 - [x] Create a function to cancel a subscription:
   - [x] Set `status = 'canceled'`.
   - [x] Update `current_period_end` and `cancel_at_period_end`.
   - [ ] Log the change in `activity_logs`. -- (not necessary)
 - [x] Create a function that reports a user:
   - [x] Insert into `reports`.
   - [x] Create a `notification` for the admin.

---

## 💬 **LEVEL 4: Relations Across Domains**

*Focus: joins, nested preloading, polymorphism patterns.*

 - [x] Get all chat threads for a given exchange, including messages and senders.
 - [x] Fetch all community threads and their messages (with author info).
 - [x] List all users who belong to at least 2 communities.
 - [x] Find communities that require paid chat (`require_paid_chat = 1`).
 - [ ] Get all notifications for a user, with type-specific payload parsing.
 - [ ] Retrieve all exchanges involving a particular user (as requester or responder).
 - [ ] For a given user, list all books they offered in exchanges that are now `completed`.

---

## ⚖️ **LEVEL 5: Aggregations, Analytics & Advanced Queries**

*Focus: group by, counts, averages, JSON, subqueries, CTEs.*

 - [ ] Find top 5 users by total number of books owned.
 - [ ] Find authors with the most books listed.
 - [ ] Calculate the average rating per user from `user_ratings`.
 - [ ] List books with more than 3 reviews and an average rating > 4.
 - [ ] Count the number of messages sent per user per month.
 - [ ] Find users who have never participated in an exchange.
 - [ ] Calculate total revenue per month from `payments`.
 - [ ] List subscription plans ranked by active subscriber count.
 - [ ] Find users who have disputed exchanges.
 - [ ] Identify the most active communities (by number of messages).

---

## 🔐 **LEVEL 6: Security, Auditing & Automation**

*Focus: triggers, hooks, middleware, and audit trails.*

 - [ ] Implement a model hook to automatically log activity (`activity_logs`) whenever a user updates their profile.
 - [ ] Automatically delete expired tokens via scheduled job.
 - [ ] Implement a trigger (or ORM event hook) to mark a message quota record when a user sends a message.
 - [ ] Add middleware that automatically preloads related data (like user → profile).
 - [ ] Automatically set `email_verified_at` when user’s email is verified.

---

## 🧠 **LEVEL 7: Complex Application Features**

*Focus: real-world logic combining multiple domains.*

 - [ ] “Recommended books” — find books that match a user’s preferred genres and are nearby (within 50km).
 - [ ] “Mutual swap candidates” — users who own books another wants (via `preferred_titles`).
 - [ ] “Exchange history” — show exchanges between two users.
 - [ ] “Top reviewers” — users who wrote the most `book_reviews`.
 - [ ] “Subscription churn” — find users who canceled in the last 30 days.
 - [ ] “Moderation dashboard” — list pending `reports` with reporter, target, and handler details.
 - [ ] “Community leaderboards” — top posters in each community.
 - [ ] “Geo heatmap” — count of books by location (grouped by country/state).

---

## 🧩 Bonus Challenges

 - [ ] Implement soft delete logic globally (e.g., for users, books).
 - [ ] Add caching layer or query batching (ORM optimization).
 - [ ] Implement pagination with cursor-based queries.
 - [ ] Write a custom ORM query for “users with no profile.”
 - [ ] Implement “lazy loading vs eager loading” experiments and measure query count differences.