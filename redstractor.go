package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/akhilcacharya/geddit"
)

var (
	flattened = []*geddit.Comment{}
)

func main() {
	//Setup flags
	sub := flag.String("sub", "", "Subreddit to extract from")
	user := flag.String("user", "", "Reddit username")
	pass := flag.String("pass", "", "Reddit password")
	flag.Parse()

	//Kill if no parameters
	if *sub == "" || *user == "" || *pass == "" {
		usage()
		return
	}

	//Extract
	extract(*sub, *user, *pass)
}

func extract(sub, user, pass string) {
	//Log in
	session, err := geddit.NewLoginSession(
		user,
		pass,
		"Redstractor",
	)

	//Handle invalid credentials
	if err != nil {
		fmt.Println("Invalid login credentials")
		return
	}

	opts := geddit.ListingOptions{}

	submissions, err := session.SubredditSubmissions(sub, geddit.HotSubmissions, opts)

	if err != nil {
		fmt.Println("Error receiving submissions")
		return
	}

	fmt.Println("Parsing through first", len(submissions), "submissions")

	for _, submission := range submissions {
		fmt.Println("=> Post:", submission.Title)
		comments, _ := session.Comments(submission)
		//Recursively add comments in
		for _, comment := range comments {
			flattenChildren(comment)
		}
	}

	//Build up the text
	corpus := ""
	for _, comment := range flattened {
		corpus += "\n" + comment.Body + "\n"
	}

	//Save to file
	err = ioutil.WriteFile(sub+".txt", []byte(corpus), 0644)

	//Kil if error
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}

	fmt.Println("==> Saved", len(flattened), "comments to comment corpus in ", sub+".txt")
}

//Flatten out the child replies and add to the list
func flattenChildren(comment *geddit.Comment) {
	//Add to flattened
	flattened = append(flattened, comment)
	//Recursively go through all child comments
	for _, child := range comment.Replies {
		flattenChildren(child)
	}
}

//Show usage
func usage() {
	fmt.Println("\nRedstractor")
	fmt.Println("==>   Extract text corpuses from subreddits")
	fmt.Println("==> $ redstractor -sub={Subreddit Name} -user={Username} -pass={Password}")
}
