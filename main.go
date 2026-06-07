package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type Post struct {
	Slug     string
	Title    string
	Excerpt  string
	Category string
	Date     string
	Body     []string
}

type Item struct {
	Slug        string
	Name        string
	Description string
	Price       string
	Includes    []string
}

type HomeData struct {
	Posts []Post
	Items []Item
}

var posts = []Post{
	{
		Slug:     "websites-that-work-harder",
		Title:    "Websites that work harder",
		Excerpt:  "A useful website should answer questions, qualify leads, and make the next step obvious.",
		Category: "Web",
		Date:     "1 June 2026",
		Body: []string{
			"A website should earn its place in the business. It should explain what you do, answer the obvious questions, and guide the right people toward a clear next step.",
			"Start with the outcome. Decide what a visitor needs to understand and what the business needs them to do next, then remove everything that gets in the way.",
			"Good websites are not necessarily complicated. Clear writing, fast pages, and sensible structure usually do more work than a long list of features.",
		},
	},
	{
		Slug:     "your-website-is-an-operating-system",
		Title:    "Your website is an operating system, not a brochure",
		Excerpt:  "The best small-business sites remove admin, qualify leads, and make the next action obvious.",
		Category: "Strategy",
		Date:     "28 May 2026",
		Body: []string{
			"A brochure tells people that a business exists. A useful website helps the business operate.",
			"It can set expectations, collect the right information, explain pricing, and send people to the right place before anyone has to answer an email.",
			"The goal is not automation for its own sake. The goal is fewer repeated tasks and a calmer path for both the customer and the business.",
		},
	},
	{
		Slug:     "five-signs-your-site-needs-a-rebuild",
		Title:    "Five signs your site needs a rebuild",
		Excerpt:  "Slow edits, vague messaging, and mobile compromises are signs the foundation has stopped serving you.",
		Category: "Web",
		Date:     "14 May 2026",
		Body: []string{
			"A rebuild makes sense when the current website creates more friction than value.",
			"Common signs include simple changes feeling risky, the mobile version being treated as an afterthought, unclear ownership, slow pages, and content that no longer reflects the business.",
			"Do not rebuild simply because a site is old. Rebuild when a clearer and simpler foundation will make the business easier to run.",
		},
	},
	{
		Slug:     "hosting-should-be-boring",
		Title:    "Hosting should be boring",
		Excerpt:  "Fast pages, valid certificates, working backups, and clear ownership are the real product.",
		Category: "Infrastructure",
		Date:     "8 May 2026",
		Body: []string{
			"Reliable hosting disappears into the background. Pages load, certificates renew, backups work, and nobody needs to think about it every morning.",
			"The minimum standard is straightforward: monitoring, backups, TLS, updates, and a documented recovery process.",
			"Complex dashboards are not a substitute for clear responsibility. Someone should know what is running, where it lives, and how to restore it.",
		},
	},
	{
		Slug:     "automation-without-the-theatre",
		Title:    "Automation without the theatre",
		Excerpt:  "Small, dependable automations often create more value than complicated platforms.",
		Category: "Automation",
		Date:     "12 April 2026",
		Body: []string{
			"Useful automation removes a repeated handoff. It does not need to transform the entire company.",
			"Start with a task that is predictable, frequent, and annoying. Make the result visible and keep manual recovery simple.",
			"A small workflow that runs every day is more valuable than an ambitious system nobody trusts.",
		},
	},
}

var items = []Item{
	{
		Slug:        "website-launch",
		Name:        "Website Launch",
		Description: "A focused small-business website that explains the offer clearly and works beautifully on every screen.",
		Price:       "From $650",
		Includes:    []string{"Up to four pages", "Responsive build", "Basic SEO setup", "Launch support"},
	},
	{
		Slug:        "website-grow",
		Name:        "Website Grow",
		Description: "A larger content-led website for businesses ready to publish, integrate, and build momentum.",
		Price:       "From $1,450",
		Includes:    []string{"Content-led page system", "Blog or resource section", "Simple integrations", "Performance review"},
	},
	{
		Slug:        "hosting-care",
		Name:        "Hosting and Care",
		Description: "Managed hosting, monitoring, backups, and updates without the mystery or dashboard overload.",
		Price:       "$20 / month",
		Includes:    []string{"Managed hosting", "Uptime monitoring", "Regular backups", "Software updates"},
	},
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func init() {
	template.Must(templates.ParseGlob("templates/partials/*.html"))
}

func render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "layout.html", HomeData{Posts: posts[:3], Items: items})
}

func homePartial(w http.ResponseWriter, r *http.Request) {
	render(w, "home.html", HomeData{Posts: posts[:3], Items: items})
}

func blog(w http.ResponseWriter, r *http.Request) {
	render(w, "blog.html", posts)
}

func postPage(w http.ResponseWriter, r *http.Request) {
	for _, post := range posts {
		if post.Slug == r.PathValue("slug") {
			render(w, "post.html", post)
			return
		}
	}
	http.NotFound(w, r)
}

func store(w http.ResponseWriter, r *http.Request) {
	render(w, "store.html", items)
}

func itemPage(w http.ResponseWriter, r *http.Request) {
	for _, item := range items {
		if item.Slug == r.PathValue("slug") {
			render(w, "item.html", item)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/home", homePartial)
	http.HandleFunc("/blog", blog)
	http.HandleFunc("/blog/{slug}", postPage)
	http.HandleFunc("/store", store)
	http.HandleFunc("/store/{slug}", itemPage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8091"
	}

	log.Printf("Aevra is running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
