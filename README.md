# goexpend

cli budget manager focused on regular expense planning

## Usage

The application is intended for monthly budgeting.

### Data Model

There are three types of expenses: **one-time**, **monthly** and **annual**. Expenses are categorized into a given **month** and either recur monthly, annually or not at all. Expenses are named and categorized, and have budgeted (accrued) amounts as well as actual (realized) amounts.

Names and categories are user-defined.

To summarize:

* `budgetItem`
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