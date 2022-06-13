# go-sample-skill

Simple skill showing how to subscribe to and transact new data.

## Files

| File                                                                       | Description                                    |
|----------------------------------------------------------------------------|------------------------------------------------|
| [datalog/schema/commit_signature.edn](datalog/schema/commit_signature.edn) | Datalog schema defining Commit signature facts |
| [datalog/subscription/on_push.edn](datalog/subscription/on_push.edn)       | Subscription for new pushes                    |
| [main.go](main.go)                                                         | Go file defining `main` entrypoint             |
| [handlers.go](handlers.go)                                                 | Go file defining the event handlers            |
| [skill.yaml](skill.yaml)                                                   | Skill descriptor (metadata and parameters etc) |
| [skill.package.yaml](skill.package.yaml)                                   | Skill packaging instructions                   |
