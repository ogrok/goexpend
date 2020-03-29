# goexpend

cli budget manager focused on regular expense planning

## Usage

The application is intended for monthly budgeting.

### Data Model

There are three types of expenses: **one-time**, **monthly** and **annual**. Expenses are categorized into a given **month** and either recur monthly, annually or not at all. Expenses are named and categorized, and have budgeted (accrued) amounts as well as actual (realized) amounts.

Names and categories are user-defined.
    
### Commands

* `goex` (no argument) - shows a report related to the current month
* `goex init` - create initial empty budget and log files in same directory
* `goex add` - create new budget item in current month with syntax `name amount`, e.g. `goex add rent 1000`
    * optional flags:
        * `-c category`
        * `-d description`
        * `-m mutable` (defaults to true)
        * `-r recurrence` (defaults to monthly)
* `goex delete` - delete budget item with ID passed as argument, e.g. `goex delete 3`
* `goex modify` - modify budget item with ID passed as argument and properties passed as flags:
    * `-a` / `--accrued` / `--amount` - amount
    * `-c` / `--category` - category
    * `--mutable` - mutable
    * `-n` / `--name` - name
    * `-r` / `--realized` - realized amount
    * `--recur` - recurrence
* `goex accrue` - accrue additional amount of budget item (passed by id) for the current month, e.g. `goex accrue 4 10.83`
    * does not affect future months
* `goex realize` - realize amount of budget item (passed by id) for the current month, e.g. `goex realize 4 10.83`
    * negative amounts are possible and will reduce realized amount
    * with no argument, realized amount defaults to entire remaining amount
* `goex month` has a number of subcommands related to the month state:
    * `goex month` by itself returns the current month, e.g. `2020-04`
    * `goex month close` asks for confirmation and then closes the current month, permanently logging the month and resetting to the current month
    * `goex month reset` asks for confirmation and then updates current month based on actual current time. No logs are created.
        * designed to reset app after a period of non-use, or for testing purposes
* `goex all` lists all budgeted items with details
* `goex purge` completely deletes all logs and then executes `goex month reset`, essentially removing all state from the application. 

## Upcoming Features
Project is WIP and very little functionality exists. Once above spec is realized, I anticipate implementing the following:
* Suggest increase/reduction of variable expenses upon month-end close if they have been high/low for 3+ consecutive months