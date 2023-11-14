MAN_PAGE := doc/swextreload.1


.PHONY: doc
doc: $(MAN_PAGE)

$(MAN_PAGE): doc/swextreload.1.txt
	a2x --no-xmllint --format manpage $<
