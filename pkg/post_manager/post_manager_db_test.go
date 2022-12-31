package post_manager

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/ani5msr/microservices-project/pkg/db_utils"
	om "github.com/ani5msr/microservices-project/pkg/object_model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DB post store tests", func() {
	var postStore *DbPostStore
	var deleteAll = func() {
		sq.Delete("posts").RunWith(postStore.db).Exec()
		sq.Delete("tags").RunWith(postStore.db).Exec()
	}
	BeforeSuite(func() {
		var err error
		dbHost, dbPort, err := db_utils.GetDbEndpoint("post_manager")
		Ω(err).Should(BeNil())

		postStore, err = NewDbPostStore(dbHost, dbPort, "postgres", "postgres")
		if err != nil {
			_, err = db_utils.RunLocalDB("postgres")
			Ω(err).Should(BeNil())
			if err != nil {
				log.Fatal(err)
			}

			postStore, err = NewDbPostStore(dbHost, dbPort, "postgres", "postgres")
			Ω(err).Should(BeNil())
			if err != nil {
				log.Fatal(err)
			}
		}

		Ω(err).Should(BeNil())
		Ω(postStore).ShouldNot(BeNil())
		Ω(postStore.db).ShouldNot(BeNil())
	})

	BeforeEach(deleteAll)
	AfterSuite(deleteAll)

	It("should add and get posts", func() {
		// No posts initially
		r := om.GetPostRequest{
			Username: "gigi",
		}
		res, err := postStore.GetPost(r)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(0))

		// Add a post
		r2 := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err = postStore.AddPost(r2)
		Ω(err).Should(BeNil())

		res, err = postStore.GetPost(r)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		post := res.Posts[0]
		Ω(post.Url).Should(Equal(r2.Url))
		Ω(post.Title).Should(Equal(r2.Title))
		Ω(post.Status).Should(Equal(om.PostStatusPending))

	})

	It("should update a post", func() {
		// Add a post
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err := postStore.AddPost(r)
		Ω(err).Should(BeNil())

		r2 := om.UpdatePostRequest{
			Username:    r.Username,
			Url:         r.Url,
			Description: "The main web site for the Go programming language",
			RemoveTags:  map[string]bool{"programming": true},
		}
		_, err = postStore.UpdatePost(r2)
		Ω(err).Should(BeNil())

		r3 := om.GetPostRequest{Username: "gigi"}
		res, err := postStore.GetPost(r3)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		post := res.Posts[0]
		Ω(post.Url).Should(Equal(r.Url))
		Ω(post.Description).Should(Equal(r2.Description))
	})

	It("should delete a post", func() {
		// Add a post
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err := postStore.AddPost(r)
		Ω(err).Should(BeNil())

		// Should have 1 post
		r2 := om.GetPostRequest{Username: "gigi"}
		res, err := postStore.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))

		// Delete the post
		err = postStore.DeletePost("gigi", r.Url)
		Ω(err).Should(BeNil())

		// There should be no more posts
		res, err = postStore.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(0))
	})

	It("should set post status", func() {
		// Add a post
		r := om.AddPostRequest{
			Username: "gigi",
			Url:      "https://golang.org/",
			Title:    "Golang",
			Tags:     map[string]bool{"programming": true},
		}
		_, err := postStore.AddPost(r)
		Ω(err).Should(BeNil())

		// Should have 1 post
		r2 := om.GetPostRequest{Username: "gigi"}
		res, err := postStore.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusPending))

		// Set post status
		err = postStore.SetPostStatus("gigi", r.Url, om.PostStatusValid)
		Ω(err).Should(BeNil())

		// The post status should be valid now instead of pending
		res, err = postStore.GetPost(r2)
		Ω(err).Should(BeNil())
		Ω(res.Posts).Should(HaveLen(1))
		Ω(res.Posts[0].Status).Should(Equal(om.PostStatusValid))

	})

})
