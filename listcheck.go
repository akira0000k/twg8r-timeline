package main

import (
	"encoding/json"
	"context"
	"fmt"
	"strings"
	"os"
	
	g8rv2 "github.com/g8rswimmer/go-twitter/v2"
)

func listIDCheck(userID string, listid string, listname string) (returnID string) {
	fmt.Printf("userID=[%v] listid=[%v] listname=[%v]\n", userID, listid, listname)
	returnID = "0"
	if userID == "0" {
		if listid != "0" {
			return listid
		}
		fmt.Fprintln(os.Stderr, "no userid")
		return
	}
	var lists = []*g8rv2.ListObj{}
	var onetime = 100
	var pagtoken = ""
	for {
		var res *g8rv2.UserListLookupResponse
		res, err := twapi.client.UserListLookup(context.Background(), userID, g8rv2.UserListLookupOpts {
			ListFields: []g8rv2.ListField{g8rv2.ListFieldPrivate},
			MaxResults: onetime,
			PaginationToken: pagtoken,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			return
		}
		if twapi.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Meta != nil && res.Meta.ResultCount == 0 {
			break
		}
		if res.Raw.Errors != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
			}
			return
		}
		lists = append(lists, res.Raw.Lists...)
		if pagtoken = res.Meta.NextToken; pagtoken == "" {
			break
		}
	}
	if len(lists) <= 0 {
		fmt.Fprintln(os.Stderr, "no list in this user.")
		return
	}
	matchcount := 0
	for _, list := range lists {
		if listid != "0" && list.ID == listid ||
			listname != "" && strings.HasPrefix(list.Name, listname) {
			returnID = list.ID
			fmt.Fprintln(os.Stderr, "listId: ", list.ID, " Name: ", list.Name)
			matchcount += 1
		}
	}
	if matchcount == 1 {
		return returnID
	} else if matchcount > 1 {
		fmt.Fprintln(os.Stderr, "choose list id.")
	} else {
		if listid == "0" && listname == "" {
			fmt.Fprintln(os.Stderr, "need -listid or -listname.")
		} else {
			fmt.Fprintln(os.Stderr, "list id or list name unmatch.")
		}
		for _, list := range lists {
			fmt.Fprintln(os.Stderr, "listId: ", list.ID, " Name: ", list.Name)
		}
	}
	return "0"
}

// // UserListLookupResponse is the raw response with meta
// type UserListLookupResponse struct {
//  	Raw       *UserListRaw
//  	Meta      *UserListLookupMeta `json:"meta"`
//  	RateLimit *RateLimit
// }
//  
// // UserListRaw is the raw response
// type UserListRaw struct {
//  	Lists    []*ListObj       `json:"data"`
//  	Includes *ListRawIncludes `json:"includes,omitempty"`
//  	Errors   []*ErrorObj      `json:"errors,omitempty"`
// }
//  
// // UserListLookupMeta is the meta data for the lists
// type UserListLookupMeta struct {
//  	ResultCount   int    `json:"result_count"`
//  	PreviousToken string `json:"previous_token"`
//  	NextToken     string `json:"next_token"`
// }
