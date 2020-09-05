//
// This program is a prototype to learm more about jir webhooks
//
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
)

const (
	FROM     = "m.darin@email.su"
	PASSWORD = "x_wK,(Tf6ff)2L_"
)

/*
Игнорировать пустое поле
Чтобы предотвратить включение Location в JSON, когда оно установлено на его нулевое значение, добавьте ,omitempty в тег json .

type Company struct {
    Name     string `json:"name"`
    Location string `json:"location,omitempty"`
}
*/
//-------------------------------
// type defenitios
//------------------------------

// JiraWebhookEvent
type JiraWebhookEvent struct {
	// "timestamp": 1598980214900,
	Timestamp int `json:"timestamp"`
	// "webhookEvent": "jira:issue_updated",
	WebhookEvent string `json:"webhookEvent"`
	// "issue_event_type_name": "issue_updated",
	IssueEventTypeName string `json:"issue_event_type_name"`
	// "user": {},
	User JiraUser `json:"user"`
	// "issue": {},
	Issue JiraIssue `json:"issue"`
	// "changelog": {}
	Changelog JiraChangelog `json:"changelog"`
}

// "user"
type JiraUser struct {
	//   "self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5d0b3e123e70300bc975860e",
	//   "accountId": "5d0b3e123e70300bc975860e",
	AccountID string `json:"accountId"`
	//   "displayName": "Michael DARIN",
	DisplayName string `json:"displayName"`
	//   "active": true,
	Active bool `json:"active"`
	//   "timeZone": "Europe/Moscow",
	//   "accountType": "atlassian"
}

// "issue"
type JiraIssue struct {
	//   "id": "10059",
	Id string `json:"id"`
	//   "self": "https://aeonmeta.atlassian.net/rest/api/2/10059",
	Self string `json:"self"`
	//   "key": "VZQO-34",
	Key string `json:"key"`
	//   "fields": {}
	Fields JiraFields `json:"fields"`
}

// "status"
type JiraStatus struct {
	// "self": "https://aeonmeta.atlassian.net/rest/api/2/status/10011",
	// "description": "This issue is being actively worked on at the moment by the assignee.",
	Desc string `json:"description"`
	// "iconUrl": "https://aeonmeta.atlassian.net/",
	// "name": "In Progress",
	Name string `json:"name"`
	// "id": "10011",
	Id string `json:"id"`
	// "statusCategory": Object
	StatusCategory JiraStatusCategory `json:"statusCategory"`
}

// "statusCategory"
type JiraStatusCategory struct {
	// "self": "https://aeonmeta.atlassian.net/rest/api/2/statuscategory/4",
	// "id": 4,
	Id int `json:"id"`
	// "key": "indeterminate",
	Key string `json:"key"`
	// "colorName": "yellow",
	// "name": "In Progress"
	Name string `json:"name"`
}

// "project"
type JiraProject struct {
	//	"self": "https://aeonmeta.atlassian.net/rest/api/2/project/10004",
	//	"id": "10004",
	Id string `json:"id"`
	//	"key": "VZQO",
	Key string `json:"key"`
	//	"name": "ДЕМО",
	Name string `json:"name"`
	//	"projectTypeKey": "software",
	//	"simplified": true,
}

// "issuetype"
type JiraIssueType struct {
	// "self": "https://aeonmeta.atlassian.net/rest/api/2/issuetype/10017",
	// "id": "10017",
	Id string `json:"id"`
	// "description": "Tasks track small, distinct pieces of work.",
	Desc string `json:"description"`
	// "iconUrl": "https://aeonmeta.atlassian.net/secure/viewavatar?size=medium&avatarId=10318&avatarType=issuetype",
	// "name": "Task",
	Name string `json:"name"`
	// "subtask": false,
	Subtask bool `json:"subtask"`
	// "avatarId": 10318,
	// "entityId": "688ddaac-f832-42ed-8903-0cafa65bf49f"
}

