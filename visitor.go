package main

import (
	"net/http"
	"strings"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
)

// The Visitor takes a link and scrapes the job description, title and all the text nodes on the page.
// If the body matches the atleast one of the provided keywords, the Visitor sends the webiste link back to
// the Finder.
type Visitor struct {
	FPID *actor.PID
	Link string
	Meta *JobRequest
}

func NewVisitor(link string, mpid *actor.PID, meta *JobRequest) actor.Producer {
	return func() actor.Receiver {
		return &Visitor{
			FPID: mpid,
			Link: link,
			Meta: meta,
		}
	}
}

func (v *Visitor) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		if v.filterJobPosting() {
			ctx.Send(v.FPID, &FilteredJobPost{
				Link: v.Link,
			})
		}
		ctx.Engine().Poison(ctx.PID())
	case actor.Stopped:
		ctx.Send(v.FPID, &KillVisitor{
			PID: ctx.PID(),
		})
	}
}

// Needs access to Meta object -> method instead of function
func (v *Visitor) jobListingKeywordMatcher(val string) bool {
	var valid bool
	if len(v.Meta.Keywords) == 0 {
		return true
	}
	for _, k := range v.Meta.Keywords {
		if strings.Contains(strings.ToLower(val), strings.ToLower(k)) {
			valid = true
		}
	}
	return valid
}

func (v *Visitor) filterJobPosting() bool {
	var (
		f func(*html.Node) bool
	)

	res, err := http.Get(v.Link)
	if err != nil {
		return false
	}
	
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return false
	}

	f = func(n *html.Node) bool {
        if n.Type == html.TextNode {
            if v.jobListingKeywordMatcher(strings.TrimSpace(n.Data)) {
                return true
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            if f(c) {
                return true
            }
        }
        return false
    }
	return f(doc)
}

