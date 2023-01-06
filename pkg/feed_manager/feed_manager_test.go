package feed_manager

import (
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("In-memory post manager tests", func() {
	var feedManager *FeedManager

	BeforeEach(func() {
		nm, err := NewFeedManager(NewInMemoryFeedStore(), "", "")
		Ω(err).Should(BeNil())
		feedManager = nm.(*FeedManager)
		Ω(feedManager).ShouldNot(BeNil())
	})

	It("should get feed", func() {
		// No feed initially
		r := om.GetFeedRequest{
			Username: "ani5msr",
		}
		res, err := feedManager.GetFeed(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(0))

		// Add a post
		post := &om.Post{
			Url:   "http://123.com",
			Title: "123",
		}
		feedManager.OnPostAdded("ani5msr", post)
		res, err = feedManager.GetFeed(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(1))
		event := res.Events[0]
		Ω(event.EventType).Should(Equal(om.PostAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		// Update a post
		post.Title = "New Title"
		feedManager.OnPostUpdated("ani5msr", post)
		res, err = feedManager.GetFeed(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(2))
		event = res.Events[0]
		Ω(event.EventType).Should(Equal(om.PostAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[1]
		Ω(event.EventType).Should(Equal(om.PostUpdated))
		Ω(event.Url).Should(Equal("http://123.com"))

		// Delete a post
		feedManager.OnPostDeleted("ani5msr", post.Url)
		res, err = feedManager.GetFeed(r)
		Ω(err).Should(BeNil())
		Ω(res.Events).Should(HaveLen(3))
		event = res.Events[0]
		Ω(event.EventType).Should(Equal(om.PostAdded))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[1]
		Ω(event.EventType).Should(Equal(om.PostUpdated))
		Ω(event.Url).Should(Equal("http://123.com"))

		event = res.Events[2]
		Ω(event.EventType).Should(Equal(om.PostDeleted))
		Ω(event.Url).Should(Equal("http://123.com"))
	})
})
