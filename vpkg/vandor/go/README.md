# vandor/go

Maintainers may preview templates with `vx view ... --plan`. Application
developers use only the Taskfile interface: `task add:usecase`,
`task add:http-handler`, `task add:job`, `task add:scheduler`, and the related
`task sync:*` commands.

`add:usecase`, `add:job`, and `add:scheduler` require `name`. `add:http-handler`
requires `name`, `method`, and a full versioned `path` such as
`/api/v1/examples/{id}`. HTTP methods are limited to GET, POST, PUT, PATCH, and
DELETE. Generation fails before writing when required values are absent or a
target path already exists.
