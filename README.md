yak
===

Yet Another Krawler

I found to some dismay that there don't seem to be many link crawlers written 
in golang so maybe this project is mistitled.  But this is in fact, a link
crawler that accepts a custom `walkFn` in the spirit of [filepath.Walk](http://golang.org/pkg/path/filepath/#Walk)
which gets each HTML node in the parse tree.  In this way, one can easily decide
what consistutes an "interesting node" for some particular definition of 
interesting.

At the core it has the following set of events:

1. Start with a single URL, add it to a queue
2. Dequeue a work item
3. Grab its content, build a parse tree of the content
4. Scan the parse tree for interesting HTML nodes.  In this case interesting consists of
    * images
    * javascript
    * external links (links that point to somewhere other than the host of the original URL in step 1)
    * stylesheet(s)
    * child links (links that are hosted on the same URL as the URL in step 1)
5. If there are any child links, enqueue them for later processing
6. Start over at step 2 until the queue is empty
7. Output each page's url and its assets, by walking the page tree from top to bottom, left to right.


Installing
----------

     $ go get https://github.com/mrallen1/yak
     $ cd $GOPATH/src/github.com/mrallen1/yak
     $ go test -cover
     $ go install
     $ cd $GOPATH/bin
     $ ./yak http://www.example.com
     
Known issues
------------

### Head of line blocking ###
The current implementation is a depth-first recursive tree walk. It suffers from excessive 
head-of-line blocking. A way to solve that would be implementing some workers to concurrently
dequeue work items, but I haven't mastered golang channels and concurrency idioms 
well enough for that yet.  That's just a matter of further practice and mentoring though.

This [example code](http://golang.org/doc/codewalk/sharemem/) shows what could be implemented.

### URLs with a `/` ###
Some times URL nodes have a `/` on the end and sometimes they don't.  It confuses
my map of visited pages so sometimes I revisit pages I have in fact already parsed.  I don't
know enough about the internals of the `net/url` package to say why this is.

### No maximum queue depth ###
At the moment, there is no maximum depth to the work queue.

Better implementations (that someone else wrote)
------------------------------------------------
He wrote two of them!

* [fetchbot](https://github.com/PuerkitoBio/fetchbot)
* [gocrawl](https://github.com/PuerkitoBio/gocrawl)
