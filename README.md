xml-diff Compare XML Files
================

Compare XML files after sorting attributes and fields.

A common database  vendors XML dump utility will regularly produce output where the order of the attributes changes.
This makes using the command line "diff" utility completely useless in comparing two XML files.

xml-diff will sort the attributes and values and then perform a diff on the results.


### Install

    go get -u github.com/pschlump/xml-diff

### Importing

The command line is in this directory.  The package that performs most of the work is xmllib.

    import github.com/pschlump/xml-diff/xmllib


