// The ghissues package provides simple hooks into Github's Issues API
package ghissues

import (
	"http"
	"fmt"
	"json"
	"io/ioutil"
        "os"
)

const base_api_url = "http://github.com/api/v2/json"

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

type PullRequest struct { // @@@ Unimplemented
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

func (ic *IssuesClient) post(url string, data map[string]string) (*http.Response, os.Error) {
	if _, username_exists := data["login"]; !username_exists {
		data["login"] = ic.username
	}
	if _, token_exists := data["token"]; !token_exists {
		data["token"] = ic.token
	}
	return ic.client.PostForm(url, data)
}

func (ic *IssuesClient) get(url string) (*http.Response, os.Error) {
	response, _, err := ic.client.Get(url)
	if response.StatusCode != 200 {
		return response, os.NewError(
			fmt.Sprintf("Got a %v status code on fetch of %v.", response.StatusCode, url))
	}
	return response, err
}

func (ic *IssuesClient) parseJson(response *http.Response, toStructure interface{}) (interface{}, os.Error) {
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
                return toStructure, err
	}
	err2 := json.Unmarshal(b, toStructure)
	return toStructure, err2
}

func (ic *IssuesClient) Search(user, repo, state, term string) ([]Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/search/%v/%v/%v/%v/", base_api_url, user, repo, state, term)
	response, err := ic.get(url_string)
	if err != nil {
                return nil, err
	}
	json, err2 := ic.parseJson(response, new(multipleIssueResponse))
	return json.(*multipleIssueResponse).Issues, err2
}

func (ic *IssuesClient) List(user, repo, state string) ([]Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/list/%v/%v/%v/", base_api_url, user, repo, state)
        response, err := ic.get(url_string)
	if err != nil {
                return nil, err
	}
	json, err2 := ic.parseJson(response, new(multipleIssueResponse))
	return json.(*multipleIssueResponse).Issues, err2
}

func (ic *IssuesClient) Create(user, repo, title, body string) (Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/open/%v/%v/", base_api_url, user, repo)
	post_data := make(map[string]string){
        "title":title,
        "body":body,
    }
	response, err := ic.post(url_string, post_data)
	if err != nil {
                return Issue{}, err
	}
	json, err2 := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue, err2
}

func (ic *IssuesClient) Detail(user, repo string, issueNumber int) (Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/show/%v/%v/%v", base_api_url, user, repo, issueNumber)
	response, err := ic.get(url_string)
	if err != nil {
                return Issue{}, err
	}
	json, err2 := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue, err2
}

func (ic *IssuesClient) Edit(user, repo string, issueNumber int, title, body string) (Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/edit/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string){
        "title":title,
        "body":body,
    }
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return Issue{}, err
        }
	json, err2 := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue, err2
}

func (ic *IssuesClient) Close(user, repo string, issueNumber int) (Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/close/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return Issue{}, err
        }
	json, err2 := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue, err2
}

func (ic *IssuesClient) Reopen(user, repo string, issueNumber int) (Issue, os.Error) {
	url_string := fmt.Sprintf("%v/issues/reopen/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return Issue{}, err
        }
	json, err2 := ic.parseJson(response, new(singleIssueResponse))
	return json.(*singleIssueResponse).Issue, err2
}

func (ic *IssuesClient) ListComments(user, repo string, issueNumber int) ([]Comment, os.Error) {
	url_string := fmt.Sprintf("%v/issues/comments/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	response, err := ic.get(url_string)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleCommentResponse))
	return json.(*multipleCommentResponse).Comments, err2
}

func (ic *IssuesClient) AddComment(user, repo string, issueNumber int, comment string) (Comment, os.Error) {
	url_string := fmt.Sprintf("%v/issues/comment/%v/%v/%v/", base_api_url, user, repo, issueNumber)
	post_data := make(map[string]string){
        "comment":comment,
    }
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return Comment{}, err
        }
	json, err2 := ic.parseJson(response, new(singleCommentResponse))
	return json.(*singleCommentResponse).Comment, err2
}

func (ic *IssuesClient) ListLabels(user, repo string) ([]Label, os.Error) {
	url_string := fmt.Sprintf("%v/issues/labels/%v/%v/", base_api_url, user, repo)
	response, err := ic.get(url_string)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels, err2
}

func (ic *IssuesClient) AddLabelToRepo(user, repo, label string) ([]Label, os.Error) {
	url_string := fmt.Sprintf("%v/issues/label/add/%v/%v/%v/", base_api_url, user, repo, label)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels, err2
}

func (ic *IssuesClient) AddLabelToIssue(user, repo string, issueNumber int, label string) ([]Label, os.Error) {
	url_string := fmt.Sprintf("%v/issues/label/add/%v/%v/%v/%v/", base_api_url, user, repo, label, issueNumber)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels, err2
}

func (ic *IssuesClient) RemoveLabelFromRepo(user, repo, label string) ([]Label, os.Error) {
	url_string := fmt.Sprintf("%v/issues/label/remove/%v/%v/%v/", base_api_url, user, repo, label)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels, err2
}

func (ic *IssuesClient) RemoveLabelFromIssue(user, repo string, issueNumber int, label string) ([]Label, os.Error) {
	url_string := fmt.Sprintf("%v/issues/label/remove/%v/%v/%v/%v/", base_api_url, user, repo, label, issueNumber)
	post_data := make(map[string]string)
	response, err := ic.post(url_string, post_data)
        if err != nil {
                return nil, err
        }
	json, err2 := ic.parseJson(response, new(multipleLabelResponse))
	return json.(*multipleLabelResponse).Labels, err2
}
