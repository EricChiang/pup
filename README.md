# pup

pup is a command line tool for processing HTML. It reads from stdin,
prints to stdout, and allows the user to filter parts of the page using
[CSS selectors](https://developer.mozilla.org/en-US/docs/Web/Guide/CSS/Getting_started/Selectors).

Inspired by [jq](http://stedolan.github.io/jq/), pup aims to be a
fast and flexible way of exploring HTML from the terminal.

## Install

Direct downloads are available through the [releases page](https://github.com/EricChiang/pup/releases/latest).

If you have Go installed on your computer just run `go install`.

    go install github.com/ericchiang/pup@latest

If you're on OS X, use [Homebrew](http://brew.sh/) to install (no Go required).

    brew install https://raw.githubusercontent.com/EricChiang/pup/master/pup.rb

## Quick start

```bash
$ curl -s https://news.ycombinator.com/
```

Ew, HTML. Let's run that through some pup selectors:

```bash
$ curl -s https://news.ycombinator.com/ | pup 'table table tr:nth-last-of-type(n+2) td.title a'
```

Okay, how about only the links?

```bash
$ curl -s https://news.ycombinator.com/ | pup 'table table tr:nth-last-of-type(n+2) td.title a attr{href}'
```

Even better, let's grab the titles too:

```bash
$ curl -s https://news.ycombinator.com/ | pup 'table table tr:nth-last-of-type(n+2) td.title a json{}'
```

## Basic Usage

```bash
$ cat index.html | pup [flags] '[selectors] [display function]'
```

## Examples

Download a webpage with wget.

```bash
$ wget http://en.wikipedia.org/wiki/Robots_exclusion_standard -O robots.html
```

#### Clean and indent

By default pup will fill in missing tags and properly indent the page.

```bash
$ cat robots.html
# nasty looking HTML
$ cat robots.html | pup --color
# cleaned, indented, and colorful HTML
```

#### Filter by tag

```bash
$ cat robots.html | pup 'title'
<title>
 Robots exclusion standard - Wikipedia, the free encyclopedia
</title>
```

#### Filter by id

```bash
$ cat robots.html | pup 'span#See_also'
<span class="mw-headline" id="See_also">
 See also
</span>
```

#### Filter by attribute

```bash
$ cat robots.html | pup 'th[scope="row"]'
<th scope="row" class="navbox-group">
 Exclusion standards
</th>
<th scope="row" class="navbox-group">
 Related marketing topics
</th>
<th scope="row" class="navbox-group">
 Search marketing related topics
</th>
<th scope="row" class="navbox-group">
 Search engine spam
</th>
<th scope="row" class="navbox-group">
 Linking
</th>
<th scope="row" class="navbox-group">
 People
</th>
<th scope="row" class="navbox-group">
 Other
</th>
```

#### Pseudo Classes

CSS selectors have a group of specifiers called ["pseudo classes"](
https://developer.mozilla.org/en-US/docs/Web/CSS/Pseudo-classes)  which are pretty
cool. pup implements a majority of the relevant ones them.

Here are some examples.

```bash
$ cat robots.html | pup 'a[rel]:empty'
<a rel="license" href="//creativecommons.org/licenses/by-sa/3.0/" style="display:none;">
</a>
```

```bash
$ cat robots.html | pup ':contains("History")'
<span class="toctext">
 History
</span>
<span class="mw-headline" id="History">
 History
</span>
```

```bash
$ cat robots.html | pup ':parent-of([action="edit"])'
<span class="wb-langlinks-edit wb-langlinks-link">
 <a action="edit" href="//www.wikidata.org/wiki/Q80776#sitelinks-wikipedia" text="Edit links" title="Edit interlanguage links" class="wbc-editpage">
  Edit links
 </a>
</span>
```

For a complete list, view the [implemented selectors](#implemented-selectors)
section.


#### `+`, `>`, and `,`

These are intermediate characters that declare special instructions. For
instance, a comma `,` allows pup to specify multiple groups of selectors.

```bash
$ cat robots.html | pup 'title, h1 span[dir="auto"]'
<title>
 Robots exclusion standard - Wikipedia, the free encyclopedia
</title>
<span dir="auto">
 Robots exclusion standard
</span>
```

#### Chain selectors together

When combining selectors, the HTML nodes selected by the previous selector will
be passed to the next ones.

```bash
$ cat robots.html | pup 'h1#firstHeading'
<h1 id="firstHeading" class="firstHeading" lang="en">
 <span dir="auto">
  Robots exclusion standard
 </span>
</h1>
```

```bash
$ cat robots.html | pup 'h1#firstHeading span'
<span dir="auto">
 Robots exclusion standard
</span>
```

## Implemented Selectors

For further examples of these selectors head over to [MDN](
https://developer.mozilla.org/en-US/docs/Web/CSS/Reference).

```bash
pup '.class'
pup '#id'
pup 'element'
pup 'selector + selector'
pup 'selector > selector'
pup '[attribute]'
pup '[attribute="value"]'
pup '[attribute*="value"]'
pup '[attribute~="value"]'
pup '[attribute^="value"]'
pup '[attribute$="value"]'
pup ':empty'
pup ':first-child'
pup ':first-of-type'
pup ':last-child'
pup ':last-of-type'
pup ':only-child'
pup ':only-of-type'
pup ':contains("text")'
pup ':nth-child(n)'
pup ':nth-of-type(n)'
pup ':nth-last-child(n)'
pup ':nth-last-of-type(n)'
pup ':not(selector)'
pup ':parent-of(selector)'
```

You can mix and match selectors as you wish.

```bash
cat index.html | pup 'element#id[attribute="value"]:first-of-type'
```

## Display Functions

Non-HTML selectors which effect the output type are implemented as functions
which can be provided as a final argument.

#### `text{}`

Print all text from selected nodes and children in depth first order.

```bash
$ cat robots.html | pup '.mw-headline text{}'
History
About the standard
Disadvantages
Alternatives
Examples
Nonstandard extensions
Crawl-delay directive
Allow directive
Sitemap
Host
Universal "*" match
Meta tags and headers
See also
References
External links
```

#### `attr{attrkey}`

Print the values of all attributes with a given key from all selected nodes.

```bash
$ cat robots.html | pup '.catlinks div attr{id}'
mw-normal-catlinks
mw-hidden-catlinks
```

#### `json{}`

Print HTML as JSON.

```bash
$ cat robots.html  | pup 'div#p-namespaces a'
<a href="/wiki/Robots_exclusion_standard" title="View the content page [c]" accesskey="c">
 Article
</a>
<a href="/wiki/Talk:Robots_exclusion_standard" title="Discussion about the content page [t]" accesskey="t">
 Talk
</a>
```

```bash
$ cat robots.html | pup 'div#p-namespaces a json{}'
[
 {
  "accesskey": "c",
  "href": "/wiki/Robots_exclusion_standard",
  "tag": "a",
  "text": "Article",
  "title": "View the content page [c]"
 },
 {
  "accesskey": "t",
  "href": "/wiki/Talk:Robots_exclusion_standard",
  "tag": "a",
  "text": "Talk",
  "title": "Discussion about the content page [t]"
 }
]
```

Use the `-i` / `--indent` flag to control the intent level.

```bash
$ cat robots.html | pup -i 4 'div#p-namespaces a json{}'
[
    {
        "accesskey": "c",
        "href": "/wiki/Robots_exclusion_standard",
        "tag": "a",
        "text": "Article",
        "title": "View the content page [c]"
    },
    {
        "accesskey": "t",
        "href": "/wiki/Talk:Robots_exclusion_standard",
        "tag": "a",
        "text": "Talk",
        "title": "Discussion about the content page [t]"
    }
]
```

If the selectors only return one element the results will be printed as a JSON
object, not a list.

```bash
$ cat robots.html  | pup --indent 4 'title json{}'
{
    "tag": "title",
    "text": "Robots exclusion standard - Wikipedia, the free encyclopedia"
}
```

Because there is no universal standard for converting HTML/XML to JSON, a
method has been chosen which hopefully fits. The goal is simply to get the
output of pup into a more consumable format.

## Flags

Run `pup --help` for a list of further options
