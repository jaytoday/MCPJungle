## Contributing 💻

This document contains notes for Developers and Contributors of MCPJungle.

If you're simply a user of MCPJungle, you can skip this doc.

### Build for local testing
```bash
# Single binary for your current system
$ goreleaser build --single-target --clean --snapshot

# Test the full release assets (binaries, docker image) without publishing
goreleaser release --clean --snapshot --skip publish

# Binaries for all supported platforms
$ goreleaser release --snapshot --clean
```

### Create a new release
1. Create a Git Tag with the new version

```bash
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

2. Release
```bash
# Make sure GPG is present on your system and you have a default key which is added to Github.

# set your github access token
export GITHUB_TOKEN="<your GH token>"

goreleaser release --clean
```

This will create a new release under Releases and also make it available via Homebrew.


### Accessing SQLite
When running MCPJungle with SQLite, you can access the database using the `sqlite3` command line tool.

```bash
sqlite3 mcp.db

> .tables
> SELECT * FROM mcp_servers;
> SELECT * FROM tools;

# and so on...
```

### Accessing PostgreSQL
When running MCPJungle with docker-compose, you can access the PostgreSQL database using the `pgadmin` utility.

Open `http://localhost:5050` in your browser and log in with the Username `admin@admin.com` and Password `admin`.

Add a new DB Server with the following settings:
- Host: `db`
- Port: `5432`
- Username: `mcpjungle`
- Password: `mcpjungle`
- Database: `mcpjungle`

Then you can open up tables and run queries.