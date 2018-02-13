docker-aws-info
===============

This is a quick and dirty hack written for a friend, to dump some basic
information from AWS's metadata service into a web-page, as a tiny Docker
image.

Source code hereby placed into the public domain.

No warranty.  You get to keep all the pieces and shards if it breaks.

There's an integration setup to auto-build this on Quay, but they don't
support ARG in Dockerfile so it doesn't work as present.  If you see an image
available as `quay.io/pennocktech/docker-aws-info` but this paragraph still
says it's broken, please report that to us!  (It means Quay now handle newer
Dockerfile features).

A manual version was setup on Docker Hub, using a different name.
(See above re quick and dirty hack).  `pennocktech/aws-basic-info` is manually
built and pushed.
