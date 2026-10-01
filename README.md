# hellnet-worker-template

> GitHub template for Go Worker Service projects (Kafka consumers, background workers, etc.).

## Initialize from this template

After **Use this template**, clone the new repository and run:

```bash
scripts/init-from-template.sh <repo-name> [service-name]   # renames the module, imports and cmd/ (services)
scripts/setup-repo.sh                                      # repo settings, "main" ruleset and CI variable
```

Then create the `HELLNET_ACTIONS_PRIVATE_KEY` secret (the script prints the exact command) and make sure the
`hellnet-actions` GitHub App is installed on the repository.

## Usage

Click **"Use this template"** to create a new repository.

## CI/CD

| Workflow | Trigger | Description |
|---|---|---|
| `pipeline` | Push to `main` | Release → build → container image |

Uses reusable workflows from [ci-templates](https://github.com/guilhermelinosp/ci-templates).

## License

[Apache 2.0](LICENSE)
