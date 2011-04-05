// The ghissues package provides simple hooks into Github's Issues API
package ghissues

import (
	"http"
	"fmt"
	"json"
	"io/ioutil"
)

const base_api_url string = "http://github.com/api/v2/json"

type JsonError struct {
	error string
}

type Issue struct {
	Gravatar_id string
	Position float32
	Number int
	Votes int
	Created_at string
	Comments int
	Body string
	Title string
	Updated_at string
	Html_url string
	User string
	Labels []Label
	State string
}
type Comment struct {
	Gravatar_id string
	Created_at string
	Body string
	Updated_at string
	Id int
	User string
}
type Label string
// @@@ TODO
type PullRequest struct {
	issue Issue
	pull_request_url string
	html_url string
	patch_url string
}
type IssuesClient struct {
	username string
	token string
	client *http.Client
}

// Responses
type multipleIssueResponse struct {
	Issues []Issue
}
type multipleCommentResponse struct {
	Comments []Comment
}
type singleIssueResponse struct {
	Issue Issue
}
type singleCommentResponse struct {
	Comment Comment
}
type multipleLabelResponse struct {
	Labels []Label
}

func NewClient(username, token string) *IssuesClient {
	return &IssuesClient{username, token, new(http.Client)}
}
func (ic *IssuesClient) post(url string, data map[string]string) *http.Response {
	if _, username_exists := data["login"]; !username_exists {
		data["login"] = ic.username
	}
	if _, token_exists := data["token"]; !token_exists {
		data["token"] = ic.token
	}
	response, err := ic.client.PostForm(url, data)
	if (err != nil) {
		fmt.Printf("Fetch error: %v", err.String())
	}
	return response
}
func (ic *IssuesClient) parseJson(response *http.Response, toStructure interface{}) interface{} {
	b, err := ioutil.ReadAll(response.Body)
	if (err != nil) {
		// @@@ what to do here? Panic? Possibly return an error type.
		fmt.Printf("Error reading response body. " + err.String())
	}
	err2 := json.Unmarshal(b, toStructure)
	if (err2 != nil) {
		fmt.Printf("Unmarshalling Error: " + err2.String())
	}
	return toStructure
}
func (ic *IssuesClient) Search(user, repo, state, term string) []Issue {
	url_string := fmt.Sprintf("%v/issues/search/%v/%v/%v/%v/", base_api_url, user, repo, state, term)
	response, _, err := ic.client.Get(url_string)
	if (err != nil) {
		fmt.Printf("Fetch error: " + err.String())
	}
	json := ic.parseJson(response, new(multipleIssueResponse))
	fmt.Printf("%v", json)
	return json.(*multipleIssueResponse).Issues
}
func (ic *IssuesClient) List(user, repo, state string) []Issue {
	url_string := fmt.Sprintf("%v/issues/list/%v/%v/%v/", base_api_url, user, repo, state)
	response, _, err := ic.client.Get(url_string)
	if (err != nil) {
		fmt.Printf("Fetch error: " + err.String())
	}
	json := ic.parseJson(response, new(multipleIssueResponse))
	return json.(*multipleIssueResponse).Issues
}
func (ic *IssuesClient) Create(user, repo, title, body string) Issue {
	url_string := fmt.Sprintf("%v/issues/open/%v/%v/", base_api_url, user, repo)
	post_data := make(map[string]string)
	post_data["title"] = title
	post_data["body"] = body
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue
}
func (ic *IssuesClient) Detail(user, repo string, issueNumber int) Issue {
	url_string := fmt.Sprintf("%v/issues/show/%v/%v/%v", base_api_url, user, repo, issueNumber)
	response, _, err := ic.client.Get(url_string)
	if (err != nil) {
		fmt.Printf("Fetch error: " + err.String())
	}
	json := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue
}
func (ic *IssuesClient) Edit(user, repo string, issueNumber int, title, body string) Issue {
	url_string := fmt.Sprintf("%v/issues/edit/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	post_data["title"] = title
	post_data["body"] = body
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue
}
func (ic *IssuesClient) Close(user, repo string, issueNumber int) Issue {
	url_string := fmt.Sprintf("%v/issues/close/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue
}
func (ic *IssuesClient) Reopen(user, repo string, issueNumber int) Issue {
	url_string := fmt.Sprintf("%v/issues/reopen/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue
}

func (ic *IssuesClient) ListComments(user, repo string, issueNumber int) []Comment {
	url_string := fmt.Sprintf("%v/issues/comments/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	response, _, err := ic.client.Get(url_string)
	if (err != nil) {
		fmt.Printf("Fetch error: " + err.String())
	}
	json := ic.parseJson(response, new(multipleCommentResponse))
	return json.(*multipleCommentResponse).Comments
}
func (ic *IssuesClient) AddComment(user, repo string, issueNumber int, comment string) Comment {
	url_string := fmt.Sprintf("%v/issues/comment/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	post_data["comment"] = comment
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(singleCommentResponse))
	return json.(*singleCommentResponse).Comment
}

func (ic *IssuesClient) ListLabels(user, repo string) []Label {
	url_string := fmt.Sprintf("%v/issues/labels/%v/%v/", base_api_url, user, repo)
	response, _, err := ic.client.Get(url_string)
	if (err != nil) {
		fmt.Printf("Fetch error: " + err.String())
	}
	json := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels
}

func (ic *IssuesClient) AddLabelToRepo(user, repo, label string) []Label {
	url_string := fmt.Sprintf("%v/issues/label/add/%v/%v/%v/", base_api_url, user, repo, label)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels
}

func (ic *IssuesClient) AddLabelToIssue(user, repo string, issueNumber int, label string) []Label {
	url_string := fmt.Sprintf("%v/issues/label/add/%v/%v/%v/%v/", base_api_url, user, repo, label, issueNumber)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels
}

func (ic *IssuesClient) RemoveLabelFromRepo(user, repo, label string) []Label {
	url_string := fmt.Sprintf("%v/issues/label/remove/%v/%v/%v/", base_api_url, user, repo, label)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels
}

func (ic *IssuesClient) RemoveLabelFromIssue(user, repo string, issueNumber int, label string) []Label {
	url_string := fmt.Sprintf("%v/issues/label/remove/%v/%v/%v/%v/", base_api_url, user, repo, label, issueNumber)
	post_data := make(map[string]string)
	response := ic.post(url_string, post_data)
	json := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels
}

func main() {
	c := NewClient("justinlilly", "ce87a1af897da128ac9a98059bfe2a41")
	// LIST
	// list := c.List("justinlilly", "justinlilly.github.com", "open")
	// SEARCH
	// list := c.Search("justinlilly", "justinlilly.github.com", "open", "curl")
	// for i := range list {
	// 	item := list[i]
	// 	fmt.Printf("%v: %s", item.Number, item.Title)
	// }
	// CREATE
	// created := c.Create("justinlilly", "justinlilly.github.com", "another.", "wee.")
	// fmt.Println("Got: " + fmt.Sprintf("%v", created))
	// DETAIL
	// issue := c.Detail("justinlilly", "justinlilly.github.com", 5)
	// fmt.Printf("#%v %v: %v", issue.Number, issue.Title, issue.Body))
	// LIST COMMENTS
	// list := c.ListComments("justinlilly", "justinlilly.github.com", 5)
	// for i := range list {
	// 	item := list[i]
	// 	fmt.Printf("%v: %s", item.User, item.Body)
	// }
	// CREATE COMMENT
	// item := c.AddComment("justinlilly", "justinlilly.github.com", 5, "This is my test API comment.")
	// fmt.Printf("%v: %v", item.User, item.Body)
	// LIST LABELS
	// list := c.ListLabels("justinlilly", "justinlilly.github.com")
	// for i := range list {
	// 	item := list[i]
	// 	fmt.Printf("%v", item)
	// }
	// ADD LABEL TO REPO
	c.RemoveLabelFromRepo("justinlilly", "justinlilly.github.com", "viaAPI")
}
