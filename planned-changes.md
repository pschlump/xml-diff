
Tree Comparison
------------------------

Q: Why read-XML, process it, then serialize back to a document and do the diff.    
What about doing the comparison (diff) on the trees themselves?

A: This is a question of what changed in the XML file and how to report it to
the user.  For Example

Original:

```xml
	<outside>
		<lev1 attr="orig">		
			<lev2>		
				<inner>			
					Some Data
					Some Data
					Some Data
				</inner>		
			</lev2>		
		</lev1>		
	</outside>
```

Changed:


```xml
	<outside>
		<lev1 attr="modified">		
			<lev2>		
				<inner>		
					Some Data
					Some Data
					Some Data
				</inner>		
			</lev2>		
		</lev1>		
	</outside>
```

From a tree perspective do you report that all of `lev1` has changed?   Do you report that 
only the attribute on the tag `lev1` has changed?  Both are correct in some contexts.

Ordering of attributes presents an additional problem.  For example


Original:

```xml
	<outside>
		<xyzzy rowId="100" />		
		<mike rowId="1102" />		
		<alpha rowId="1102" />		
	</outside>
```

Changed:

```xml
	<outside>
		<mike rowId="1102" />		
		<alpha rowId="1102" />		
		<xyzzy rowId="100" />		
	</outside>
```

If you sort the data inside the `outside` tag based on tag name, so that you get

```xml
	<outside>
		<alpha rowId="1102" />		
		<mike rowId="1102" />		
		<xyzzy rowId="100" />		
	</outside>
```

Then the two chunks of XML are the same and you can report this.  The question is how
you show the user what the sorted XML is?   

Serializing the XML allows for the tool to output the modified/sorted XML and interact
with other tools.  For example, you could take the modified XML and send it through
command line `diff` with appropriate options to create a patch file.  In some contexts
this would be a very useful ability.  If you are updating a database and the data
arrives in XML then you can easily translate the patch into SQL inserts and updates.



Planned Changes
------------------------

### Some Simple Data Maping

I plan on implementing a system that allows you to take attributes and data values and specify swapping them.  This will be a tree modification before the XML is re-serialized.
Given a confirmation that says take `lev1` `attr` and convert it into data:

```xml
	<outside>
		<lev1 attr="orig" />		
	</outside>
```

Should match with:

```xml
	<outside>
		<lev1>
			<attr>orig</attr>
		</lev1>		
	</outside>
```

A configuration that works in the opposite director would also be useful so that:

```xml
	<outside>
		<lev1>
			<attr>orig</attr>
		</lev1>		
	</outside>
```

Should match with:

```xml
	<outside>
		<lev1 attr="orig" />		
	</outside>
```

### Configuration on what gets sorted

Sometimes the ordering of tags is important in XML.  For example in SVG this determines the order that items are drawn on the canvas.
Other times the order is not relevant.    A configuration that sys inside tag `yyy` you should or should not sort would be useful.


