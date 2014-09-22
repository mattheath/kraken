# Kraken

Kraken is a parallelised web crawler written in Go

## Usage

	* go get github.com/mattheath/kraken
	* kraken -target="http://example.com"

### Options

Kraken takes a number of command line flags:

	* -target="http://example.com" - The site to crawl
	* -depth=4                     - Depth of links to follow
	* -v                           - Enable verbose logging
