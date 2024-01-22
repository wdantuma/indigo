package rules

import (
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/automod"
	"github.com/bluesky-social/indigo/automod/keyword"
)

// looks for specific hashtags from known lists
func BadHashtagsPostRule(c *automod.RecordContext, post *appbsky.FeedPost) error {
	for _, tag := range ExtractHashtagsPost(post) {
		tag = NormalizeHashtag(tag)
		if c.InSet("bad-hashtags", tag) || c.InSet("bad-words", tag) || c.InSet("worst-words", tag) {
			c.AddRecordFlag("bad-hashtag")
			c.Notify("slack")
			break
		}
		word := keyword.SlugContainsExplicitSlur(keyword.Slugify(tag))
		if word != "" {
			c.AddAccountFlag("bad-hashtag")
		}
	}
	return nil
}

var _ automod.PostRuleFunc = BadHashtagsPostRule

// if a post is "almost all" hashtags, it might be a form of search spam
func TooManyHashtagsPostRule(c *automod.RecordContext, post *appbsky.FeedPost) error {
	tags := ExtractHashtagsPost(post)
	tagChars := 0
	for _, tag := range tags {
		tagChars += len(tag)
	}
	tagTextRatio := float64(tagChars) / float64(len(post.Text))
	// if there is an image, allow some more tags
	if len(tags) > 4 && tagTextRatio > 0.6 && post.Embed.EmbedImages == nil {
		c.AddRecordFlag("many-hashtags")
		c.Notify("slack")
	} else if len(tags) > 7 && tagTextRatio > 0.8 {
		c.AddRecordFlag("many-hashtags")
		c.Notify("slack")
	}
	return nil
}

var _ automod.PostRuleFunc = TooManyHashtagsPostRule
