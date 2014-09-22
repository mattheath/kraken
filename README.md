# Kraken

Kraken is a parallelised web crawler written in Go

## Usage

	* go get github.com/mattheath/kraken
	* kraken -target="http://example.com"

Kraken will output a standard XML sitemap and a richer JSON description with links and static assets per page into the current directory (output directory can be overridden, see below).

### Options

Kraken takes a number of command line flags:

	* -target="http://example.com" - The site to crawl
	* -depth=4                     - Depth of links to follow
	* -v                           - Enable verbose logging
	* -o                           - Specify output directory

## Todo

 - [ ] Limit the number of concurrent goroutines, currently this runs as fast as possible
 - [ ] Retry failed page loads with exponential backoff
 - [ ] Listen on HTTP port and serve back site description