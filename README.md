# Go Sample Skill

Simple skill showing how to subscribe to and transact new data written in Go.

## Files

| File                                                                       | Description                                    |
| -------------------------------------------------------------------------- | ---------------------------------------------- |
| [datalog/schema/commit_signature.edn](datalog/schema/commit_signature.edn) | Datalog schema defining Commit signature facts |
| [datalog/subscription/on_push.edn](datalog/subscription/on_push.edn)       | Subscription for new pushes                    |
| [datalog/subscription/on_commit_signature.edn](datalog/subscription/on_commit_signature.edn)       | Subscription for new commit signatures added to the database                    |
| [Dockerfile](Dockerfile)                                                   | Dockerfile to build the Skill runtime container image |
| [main.go](main.go)                                                         | Go file defining `main` entrypoint             |
| [handlers.go](handlers.go)                                                 | Go file defining the event handlers            |
| [types.go](main.go)                                                        | Go file defining mapping for in and outgoing structs |
| [skill.yaml](skill.yaml)                                                   | Skill descriptor (metadata and parameters etc) |

## Contributing

Contributions to this project from community members are encouraged and
appreciated. Please review the [Contributing Guidelines](CONTRIBUTING.md) for
more information. Also see the [Development](#development) section in this
document.

## Code of conduct

This project is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). You are
expected to act in accordance with this code by participating. Please report any
unacceptable behavior to code-of-conduct@atomist.com.

## Connect

Follow [@atomist][atomist-twitter] on Twitter and [The Atomist
Blog][atomist-blog].

[atomist-twitter]: https://twitter.com/atomist "Atomist on Twitter"
[atomist-blog]: https://blog.atomist.com/ "The Atomist Blog"


## Support

General support questions should be discussed in the `#help` channel in the
[Atomist community Slack workspace][slack].

If you find a problem, please create an [issue](../../issues).

## Development

You will need to install [Golang][go] to build and test this project.

[go]: https://go.dev/ "Go"

### Build and test

Install dependencies.

```
$ go get
```

Use the following command to build the skill.

```
$ go build
```

### Release

Releases are created by pushing a release [semantic version][semver] tag to the
repository, Atomist Skills take care of the rest.

To make this skill globally available, set its maturity to "stable" via the set
maturity drop-down in its Atomist Community Slack channel.

[semver]: https://semver.org/ "Semantic Version"

---

Created by [Atomist][atomist]. Need Help? [Join our Slack workspace][slack].

[atomist]: https://atomist.com/ "Atomist"
[slack]: https://join.atomist.com/ "Atomist Community Slack"

