# FFX Blue Technical Test

This is my response to the technical test https://ffxblue.github.io/interview-tests/test/article-api/

## Deliverables

### 1. Source Code

Please see other files in repository

### 2. Setup/Installation Instructions

*Requirements:*
* Go 1.15

Run:
```
make run
```


Tests:
```
make test
```

### 3. Description of Solution

Before diving into the detail I'm going to take you through my thought process of some of the big architectural decisions I made. What other options I considered and why I didn’t choose them.

First, language. I picked Go because while I know other languages, it is a Golang developer role I’m applying for :)

Next, I considered how I would want this code would run. I considered 2 main options:
AWS Lambdas backing an AWS API Gateway
A more traditional container-based service

I have more familiarity with AWS Lambas so I was tempted to go down that path. But given the size of the organisation this test is for, I can imagine that if this was a real-world case then Lambdas would probably be expensive compared to a container-based service. There were these additional befits:
Easier to have something running on a local machine
I would learn something new
So in the end I decided to build something that would set better on a container.

Given that decision, I next thought about frameworks. I did a bit of Googling and settled on the [Gin Framework](https://gin-gonic.com/). It seemed pretty straightforward and so I familiarized myself with some examples.

Then came choosing a data store. Once again I asked myself these two questions:
What would be a good choice if this was a real-world example, with the kind of loads that the organisation the test is for experiences?
What will be easy to get running on someone’s laptop?

The best answers to these questions separately are not the same. To the first one, I think AWS DynamoDB, may be paired with AWS S3. It has, I think, the right mix of functionality, performance and cost (both in terms of AWS costs and time administrating it).

But running a version of AWS DynamoDB locally on someone’s laptop is notoriously buggy. I’ve yet to meet someone that has done so and thought it was fun :)

I seriously considered these options:
Including in the Setup/Installation Instructions using Docker Compose to MySQL instance running
Doing something similar with a NoSQL database engine
Using some in-memory relational datastore
Using some in-memory NoSQL solution, like [go-cache](https://github.com/patrickmn/go-cache)

I dismissed all the ones involving a relation datastore because it’s not what I would pick for a production workload I would expect for the size of the organisation this test.

In the end, I decided to just create a quick and dirty in-memory store myself. Not generally my favoured option, I 
don’t like [NIH](https://en.wikipedia.org/wiki/Not_invented_here). But it allowed me to very easily structure the 
rest of my solution as if I had used DynamoDB. And it allowed me to demonstrate some knowledge/awareness of how to 
use the sync package [sync](https://pkg.go.dev/sync).

The last big questions I considered were: How much to treat this example as a service with:
A lot larger, feature wise?
Deployed to production and handling customer workloads?

That is, pretend this was not just a service that only needs to meet specs provided and is not deployed to Production. A strange question I know given the other things I considered above, but still worthwhile to ask.

In the end, I decided to [YAGNI](https://en.wikipedia.org/wiki/You_aren't_gonna_need_it) principle. Build what I 
needed to meet the needs of the service in specifications and also not worry about things like logging, resilience, etc. My reason for the former is because YAGNI is the practice I like to try to follow. Build a solution for the problem at hand, not some potential problem at some point in the future. And the latter (logging, resilience, etc) because I think my years of DevOps/SRE work I think demonstrates that I do in fact care about and know about these things.

#### Run through the code

OK, so now that I’ve covered the things I thought about before I started coding, let me cover what I ended up creating :)

At the core of the service are models around to concepts:
Articles
Tags

Both of these are in the `/models` directory.

There are two “services” related to these 2 domain concepts:
article_service
tag_service

These are under `/services`. They take care of persisting and retrieving data to/from the data store.

Under `/router` is the code to configure the router and the handlers for the 3 endpoints exposed.

`/pkg` has 2 parts:
app
util

`/pkg/app` has some functions that help with creating some consistency between microservices. In practice, this is the kind of code that I would be centralised in an organisation as a package that was consumed by all services in that organisation.

`/pkg/util` is just some date and map/slice manipulation functions. Once again things that would probably be reused across multiple code repositories in an organisation.

#### Testing

Because the service is so simple having a few tests on the endpoints. Indeed 2 test files, one on the routes and another on one of the functions in `/pkg/util` covers 98% of lines. I understand that lines covered does not mean the logic is covered. With a more complex service, I would of course be making more use of unit tests to cover specific parts of logic in more detail.

#### What I did change about the specs for the service
I added a `/v1` prefix to the endpoints. I didn’t do much in my code to support the adding of new versions of the endpoints. But I’ve found it good to assume that at some later point new versions of the endpoints will be required. This breaks with YAGNI but comes at very little cost. And failing to do it can make adding it later a real pain.

#### What I would have changed about the service if I could have changed the specs.
I would have made the service idempotent. That is instead of exposing a POST method for `/articles` exposes a PUT method. Making services idempotent is a huge step to making a microserice architecture resilient.

### 4. Assumptions

I found a bug, at one point in the specs the endpoint is `/tag/{tagName}/{date}` but all other points it is `/tags/{tagName}/{date}`. I used `/tags/{tagName}/{date}`.

It’s unclear if the payload for POST `/articles` contains an ID or not. I assumed so.

I really didn’t understand this line:

> The `count` field shows the number of tags for the tag for that day.

I just used the number of articles found, but really wasn’t sure if I was close to what was required. The example result didn’t help me.


### 5. What I thought of the test

I have to say I really enjoyed doing this. I was a bit daunted at first, it was a lot more involved than other technical tests I’ve had to do. But it also meant that I was forced to think about problems I’ve not yet tackled and so I learnt some things along the way.

### 6. What I would have done if I had more time

There are a couple of TODOs in the comment. Things I want to refactor and/or explain.

I'm not happy with the cyclomatic complexity of some of the functions.

100% would have added more about making the service "production ready". Like structured logging and retries for 
accessing the datastore.
