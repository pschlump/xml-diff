xml-diff Compare XML Files
================

Compare XML files after sorting attributes and fields.

A common database vendor's XML dump utility will regularly produce output where the order of the attributes changes.
This makes using the command line `diff` utility completely useless in comparing two XML files.

xml-diff will sort the attributes and values and then perform a `diff` on the results.


## Install

    go get -u github.com/pschlump/xml-diff

## To Build

The command line utility is in the top level directory.

	cd ~/.../github.com/pschlump/xml-diff
	go build
	cp xml-diff ~/bin

You can run tests on the command line with

	make test

## Importing Package

The command line is in this directory.  The package that performs most of the work is `xmllib`.

    import "github.com/pschlump/xml-diff/xmllib"

## Example

For XML inputs

```xml
<?xml version="1.0" encoding="UTF-8"?>
<ConnectedApp xmlns="http://soap.sforce.com/2006/04/metadata">
	<contactEmail>foo@example.org</contactEmail>
	<label>WooCommerce</label>
	<oauthConfig>
		<scopes>Basic</scopes>
		<scopes>Api</scopes>
		<scopes>Web</scopes>
		<scopes>Full</scopes>
		<callbackUrl>https://login.salesforce.com/services/oauth2/callback</callbackUrl>
		<consumerKey>CLIENTID</consumerKey>
	</oauthConfig>
</ConnectedApp>
```

and 

```xml
<?xml version="1.0" encoding="UTF-8"?>
<ConnectedApp xmlns="http://soap.sforce.com/2006/04/metadata">
	<contactEmail>foo@example.org</contactEmail>
	<label>WooCommerce</label>
	<oauthConfig>
		<callbackUrl>https://login.salesforce.com/services/oauth2/callback</callbackUrl>
		<consumerKey>OTHER</consumerKey>
		<scopes>Full</scopes>
		<scopes>Basic</scopes>
	</oauthConfig>
</ConnectedApp>
```

You can run:

	./xml-diff -l ./testdata/left.xml -r ./testdata/right.xml 

The output is:

![Output From Diff](https://github.com/pschlump/xml-diff/raw/master/out/test01.png "Output from xml-diff")

If you add the `-byLine` flag the diff will be shown by lines.

	./xml-diff -l ./testdata/left.xml -r ./testdata/right.xml -byLine

The output is:

![Output From Diff](https://github.com/pschlump/xml-diff/raw/master/out/test02.png "Output from xml-diff with byLine flag")


### Algorithm

The diff is based on the [Myers](https://neil.fraser.name/software/diff_match_patch/myers.pdf) algorithm.  This is the most common
approach to comparing differences between files.  

An alternative approach would be to perform the difference on the XML node-tree in memory.   Because of my plan to be able to move
attributes to values and back this is an undesirable way to express the differences.

## Performance

### Memory

It takes about 4.2 times the heap size as the size of the XML file to run.  This means that if you have a 128 Mb of memory you should be
able to xml-diff files of up to 30 Mb in size.

### Performance

The XML read/parse and generate will run about 10 MB of XML in a second.  100ms should compare about 1 MB of XML.  Performance is
heavily dependent on how much data has to be sorted.   If there are lots of XML nodes that have to be built and then sorted it will
take longer to process.

## TODO

1. Partial support is in place for moving attributes to values or values back to attributes. The same database vendor seems to arbitrarily swap these in its XML dump.
2. Better documentation.

test

