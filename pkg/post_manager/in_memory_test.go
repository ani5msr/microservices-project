package post_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("In-memory link manager tests", func() {
	var err error
	var linkManager om.PostManager
	var socialGraphManager *mockSocialGraphManager
	var eventSink *testEventsSink

	BeforeEach(func() {
		socialGraphManager = newMockSocialGraphManager([]string{"liat"})
		eventSink = newPostManagerEventsSink()
		linkManager, err = NewPostManager(NewInMemoryPostStore(),
			socialGraphManager,
			"",
			eventSink,
			10)
		Ω(err).Should(BeNil())
	})

	It("should add and get links", func() {
		// No links initially
		r := om.GetPostRequest{
			Username: "gigi",
		}
		res, err := linkManager.GetPost(r)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(0))

		// Add a link
		r2 := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err = linkManager.AddPost(r2)
		Ω(err).Should(BeNil())

		res, err = linkManager.GetPost(r)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		link := res.Posts[0]
		Ω(link.Url).Should(Equal(r2.Url))
		Ω(link.Title).Should(Equal(r2.Title))

		// Verify link manager notified the event sink about a single added event for the follower "liat"
		Ω(eventSink.addPostEvents).Should(HaveLen(1))
		Ω(eventSink.addPostEvents["liat"]).Should(HaveLen(1))
		Ω(*eventSink.addPostEvents["liat"][0]).Should(Equal(link))
		Ω(eventSink.updatePostEvents).Should(HaveLen(0))
		Ω(eventSink.deletedPostEvents).Should(HaveLen(0))
	})

	It("should update a link", func() {
		// Add a link
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err := linkManager.AddPost(r)
		Ω(err).Should(BeNil())

		r2 := om.UpdatePostRequest{
			Username:    r.Username,
			Url:         r.Url,
			Description: "The main web site for the Go programming language",
			RemoveTags:  map[string]bool{"programming": true},
		}
		err = linkManager.UpdatePost(r2)
		Ω(err).Should(BeNil())

		r3 := om.GetPostRequest{Username: "gigi"}
		res, err := linkManager.GetPost(r3)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		link := res.Posts[0]
		Ω(link.Url).Should(Equal(r.Url))
		Ω(link.Description).Should(Equal(r2.Description))
	})

	It("should delete a link", func() {
		// Add a link
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err := linkManager.AddPost(r)
		Ω(err).Should(BeNil())

		// Should have 1 link
		r2 := om.GetPostRequest{Username: "gigi"}
		res, err := linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))

		// Delete the link
		err = linkManager.DeletePost("gigi", r.Url)
		Ω(err).Should(BeNil())

		// There should be no more links
		res, err = linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(0))
	})

	It("should update link status when receiving OnPostChecked() calls", func() {
		// Add a link
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		err := linkManager.AddPost(r)
		Ω(err).Should(BeNil())

		// Should have 1 link in pending status
		r2 := om.GetPostRequest{Username: "gigi"}
		res, err := linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusPending))

		// Call OnPostChecked() with valid status on link manager (after type asserting to PostCheckerEvents)
		linkCheckSink := linkManager.(om.PostCheckerEvents)
		linkCheckSink.OnPostChecked("gigi", r.Url, om.PostStatusValid)

		// The link should have valid status
		res, err = linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusValid))

		// Call OnPostChecked() with valid status again
		linkCheckSink.OnPostChecked("gigi", r.Url, om.PostStatusValid)

		// The link should still have valid status
		res, err = linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusValid))

		// Call OnPostChecked() with invalid status
		linkCheckSink.OnPostChecked("gigi", r.Url, om.PostStatusInvalid)
		// The link should have invalid status now
		res, err = linkManager.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusInvalid))
	})

})
