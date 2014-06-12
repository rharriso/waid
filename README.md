waid
====

What Am I doing? Command line activity tracker in go.

###Usage

```
waid [command] [options]
```


###Commands

* **start** - start a new entry
* **stop** - stop the current entry
* **add** - add an entry
* **list** - list all the entries
* **clear** - remove all tasks in the list

###Options

* **-m** - Add a message to the current task on start or stop
* **-t** - Set the time for a entry (used on add)

### Authentication

Add a .env file to the root of the project

```
#!bash
# .env
USERNAME=your_username
PASSWORD=super-secret-password

```
