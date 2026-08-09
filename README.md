# EZ2 CLI

A simple CLI tool for listing, starting, stopping, and rebooting EC2 instances.

## Why?

There are already great tools for managing AWS infrastructure — CloudFormation, Terraform, Pulumi — for quickly and reliably creating and destroying resources. There's also AWS's own CLI, which can manage pretty much any cloud resource you throw at it.

But for my use case, none of that fit well. Most of the time, my EC2 instances are already created (manually or via automation), and the *only* thing I actually need day-to-day is to start and stop them. Opening a browser and logging into the AWS console just to click "Start" or "Stop" felt like overkill, and while the AWS CLI absolutely can do this, it's a much bigger tool than I need for such a small, repetitive task.

So I built EZ2 CLI: a small, focused tool that does a few things — list, start, stop, restart, and connect to EC2 instances. Nothing more.

## Configuration

EZ2 CLI uses the AWS SDK, so it picks up credentials the same way the AWS CLI does. Set the following environment variables before running it:

```
AWS_ACCESS_KEY_ID=<your-access-key-id>
AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
AWS_REGION=<your-aws-region>
```

If you're using temporary credentials (e.g. from an assumed role or SSO), also set:

```
AWS_SESSION_TOKEN=<your-session-token>
```

Alternatively, if you already have a profile configured in `~/.aws/credentials`, you can skip the environment variables and instead set:

```
AWS_PROFILE=<your-profile-name>
```

## Usage

```
Available subcommands: list, ls, start, stop, restart, connect

ez2cli [subcommand] [args...]

ez2cli ls name=<tag:Name>
ez2cli ls id=<instance-id>
ez2cli ls state=<instance-state>

ez2cli list name=<tag:Name>
ez2cli list id=<instance-id>
ez2cli list state=<instance-state>

ez2cli start id=<instance-id> name=<instance-name> [id=<id2> name=<name2> ...]
ez2cli stop id=<instance-id> name=<instance-name> [id=<id2> name=<name2> ...]
ez2cli restart id=<instance-id> name=<instance-name> [id=<id2> name=<name2> ...]

ez2cli connect <instance-id|instance-name> <user>
```

## Note

This is a personal project I built mainly to learn Go. I wanted a project simple enough to start right away, and the fact that Go compiles to a single native binary made it a great fit for a small CLI tool like this.
