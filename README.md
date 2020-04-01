# goexpend

cli budget manager focused on regular expense planning

## Usage

The application is intended for monthly budgeting.

### Data Model

There are three types of expenses: **one-time**, **monthly** and **yearly**. Expenses are categorized into a given **month** and either recur monthly, annually or not at all. Expenses are named and categorized, and have budgeted (accrued) amounts as well as actual (realized) amounts.

Names and categories are user-defined.
    
### Commands

* `goex` (no argument) - shows a report related to the current month
* `goex init` - create initial empty budget and log files in same directory
* `goex add` - create new budget item in current month with syntax `-n name -a amount`, e.g. `goex add -n rent -a 1000`
    * required flags:
        * `-n name`
        * `-a amount`
    * optional flags:
        * `-c category`
        * `-d description`
        * `-m mutable` (defaults to true)
        * `-r recurrence` (defaults to monthly)
* `goex delete` - delete budget item with ID passed as argument, e.g. `goex delete 3`
* `goex modify` - modify budget item with ID and properties passed as flags, e.g. `goex modify 4 -a 24.4`
    * `-a` - amount
    * `-c` - category
    * `-d` - description
    * `-n` - name
    * `-r` - realized amount
    * `mutable` and `recurrence` characteristics cannot be modified. To edit these, delete and recreate the budget item
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
Project is WIP and very little functionality exists. Once above spec is realized, I anticipate implementing the following:
* Suggest increase/reduction of variable expenses upon month-end close if they have been high/low for 3+ consecutive months