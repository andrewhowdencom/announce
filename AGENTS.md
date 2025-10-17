# Agent Instructions

## Configuration
### Path

Where applications have configuration, they should store that configuration in a path that follows the XDG standards. The
Golang library from [ardg/xdg] is a particularly excellent implementation finding or writing to those files.

[ardg/xdg]: https://github.com/ardg/xdg

## Git
### Staging Changes

If possible, stage changes by reviewing specific changes applied in files rather than staging files directly. In
practice, use:

```bash
git add --patch ./path/to/file
```

Rather than just `git add ./path to file`.

### Commit Styles

In general, commits should follow the styleguide set out in "[What constitutes a 'Good Commit']".

#### Title

For the "title" of a commit it should be maximally 72 characters long, be descriptive and be
expressed as though the change is being applied. The commit should be useful to view with commands such as:

```
$ git log --pretty=oneline
```

Some examples include:

* Deploy changes automatically on updates to the main branch
* Update the port (8083 → 9093) used for Prometheus connections
* Introduce the widget to handle image creation

[What constitutes a 'Good Commit']: https://www.andrewhowden.com/p/anatomy-of-a-good-commit-message

#### Body

The body should include primarily the justification for the changes, rather than a description of the changes themselves.
This allows understanding why this change was made by either humans or large language models when the change is reviewed
later with `git blame` or `git log --patch`.

It should be written as a paragraph, also breaking at 72 characters long.

#### Coauthor

Where you create commits, be sure to mark yourself as a co-author of the changes, if not the primary author. The syntax for
doing so is:

```

Co-authored-by: NAME <NAME@EXAMPLE.COM>
Co-authored-by: ANOTHER-NAME <ANOTHER-NAME@EXAMPLE.COM>"
```

Complete with the line break before the "Co-authored-by" key/value pair. The format follows RFC 5322 for defining the
display name / email pair.
