# Research

## Reactions Implementations

### Telegram

* Criticism (on early stage)
  * No settings to manage reactions behaviour (disable quick-reactions, bring back double-tap context menu, disable long animations)
  * Inconsistency across platforms
* Implementation (Bard):
  * Centralized storage
  * Real-time updates with WebSocket connection
  * Polling mechanism for the latest reaction data for each message
* Implementation (Telegram API doc):
  * 1+ reactions from signle user
  * 2 flows: normal + recent_reactions (local cache; https://core.telegram.org/api/reactions#recent-reactions)
  * Polls reactions for visible messages every 15-30 seconds
  * Uses lottie animations (https://core.telegram.org/api/reactions)
  * Has [reactions_uniq_max](https://core.telegram.org/api/config#reactions-uniq-max) configuration (max uniq reactions for message)
  * Normal and custom emojis
  * [getMessagesReactions](https://core.telegram.org/method/messages.getMessagesReactions)
  * [Reaction](https://core.telegram.org/type/Reaction): empty, normal, custom
  * [messagePeerReaction](messagePeerReaction)
  * [sendReaction](https://core.telegram.org/method/messages.sendReaction): add_to_recent, peer, msg_id, list of reactions

### Discord

* Criticism
  * Mass reactions aren't supported ([github issue](https://github.com/discord/discord-api-docs/issues/1301))
  * Ghost reactions ([reddit](https://www.reddit.com/r/discordapp/comments/wn5ox3/ghost_spotted_in_dm/))
  * Frontend inconsistency ([github issue](https://github.com/discord/discord-api-docs/issues/6028))
* Implementation ([Disocrd JS client doc](https://discordjs.guide/popular-topics/reactions.html#reacting-to-messages)):
  * React using unicode emoji or id of custom emoji
  * Client has cache with reactions
  * Allows to listen for reactions on messages
* Implementation ([Discord API doc](https://discord.com/developers/docs/resources/channel#create-reaction)):
  * Custom emoji has name + id
  * GetReactions has after and limit params
    * after — reactions of users with ID $\ge$ given (snowflake ID is used)
    * liimt — max number of users to return (1-100, default: 25)
  * Different max unique reactions depending on server level

### Facebook

* Criticism
  * Limited set of reactions
  * Toxicity (reactions like SAD, HAHA, SHOCK affect author's psychological well-being)
* Implementation animation example: [medium post](https://medium.com/@huydotnet/implementing-facebooks-reaction-animation-9ab05460d9f7)
* Implementation ([Facebook API doc](https://developers.facebook.com/docs/graph-api/reference/v18.0/object/reactions)):
  * GetReactions by post_id -> [{user_id, user_name, reaction_type}]
    * GET /post_id/reactions?access_token=
  * GetReactions has paging with user_id (before, after)
  * reaction_type is a string

### Github

* Criticism
  * Lack of expressiveness
  * Toxicity
* Implementation ([GitHub Docs](https://docs.github.com/en/rest/reactions/reactions?apiVersion=2022-11-28#list-reactions-for-a-team-discussion-comment)):
  * Mapping content -> emoji
  * API (CRD for every entity)
    * Get reactions for post -> [{reaction_id, post_id, user_info}]
      * GET /organizations/org_id/team/team_id/discussions/discussion_id/comments/comment_id/reactions
      * Query_params: content, per_page, page
    * Add reaction for post
      * POST /organizations/org_id/team/team_id/discussions/discussion_id/comments/comment_id/reactions
      * Body scheme: {reaction_type}
    * Remove reaction for post

### Instagram

* Implementation:
  * Utf8 reactions only

### Linkedin

* Implementation ([Linkedin doc](https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/reactions-api?view=li-lms-2024-01&viewFallbackFrom=li-lms-2022-12&tabs=http))
  * Reaction schema: {reaction_type, root: urn, create_ts, last_modified_ts, delete_ts}
    * Timestamps scheme: {actor, time} , where actor is like "urn:li:person:rboDhL7Xsf", time - unix ts
  * API
    * Get reaction by actor_id, entity_id -> Reaction
      * GET /rest/reactions/(actor_id, entity_id)
    * Get reactions by actor_id, entity_id in batch -> [Reaction]
      * GET /rest/reactions?ids=[(actor_id, entity_id)]
    * Get reactions by entity_id -> {elements: [Reaction], paging: {count, start, total, links}}
      * GET /rest/reactiosn/entity_id?q=entity&sort=sort_type
      * sort_types: CHRONOLOGICAL, REVERSE_CHRONOLOGICAL (default), RELEVANCE (relevance to the viewer)
    * Add reaction by actor_id, entity_id <- {root: urn, reaction_type}
      * POST /rest/reactions?actor=actor_id
    * Delete rection by actor_id, entity_id
      * DELETE /rest/reactions/(actor_id, entity_id)
* Implementation ([Linkedin help page](https://www.linkedin.com/help/linkedin/answer/a528190))
  * Interface: hold like icon to see available reactions
  * Limited set of custom reactions

## Technology Stack

### Language

* While language such as C++ may be faster, it usually requires non-trivial dependency management that usually becomes yet an another challenge.

### Golang Framework

* All options look similar
* Research, Table 3: https://dspace.lib.uom.gr/bitstream/2159/24536/5/SochopoulosEugeniosMsc2020.pdf
* Chose Gin, becuase
  * I'm against ORM
    * "It's time for ORM retirement" — https://ieeexplore.ieee.org/document/1708563
      * UML for database modeling: https://sparxsystems.com/resources/tutorials/uml/datamodel.html
  * Most important requirements: simplicity, support, performance
  * The most popular framework (the most of stars on github)
  * Decent performance (according to https://www.techempower.com/benchmarks/#section=data-r21&l=zijocf-6bj&hw=ph&test=fortune and research)
  * Good documentation
  * Still supported with the biggest community
  * Best choice according to research table 3 (learning curve, code structure, necessary functionality, such as caching)

#### Fiber

* Benefits
  * Documentation
  * Functionality (routing, hooks, a various of middleware: CORS, Cache, limiter, etc)
  * Conciseness
* Drawbacks

#### Gin

* Benefits
  * Performance
  * Minimal memory footprint
  * Functionality (routing, middleware, etc)
  * Documentation
  * Support
* Drawbacks
  * Not suitable for large-scale applications
  * Lack of advanced features
  * Testing (encourages a global state rather than dependency injection)

#### Echo

* Benefits
  * Performance
  * Functionality (routing, middleware, etc)
  * Built-in validation support
  * Built-in support for WebSockets
  * Graceful shutdown
* Drawbacks
  * Don't cover all cases (may require additional libraries)

#### Beego

* Benefits
  * ORM
  * MVC structure
  * Functionality (caching, logging, etc)
  * Built-in internationalization support
* Drawbacks
  * ORM forces SQL-bsed design
  * Monolith app preference
  * Lack of code transparency (magic code)
  * Admin UI customization
  * Flexibility
  * Routing features

#### Revel

* Benefits
  * Full-stack
  * Modular architecture
  * Built-in support for WebSocket
  * Built-in support for reverse-proxying
  * Hot code reloading for faster development
* Drawbacks
  * Heavy framework
  * Steep learning curve

#### Buffalo

* Benefits
  * Complete developmetn ecosystem
  * Built-in support for database migrations and ORM
  * Integrated templating and asset management
  * Support
* Drawbacks
  * Less suitable for small projects
  * Steep learning curve

#### CHI

*

### Golang Libraries

* fx OR dig — for dependency management (dependency injection based application framework)
  * Reasoning: any dependency management library is fine, both from Uber
* TODO: postgres library