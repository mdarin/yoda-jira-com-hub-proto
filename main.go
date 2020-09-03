//
// This program is a prototype to learm more about jir webhooks
//
package main

import (
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

//-------------------------------
// type defenitios
//------------------------------
type JiraWebhookEvent struct {
	// "timestamp": 1598980214900,
	Timestamp string `json:"timestamp"`
	// "webhookEvent": "jira:issue_updated",
	WebhookEvent string `json:"webhookEvent"`
	// "issue_event_type_name": "issue_updated",
	IssueEventTypeName string `json:"issue_event_type_name"`
	// "user": {
	//   "self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5d0b3e123e70300bc975860e",
	//   "accountId": "5d0b3e123e70300bc975860e",
	//   "displayName": "Michael DARIN",
	//   "active": true,
	//   "timeZone": "Europe/Moscow",
	//   "accountType": "atlassian"
	// },
	User JiraUser `json:"user"`
	// "issue": {
	//   "id": "10059",
	//   "self": "https://aeonmeta.atlassian.net/rest/api/2/10059",
	//   "key": "VZQO-34",
	//   "fields": {
	// 		"statuscategorychangedate": "2020-03-23T13:57:55.935+0300",
	// 		"issuetype": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/issuetype/10017",
	// 	  		"id": "10017",
	// 	  		"description": "Tasks track small, distinct pieces of work.",
	// 	  		"iconUrl": "https://aeonmeta.atlassian.net/secure/viewavatar?size=medium&avatarId=10318&avatarType=issuetype",
	// 	  		"name": "Task",
	// 	  		"subtask": false,
	// 	  		"avatarId": 10318,
	// 	  		"entityId": "688ddaac-f832-42ed-8903-0cafa65bf49f"
	// 		},
	// 		"parent": {
	// 	  		"id": "10038",
	// 	  		"key": "VZQO-14",
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/issue/10038",
	// 	  		"fields": {
	// 				"summary": "Авторизация",
	// 				"status": {
	// 		  			"self": "https://aeonmeta.atlassian.net/rest/api/2/status/10011",
	// 		  			"description": "This issue is being actively worked on at the moment by the assignee.",
	// 		  			"iconUrl": "https://aeonmeta.atlassian.net/",
	// 		  			"name": "In Progress",
	// 		  			"id": "10011",
	// 		  			"statusCategory": {
	// 						"self": "https://aeonmeta.atlassian.net/rest/api/2/statuscategory/4",
	// 						"id": 4,
	// 						"key": "indeterminate",
	// 						"colorName": "yellow",
	// 						"name": "In Progress"
	// 		 			}
	// 				},
	// 				"priority": {
	// 		  			"self": "https://aeonmeta.atlassian.net/rest/api/2/priority/3",
	// 		  			"iconUrl": "https://aeonmeta.atlassian.net/images/icons/priorities/medium.svg",
	// 		  			"name": "Medium",
	// 		  			"id": "3"
	// 				},
	// 				"issuetype": {
	// 		  			"self": "https://aeonmeta.atlassian.net/rest/api/2/issuetype/10018",
	// 		  			"id": "10018",
	// 		  			"description": "Эпики позволяют отслеживать похожие истории, задания и связанные с ними ошибки.",
	// 		  			"iconUrl": "https://aeonmeta.atlassian.net/secure/viewavatar?size=medium&avatarId=10307&avatarType=issuetype",
	// 		  			"name": "Эпик",
	// 		  			"subtask": false,
	// 		  			"avatarId": 10307,
	// 		  			"entityId": "856c090e-3ad0-4622-864d-d3e6e35d1be1"
	// 				}
	// 	  		}
	// 		},
	// 		"timespent": null,
	// 		"project": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/project/10004",
	// 	  		"id": "10004",
	// 	  		"key": "VZQO",
	// 	  		"name": "ДЕМО",
	// 	  		"projectTypeKey": "software",
	// 	  		"simplified": true,
	// 		},
	// 		"aggregatetimespent": null,
	// 		"resolution": null,
	// 		"resolutiondate": null,
	// 		"workratio": -1,
	// 		"lastViewed": null,
	// 		"watches": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/issue/VZQO-34/watchers",
	// 	  		"watchCount": 2,
	// 	  		"isWatching": false
	// 		},
	// 		"issuerestriction": {
	// 	  		"issuerestrictions": {},
	// 	  		"shouldDisplay": true
	// 		},
	// 		"created": "2020-03-11T16:21:37.654+0300",
	// 		"priority": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/priority/3",
	// 	  		"iconUrl": "https://aeonmeta.atlassian.net/images/icons/priorities/medium.svg",
	// 	  		"name": "Medium",
	// 	  		"id": "3"
	// 		},
	// 		"labels": [
	// 	  		"Авторизация"
	// 		],
	// 		"timeestimate": null,
	// 		"aggregatetimeoriginalestimate": null,
	// 		"versions": [],
	// 		"issuelinks": [],
	// 		"assignee": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5b22826d58be252d7452ef66",
	// 	  		"accountId": "5b22826d58be252d7452ef66",
	// 	  		"displayName": "Andre Orlov",
	// 	  		"active": false,
	// 	  		"timeZone": "Europe/Moscow",
	// 	  		"accountType": "atlassian"
	// 		},
	// 		"updated": "2020-09-01T20:10:14.894+0300",
	// 		"status": {
	// 			"self": "https://aeonmeta.atlassian.net/rest/api/2/status/10012",
	// 	  		"description": "",
	// 	  		"iconUrl": "https://aeonmeta.atlassian.net/",
	// 	  		"name": "Done",
	// 	  		"id": "10012",
	// 	  		"statusCategory": {
	// 				"self": "https://aeonmeta.atlassian.net/rest/api/2/statuscategory/4",
	// 				"id": 4,
	// 				"key": "indeterminate",
	// 				"colorName": "yellow",
	// 				"name": "In Progress"
	// 	 		}
	// 		},
	// 		"components": [],
	// 		"description": "Пофиксить",
	// 		"timetracking": {},
	// 		"security": null,
	// 		"attachment": [],
	// 		"aggregatetimeestimate": null,
	// 		"summary": "modelPhone возвращает номер и модель даже когда не одобрили в админке",
	// 		"creator": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5d8cce09e4b7210dd0a676e7",
	// 	  		"accountId": "5d8cce09e4b7210dd0a676e7",
	// 	  		"displayName": "Сергей Соловьев",
	// 	  		"active": false,
	// 	  		"timeZone": "Europe/Moscow",
	// 	  		"accountType": "atlassian"
	// 		},
	// 		"subtasks": [],
	// 		"reporter": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5d8cce09e4b7210dd0a676e7",
	// 	  		"accountId": "5d8cce09e4b7210dd0a676e7",
	// 	  		"displayName": "Сергей Соловьев",
	// 	  		"active": false,
	// 	  		"timeZone": "Europe/Moscow",
	// 	  		"accountType": "atlassian"
	// 		},
	// 		"aggregateprogress": {
	// 	  		"progress": 0,
	// 	  		"total": 0
	// 		},
	// 		"environment": null,
	// 		"duedate": "2020-03-31",
	// 		"progress": {
	// 	  		"progress": 0,
	// 	  		"total": 0
	// 		},
	// 		"votes": {
	// 	  		"self": "https://aeonmeta.atlassian.net/rest/api/2/issue/VZQO-34/votes",
	// 	  		"votes": 0,
	// 	  		"hasVoted": false
	// 		}
	//   }
	// },
	// "changelog": {
	//   "id": "28403",
	//   "items": [{
	// 	  	"field": "Rank",
	// 		"fieldtype": "custom",
	// 	  	"fieldId": "customfield_10019",
	// 	  	"from": "",
	// 	  	"fromString": "",
	// 	  	"to": "",
	// 	  	"toString": "Рейтинг понижен"
	// 	}]
	// }
	// }
}

// "user": {
type JiraUser struct {
	//   "self": "https://aeonmeta.atlassian.net/rest/api/2/user?accountId=5d0b3e123e70300bc975860e",
	//   "accountId": "5d0b3e123e70300bc975860e",
	AccountID string `json:"accountId"`
	//   "displayName": "Michael DARIN",
	DislayName string `json:"displayName"`
	//   "active": true,
	Active bool `json:"active"`
	//   "timeZone": "Europe/Moscow",
	//   "accountType": "atlassian"
	// }
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
	send_email(string(body))
	fmt.Fprintf(w, "webhook machined!")
}

func handle_gitlab_push_webhook(w http.ResponseWriter, r *http.Request) {
	// get body
	body, _ := ioutil.ReadAll(r.Body)

	log.Println(string(body))

	send_email(string(body))

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
		"BODY\r\n" + message +
		"\r\n")
	err := smtp.SendMail("smtp.yandex.ru:587", auth, FROM, to, msg)
	if err != nil {
		log.Print("Error: ")
		log.Fatal(err)
	} else {
		log.Println("Successfull done!")
	}
}
