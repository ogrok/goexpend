# goexpend

cli budget manager focused on regular expense planning

## Usage

The application is intended for monthly budgeting.

### Data Model

There are three types of expenses: **one-time**, **monthly** and **annual**. Expenses are categorized into a given **month** and either recur monthly, annually or not at all. Expenses are named and categorized, and have budgeted (accrued) amounts as well as actual (realized) amounts.

Names and categories are user-defined.

To summarize:

* `budgetItem`
    * ID `int` (auto-generated, 1-indexed)
    * Name `string`
    * Category - `string`
    * Month - `int`
        * for `monthly` recurrence, set to `0`
        * for `yearly` recurrence, set to month of recurrence
        * for `none` recurrence, set to month of occurrence
    * Recurrence - `string`
        * `monthly`, `yearly`, `none`
    * accrued - `float64`
    * realized - `float64`
    * mutable - `bool` (if true, expense can be trivially avoided / canceled, like a subscription service; if false, it can't, like a mortgage)
        
* `month`
    * `month` - `int`
    * `year` - `int`
    * `open` - `bool` (indicates whether month is currently open for edits)
    
### Commands

* `goex` (no argument) - shows a report related to the current month
* `goex init` - create initial empty budget and config files in same directory
* `goex add` - create new budget item in current month with syntax `name amount`
    * optional flags:
        * `-c category`
        * `-d description`
        * `-m mutable` (defaults to true)
        * `-r recurrence` (defaults to monthly)
* `goex delete` - delete budget item with ID passed as argument, e.g. `goex delete 3`
* `goex modify` - modify budget item with ID passed as argument and properties passed as flags:
    * `-a` / `--accrued` / `--amount` - amount
    * `-c` / `--category` - category
    * `-m` / `--month` - month
    * `--mutable` - mutable
    * `-n` / `--name` - name
    * `-r` / `--realized` - realized amount
    * `--recur` - recurrence
    * `-y` / `--year` - year
* `goex realize` - realize amount of budget item (passed by id) for the current month, e.g. `goex realize 4 10.83`
    * negative amounts are possible and will reduce realized amount
    * with no argument, realized amount defaults to entire remaining amount
* `goex month` has a number of subcommands related to the month state:
    * `goex month` by itself returns the current month, e.g. `2020-04`
    * `goex month 2020 4` or `goex month change 2020 4` changes the active month to the specified month (year is also required)
    * `goex month open` opens the current month, allowing changes to realized amounts until closed. Months default to open
    * `goex month close` closes the current month, preventing changes to realized amounts until month is reopened
    * `goex month reset` asks for confirmation and then returns realized amounts for the current month to 0.
* `goex all` lists all budgeted items with details

## Upcoming Features
Project is WIP and very little functionality exists. Once above spec is realized, I anticipate implementing the following:
* Suggest increase/reduction of variable expenses upon month-end close if they have been high/low for 3+ consecutive months