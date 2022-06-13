# go-sample-skill

Simple skill showing how to subscribe to and transact new data.
                    
## Files

| File                                    | Description                                    |
|-----------------------------------------|------------------------------------------------|
| [datalog/schema/commit_signature.edn]() | Datalog schema defining Commit signature facts |
| [datalog/subscription/on_push.edn]()    | Subscription for new pushes                    |
| [main.go]()                             | Go file defining `main` entrypoint             |
| [handlers.go]()                         | Go file defining the event handlers            |
