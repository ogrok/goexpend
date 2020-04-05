# goexpend

cli budget manager focused on regular expense planning

## Usage

The application is intended for monthly budgeting.

The build script simply builds the app and moves the executable to somewhere likely to be in your $PATH. Requires a working Go environment. 

### Data Model

There are three types of expenses: **one-time**, **monthly** and **yearly**. Expenses are categorized into a given **month** and either recur monthly, annually or not at all. Expenses are named and categorized, and have budgeted (accrued) amounts as well as actual (realized) amounts.

Names and categories are user-defined.

All the data is stored as JSON, so your data is truly yours! Parse, export, manipulate as you see fit.
    
### Commands

* `goex` (no argument) - shows a report related to the current month
* `goex init` - create initial empty budget and log files in same directory
* `goex income` - edit monthly income figure (used for reporting), e.g. `goex income 3000`
* `goex add` - create new budget item in current month with syntax `-n name -a amount`, e.g. `goex add -n rent -a 1000`
    * required flags:
        * `-n name`
        * `-a amount`
    * optional flags:
        * `-c category`
        * `-d description`
        * `-i immutable` (defaults to false; no argument required)
        * `-r recurrence` (defaults to monthly)
* `goex delete` - delete budget item with ID passed as argument, e.g. `goex delete 3`
* `goex modify` - modify budget item with ID and properties passed as flags, e.g. `goex modify 4 -a 24.4`
    * `-a` - amount
    * `-c` - category
    * `-d` - description
    * `-n` - name
    * `-r` - realized amount
    * `immutable` and `recurrence` characteristics cannot be modified. To edit these, delete and recreate the budget item
    * all edits affect current month and future months 
* `goex accrue` - accrue additional amount of budget item (passed by id) for the current month, e.g. `goex accrue 4 -a 10.83`
    * does not affect future months
* `goex realize` - realize amount of budget item (passed by id) for the current month, e.g. `goex realize 4 10.83`
    * negative amounts are possible and will reduce realized amount
    * with no argument, realized amount defaults to entire remaining amount
* `goex month` has a number of subcommands related to the month state:
    * `goex month` by itself returns the current month, e.g. `2020-04`
    * `goex month close` asks for confirmation and then closes the current month, permanently logging the month and resetting to the current month
    * `goex month reset` asks for confirmation and then resets month value to current and removes all realized values. No logs are created.
        * all non-recurring items are purged and not logged
        * designed to reset app after a period of non-use, or for testing purposes
* `goex info` lists details for a specific budget item (by id), e.g. `goex info 4`
* `goex all` lists all budgeted items with details
* `goex report` lists report intended for regular viewing
* `goex purge` completely deletes all logs and then executes `goex month reset`, essentially removing all state from the application. 

## Upcoming Features
Project is a work in progress, but the above functionality currently works. Reporting is highly incomplete. I anticipate the following further development:
* Complete intended reports functionality: robustness, usefulness, summary tables, filtering by category, etc. Specifically:
    * Creating list of expenses that includes ID
    * Relating stored income figure to report results
    * Sort expenses by amount in reports, filter by category, etc.
* Suggest increase/reduction of variable expenses upon month-end close if they have been high/low for 3+ consecutive months
    * This will most likely involve income figure and result in a new command such as `goex advice` or something like that
