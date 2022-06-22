package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	g8rv2 "github.com/g8rswimmer/go-twitter/v2"
)

func name2id(screen_name string) (string, error) {
	id := ""
	res, err := twapi.client.UserNameLookup(context.Background(), []string{screen_name},
		g8rv2.UserLookupOpts{},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	if res == nil {
		fmt.Fprintln(os.Stderr, "res == nil")
		return id, err
	}
	if twapi.jsonp {
		jsonraw, _ := json.MarshalIndent(res, "", "    ")
		fmt.Println(string(jsonraw))
	}
	if res.Raw.Errors != nil {
		for _, e := range res.Raw.Errors {
			fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
		}
		return id, err
	}
	user := res.Raw.Users
	if user != nil {
		id = user[0].ID
	}
	return id, err
}

// // UserLookupResponse contains all of the information from an user lookup callout
// type UserLookupResponse struct {
//  	Raw       *UserRaw
//  	RateLimit *RateLimit
// }
//  
// // UserRaw is the raw response from the user lookup endpoint
// type UserRaw struct {
//  	Users        []*UserObj       `json:"data"`
//  	Includes     *UserRawIncludes `json:"includes,omitempty"`
//  	Errors       []*ErrorObj      `json:"errors,omitempty"`
//  	dictionaries map[string]*UserDictionary
// }
//  
// // UserObj contains Twitter user account metadata describing the referenced user
// type UserObj struct {
//  	ID              string          `json:"id"`
//  	Name            string          `json:"name"`
//  	UserName        string          `json:"username"`
//  	CreatedAt       string          `json:"created_at,omitempty"`
//  	Description     string          `json:"description,omitempty"`
//  	Entities        *EntitiesObj    `json:"entities,omitempty"`
//  	Location        string          `json:"location,omitempty"`
//  	PinnedTweetID   string          `json:"pinned_tweet_id,omitempty"`
//  	ProfileImageURL string          `json:"profile_image_url,omitempty"`
//  	Protected       bool            `json:"protected,omitempty"`
//  	PublicMetrics   *UserMetricsObj `json:"public_metrics,omitempty"`
//  	URL             string          `json:"url,omitempty"`
//  	Verified        bool            `json:"verified,omitempty"`
//  	WithHeld        *WithHeldObj    `json:"withheld,omitempty"`
// }