// "watches"
type JiraWatches struct {
	//	"self": "https://aeonmeta.atlassian.net/rest/api/2/issue/VZQO-34/watchers",
	//	"watchCount": 2,
	WatchCount int `json:"watchCount"`
	//	"isWatching": false
	IsWatching bool `json:"isWatching"`
}

// "progress"
type JiraProgress struct {
	//	"progress": 0,
	Progress int `json:"progress"`
	//	"total": 0
	Total int `json:"total"`
}

// "priority"
type JiraPriority struct {
	// 	"self": "https://aeonmeta.atlassian.net/rest/api/2/priority/3",
	// 	"iconUrl": "https://aeonmeta.atlassian.net/images/icons/priorities/medium.svg",
	// 	"name": "Medium",
	Name string `json:"name"`
	// 	"id": "3"
	Id string `json:"id"`
}

// "fields"
type JiraFields struct {
	// "statuscategorychangedate": "2020-03-23T13:57:55.935+0300",
	StatusCategoryChangeDate string `json:"statuscategorychangedate"`
	// "issuetype": {},
	IssueType JiraIssueType `json:"issuetype"`
	// "parent": {},
	Parent JiraParent `json:"parent"`
	// "timespent": null,
	// "project": {},
	Project JiraProject `json:"project"`
	// "aggregatetimespent": null,
	// "resolution": null,
	// "resolutiondate": null,
	// "workratio": -1,
	// "lastViewed": null,
	// "watches": {},
	Watches JiraWatches `json:"watches"`
	// "issuerestriction": {},
	// "created": "2020-03-11T16:21:37.654+0300",
	Created string `json:"created"`
	// "priority": {},
	Priority JiraPriority `json:"priority"`
	// "labels": [
	// 	"Авторизация"
	// ],
	Labels []string `json:"labels"`
	// "timeestimate": null,
	// "aggregatetimeoriginalestimate": null,
	// "versions": [],
	// "issuelinks": [],
	// "assignee": {},
	Assignee JiraUser `json:"assignee"`
	// "updated": "2020-09-01T20:10:14.894+0300",
	Updated string `json:"updated"`
	// "status": {},
	Status JiraStatus `json:"status"`
	// "components": [],
	// "description": "Пофиксить",
	Desc string `json:"description"`
	// "timetracking": {},
	// "security": null,
	// "attachment": [],
	// "aggregatetimeestimate": null,
	// "summary": "modelPhone возвращает номер и модель даже когда не одобрили в админке",
	Summary string `json:"summary"`
	// "creator": {},
	Creator JiraUser `json:"creator"`
	// "subtasks": [],
	// "reporter": {},
	Reporter JiraUser `json:"reporter"`
	// "aggregateprogress": {},
	Aggregateprogress JiraProgress `json:"aggregateprogress"`
	// "environment": null,
	// "duedate": "2020-03-31",
	Duedate string `json:"duedate"`
	// "progress": {},
	Progress JiraProgress `json:"progress"`
	// "votes": {}
}

// "parent"
type JiraParent struct {
	// 	"id": "10038",
	Id string `json:"id"`
	// 	"key": "VZQO-14",
	Key string `json:"key"`
	// 	"self": "https://aeonmeta.atlassian.net/rest/api/2/issue/10038",
	// 	"fields": {
	// 		"summary": "Авторизация",
	// 		"status": Object
	// 		"priority": Object
	// 		"issuetype": Object
	// 	}
	// }
}

// "changelog"
type JiraChangelog struct {
	// "id": "28403",
	Id string `json:"id"`
	// "items": [{}]
	Items []JiraChangelogItem
}

type JiraChangelogItem struct {
	// "field": "Rank",
	Field string `json:"field"`
	// "fieldtype": "custom",
	FieldType string `json:"fieldtype"`
	// "fieldId": "customfield_10019",
	FieldId string `json:"fieldId"`
	// "from": "",
	From string `json:"from"`
	// "fromString": "",
	FromString string `json:"fromString"`
	// "to": "",
	To string `json:"to"`
	// "toString": "Рейтинг понижен"
	ToString string `json:"toString"`
}

