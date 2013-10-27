git-version-proxy
=================

A HTTP Git proxy that only exposes certain versions.

Setup
-----

Build and start up the proxy:

    go intsall github.com/msiebuhr/git-version-proxy
	rehash # required by some shells
	git-version-proxy

It starts up a webserver on `127.0.0.1:8080`, which will understand Github URLs on the form

    http://127.0.0.1:8080/github.com/msiebuhr/@<commitish>/git-version-proxy
    http://127.0.0.1:8080/github.com/msiebuhr/git-version-proxy@<commitish>

For example, you can `go get` a version 0.1.0 of Etcd by doing:

    go get 127.0.0.1:8080/github.com/coreos/etcd@v0.1.0

(Currently, the `@version` can go pretty much anywhere in the URL. I'll have to
test if it breaks too many things to put it at the very end.)

Goals
-----

Make a demonstrator/experimental implementation of using VCS
tagging/versioning/binding to get specific versions of go packages.

My utopia-fantasy-goal would be for Go to support something along these lines

	import "github.com/username/project" `v1.2.3`
	import "github.com/username/project" @ "v1.2.3"
	import "github.com/username/project@v1.2.3"

Simply having `go get` put it somewhere sensible (I don't care terribly about
the particulars on how it is serialized to disk). Versioning/branch/commitish
could either be embedded in the import string or somewhere nearby, as with
struct field tags.

Personally, I would like it to understand [semver](http://semver.org/) (I find
it works well in Node.js), but I might be fine without.

And hey, It doesn't break Go1.

Limitations
-----------

 * Can't use `localhost`, because Go strongly believes hostnames should have
   dots in them. `127.0.0.1` works.
 * Only git-stuff. From github.
 * The git parser isn't well tested (it will break from time to time).
 * The syntax for `@commitish` is chosen because it was the first to come to
   mind (after a brief affair with `__commitish__`, that ended when I found out
   `@` was allowed in import paths.)
 * Won't work on windows - for development purposes, it uses a non-standard
   port. Ports are set with a `:`, which is illegal in paths in Windows.
 * Has a lot of corner-cases with subtle breakage (it's pushing the VCS/go get
   implementation quite a bit).
