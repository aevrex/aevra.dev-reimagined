package main

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
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

type PageData struct {
	Content string
	Posts   []Post
	Items   []Item
	Post    Post
	Item    Item
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

func loadTemplates() *template.Template {
	t := template.New("")
	var renderFn func(name string, data any) template.HTML
	renderFn = func(name string, data any) template.HTML {
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, name, data); err != nil {
			return template.HTML("<!-- template render error: " + err.Error() + " -->")
		}
		return template.HTML(buf.String())
	}

	t = template.Must(t.Funcs(template.FuncMap{
		"render": renderFn,
	}).ParseGlob("templates/*.html"))
	template.Must(t.ParseGlob("templates/partials/*.html"))
	return t
}

var templates = loadTemplates()

func render(w http.ResponseWriter, name string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	if isHTMX(r) {
		render(w, "home.html", PageData{Posts: posts[:3], Items: items})
		return
	}

	if r.URL.Path == "/home" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	render(w, "layout.html", PageData{Content: "home.html", Posts: posts[:3], Items: items})
}

func blog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.NotFound(w, r)
		return
	}

	if isHTMX(r) {
		render(w, "blog.html", PageData{Posts: posts})
		return
	}

	render(w, "layout.html", PageData{Content: "blog.html", Posts: posts})
}

func postPage(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/blog/") {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/blog/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	for _, post := range posts {
		if post.Slug == slug {
			if isHTMX(r) {
				render(w, "post.html", PageData{Post: post})
				return
			}
			render(w, "layout.html", PageData{Content: "post.html", Post: post})
			return
		}
	}

	http.NotFound(w, r)
}

func store(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/store" {
		http.NotFound(w, r)
		return
	}

	if isHTMX(r) {
		render(w, "store.html", PageData{Items: items})
		return
	}

	render(w, "layout.html", PageData{Content: "store.html", Items: items})
}

func itemPage(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/store/") {
		http.NotFound(w, r)
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/store/")
	if slug == "" || strings.Contains(slug, "/") {
		http.NotFound(w, r)
		return
	}

	for _, item := range items {
		if item.Slug == slug {
			if isHTMX(r) {
				render(w, "item.html", PageData{Item: item})
				return
			}
			render(w, "layout.html", PageData{Content: "item.html", Item: item})
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/", home)
	http.HandleFunc("/blog", blog)
	http.HandleFunc("/blog/", postPage)
	http.HandleFunc("/store", store)
	http.HandleFunc("/store/", itemPage)

	http.ListenAndServe(":8080", nil)
}
