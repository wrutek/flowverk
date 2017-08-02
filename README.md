# Flowvek #

This project is to make our work on code consistent with each other.

Flowverk connects to jira, fetch all issues with status ToDo end present them
to developer. Developer can choose one of the issues and then this issue is
automatically moved to InProgress column with appropriate comment. Also
git branch is created with consistent naming convention for whole project

## Instalation ##

Copy config template file:

```bash
cp flowverk.yaml .flowverk.yaml
```

Configure all keys presented in your config file appropriate to you
project.

run go installation tool

```bash
go install flowverk
```

Move your config file to root directory of your project. After that you
are ready to use flowverk from yours project root directory