// GitLabWebhookEvent
type GitLabWebhookEvent struct {
	// "object_kind": "push",
	ObjectKind string `json:"object_kind"`
	// "event_name": "push",
	EventName string `json:"event_name"`
	// "before": "0000000000000000000000000000000000000000",
	Before string `json:"before"`
	// "after": "70672f9603dbbaa2f217d3b38176886f1a92ad33",
	After string `json:"after"`
	// "ref": "refs/heads/test",
	Ref string `json:"ref"`
	// "checkout_sha": "70672f9603dbbaa2f217d3b38176886f1a92ad33",
	CheckoutSHA string `json:"checkout_sha"`
	// "message": null,
	Message string `json:"message"`
	// "user_id": 36,
	UserId int `json:"user_id"`
	// "user_name": "Касумов Махач",
	UserName string `json:"user_name"`
	// "user_username": "kasumov-mk",
	UserUsername string `josn:"user_username"`
	// "user_email": "",
	UserEmail string `json:"user_email"`
	// "user_avatar": "https://gitlab.aeon.world/uploads/-/system/user/avatar/36/avatar.png",
	// "project_id": 71,
	ProjectId int `json:"project_id"`
	// "project": {},
	Project GitLabProject `json:"project"`
	// "commits": [{}],
	Commits []GitLabCommit `json:"commits"`
	// "total_commits_count": 1,
	TotalCommitsCount int `json:"total_commits_count"`
	// "push_options": {},
	// "repository": {}
	Repository GitLabRepository `json:"repository"`
}

// 	"project"
type GitLabProject struct {
	// "id": 71,
	Id int `json:"id"`
	// "name": "meta_id",
	Name string `json:"name"`
	// "description": "",
	Desc string `json:"description"`
	// "web_url": "https://gitlab.aeon.world/services/meta_id",
	WebUrl string `json:"web_url"`
	// "avatar_url": null,
	// "git_ssh_url": "git@gitlab.aeon.world:services/meta_id.git",
	GitSSHUrl string `json:"git_ssh_url"`
	// "git_http_url": "https://gitlab.aeon.world/services/meta_id.git",
	GitHTTPUrl string `json:"git_http_url"`
	// "namespace": "services",
	Namespace string `json:"namespace"`
	// "visibility_level": 10,
	// "path_with_namespace": "services/meta_id",
	// "default_branch": "dev",
	DefaultBranch string `json:"default_branch"`
	// "ci_config_path": null,
	// "homepage": "https://gitlab.aeon.world/services/meta_id",
	Homepage string `json:"homepage"`
	// "url": "git@gitlab.aeon.world:services/meta_id.git",
	Url string `json:"url"`
	// "ssh_url": "git@gitlab.aeon.world:services/meta_id.git",
	SSHUrl string `json:"ssh_url"`
	// "http_url": "https://gitlab.aeon.world/services/meta_id.git"
	HTTPUrl string `json:"http_url"`
}

// "commit"
type GitLabCommit struct {
	// "id": "70672f9603dbbaa2f217d3b38176886f1a92ad33",
	Id string `json:"id"`
	// "message": "test push\n",
	Message string `josn:"message"`
	// "timestamp": "2020-09-01T18:24:33Z",
	Timestamp string `json:"timestamp"`
	// "url": "https://gitlab.aeon.world/services/meta_id/commit/70672f9603dbbaa2f217d3b38176886f1a92ad33",
	Url string `json:"url"`
	// "author": {},
	Author GitLabAuthor `json:"author"`
	// "added": [],
	Added []string `json:"added"`
	// "modified": ["lib/meta_id.ex"],
	Modified []string `json:"modified"`
	// "removed": []
	Removed []string `json:"removed"`
}

// "author"
type GitLabAuthor struct {
	// "name": "Касумов Махач",
	Name string `json:"name"`
	// "email": "maktempgma@gmail.com"
	Email string `json:"email"`
}

