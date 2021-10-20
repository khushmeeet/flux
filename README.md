# Flux - Static Site Generator

Flux is an opinionated Static Site Generator written in Go programming language. It works as a command line utility. **(In Development, highly unstable)**

**_Project on hold_**

### Features (v1)

- Support for front matter
- Syntax highlighting
- Support for Go templates
- Support for Markdown files
- Live reloading with http server
- Support for Pug (stretch)
- Support for CSS language extensions (stretch)

### Folder Structure

```other
/my_blog
	/templates
		/partials
	/static
	/pages
		post_title.md
	/_site
```

### Markdown Front Matter

```other
---
title: Post Heading
short: Short post intro
date: 24th Jan 2021
template: index.html
---
```

### Configuration File

This file will contain data related to site metadata. This data is available throughout the project and can contain any number of key value pairs.

```other
site_title: Blog title
email: khushmeet@hey.com
twitter_username: @khushmeeet
github_username: @khushmeet
...
```

### Binary functionalities

- `flux help`
- `flux build`
- `flux serve`
- `flux init`
- `flux clean`