// 	"repository"
type GitLabRepository struct {
	// "name": "meta_id",
	Name string `json:"name"`
	// "url": "git@gitlab.aeon.world:services/meta_id.git",
	Url string `json:"url"`
	// "description": "",
	Desc string `json:"description"`
	// "homepage": "https://gitlab.aeon.world/services/meta_id",
	Homepage string `json:"homepage"`
	// "git_http_url": "https://gitlab.aeon.world/services/meta_id.git",
	GitHTTPUrl string `json:"git_http_url"`
	// "git_ssh_url": "git@gitlab.aeon.world:services/meta_id.git",
	GitSSHUrl string `json:"git_ssh_url"`
	// "visibility_level": 10
}

//-------------------------------
// main driver
//-------------------------------
func main() {
	// configure handlers
	http.HandleFunc("/", handle_hello)
	http.HandleFunc("/hi", handle_hi)
	http.HandleFunc("/jira-webhookreceiver", handle_jira_webhook)
	http.HandleFunc("/world", handle_world)
	http.HandleFunc("/gitlab-push-webhookreceiver", handle_gitlab_push_webhook)

	// start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//-------------------------------
// handlers
//-------------------------------
func handle_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path[1:]))
}

func handle_hi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi")
}

func handle_world(w http.ResponseWriter, r *http.Request) {
	// get body
	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))
	send_email(string(body))
	fmt.Fprintf(w, "[world] webhook machined!")
}

func handle_jira_webhook(w http.ResponseWriter, r *http.Request) {
	// get body
	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))
	var event JiraWebhookEvent
	// get event
	err := json.Unmarshal([]byte(body), &event)
	if err != nil {
		panic(err)
	}
	email := event.User.DisplayName + "\r\n" +
		event.IssueEventTypeName + "\r\n" +
		event.Issue.Fields.Project.Name + "\r\n" +
		event.Issue.Fields.Project.Key + "\r\n" +
		event.Issue.Key + "\r\n" +
		event.Issue.Fields.Summary + "\r\n" +
		event.Issue.Fields.Desc + "\r\n" +
		event.Issue.Fields.Status.Name + "\r\n" +
		event.Issue.Fields.Created + "\r\n" +
		event.Issue.Fields.Updated + "\r\n" +
		event.Issue.Fields.Duedate
	send_email(email)
	fmt.Fprintf(w, "webhook machined!")
}

func handle_gitlab_push_webhook(w http.ResponseWriter, r *http.Request) {
	// get body
	body, _ := ioutil.ReadAll(r.Body)
	log.Println(string(body))

	var event GitLabWebhookEvent
	// get event
	err := json.Unmarshal([]byte(body), &event)
	if err != nil {
		panic(err)
	}

	email := event.EventName + "\r\n" +
		event.Ref + "\r\n" +
		event.Message + "\r\n" +
		event.UserName + "\r\n" +
		event.UserUsername + "\r\n" +
		event.UserEmail + "\r\n" +
		event.Project.Name + "\r\n" +
		event.Project.Namespace + "\r\n" +
		event.Project.Homepage

	send_email(email)
	fmt.Fprintf(w, "gitlab push webhook machined!")
}

//-------------------------------
// internals
//-------------------------------
func send_email(message string) {
	// yandex
	// smtp.yandex.ru:465
	// http://ilyakhasanov.ru/baza-znanij/prochee/nuzhno-znat/139-nastrojki-otpravki-pochty-cherez-smtp
	// how to resolve EOF error
	// https://stackoverflow.com/questions/11662017/smtp-sendmail-fails-after-10-minutes-with-eof

	log.Println("auth")
	// Choose auth method and set it up
	auth := smtp.PlainAuth("", FROM, PASSWORD, "smtp.yandex.ru")
	log.Println("prepare and send")
	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"m.darin.comco@yandex.ru"}
	msg := []byte("To: m.darin.comco@yandex.ru\r\n" +
		"Subject: JIRA webhook machined\r\n" +
		"\r\n" +
		message +
		"\r\n")
	err := smtp.SendMail("smtp.yandex.ru:587", auth, FROM, to, msg)
	if err != nil {
		log.Print("Error: ")
		log.Fatal(err)
	} else {
		log.Println("Successfull done!")
	}
}
